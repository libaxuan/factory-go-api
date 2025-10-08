# 支持的模型列表

> 最后更新: 2025-10-08
> 测试方法: 通过 Factory AI 真实 API 测试

## ✅ 当前支持的模型

通过真实 API 测试，以下模型**确认可用**：

### Anthropic Claude 模型

| 模型 ID | 状态 | 说明 |
|--------|------|------|
| `claude-3-7-sonnet-20250219` | ✅ 可用 | Claude 3.7 Sonnet |
| `claude-sonnet-4-20250514` | ✅ 可用 | Claude Sonnet 4 |
| `claude-sonnet-4-5-20250929` | ✅ 可用 | Claude Sonnet 4.5（推荐） |

**使用示例：**
```bash
curl -X POST http://localhost:8003/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -d '{
    "model": "claude-sonnet-4-5-20250929",
    "messages": [{"role": "user", "content": "Hello"}]
  }'
```

## ❌ 测试失败的模型

以下模型在 Factory AI 测试中**不可用**：

### Claude 模型（不支持）

- `claude-3-5-sonnet-20241022` - HTTP 400: Unsupported OpenAI model ID
- `claude-3-5-sonnet-20250219` - HTTP 400: Unsupported OpenAI model ID
- `claude-sonnet-4-1-20250514` - HTTP 400: Unsupported OpenAI model ID
- `claude-3-5-haiku-20241022` - HTTP 400: Unsupported OpenAI model ID
- `claude-3-haiku-20240307` - HTTP 400: Unsupported OpenAI model ID

### OpenAI GPT 模型（端点超时）

所有 OpenAI 模型通过 Responses API (`/api/llm/o/v1/responses`) 请求超时：

- `gpt-5-2025-08-07` - 超时（30秒无响应）
- `gpt-5-mini-2025-08-07` 
- `gpt-5-nano-2025-08-07` - 超时
- `gpt-5-codex` - 超时
- `o1-2024-12-17` - HTTP 400
- `o1-mini-2024-09-12` - HTTP 400
- `o1-preview-2024-09-12` - HTTP 400
- `o3-mini` - HTTP 400
- `o4-mini` - 超时
- `gpt-4.1-preview` - 未测试
- `gpt-4o-2024-11-20` - HTTP 400
- `gpt-4o-mini-2024-07-18` - HTTP 400

### Google Gemini 模型（不支持）

- `gemini-2.5-flash` - HTTP 405: Method not allowed
- `gemini-2.5-pro` - HTTP 405: Method not allowed
- `gemini-2.0-flash-exp` - HTTP 400
- `gemini-exp-1206` - HTTP 400

### xAI Grok 模型（不支持）

- `grok-4` - HTTP 405: Method not allowed
- `grok-beta` - HTTP 400
- `grok-vision-beta` - HTTP 400

### 其他模型（不支持）

- `glm-4.6` - HTTP 400

## 📊 测试统计

- **总测试模型数**: 27
- **可用模型数**: 3
- **支持率**: 11.1%
- **测试日期**: 2025-10-08

## 🔍 问题分析

### OpenAI/GPT 模型超时原因

OpenAI 系列模型通过 Responses API 端点访问时持续超时（30秒无响应）。可能原因：

1. **API 权限限制**: 您的 Factory AI API Key 可能没有 GPT-5/Responses API 的访问权限
2. **端点配置问题**: Responses API 可能需要不同的认证或请求格式
3. **模型可用性**: GPT-5 系列可能还在内测中，需要特殊申请

**建议**: 联系 Factory AI 支持确认 OpenAI 模型的访问权限和正确的 API 端点。

### Gemini/Grok HTTP 405 错误

HTTP 405 表示 "Method Not Allowed"，说明：

1. 端点路径可能不正确
2. Factory AI 可能还未开放这些模型的 OpenAI 兼容 API

### 其他模型 HTTP 400 错误

HTTP 400 通常表示 "Unsupported OpenAI model ID"，说明：

1. 模型 ID 不正确或已废弃
2. Factory AI 不支持该模型
3. 需要使用完整的模型版本号

## 💡 使用建议

### 推荐使用

目前**强烈推荐**使用以下模型，已验证稳定可用：

```bash
# Claude Sonnet 4.5 (最新，性能最好)
"model": "claude-sonnet-4-5-20250929"

# Claude Sonnet 4
"model": "claude-sonnet-4-20250514"

# Claude 3.7 Sonnet
"model": "claude-3-7-sonnet-20250219"
```

### 完整示例

```bash
curl -X POST http://localhost:8003/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_PROXY_API_KEY" \
  -d '{
    "model": "claude-sonnet-4-5-20250929",
    "messages": [
      {"role": "user", "content": "Hello! How are you?"}
    ],
    "max_tokens": 1024,
    "stream": false
  }'
```

### 流式响应示例

```bash
curl -X POST http://localhost:8003/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_PROXY_API_KEY" \
  -d '{
    "model": "claude-sonnet-4-5-20250929",
    "messages": [
      {"role": "user", "content": "Write a short poem"}
    ],
    "stream": true
  }'
```

## 🔄 更新记录

- **2025-10-08**: 初始测试，确认 3 
个 Claude 模型可用
- **2025-10-08**: 测试发现 OpenAI/Gemini/Grok 模型暂不可用

## 📚 相关文档

- [README](../README.md) - 项目总览
- [快速开始](QUICK_START.md) - 快速上手指南
- [API 文档](README.md) - 完整 
API 参考
- [部署指南](DEPLOYMENT.md) - 生产环境部署