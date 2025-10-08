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
}

// 默认配置
var config = Config{
	Port:            getEnv("PORT", "8000"),
	AnthropicTarget: getEnv("ANTHROPIC_TARGET_URL", "https://gibuoilncyzqebelqjqz.supabase.co/functions/v1/smooth-handler/https://app.factory.ai/api/llm/a/v1/messages"),
}

var startTime = time.Now()

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
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

	// 获取API Key
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
	apiKey := parts[1]
	log.Printf("API Key已获取: %s...", apiKey[:10])

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

	// 设置头信息
	proxyReq.Header.Set("Content-Type", "application/json")
	proxyReq.Header.Set("Authorization", "Bearer "+apiKey)
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

func main() {
	log.Printf("🚀 Factory OpenAI-Compatible Proxy 启动中...")
	log.Printf("📍 端口: %s", config.Port)
	log.Printf("➡️  目标: %s", config.AnthropicTarget)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		path := r.URL.Path

		if path == "/health" || path == "/v1/health" {
			healthHandler(recorder, r)
		} else if path == "/v1/chat/completions" {
			chatCompletionsHandler(recorder, r)
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
	log.Printf("📋 OpenAI兼容接口:")
	log.Printf("   - POST /v1/chat/completions -> 需要 Authorization: Bearer <factory-api-key>")
	log.Printf("   - GET /health 或 /v1/health -> 健康检查")
	log.Printf("")
	log.Printf("💡 使用示例:")
	log.Printf("   curl -X POST http://localhost:%s/v1/chat/completions \\", config.Port)
	log.Printf("     -H 'Content-Type: application/json' \\")
	log.Printf("     -H 'Authorization: Bearer YOUR_FACTORY_API_KEY' \\")
	log.Printf("     -d '{")
	log.Printf("       \"model\": \"claude-sonnet-4-5-20250929\",")
	log.Printf("       \"messages\": [{\"role\": \"user\", \"content\": \"Hello\"}],")
	log.Printf("       \"max_tokens\": 100")
	log.Printf("     }'")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("❌ 服务器启动失败: %v", err)
	}
}
