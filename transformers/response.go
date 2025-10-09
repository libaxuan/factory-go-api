
package transformers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

// OpenAIResponse OpenAI 标准响应格式
type OpenAIResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Choices []OpenAIChoice         `json:"choices"`
	Usage   map[string]interface{} `json:"usage,omitempty"`
}

// OpenAIChoice 响应选项
type OpenAIChoice struct {
	Index        int                    `json:"index"`
	Message      *OpenAIMessageResponse `json:"message,omitempty"`
	Delta        *OpenAIMessageResponse `json:"delta,omitempty"`
	FinishReason *string                `json:"finish_reason"`
}

// OpenAIMessageResponse 消息响应
type OpenAIMessageResponse struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

// AnthropicResponseTransformer Anthropic 响应转换器
type AnthropicResponseTransformer struct {
	Model     string
	RequestID string
	Created   int64
}

// NewAnthropicResponseTransformer 创建 Anthropic 响应转换器
func NewAnthropicResponseTransformer(model, requestID string) *AnthropicResponseTransformer {
	if requestID == "" {
		requestID = fmt.Sprintf("chatcmpl-%d", time.Now().UnixNano())
	}
	return &AnthropicResponseTransformer{
		Model:     model,
		RequestID: requestID,
		Created:   time.Now().Unix(),
	}
}

// TransformNonStreamResponse 转换非流式响应
func (t *AnthropicResponseTransformer) TransformNonStreamResponse(anthropicResp map[string]interface{}) (*OpenAIResponse, error) {
	openaiResp := &OpenAIResponse{
		ID:      fmt.Sprintf("%v", anthropicResp["id"]),
		Object:  "chat.completion",
		Created: t.Created,
		Model:   t.Model,
		Choices: []OpenAIChoice{
			{
				Index: 0,
				Message: &OpenAIMessageResponse{
					Role:    "assistant",
					Content: "",
				},
				FinishReason: stringPtr("stop"),
			},
		},
	}

	// 提取内容
	// Extended Thinking 模型会返回多个 content 块：thinking + text
	// 我们需要找到 type=text 的块
	if content, ok := anthropicResp["content"].([]interface{}); ok && len(content) > 0 {
		for _, item := range content {
			if contentItem, ok := item.(map[string]interface{}); ok {
				// 检查类型
				contentType, hasType := contentItem["type"].(string)
				if hasType && contentType == "text" {
					// 提取文本内容
					if text, ok := contentItem["text"].(string); ok {
						openaiResp.Choices[0].Message.Content = text
						break
					}
				} else if !hasType {
					// 兼容旧格式：没有 type 字段的情况
					if text, ok := contentItem["text"].(string); ok {
						openaiResp.Choices[0].Message.Content = text
						break
					}
				}
			}
		}
	}

	// 转换 stop_reason
	if stopReason, ok := anthropicResp["stop_reason"].(string); ok {
		finishReason := "stop"
		if stopReason == "max_tokens" {
			finishReason = "length"
		}
		openaiResp.Choices[0].FinishReason = &finishReason
	}

	// 添加 usage 信息
	if usage, ok := anthropicResp["usage"].(map[string]interface{}); ok {
		inputTokens := 0
		outputTokens := 0
		if it, ok := usage["input_tokens"].(float64); ok {
			inputTokens = int(it)
		}
		if ot, ok := usage["output_tokens"].(float64); ok {
			outputTokens = int(ot)
		}
		openaiResp.Usage = map[string]interface{}{
			"prompt_tokens":     inputTokens,
			"completion_tokens": outputTokens,
			"total_tokens":      inputTokens + outputTokens,
		}
	}

	return openaiResp, nil
}

// TransformStreamChunk 转换流式响应块
func (t *AnthropicResponseTransformer) TransformStreamChunk(eventType string, eventData map[string]interface{}) (string, error) {
	switch eventType {
	case "message_start":
		return t.createOpenAIChunk("", "assistant", false, ""), nil

	case "content_block_delta":
		text := ""
		if delta, ok := eventData["delta"].(map[string]interface{}); ok {
			if textVal, ok := delta["text"].(string); ok {
				text = textVal
			}
		}
		return t.createOpenAIChunk(text, "", false, ""), nil

	case "message_delta":
		finishReason := "stop"
		if delta, ok := eventData["delta"].(map[string]interface{}); ok {
			if stopReason, ok := delta["stop_reason"].(string); ok && stopReason == "max_tokens" {
				finishReason = "length"
			}
		}
		return t.createOpenAIChunk("", "", true, finishReason), nil

	case "message_stop":
		return "", nil // 已经在 message_delta 中处理

	default:
		return "", nil // 忽略其他事件
	}
}

// createOpenAIChunk 创建 OpenAI 格式的流式块
func (t *AnthropicResponseTransformer) createOpenAIChunk(content, role string, finish bool, finishReason string) string {
	chunk := OpenAIResponse{
		ID:      t.RequestID,
		Object:  "chat.completion.chunk",
		Created: t.Created,
		Model:   t.Model,
		Choices: []OpenAIChoice{
			{
				Index: 0,
				Delta: &OpenAIMessageResponse{},
			},
		},
	}

	if role != "" {
		chunk.Choices[0].Delta.Role = role
	}
	if content != "" {
		chunk.Choices[0].Delta.Content = content
	}
	if finish {
		chunk.Choices[0].FinishReason = &finishReason
	}

	jsonData, _ := json.Marshal(chunk)
	return fmt.Sprintf("data: %s\n\n", string(jsonData))
}

// TransformStream 转换流式响应
func (t *AnthropicResponseTransformer) TransformStream(reader io.Reader) chan string {
	output := make(chan string, 100)

	go func() {
		defer close(output)

		scanner := bufio.NewScanner(reader)
		var currentEvent string

		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}

			if strings.HasPrefix(line, "event: ") {
				currentEvent = strings.TrimPrefix(line, "event: ")
			} else if strings.HasPrefix(line, "data: ") {
				dataStr := strings.TrimPrefix(line, "data: ")
				var eventData map[string]interface{}
				if err := json.Unmarshal([]byte(dataStr), &eventData); err == nil {
					if chunk, err := t.TransformStreamChunk(currentEvent, eventData); err == nil && chunk != "" {
						output <- chunk
					}
				}
			}
		}

		// 发送结束标记
		output <- "data: [DONE]\n\n"
	}()

	return output
}

// FactoryOpenAIResponseTransformer Factory OpenAI 响应转换器
type FactoryOpenAIResponseTransformer struct {
	Model     string
	RequestID string
	Created   int64
}

// NewFactoryOpenAIResponseTransformer 创建 Factory OpenAI 响应转换器
func NewFactoryOpenAIResponseTransformer(model, requestID string) *FactoryOpenAIResponseTransformer {
	if requestID == "" {
		requestID = fmt.Sprintf("chatcmpl-%d", time.Now().UnixNano())
	}
	return &FactoryOpenAIResponseTransformer{
		Model:     model,
		RequestID: requestID,
		Created:   time.Now().Unix(),
	}
}

// TransformNonStreamResponse 转换非流式响应
func (t *FactoryOpenAIResponseTransformer) TransformNonStreamResponse(factoryResp map[string]interface{}) (*OpenAIResponse, error) {
	// 检查是否已经是标准 OpenAI 格式（Factory 可能直接返回 OpenAI 格式）
	if choices, ok := factoryResp["choices"].([]interface{}); ok && len(choices) > 0 {
		// 已经是标准 OpenAI 格式，直接返回（只更新 model）
		openaiResp := &OpenAIResponse{
			ID:      fmt.Sprintf("%v", factoryResp["id"]),
			Object:  fmt.Sprintf("%v", factoryResp["object"]),
			Created: t.Created,
			Model:   t.Model, // 使用我们的 model ID
			Choices: []OpenAIChoice{},
		}
		
		// 转换 choices
		for i, choice := range choices {
			if choiceMap, ok := choice.(map[string]interface{}); ok {
				openaiChoice := OpenAIChoice{
					Index: i,
				}
				
				// 提取 message
				if message, ok := choiceMap["message"].(map[string]interface{}); ok {
					openaiChoice.Message = &OpenAIMessageResponse{
						Role:    fmt.Sprintf("%v", message["role"]),
						Content: fmt.Sprintf("%v", message["content"]),
					}
				}
				
				// 提取 finish_reason
				if fr, ok := choiceMap["finish_reason"].(string); ok {
					openaiChoice.FinishReason = &fr
				}
				
				openaiResp.Choices = append(openaiResp.Choices, openaiChoice)
			}
		}
		
		// 提取 usage
		if usage, ok := factoryResp["usage"].(map[string]interface{}); ok {
			openaiResp.Usage = usage
		}
		
		return openaiResp, nil
	}
	
	// Factory 自定义格式：output 数组
	openaiResp := &OpenAIResponse{
		ID:      fmt.Sprintf("%v", factoryResp["id"]),
		Object:  "chat.completion",
		Created: t.Created,
		Model:   t.Model,
		Choices: []OpenAIChoice{
			{
				Index: 0,
				Message: &OpenAIMessageResponse{
					Role:    "assistant",
					Content: "",
				},
				FinishReason: stringPtr("stop"),
			},
		},
	}

	// 提取响应内容
	// Factory OpenAI 响应格式：
	// output: [
	//   {type: "reasoning", ...},  // 推理过程
	//   {type: "message", content: [{text: "...", type: "output_text"}], ...}  // 实际回复
	// ]
	if output, ok := factoryResp["output"].([]interface{}); ok && len(output) > 0 {
		for _, item := range output {
			if outputItem, ok := item.(map[string]interface{}); ok {
				// 查找 type=message 的项
				if itemType, ok := outputItem["type"].(string); ok && itemType == "message" {
					// 提取 content 数组
					if contentArray, ok := outputItem["content"].([]interface{}); ok {
						for _, contentItem := range contentArray {
							if contentMap, ok := contentItem.(map[string]interface{}); ok {
								// 提取 text 字段
								if text, ok := contentMap["text"].(string); ok {
									openaiResp.Choices[0].Message.Content += text
								}
							}
						}
					}
					break // 找到 message 后就退出
				}
			}
		}
	}

	// 提取 finish_reason
	if status, ok := factoryResp["status"].(string); ok {
		finishReason := "stop"
		if status == "incomplete" {
			finishReason = "length"
		}
		openaiResp.Choices[0].FinishReason = &finishReason
	}

	// 添加 usage 信息
	if usage, ok := factoryResp["usage"].(map[string]interface{}); ok {
		inputTokens := 0
		outputTokens := 0
		if it, ok := usage["input_tokens"].(float64); ok {
			inputTokens = int(it)
		}
		if ot, ok := usage["output_tokens"].(float64); ok {
			outputTokens = int(ot)
		}
		openaiResp.Usage = map[string]interface{}{
			"prompt_tokens":     inputTokens,
			"completion_tokens": outputTokens,
			"total_tokens":      inputTokens + outputTokens,
		}
	}

	return openaiResp, nil
}

// TransformStreamChunk 转换流式响应块
func (t *FactoryOpenAIResponseTransformer) TransformStreamChunk(eventType string, eventData map[string]interface{}) (string, error) {
	switch eventType {
	case "response.created":
		return t.createOpenAIChunk("", "assistant", false, ""), nil

	case "response.in_progress":
		return "", nil

	// GPT Extended Thinking: 推理过程（不转发）
	case "response.reasoning_summary_text.delta":
		return "", nil

	case "response.reasoning_summary_text.done":
		return "", nil

	case "response.reasoning_summary_part.done":
		return "", nil

	// GPT 实际输出文本
	case "response.output_text.delta":
		text := ""
		if delta, ok := eventData["delta"].(string); ok {
			text = delta
		} else if textVal, ok := eventData["text"].(string); ok {
			text = textVal
		}
		return t.createOpenAIChunk(text, "", false, ""), nil

	case "response.output_text.done":
		return "", nil

	case "response.output_item.added":
		return "", nil

	case "response.output_item.done":
		return "", nil

	case "response.done":
		status := ""
		if response, ok := eventData["response"].(map[string]interface{}); ok {
			if statusVal, ok := response["status"].(string); ok {
				status = statusVal
			}
		}

		finishReason := "stop"
		if status == "incomplete" {
			finishReason = "length"
		}
		return t.createOpenAIChunk("", "", true, finishReason), nil

	case "response.incomplete":
		// GPT 因 max_output_tokens 等原因未完成
		finishReason := "length"
		if response, ok := eventData["response"].(map[string]interface{}); ok {
			if incompleteDetails, ok := response["incomplete_details"].(map[string]interface{}); ok {
				if reason, ok := incompleteDetails["reason"].(string); ok {
					if reason == "max_output_tokens" {
						finishReason = "length"
					}
				}
			}
		}
		return t.createOpenAIChunk("", "", true, finishReason), nil

	default:
		return "", nil
	}
}

// createOpenAIChunk 创建 OpenAI 格式的流式块
func (t *FactoryOpenAIResponseTransformer) createOpenAIChunk(content, role string, finish bool, finishReason string) string {
	chunk := OpenAIResponse{
		ID:      t.RequestID,
		Object:  "chat.completion.chunk",
		Created: t.Created,
		Model:   t.Model,
		Choices: []OpenAIChoice{
			{
				Index: 0,
				Delta: &OpenAIMessageResponse{},
			},
		},
	}

	if role != "" {
		chunk.Choices[0].Delta.Role = role
	}
	if content != "" {
		chunk.Choices[0].Delta.Content = content
	}
	if finish {
		chunk.Choices[0].FinishReason = &finishReason
	}

	jsonData, _ := json.Marshal(chunk)
	return fmt.Sprintf("data: %s\n\n", string(jsonData))
}

// TransformStream 转换流式响应
func (t *FactoryOpenAIResponseTransformer) TransformStream(reader io.Reader) chan string {
	output := make(chan string, 100)

	go func() {
		defer close(output)

		scanner := bufio.NewScanner(reader)
		var currentEvent string

		for scanner.Scan() {
			line := scanner.Text()
			
			if line == "" {
				continue
			}

			// 处理标准 OpenAI SSE 格式（Factory 可能直接返回）
			if strings.HasPrefix(line, "data: ") {
				dataStr := strings.TrimPrefix(line, "data: ")
				
				// 检查是否是 [DONE] 标记
				if strings.TrimSpace(dataStr) == "[DONE]" {
					output <- "data: [DONE]\n\n"
					continue
				}
				
				// 尝试解析为标准 OpenAI chunk
				var openaiChunk map[string]interface{}
				if err := json.Unmarshal([]byte(dataStr), &openaiChunk); err == nil {
					// 检查是否已经是标准 OpenAI 格式
					if _, hasChoices := openaiChunk["choices"]; hasChoices {
						// 直接转发（只更新 model）
						openaiChunk["model"] = t.Model
						if jsonData, err := json.Marshal(openaiChunk); err == nil {
							output <- fmt.Sprintf("data: %s\n\n", string(jsonData))
						}
						continue
					}
					
					// Factory 自定义事件格式
					if chunk, err := t.TransformStreamChunk(currentEvent, openaiChunk); err == nil && chunk != "" {
						output <- chunk
					}
				}
			} else if strings.HasPrefix(line, "event: ") {
				currentEvent = strings.TrimPrefix(line, "event: ")
			}
		}

		// 发送结束标记（如果还没发送）
		output <- "data: [DONE]\n\n"
	}()

	return output
}

// stringPtr 返回字符串指针
func stringPtr(s string) *string {
	return &s
}