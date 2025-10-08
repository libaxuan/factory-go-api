package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// 配置结构体
type Config struct {
	Port            string
	AnthropicTarget string
	OpenAITarget    string
	BedrockTarget   string
}

// 默认配置
var config = Config{
	Port:            getEnv("PORT", "8000"),
	AnthropicTarget: getEnv("ANTHROPIC_TARGET_URL", "https://gibuoilncyzqebelqjqz.supabase.co/functions/v1/smooth-handler/https://app.factory.ai/api/llm/a/v1/messages"),
	OpenAITarget:    getEnv("OPENAI_TARGET_URL", "https://spec.ngregersen.workers.dev/spec-114514/https://app.factory.ai/api/llm/o/v1/responses"),
	BedrockTarget:   getEnv("BEDROCK_TARGET_URL", "https://gibuoilncyzqebelqjqz.supabase.co/functions/v1/smooth-handler/https://app.factory.ai/api/llm/a/v1/messages"),
}

var startTime = time.Now()

// 获取环境变量，支持默认值
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

// 日志中间件
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(recorder, r)
		duration := time.Since(start)
		log.Printf("[%s] %s %s - %d - %v", r.Method, r.URL.Path, r.RemoteAddr, recorder.statusCode, duration)
	}
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

// 通用的代理请求处理
func proxyHandler(targetURL, serviceType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("代理请求: %s %s (服务类型: %s)", r.Method, r.URL.Path, serviceType)

		// 获取 API Key
		var apiKey string
		if serviceType == "openai" {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Printf("错误: 缺少 Authorization 头")
				http.Error(w, `{"error": "Authorization header is required"}`, http.StatusUnauthorized)
				return
			}
			// 提取 Bearer token
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				apiKey = parts[1]
			}
		} else {
			apiKey = r.Header.Get("x-api-key")
			if apiKey == "" {
				log.Printf("错误: 缺少 x-api-key 头")
				http.Error(w, `{"error": "x-api-key header is required"}`, http.StatusUnauthorized)
				return
			}
		}

		log.Printf("API Key已获取: %s...", apiKey[:10])

		// 解析目标 URL
		target, err := url.Parse(targetURL)
		if err != nil {
			log.Printf("错误: 无效的目标URL: %v", err)
			http.Error(w, fmt.Sprintf(`{"error": "Invalid target URL: %v"}`, err), http.StatusInternalServerError)
			return
		}

		log.Printf("目标URL: %s", target.String())

		// 读取并处理请求体
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("错误: 读取请求体失败: %v", err)
			http.Error(w, `{"error": "Failed to read request body"}`, http.StatusBadRequest)
			return
		}
		r.Body.Close()

		// 处理请求体（为Anthropic/Bedrock添加system prompt）
		if (serviceType == "anthropic" || serviceType == "bedrock") && len(bodyBytes) > 0 {
			var data map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &data); err == nil {
				// 处理system字段
				if _, hasSystem := data["system"]; !hasSystem {
					// 没有system字段，添加默认的
					data["system"] = []map[string]interface{}{
						{"type": "text", "text": "You are Droid, an AI software engineering agent built by Factory."},
					}
				} else if systemStr, ok := data["system"].(string); ok {
					// system是字符串，转换为数组并添加Factory prompt
					data["system"] = []interface{}{
						map[string]interface{}{"type": "text", "text": "You are Droid, an AI software engineering agent built by Factory."},
						map[string]interface{}{"type": "text", "text": systemStr},
					}
				} else if systemArr, ok := data["system"].([]interface{}); ok {
					// system是数组，在开头添加Factory prompt
					newSystem := []interface{}{
						map[string]interface{}{"type": "text", "text": "You are Droid, an AI software engineering agent built by Factory."},
					}
					data["system"] = append(newSystem, systemArr...)
				}

				// 重新编码
				if newBody, err := json.Marshal(data); err == nil {
					bodyBytes = newBody
					log.Printf("已添加 Factory system prompt")
				}
			}
		}

		log.Printf("请求体大小: %d bytes", len(bodyBytes))

		// 创建新请求
		proxyReq, err := http.NewRequest(r.Method, target.String(), bytes.NewBuffer(bodyBytes))
		if err != nil {
			log.Printf("错误: 创建代理请求失败: %v", err)
			http.Error(w, `{"error": "Failed to create proxy request"}`, http.StatusInternalServerError)
			return
		}

		// 只复制特定的头信息（与Deno版本保持一致）
		// Deno的createForwardedHeaders只保留了部分头
		allowedHeaders := map[string]bool{
			"content-type":    true,
			"content-length":  true,
			"accept":          true,
			"accept-encoding": true,
		}

		for key, values := range r.Header {
			lowerKey := strings.ToLower(key)
			if allowedHeaders[lowerKey] {
				for _, value := range values {
					proxyReq.Header.Add(key, value)
				}
			}
		}

		// 设置必要的头信息（与Deno完全一致）
		proxyReq.Header.Set("host", target.Host)
		proxyReq.Header.Set("user-agent", "Factory-Proxy/1.0.0")

		// 设置转发头
		if xForwardedFor := r.Header.Get("x-forwarded-for"); xForwardedFor != "" {
			proxyReq.Header.Set("x-forwarded-for", xForwardedFor)
		} else {
			proxyReq.Header.Set("x-forwarded-for", "unknown")
		}

		if xForwardedProto := r.Header.Get("x-forwarded-proto"); xForwardedProto != "" {
			proxyReq.Header.Set("x-forwarded-proto", xForwardedProto)
		} else {
			proxyReq.Header.Set("x-forwarded-proto", "http")
		}

		// 设置认证头
		proxyReq.Header.Set("authorization", "Bearer "+apiKey)

		// 对于 Bedrock 添加特定头
		if serviceType == "bedrock" {
			proxyReq.Header.Set("x-model-provider", "bedrock")
		}

		// 打印所有发送的头信息（调试用）
		log.Printf("发送的请求头:")
		for key, values := range proxyReq.Header {
			for _, value := range values {
				log.Printf("  %s: %s", key, value)
			}
		}

		// 发送请求
		client := &http.Client{
			Timeout: 30 * time.Second,
		}

		log.Printf("发送请求到: %s", target.String())
		resp, err := client.Do(proxyReq)
		if err != nil {
			log.Printf("错误: 请求失败: %v", err)
			http.Error(w, `{"error": "Proxy request failed"}`, http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		log.Printf("收到响应: 状态码 %d", resp.StatusCode)

		// 复制响应头
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		// 设置状态码
		w.WriteHeader(resp.StatusCode)

		// 复制响应体
		written, err := io.Copy(w, resp.Body)
		if err != nil {
			log.Printf("错误: 复制响应体失败: %v", err)
			return
		}

		log.Printf("响应完成: 写入 %d bytes", written)
	}
}

// 修改请求体（如果需要）
func modifyRequestBody(req *http.Request, serviceType string) {
	// 读取原始请求体
	body, err := io.ReadAll(req.Body)
	if err != nil || len(body) == 0 {
		return
	}
	req.Body.Close()

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return
	}

	// 根据服务类型修改请求体
	switch serviceType {
	case "anthropic", "bedrock":
		// 处理 system 参数
		if system, exists := data["system"]; exists {
			if system == nil {
				data["system"] = []map[string]interface{}{
					{"type": "text", "text": "You are Droid, an AI software engineering agent built by Factory."},
				}
			} else if systemStr, ok := system.(string); ok {
				data["system"] = []map[string]interface{}{
					{"type": "text", "text": "You are Droid, an AI software engineering agent built by Factory."},
					{"type": "text", "text": systemStr},
				}
			}
		}
	case "openai":
		// 模型替换
		if model, exists := data["model"]; exists {
			if model == "gpt-5" {
				data["model"] = "gpt-5-2025-08-07"
			}
		}
		// 添加 instructions
		data["instructions"] = "You are Droid, an AI software engineering agent built by Factory.\n"
	}

	// 重新编码请求体
	newBody, err := json.Marshal(data)
	if err != nil {
		return
	}

	req.Body = io.NopCloser(bytes.NewReader(newBody))
	req.ContentLength = int64(len(newBody))
}

func main() {
	// 显示启动信息
	log.Printf("🚀 Factory Go Proxy 启动中...")
	log.Printf("📍 端口: %s", config.Port)
	log.Printf("➡️  Anthropic 目标: %s", config.AnthropicTarget)
	log.Printf("➡️  OpenAI 目标: %s", config.OpenAITarget)
	log.Printf("➡️  Bedrock 目标: %s", config.BedrockTarget)

	// 创建自定义路由处理器
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		// 路由匹配
		path := r.URL.Path

		// 定义实际处理器
		var actualHandler http.HandlerFunc

		if path == "/health" {
			actualHandler = healthHandler
		} else if strings.HasPrefix(path, "/anthropic/") {
			actualHandler = proxyHandler(config.AnthropicTarget, "anthropic")
		} else if strings.HasPrefix(path, "/openai/") {
			actualHandler = proxyHandler(config.OpenAITarget, "openai")
		} else if strings.HasPrefix(path, "/bedrock/") {
			actualHandler = proxyHandler(config.BedrockTarget, "bedrock")
		} else if path == "/" {
			actualHandler = func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]string{
					"service": "Factory Go Proxy",
					"version": "1.0",
				})
			}
		} else {
			actualHandler = func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "无效的端点。请使用 /anthropic/, /openai/, /bedrock/ 或 /health",
				})
			}
		}

		// 执行处理器
		if actualHandler != nil {
			actualHandler.ServeHTTP(recorder, r)
		}

		duration := time.Since(start)
		log.Printf("[%s] %s %s - %d - %v", r.Method, path, r.RemoteAddr, recorder.statusCode, duration)
	})

	// 启动服务器
	server := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("✅ 服务器已启动，监听于 http://localhost:%s", config.Port)
	log.Printf("📋 使用方法:")
	log.Printf("   - /anthropic/* -> 需要 x-api-key 头")
	log.Printf("   - /openai/* -> 需要 Authorization: Bearer <token> 头")
	log.Printf("   - /bedrock/* -> 需要 x-api-key 头")
	log.Printf("   - /health -> 健康检查")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("❌ 服务器启动失败: %v", err)
	}
}
