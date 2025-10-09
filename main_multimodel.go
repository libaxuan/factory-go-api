
package main

import (
	"bytes"
	"encoding/json"
	"factory-go-api/config"
	"factory-go-api/transformers"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var startTime = time.Now()

// 获取环境变量，支持默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// min 函数
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

// 健康检查端点
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"uptime":    time.Since(startTime).Seconds(),
	})
}

// 模型列表端点
func modelsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	models := config.GetAllModels()
	openaiModels := make([]map[string]interface{}, 0, len(models))
	
	for _, model := range models {
		openaiModels = append(openaiModels, map[string]interface{}{
			"id":      model.ID,
			"object":  "model",
			"created": time.Now().Unix(),
			"owned_by": "factory",
		})
	}
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"object": "list",
		"data":   openaiModels,
	})
}

// API 文档端点
func docsHandler(w http.ResponseWriter, r *http.Request) {
	htmlContent, err := os.ReadFile("docs.html")
	if err != nil {
		http.Error(w, "Documentation not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(htmlContent)
}

// OpenAI 兼容的聊天端点
func chatCompletionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// 获取客户端 Authorization 头
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, `{"error": {"message": "Authorization header is required", "type": "invalid_request_error"}}`, http.StatusUnauthorized)
		return
	}

	// 验证 PROXY_API_KEY（如果配置了）
	proxyAPIKey := getEnv("PROXY_API_KEY", "")
	if proxyAPIKey != "" {
		// 提取客户端 API Key
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"error": {"message": "Invalid authorization header format", "type": "invalid_request_error"}}`, http.StatusUnauthorized)
			return
		}
		clientAPIKey := parts[1]

		// 验证客户端 Key 是否匹配代理 Key
		if clientAPIKey != proxyAPIKey {
			log.Printf("❌ API Key 验证失败")
			http.Error(w, `{"error": {"message": "Invalid API key", "type": "authentication_error"}}`, http.StatusUnauthorized)
			return
		}
	}

	// 使用源头 FACTORY_API_KEY 替换 Authorization 头
	factoryAPIKey := getEnv("FACTORY_API_KEY", "")
	if factoryAPIKey == "" {
		log.Printf("❌ FACTORY_API_KEY 未配置")
		http.Error(w, `{"error": {"message": "Server configuration error", "type": "server_error"}}`, http.StatusInternalServerError)
		return
	}
	authHeader = "Bearer " + factoryAPIKey

	// 读取请求体
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("错误: 读取请求体失败: %v", err)
		http.Error(w, `{"error": {"message": "Failed to read request body", "type": "invalid_request_error"}}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 解析 OpenAI 请求
	var openaiReq transformers.OpenAIRequest
	if err := json.Unmarshal(bodyBytes, &openaiReq); err != nil {
		log.Printf("错误: 解析请求体失败: %v", err)
		http.Error(w, `{"error": {"message": "Invalid JSON", "type": "invalid_request_error"}}`, http.StatusBadRequest)
		return
	}

	// 检查模型是否支持
	model := config.GetModelByID(openaiReq.Model)
	if model == nil {
		log.Printf("❌ 不支持的模型: %s", openaiReq.Model)
		http.Error(w, fmt.Sprintf(`{"error": {"message": "Model '%s' not found", "type": "invalid_request_error"}}`, openaiReq.Model), http.StatusNotFound)
		return
	}

	log.Printf("✅ %s [%s] stream=%v", openaiReq.Model, model.Type, openaiReq.Stream)

	// 根据模型类型路由请求
	switch model.Type {
	case "anthropic":
		handleAnthropicRequest(w, r, &openaiReq, model, authHeader)
	case "openai":
		handleFactoryOpenAIRequest(w, r, &openaiReq, model, authHeader)
	default:
		http.Error(w, `{"error": {"message": "Unsupported model type", "type": "invalid_request_error"}}`, http.StatusBadRequest)
	}
}

// 处理 Anthropic 类型请求
func handleAnthropicRequest(w http.ResponseWriter, r *http.Request, openaiReq *transformers.OpenAIRequest, model *config.Model, authHeader string) {
	// 转换请求
	anthropicReq := transformers.TransformToAnthropic(openaiReq)
	
	// 获取端点
	endpoint := config.GetEndpointByType("anthropic")
	if endpoint == nil {
		http.Error(w, `{"error": {"message": "Anthropic endpoint not configured", "type": "configuration_error"}}`, http.StatusInternalServerError)
		return
	}

	// 序列化请求
	reqBody, err := json.Marshal(anthropicReq)
	if err != nil {
		log.Printf("错误: 序列化请求失败: %v", err)
		http.Error(w, `{"error": {"message": "Failed to serialize request", "type": "server_error"}}`, http.StatusInternalServerError)
		return
	}


	// 创建 HTTP 请求
	proxyReq, err := http.NewRequest(http.MethodPost, endpoint.BaseURL, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("错误: 创建请求失败: %v", err)
		http.Error(w, `{"error": {"message": "Failed to create request", "type": "server_error"}}`, http.StatusInternalServerError)
		return
	}

	// 设置请求头
	clientHeaders := extractClientHeaders(r)
	headers := transformers.GetAnthropicHeaders(authHeader, clientHeaders, openaiReq.Stream, model.ID)
	for key, value := range headers {
		proxyReq.Header.Set(key, value)
	}

	// 发送请求
	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(proxyReq)
	if err != nil {
		log.Printf("错误: 请求失败: %v", err)
		http.Error(w, `{"error": {"message": "Request to upstream failed", "type": "upstream_error"}}`, http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	log.Printf("📥 Anthropic 响应: %d", resp.StatusCode)

	// 处理响应
	if openaiReq.Stream {
		// 流式响应
		handleAnthropicStreamResponse(w, resp, model.ID)
	} else {
		// 非流式响应
		handleAnthropicNonStreamResponse(w, resp, model.ID)
	}
}

// 处理 Factory OpenAI 类型请求
func handleFactoryOpenAIRequest(w http.ResponseWriter, r *http.Request, openaiReq *transformers.OpenAIRequest, model *config.Model, authHeader string) {
	// 转换请求
	factoryReq := transformers.TransformToFactoryOpenAI(openaiReq)
	
	// 获取端点
	endpoint := config.GetEndpointByType("openai")
	if endpoint == nil {
		http.Error(w, `{"error": {"message": "OpenAI endpoint not configured", "type": "configuration_error"}}`, http.StatusInternalServerError)
		return
	}

	// 序列化请求
	reqBody, err := json.Marshal(factoryReq)
	if err != nil {
		log.Printf("错误: 序列化请求失败: %v", err)
		http.Error(w, `{"error": {"message": "Failed to serialize request", "type": "server_error"}}`, http.StatusInternalServerError)
		return
	}


	// 创建 HTTP 请求
	proxyReq, err := http.NewRequest(http.MethodPost, endpoint.BaseURL, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("错误: 创建请求失败: %v", err)
		http.Error(w, `{"error": {"message": "Failed to create request", "type": "server_error"}}`, http.StatusInternalServerError)
		return
	}

	// 设置请求头
	clientHeaders := extractClientHeaders(r)
	headers := transformers.GetFactoryOpenAIHeaders(authHeader, clientHeaders)
	for key, value := range headers {
		proxyReq.Header.Set(key, value)
	}

	// 发送请求
	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(proxyReq)
	if err != nil {
		log.Printf("错误: 请求失败: %v", err)
		http.Error(w, `{"error": {"message": "Request to upstream failed", "type": 
"upstream_error"}}`, http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	log.Printf("📥 Factory OpenAI 响应: %d", resp.StatusCode)

	// 处理响应
	if openaiReq.Stream {
		// 流式响应
		handleFactoryOpenAIStreamResponse(w, resp, model.ID)
	} else {
		// 非流式响应
		handleFactoryOpenAINonStreamResponse(w, resp, model.ID)
	}
}

// 处理 Anthropic 非流式响应
func handleAnthropicNonStreamResponse(w http.ResponseWriter, resp *http.Response, modelID string) {
	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("错误: 读取响应失败: %v", err)
		http.Error(w, `{"error": {"message": "Failed to read response", "type": "server_error"}}`, http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		// 错误响应直接转发
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
		return
	}

	// 解析 Anthropic 响应
	var anthropicResp map[string]interface{}
	if err := json.Unmarshal(body, &anthropicResp); err != nil {
		log.Printf("错误: 解析响应失败: %v", err)
		http.Error(w, `{"error": {"message": "Failed to parse response", "type": "server_error"}}`, http.StatusInternalServerError)
		return
	}

	// 转换为 OpenAI 格式
	transformer := transformers.NewAnthropicResponseTransformer(modelID, "")
	openaiResp, err := transformer.TransformNonStreamResponse(anthropicResp)
	if err != nil {
		log.Printf("错误: 转换响应失败: %v", err)
		http.Error(w, `{"error": {"message": "Failed to transform response", "type": "server_error"}}`, http.StatusInternalServerError)
		return
	}

	// 返回 OpenAI 格式响应
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(openaiResp)
}

// 处理 Anthropic 流式响应
func handleAnthropicStreamResponse(w http.ResponseWriter, resp *http.Response, modelID string) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, `{"error": {"message": "Streaming not supported", "type": "server_error"}}`, http.StatusInternalServerError)
		return
	}

	// 创建转换器
	transformer := transformers.NewAnthropicResponseTransformer(modelID, "")
	
	// 转换流式响应
	outputChan := transformer.TransformStream(resp.Body)
	
	for chunk := range outputChan {
		fmt.Fprint(w, chunk)
		flusher.Flush()
	}
}

// 处理 Factory OpenAI 非流式响应
func handleFactoryOpenAINonStreamResponse(w http.ResponseWriter, resp *http.Response, modelID string) {
	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("错误: 读取响应失败: %v", err)
		http.Error(w, `{"error": {"message": "Failed to read response", "type": "server_error"}}`, http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		// 错误响应直接转发
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
		return
	}

	// 解析 Factory OpenAI 响应
	var factoryResp map[string]interface{}
	if err := json.Unmarshal(body, &factoryResp); err != nil {
		log.Printf("错误: 解析响应失败: %v", err)
		http.Error(w, `{"error": {"message": "Failed to parse response", "type": "server_error"}}`, http.StatusInternalServerError)
		return
	}

	// 转换为 OpenAI 格式
	transformer := transformers.NewFactoryOpenAIResponseTransformer(modelID, "")
	openaiResp, err := transformer.TransformNonStreamResponse(factoryResp)
	if err != nil {
		log.Printf("错误: 转换响应失败: %v", err)
		http.Error(w, `{"error": {"message": "Failed to transform response", "type": "server_error"}}`, http.StatusInternalServerError)
		return
	}

	// 返回 OpenAI 格式响应
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(openaiResp)
}

// 处理 Factory OpenAI 流式响应
func handleFactoryOpenAIStreamResponse(w http.ResponseWriter, resp *http.Response, modelID string) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, `{"error": {"message": "Streaming not supported", "type": "server_error"}}`, http.StatusInternalServerError)
		return
	}

	// 创建转换器
	transformer := transformers.NewFactoryOpenAIResponseTransformer(modelID, "")
	
	// 转换流式响应
	outputChan := transformer.TransformStream(resp.Body)
	
	for chunk := range outputChan {
		fmt.Fprint(w, chunk)
		flusher.Flush()
	}
}

// 提取客户端请求头
func extractClientHeaders(r *http.Request) map[string]string {
	headers := make(map[string]string)
	
	// 提取需要转发的头
	forwardHeaders := []string{
		"x-session-id",
		"x-assistant-message-id",
		"x-factory-client",
		"x-stainless-arch",
		"x-stainless-lang",
		"x-stainless-os",
		"x-stainless-runtime",
		"x-stainless-retry-count",
		"x-stainless-package-version",
		"x-stainless-runtime-version",
	}
	
	for _, header := range forwardHeaders {
		if value := r.Header.Get(header); value != "" {
			headers[header] = value
		}
	}
	
	return headers
}

func main() {
	// 验证必需的环境变量
	factoryAPIKey := getEnv("FACTORY_API_KEY", "")
	if factoryAPIKey == "" {
		log.Fatalf("❌ 错误: 必须设置 FACTORY_API_KEY 环境变量")
	}

	proxyAPIKey := getEnv("PROXY_API_KEY", "")
	if proxyAPIKey != "" {
		log.Printf("🔐 代理模式: 已启用")
		log.Printf("   • 对外 Key: %s", proxyAPIKey)
		log.Printf("   • 源头 Key: %s***", factoryAPIKey[:min(8, len(factoryAPIKey))])
	} else {
		log.Printf("🔐 直连模式: 使用 FACTORY_API_KEY 直接访问")
	}

	// 加载配置
	configPath := getEnv("CONFIG_PATH", "config.json")
	log.Printf("📖 加载配置文件: %s", configPath)
	
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("❌ 加载配置失败: %v", err)
	}

	log.Printf("✅ 配置加载成功")
	log.Printf("📍 支持的模型 (%d):", len(cfg.Models))
	for _, model := range cfg.Models {
		log.Printf("   • %s [%s]", model.ID, model.Type)
	}

	// 设置路由
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/v1/models", modelsHandler)
	http.HandleFunc("/v1/chat/completions", chatCompletionsHandler)
	http.HandleFunc("/docs", docsHandler)
	
	// 根路径
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, `{"error": {"message": "Not found", "type": "invalid_request_error"}}`, http.StatusNotFound)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"service": "Factory Go API - Multi-Model Support",
			"version": "2.0",
			"endpoints": []string{
				"/health",
				"/v1/models",
				"/v1/chat/completions",
			},
		})
	})

	// 启动服务器
	port := fmt.Sprintf(":%d", cfg.Port)
	server := &http.Server{
		Addr:         port,
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 300 * time.Second,
		
IdleTimeout:  120 * time.Second,
	}

	log.Printf("🚀 服务启动于 http://localhost%s", port)
	log.Printf("📖 文档: http://localhost%s/docs", port)
	
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("❌ 服务器启动失败: %v", err)
	}
}