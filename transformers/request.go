
package transformers

import (
	"factory-go-api/config"

	"github.com/google/uuid"
)

// OpenAIMessage OpenAI 格式的消息
type OpenAIMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"` // 可以是 string 或 []ContentPart
}

// ContentPart 消息内容部分
type ContentPart struct {
	Type     string                 `json:"type"`
	Text     string                 `json:"text,omitempty"`
	ImageURL map[string]interface{} `json:"image_url,omitempty"`
}

// OpenAIRequest OpenAI 标准请求格式
type OpenAIRequest struct {
	Model            string          `json:"model"`
	Messages         []OpenAIMessage `json:"messages"`
	MaxTokens        int             `json:"max_tokens,omitempty"`
	Temperature      float64         `json:"temperature,omitempty"`
	TopP             float64         `json:"top_p,omitempty"`
	Stream           bool            `json:"stream,omitempty"`
	Tools            []interface{}   `json:"tools,omitempty"`
	PresencePenalty  float64         `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64         `json:"frequency_penalty,omitempty"`
}

// AnthropicMessage Anthropic 格式的消息
type AnthropicMessage struct {
	Role    string                   `json:"role"`
	Content []map[string]interface{} `json:"content"`
}

// AnthropicRequest Anthropic 请求格式
type AnthropicRequest struct {
	Model       string                   `json:"model"`
	Messages    []AnthropicMessage       `json:"messages"`
	System      []map[string]interface{} `json:"system,omitempty"`
	MaxTokens   int                      `json:"max_tokens"`
	Temperature float64                  `json:"temperature,omitempty"`
	Stream      bool                     `json:"stream,omitempty"`
	Thinking    *ThinkingConfig          `json:"thinking,omitempty"`
}

// ThinkingConfig Anthropic 的思考配置
type ThinkingConfig struct {
	Type         string `json:"type"`
	BudgetTokens int    `json:"budget_tokens"`
}

// FactoryOpenAIMessage Factory OpenAI 格式的消息
type FactoryOpenAIMessage struct {
	Role    string                   `json:"role"`
	Content []map[string]interface{} `json:"content"`
}

// FactoryOpenAIRequest Factory OpenAI 请求格式 (用于 /v1/responses 端点)
type FactoryOpenAIRequest struct {
	Model              string                 `json:"model"`
	Input              []FactoryOpenAIMessage `json:"input"`
	Instructions       string                 `json:"instructions,omitempty"`
	MaxOutputTokens    int                    `json:"max_output_tokens,omitempty"`
	Temperature        float64                `json:"temperature,omitempty"`
	TopP               float64                `json:"top_p,omitempty"`
	Stream             bool                   `json:"stream,omitempty"`
	Store              bool                   `json:"store"`
	Tools              []interface{}          `json:"tools,omitempty"`
	Reasoning          *ReasoningConfig       `json:"reasoning,omitempty"`
	PresencePenalty    float64                `json:"presence_penalty,omitempty"`
	FrequencyPenalty   float64                `json:"frequency_penalty,omitempty"`
	ParallelToolCalls  bool                   `json:"parallel_tool_calls,omitempty"`
}

// ReasoningConfig OpenAI 的推理配置
type ReasoningConfig struct {
	Effort  string `json:"effort"`
	Summary string `json:"summary"`
}

// TransformToAnthropic 将 OpenAI 格式转换为 Anthropic 格式
func TransformToAnthropic(req *OpenAIRequest) *AnthropicRequest {
	anthropicReq := &AnthropicRequest{
		Model:    req.Model,
		Messages: []AnthropicMessage{},
		Stream:   req.Stream,
	}

	// 设置 max_tokens
	if req.MaxTokens > 0 {
		anthropicReq.MaxTokens = req.MaxTokens
	} else {
		anthropicReq.MaxTokens = 64000 // 默认值
	}

	// 设置 temperature
	if req.Temperature > 0 {
		anthropicReq.Temperature = req.Temperature
	}

	// 转换消息并提取 system
	var systemPrompts []string
	systemPrompt := config.GetSystemPrompt()
	if systemPrompt != "" {
		systemPrompts = append(systemPrompts, systemPrompt)
	}

	for _, msg := range req.Messages {
		if msg.Role == "system" {
			// 提取 system 消息
			if text, ok := msg.Content.(string); ok {
				systemPrompts = append(systemPrompts, text)
			}
			continue
		}

		// 转换 user 和 assistant 消息
		anthropicMsg := AnthropicMessage{
			Role:    msg.Role,
			Content: []map[string]interface{}{},
		}

		if text, ok := msg.Content.(string); ok {
			anthropicMsg.Content = append(anthropicMsg.Content, map[string]interface{}{
				"type": "text",
				"text": text,
			})
		} else if parts, ok := msg.Content.([]interface{}); ok {
			for _, part := range parts {
				if partMap, ok := part.(map[string]interface{}); ok {
					anthropicMsg.Content = append(anthropicMsg.Content, partMap)
				}
			}
		}

		anthropicReq.Messages = append(anthropicReq.Messages, anthropicMsg)
	}

	// 设置 system 字段
	if len(systemPrompts) > 0 {
		anthropicReq.System = []map[string]interface{}{}
		for _, prompt := range systemPrompts {
			anthropicReq.System = append(anthropicReq.System, map[string]interface{}{
				"type": "text",
				"text": prompt,
			})
		}
	}

	// 处理 thinking 字段
	reasoning := config.GetModelReasoning(req.Model)
	if reasoning != "" {
		budgetTokens := map[string]int{
			"low":    4096,
			"medium": 12288,
			"high":   24576,
		}
		
		// 确保 max_tokens 大于 budget_tokens
		requiredBudget := budgetTokens[reasoning]
		if anthropicReq.MaxTokens <= requiredBudget {
			// 增加 max_tokens 以满足要求
			anthropicReq.MaxTokens = requiredBudget + 4000
		}
		
		anthropicReq.Thinking = &ThinkingConfig{
			Type:         "enabled",
			BudgetTokens: requiredBudget,
		}
	}

	return anthropicReq
}

// TransformToFactoryOpenAI 将 OpenAI 格式转换为 Factory OpenAI 格式
func TransformToFactoryOpenAI(req *OpenAIRequest) *FactoryOpenAIRequest {
	factoryReq := &FactoryOpenAIRequest{
		Model:  req.Model,
		Input:  []FactoryOpenAIMessage{},
		Stream: req.Stream,
		Store:  false,
	}

	// 转换 max_tokens
	if req.MaxTokens > 0 {
		factoryReq.MaxOutputTokens = req.MaxTokens
	}

	// 转换其他参数
	if req.Temperature > 0 {
		factoryReq.Temperature = req.Temperature
	}
	if req.TopP > 0 {
		factoryReq.TopP = req.TopP
	}
	if req.PresencePenalty != 0 {
		factoryReq.PresencePenalty = req.PresencePenalty
	}
	if req.FrequencyPenalty != 0 {
		factoryReq.FrequencyPenalty = req.FrequencyPenalty
	}

	// 转换工具
	if len(req.Tools) > 0 {
		factoryReq.Tools = req.Tools
	}

	// 提取 system 消息作为 instructions
	systemPrompt := config.GetSystemPrompt()
	var userSystemMessages []string

	for _, msg := range req.Messages {
		if msg.Role == "system" {
			if text, ok := msg.Content.(string); ok {
				userSystemMessages = append(userSystemMessages, text)
			}
			continue
		}

		// 转换消息内容
		factoryMsg := FactoryOpenAIMessage{
			Role:    msg.Role,
			Content: []map[string]interface{}{},
		}

		// 根据角色确定内容类型
		textType := "input_text"
		imageType := "input_image"
		if msg.Role == "assistant" {
			textType = "output_text"
			imageType = "output_image"
		}

		if text, ok := msg.Content.(string); ok {
			factoryMsg.Content = append(factoryMsg.Content, map[string]interface{}{
				"type": textType,
				"text": text,
			})
		} else if parts, ok := msg.Content.([]interface{}); ok {
			for _, part := range parts {
				if partMap, ok := part.(map[string]interface{}); ok {
					partType, _ := partMap["type"].(string)
					if partType == "text" {
						factoryMsg.Content = append(factoryMsg.Content, map[string]interface{}{
							"type": textType,
							"text": partMap["text"],
						})
					} else if partType == "image_url" {
						factoryMsg.Content = append(factoryMsg.Content, map[string]interface{}{
							"type":      imageType,
							"image_url": partMap["image_url"],
						})
					} else {
						// 其他类型直接传递
						factoryMsg.Content = append(factoryMsg.Content, partMap)
					}
				}
			}
		}

		factoryReq.Input = append(factoryReq.Input, factoryMsg)
	}

	// 设置 instructions
	if systemPrompt != "" || len(userSystemMessages) > 0 {
		instructions := systemPrompt
		for _, msg := range userSystemMessages {
			instructions += msg
		}
		factoryReq.Instructions = instructions
	}

	// 处理 reasoning 字段
	// 注意：GPT 模型的 reasoning 模式会导致只输出推理过程而无实际答案
	// 暂时禁用 reasoning 配置，让模型直接输出答案
	reasoning := config.GetModelReasoning(req.Model)
	if reasoning != "" {
		// 暂时禁用：GPT Extended Thinking 导致只有推理无答案
		// factoryReq.Reasoning = &ReasoningConfig{
		// 	Effort:  reasoning,
		// 	Summary: "auto",
		// }
		
		// 如果需要推理，增加 max_output_tokens 以确保有足够空间输出答案
		if factoryReq.MaxOutputTokens == 0 {
			factoryReq.MaxOutputTokens = 4000
		} else if factoryReq.MaxOutputTokens < 1000 {
			// 确保至少有足够的 token 用于实际输出
			factoryReq.MaxOutputTokens = factoryReq.MaxOutputTokens + 2000
		}
	}

	return factoryReq
}

// GetAnthropicHeaders 获取 Anthropic 请求头
func GetAnthropicHeaders(authHeader string, clientHeaders map[string]string, isStreaming bool, modelID string) map[string]string {
	headers := map[string]string{
		"content-type":      "application/json",
		"authorization":     authHeader,
		"anthropic-version": "2023-06-01",
		"user-agent":        config.GetUserAgent(),
	}

	if isStreaming {
		headers["accept"] = "text/event-stream"
	}

	// 传递客户端头
	if clientID, ok := clientHeaders["x-factory-client"]; ok {
		headers["x-factory-client"] = clientID
	} else {
		headers["x-factory-client"] = "cli"
	}

	if sessionID, ok := clientHeaders["x-session-id"]; ok {
		headers["x-session-id"] = sessionID
	}

	if msgID, ok := clientHeaders["x-assistant-message-id"]; ok {
		headers["x-assistant-message-id"] = msgID
	}

	return headers
}

// GetFactoryOpenAIHeaders 获取 Factory OpenAI 请求头
func GetFactoryOpenAIHeaders(authHeader string, clientHeaders map[string]string) map[string]string {
	// 生成唯一 ID
	sessionID := clientHeaders["x-session-id"]
	if sessionID == "" {
		sessionID = uuid.New().String()
	}

	messageID := clientHeaders["x-assistant-message-id"]
	if messageID == "" {
		messageID = uuid.New().String()
	}

	headers := map[string]string{
		"content-type":             "application/json",
		"authorization":            authHeader,
		"x-api-provider":           "azure_openai",
		"x-factory-client":         "cli",
		"x-session-id":             sessionID,
		"x-assistant-message-id":   messageID,
		"user-agent":               config.GetUserAgent(),
		"connection":               "keep-alive",
		"x-stainless-arch":         "x64",
		"x-stainless-lang":         "js",
		"x-stainless-os":           "MacOS",
		"x-stainless-runtime":      "node",
		"x-stainless-retry-count":  "0",
		"x-stainless-package-version": "5.23.2",
		"x-stainless-runtime-version": "v24.3.0",
	}

	// 从客户端头中覆盖
	for key, value := range clientHeaders {
		if key == "x-factory-client" || key == "x-session-id" || key == "x-assistant-message-id" {
			headers[key] = value
		}
	}

	return headers
}