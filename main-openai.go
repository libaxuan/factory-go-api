package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// 配置结构体
type Config struct {
	Port            string
	AnthropicTarget string
	FactoryAPIKey   string // 源头 Factory API Key（用于调用上游 API）
	ProxyAPIKey     string // 对外代理 API Key（客户端使用此 Key）
}

// 默认配置
var config = Config{
	Port:            getEnv("PORT", "8000"),
	AnthropicTarget: getEnv("ANTHROPIC_TARGET_URL", "https://gibuoilncyzqebelqjqz.supabase.co/functions/v1/smooth-handler/https://app.factory.ai/api/llm/a/v1/messages"),
	FactoryAPIKey:   os.Getenv("FACTORY_API_KEY"),   // 必须配置
	ProxyAPIKey:     os.Getenv("PROXY_API_KEY"),     // 必须配置
}

var startTime = time.Now()

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 响应记录器
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

// OpenAI格式转Anthropic格式
func convertOpenAIToAnthropic(openaiBody map[string]interface{}) map[string]interface{} {
	anthropicBody := make(map[string]interface{})

	// 转换model
	if model, ok := openaiBody["model"].(string); ok {
		anthropicBody["model"] = model
	}

	// 转换messages
	if messages, ok := openaiBody["messages"].([]interface{}); ok {
		anthropicMessages := make([]map[string]interface{}, 0)
		var systemPrompts []string

		for _, msg := range messages {
			if msgMap, ok := msg.(map[string]interface{}); ok {
				role, _ := msgMap["role"].(string)
				content, _ := msgMap["content"].(string)

				if role == "system" {
					systemPrompts = append(systemPrompts, content)
				} else if role == "user" || role == "assistant" {
					anthropicMessages = append(anthropicMessages, map[string]interface{}{
						"role":    role,
						"content": content,
					})
				}
			}
		}

		anthropicBody["messages"] = anthropicMessages

		// 设置system字段
		systemBlocks := []map[string]interface{}{
			{"type": "text", "text": "You are Droid, an AI software engineering agent built by Factory."},
		}
		for _, sp := range systemPrompts {
			systemBlocks = append(systemBlocks, map[string]interface{}{
				"type": "text",
				"text": sp,
			})
		}
		anthropicBody["system"] = systemBlocks
	}

	// 转换max_tokens
	if maxTokens, ok := openaiBody["max_tokens"].(float64); ok {
		anthropicBody["max_tokens"] = int(maxTokens)
	} else {
		anthropicBody["max_tokens"] = 1024
	}

	// 转换temperature
	if temp, ok := openaiBody["temperature"].(float64); ok {
		anthropicBody["temperature"] = temp
	}

	// 转换stream
	if stream, ok := openaiBody["stream"].(bool); ok {
		anthropicBody["stream"] = stream
	}

	return anthropicBody
}

// Anthropic响应转OpenAI格式
func convertAnthropicToOpenAI(anthropicResp map[string]interface{}) map[string]interface{} {
	openaiResp := map[string]interface{}{
		"id":      anthropicResp["id"],
		"object":  "chat.completion",
		"created": time.Now().Unix(),
		"model":   anthropicResp["model"],
		"choices": []map[string]interface{}{
			{
				"index": 0,
				"message": map[string]interface{}{
					"role":    "assistant",
					"content": "",
				},
				"finish_reason": "stop",
			},
		},
	}

	// 提取content
	if content, ok := anthropicResp["content"].([]interface{}); ok && len(content) > 0 {
		if firstContent, ok := content[0].(map[string]interface{}); ok {
			if text, ok := firstContent["text"].(string); ok {
				openaiResp["choices"].([]map[string]interface{})[0]["message"].(map[string]interface{})["content"] = text
			}
		}
	}

	// 转换stop_reason
	if stopReason, ok := anthropicResp["stop_reason"].(string); ok {
		finishReason := "stop"
		if stopReason == "max_tokens" {
			finishReason = "length"
		}
		openaiResp["choices"].([]map[string]interface{})[0]["finish_reason"] = finishReason
	}

	// 添加usage信息
	if usage, ok := anthropicResp["usage"].(map[string]interface{}); ok {
		inputTokens := 0
		outputTokens := 0
		if it, ok := usage["input_tokens"].(float64); ok {
			inputTokens = int(it)
		}
		if ot, ok := usage["output_tokens"].(float64); ok {
			outputTokens = int(ot)
		}
		openaiResp["usage"] = map[string]interface{}{
			"prompt_tokens":     inputTokens,
			"completion_tokens": outputTokens,
			"total_tokens":      inputTokens + outputTokens,
		}
	}

	return openaiResp
}

// OpenAI兼容的chat completions端点
func chatCompletionsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("收到OpenAI格式请求: %s %s", r.Method, r.URL.Path)

	// 获取客户端提供的 API Key
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Printf("错误: 缺少 Authorization 头")
		http.Error(w, `{"error": {"message": "Authorization header is required", "type": "invalid_request_error"}}`, http.StatusUnauthorized)
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, `{"error": {"message": "Invalid authorization header format", "type": "invalid_request_error"}}`, http.StatusUnauthorized)
		return
	}
	clientAPIKey := parts[1]

	// 验证客户端 API Key 是否匹配代理 Key
	if config.ProxyAPIKey != "" && clientAPIKey != config.ProxyAPIKey {
		log.Printf("错误: API Key 验证失败")
		http.Error(w, `{"error": {"message": "Invalid API key", "type": "authentication_error"}}`, http.StatusUnauthorized)
		return
	}

	log.Printf("API Key 验证通过")

	// 读取请求体
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("错误: 读取请求体失败: %v", err)
		http.Error(w, `{"error": {"message": "Failed to read request body", "type": "invalid_request_error"}}`, http.StatusBadRequest)
		return
	}
	r.Body.Close()

	// 解析OpenAI格式请求
	var openaiReq map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &openaiReq); err != nil {
		log.Printf("错误: 解析JSON失败: %v", err)
		http.Error(w, `{"error": {"message": "Invalid JSON", "type": "invalid_request_error"}}`, http.StatusBadRequest)
		return
	}

	log.Printf("OpenAI请求: model=%v, messages数量=%d", openaiReq["model"], len(openaiReq["messages"].([]interface{})))

	// 转换为Anthropic格式
	anthropicReq := convertOpenAIToAnthropic(openaiReq)
	anthropicBody, err := json.Marshal(anthropicReq)
	if err != nil {
		log.Printf("错误: 序列化Anthropic请求失败: %v", err)
		http.Error(w, `{"error": {"message": "Internal error", "type": "server_error"}}`, http.StatusInternalServerError)
		return
	}

	log.Printf("已转换为Anthropic格式，请求体大小: %d bytes", len(anthropicBody))

	// 创建代理请求
	proxyReq, err := http.NewRequest("POST", config.AnthropicTarget, bytes.NewBuffer(anthropicBody))
	if err != nil {
		log.Printf("错误: 创建代理请求失败: %v", err)
		http.Error(w, `{"error": {"message": "Internal error", "type": "server_error"}}`, http.StatusInternalServerError)
		return
	}

	// 设置头信息 - 使用源头 Factory API Key
	proxyReq.Header.Set("Content-Type", "application/json")
	proxyReq.Header.Set("Authorization", "Bearer "+config.FactoryAPIKey)
	proxyReq.Header.Set("Host", "gibuoilncyzqebelqjqz.supabase.co")
	proxyReq.Header.Set("User-Agent", "Factory-Proxy/1.0.0")
	proxyReq.Header.Set("x-forwarded-for", "unknown")
	proxyReq.Header.Set("x-forwarded-proto", "http")

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	log.Printf("发送请求到Anthropic API...")
	resp, err := client.Do(proxyReq)
	if err != nil {
		log.Printf("错误: 请求失败: %v", err)
		http.Error(w, `{"error": {"message": "Proxy request failed", "type": "server_error"}}`, http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	log.Printf("收到响应: 状态码 %d", resp.StatusCode)

	// 读取Anthropic响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("错误: 读取响应体失败: %v", err)
		http.Error(w, `{"error": {"message": "Failed to read response", "type": "server_error"}}`, http.StatusInternalServerError)
		return
	}

	// 如果不是200，直接返回错误
	if resp.StatusCode != http.StatusOK {
		log.Printf("Anthropic API返回错误: %d, %s", resp.StatusCode, string(respBody))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)

		// 尝试转换错误格式
		var anthropicError map[string]interface{}
		if json.Unmarshal(respBody, &anthropicError) == nil {
			openaiError := map[string]interface{}{
				"error": map[string]interface{}{
					"message": fmt.Sprintf("%v", anthropicError),
					"type":    "api_error",
				},
			}
			json.NewEncoder(w).Encode(openaiError)
		} else {
			w.Write(respBody)
		}
		return
	}

	// 解析Anthropic响应
	var anthropicResp map[string]interface{}
	if err := json.Unmarshal(respBody, &anthropicResp); err != nil {
		log.Printf("错误: 解析Anthropic响应失败: %v", err)
		http.Error(w, `{"error": {"message": "Failed to parse response", "type": "server_error"}}`, http.StatusInternalServerError)
		return
	}

	// 转换为OpenAI格式
	openaiResp := convertAnthropicToOpenAI(anthropicResp)

	log.Printf("已转换为OpenAI格式，返回响应")

	// 返回OpenAI格式响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(openaiResp)
}

// 健康检查端点
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"uptime":    time.Since(startTime).Seconds(),
	})
}

// API 文档端点
func docsHandler(w http.ResponseWriter, r *http.Request) {
	// 读取 docs.html 文件
	htmlContent, err := os.ReadFile("docs.html")
	if err != nil {
		// 如果文件不存在，返回内嵌的文档
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, getEmbeddedDocs())
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(htmlContent)
}

// 内嵌的 API 文档
func getEmbeddedDocs() string {
	return `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Factory Proxy API - 文档</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', sans-serif;
            line-height: 1.6;
            color: #333;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            border-radius: 12px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            overflow: hidden;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 40px;
            text-align: center;
        }
        .header h1 { font-size: 2.5em; margin-bottom: 10px; }
        .content { padding: 40px; }
        .section { margin-bottom: 40px; }
        .section h2 {
            color: #667eea;
            font-size: 1.8em;
            margin-bottom: 20px;
            padding-bottom: 10px;
            border-bottom: 3px solid #667eea;
        }
        .endpoint {
            background: #f8f9fa;
            border-left: 4px solid #667eea;
            padding: 20px;
            margin: 15px 0;
            border-radius: 4px;
        }
        .method {
            display: inline-block;
            padding: 4px 12px;
            border-radius: 4px;
            font-weight: bold;
            margin-right: 10px;
        }
        .method.post { background: #10b981; color: white; }
        .method.get { background: #3b82f6; color: white; }
        code {
            background: #f1f5f9;
            padding: 2px 6px;
            border-radius: 3px;
            font-family: Monaco, monospace;
            color: #e74c3c;
        }
        pre {
            background: #1e293b;
            color: #e2e8f0;
            padding: 20px;
            border-radius: 6px;
            overflow-x: auto;
            margin: 15px 0;
            font-family: Monaco, monospace;
        }
        .footer {
            background: #f8f9fa;
            padding: 30px;
            text-align: center;
            color: #666;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚀 Factory Proxy API</h1>
            <p>OpenAI 兼容格式 | 支持 25+ AI 模型</p>
        </div>
        <div class="content">
            <div class="section">
                <h2>📖 快速开始</h2>
                <p>Factory Proxy API 提供 OpenAI 兼容的接口，让您可以使用标准的 OpenAI SDK 访问 Factory AI 的强大模型。</p>
            </div>

            <div class="section">
                <h2>🔌 API 端点</h2>
                
                <div class="endpoint">
                    <div><span class="method post">POST</span><code>/v1/chat/completions</code></div>
                    <p>创建对话补全</p>
                    <pre>curl -X POST http://localhost:8003/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_FACTORY_API_KEY" \
  -d '{
    "model": "claude-sonnet-4-5-20250929",
    "messages": [{"role": "user", "content": "Hello!"}],
    "max_tokens": 100
  }'</pre>
                </div>

                <div class="endpoint">
                    <div><span class="method get">GET</span><code>/v1/health</code></div>
                    <p>健康检查</p>
                    <pre>curl http://localhost:8003/v1/health</pre>
                </div>
            </div>

            <div class="section">
                <h2>🤖 支持的模型</h2>
                <p>支持 <strong>25+</strong> 种 AI 模型，包括：</p>
                <ul style="margin-left: 20px; margin-top: 10px;">
                    <li><code>claude-sonnet-4-5-20250929</code> - Claude 4.5 Sonnet (推荐)</li>
                    <li><code>claude-opus-4-1-20250805</code> - Claude 4 Opus</li>
                    <li><code>gpt-5-2025-08-07</code> - GPT-5</li>
                    <li><code>gpt-5-codex</code> - GPT-5 Codex</li>
                    <li><code>gemini-2.5-pro</code> - Gemini 2.5 Pro</li>
                    <li>更多模型请查看 <a href="https://github.com/libaxuan/factory-go-api/blob/main/MODELS.md" target="_blank">MODELS.md</a></li>
                </ul>
            </div>

            <div class="section">
                <h2>🔑 认证</h2>
                <p>使用 Factory API Key 进行认证：</p>
                <pre>Authorization: Bearer YOUR_FACTORY_API_KEY</pre>
            </div>

            <div class="section">
                <h2>📝 请求示例</h2>
                <h3>Python (OpenAI SDK)</h3>
                <pre>from openai import OpenAI

client = OpenAI(
    base_url="http://localhost:8003/v1",
    api_key="YOUR_FACTORY_API_KEY"
)

response = client.chat.completions.create(
    model="claude-sonnet-4-5-20250929",
    messages=[{"role": "user", "content": "Hello!"}]
)
print(response.choices[0].message.content)</pre>

                <h3>Node.js</h3>
                <pre>const OpenAI = require('openai');

const client = new OpenAI({
    baseURL: 'http://localhost:8003/v1',
    apiKey: 'YOUR_FACTORY_API_KEY'
});

const response = await client.chat.completions.create({
    model: 'claude-sonnet-4-5-20250929',
    messages: [{ role: 'user', content: 'Hello!' }]
});
console.log(response.choices[0].message.content);</pre>

                <h3>cURL</h3>
                <pre>curl -X POST http://localhost:8003/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_FACTORY_API_KEY" \
  -d '{
    "model": "claude-sonnet-4-5-20250929",
    "messages": [{"role": "user", "content": "Hello!"}],
    "max_tokens": 100,
    "temperature": 0.7
  }'</pre>
            </div>

            <div class="section">
                <h2>⚙️ 参数说明</h2>
                <table style="width:100%; border-collapse: collapse;">
                    <tr style="background: #f8f9fa;">
                        <th style="padding: 12px; text-align: left; border-bottom: 1px solid #e2e8f0;">参数</th>
                        <th style="padding: 12px; text-align: left; border-bottom: 1px solid #e2e8f0;">类型</th>
                        <th style="padding: 12px; text-align: left; border-bottom: 1px solid #e2e8f0;">说明</th>
                    </tr>
                    <tr>
                        <td style="padding: 12px; border-bottom: 1px solid #e2e8f0;"><code>model</code></td>
                        <td style="padding: 12px; border-bottom: 1px solid #e2e8f0;">string</td>
                        <td style="padding: 12px; border-bottom: 1px solid #e2e8f0;">模型名称（必填）</td>
                    </tr>
                    <tr>
                        <td style="padding: 12px; border-bottom: 1px solid #e2e8f0;"><code>messages</code></td>
                        <td style="padding: 12px; border-bottom: 1px solid #e2e8f0;">array</td>
                        <td style="padding: 12px; border-bottom: 1px solid #e2e8f0;">对话消息数组（必填）</td>
                    </tr>
                    <tr>
                        <td style="padding: 12px; border-bottom: 1px solid #e2e8f0;"><code>max_tokens</code></td>
                        <td style="padding: 12px; border-bottom: 1px solid #e2e8f0;">integer</td>
                        <td style="padding: 12px; border-bottom: 1px solid #e2e8f0;">最大生成 token 数</td>
                    </tr>
                    <tr>
                        <td style="padding: 12px; border-bottom: 1px solid #e2e8f0;"><code>temperature</code></td>
                        <td style="padding: 12px; border-bottom: 1px solid #e2e8f0;">float</td>
                        <td style="padding: 12px; border-bottom: 1px solid #e2e8f0;">温度参数 (0-2)</td>
                    </tr>
                    <tr>
                        <td style="padding: 12px; border-bottom: 1px solid #e2e8f0;"><code>stream</code></td>
                        <td style="padding: 12px; border-bottom: 1px solid #e2e8f0;">boolean</td>
                        <td style="padding: 12px; border-bottom: 1px solid #e2e8f0;">是否流式输出</td>
                    </tr>
                </table>
            </div>
        </div>
        <div class="footer">
            <p><strong>Factory Proxy API</strong> | <a href="https://github.com/libaxuan/factory-go-api" target="_blank" style="color: #667eea;">GitHub</a> | <a href="https://github.com/libaxuan/factory-go-api/blob/main/README.md" target="_blank" style="color: #667eea;">文档</a></p>
        </div>
    </div>
</body>
</html>`
}

func main() {
	// 验证必需的环境变量
	if config.FactoryAPIKey == "" {
		log.Fatalf("❌ 错误: 必须设置 FACTORY_API_KEY 环境变量")
	}
	if config.ProxyAPIKey == "" {
		log.Fatalf("❌ 错误: 必须设置 PROXY_API_KEY 环境变量")
	}

	log.Printf("🚀 Factory OpenAI-Compatible Proxy 启动中...")
	log.Printf("📍 端口: %s", config.Port)
	log.Printf("➡️  目标: %s", config.AnthropicTarget)
	log.Printf("🔐 API Key 代理: 已启用")
	log.Printf("   - 对外 Key: %s***", config.ProxyAPIKey[:min(8, len(config.ProxyAPIKey))])
	log.Printf("   - 源头 Key: %s***", config.FactoryAPIKey[:min(8, len(config.FactoryAPIKey))])

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		path := r.URL.Path

		if path == "/health" || path == "/v1/health" {
			healthHandler(recorder, r)
		} else if path == "/v1/chat/completions" {
			chatCompletionsHandler(recorder, r)
		} else if path == "/docs" || path == "/v1/docs" {
			docsHandler(recorder, r)
		} else if path == "/" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"service": "Factory OpenAI-Compatible Proxy",
				"version": "1.0",
			})
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{
					"message": "Invalid endpoint. Use /v1/chat/completions or /health",
					"type":    "invalid_request_error",
				},
			})
		}

		duration := time.Since(start)
		log.Printf("[%s] %s %s - %d - %v", r.Method, path, r.RemoteAddr, recorder.statusCode, duration)
	})

	server := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("✅ 服务器已启动，监听于 http://localhost:%s", config.Port)
	log.Printf("📋 API 端点:")
	log.Printf("   - POST /v1/chat/completions")
	log.Printf("   - GET  /v1/health")
	log.Printf("   - GET  /docs")
	log.Printf("")
	log.Printf("📖 API 文档: http://localhost:%s/docs", config.Port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("❌ 服务器启动失败: %v", err)
	}
}
