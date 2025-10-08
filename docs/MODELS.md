# 支持的模型列表

## ⚠️ 重要说明

**经过 2025-10-08 真实 API 测试验证，Factory AI 目前仅支持以下 3 个 Claude 模型。**

通过 OpenAI 兼容格式的代理，您可以使用标准的 OpenAI SDK 访问这些 Claude 模型。

---

## ✅ 支持的模型（已验证可用）

### 🤖 Claude 模型

| 模型 ID | 版本 | 上下文 | 状态 | 说明 |
|---------|------|--------|------|------|
| **`claude-sonnet-4-5-20250929`** ⭐ | Claude Sonnet 4.5 | 200K | ✅ 已验证 | **强烈推荐** - 最新最强 |
| **`claude-sonnet-4-20250514`** | Claude Sonnet 4 | 200K | ✅ 已验证 | 稳定可靠 |
| **`claude-3-7-sonnet-20250219`** | Claude 3.7 Sonnet | 200K | ✅ 已验证 | 性能优秀 |

> 📝 **测试日期**: 2025-10-08  
> 🔍 **测试方法**: 真实 API 调用  
> ✅ **可用模型**: 3 个

---

## 🎯 模型选择建议

### 🏆 推荐使用：Claude Sonnet 4.5

```python
model = "claude-sonnet-4-5-20250929"  # 强烈推荐
```

**适用场景**:
- ✅ 日常对话
- ✅ 代码生成与分析
- ✅ 文档处理
- ✅ 长文本分析
- ✅ 复杂推理任务

**优势**:
- 🚀 最新版本，性能最强
- 📊 200K tokens 上下文
- 💡 全能型，适合 90% 的场景

---

## 📖 使用示例

### Python (OpenAI SDK)

```python
from openai import OpenAI

client = OpenAI(
    api_key="YOUR_PROXY_API_KEY",
    base_url="http://localhost:8003/v1"
)

# 使用 Claude Sonnet 4.5 (推荐)
response = client.chat.completions.create(
    model="claude-sonnet-4-5-20250929",
    messages=[
        {"role": "user", "content": "Hello! 用中文介绍一下你自己"}
    ],
    max_tokens=1024
)

print(response.choices[0].message.content)
```

### 流式响应

```python
# 流式输出 - 实时显示生成过程
stream = client.chat.completions.create(
    model="claude-sonnet-4-5-20250929",
    messages=[{"role": "user", "content": "写一首关于编程的诗"}],
    stream=True
)

for chunk in stream:
    if chunk.choices[0].delta.content:
        print(chunk.choices[0].delta.content, end="", flush=True)
```

### cURL

```bash
# 非流式请求
curl -X POST http://localhost:8003/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_PROXY_API_KEY" \
  -d '{
    "model": "claude-sonnet-4-5-20250929",
    "messages": [{"role": "user", "content": "Hello!"}],
    "max_tokens": 1024,
    "stream": false
  }'

# 流式请求
curl -N -X POST http://localhost:8003/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_PROXY_API_KEY" \
  -d '{
    "model": "claude-sonnet-4-5-20250929",
    "messages": [{"role": "user", "content": "讲个笑话"}],
    "stream": true
  }'
```

### Node.js

```javascript
import OpenAI from 'openai';

const client = new OpenAI({
  apiKey: process.env.PROXY_API_KEY,
  baseURL: 'http://localhost:8003/v1'
});

// 非流式
const response = await client.chat.completions.create({
  model: 'claude-sonnet-4-5-20250929',
  messages: [{ role: 'user', content: 'Hello!' }]
});

console.log(response.choices[0].message.content);

// 流式
const stream = await client.chat.completions.create({
  model: 'claude-sonnet-4-5-20250929',
  messages: [{ role: 'user', content: 'Hello!' }],
  stream: true
});

for await (const chunk of stream) {
  if (chunk.choices[0]?.delta?.content) {
    process.stdout.write(chunk.choices[0].delta.content);
  }
}
```

---

## 🔄 快速切换模型

所有模型使用相同的 API 格式，只需更改 `model` 参数：

```python
# 测试不同模型
models_to_try = [
    "claude-sonnet-4-5-20250929",   # 推荐 ⭐
    "claude-sonnet-4-20250514",     # Claude 4
    "claude-3-7-sonnet-20250219"    # Claude 3.7
]

for model in models_to_try:
    response = client.chat.completions.create(
        model=model,
        messages=[{"role": "user", "content": "你好"}],
        max_tokens=200
    )
    print(f"\n模型: {model}")
    print(response.choices[0].message.content)
```

---

## ⚙️ 请求参数

所有模型支持以下标准参数：

```python
response = client.chat.completions.create(
    model="claude-sonnet-4-5-20250929",
    messages=[
        {"role": "system", "content": "You are a helpful assistant."},
        {"role": "user", "content": "Hello!"}
    ],
    max_tokens=4096,       # 最大输出长度
    temperature=0.7,       # 随机性 (0-2)
    top_p=1.0,            # 核采样
    stream=False,         # 是否流式输出
    stop=None             # 停止词
)
```

### 推荐的 max_tokens 设置

```python
# 简短对话
max_tokens = 1024      
# ~1K tokens

# 标准响应 (推荐)
max_tokens = 4096      # ~4K tokens

# 长文本生成
max_tokens = 8192      # ~8K tokens
```

---

## 📝 注意事项

1. **模型可用性**: 目前仅支持 3 个 Claude 模型（经真实测试验证）
2. **API 兼容性**: 完全兼容 OpenAI SDK 和 API 格式
3. **上下文长度**: 所有模型支持 200K tokens 上下文
4. **流式支持**: 所有模型都支持流式（`stream: true`）和非流式响应
5. **成本优化**: 使用合适的 `max_tokens` 参数控制成本

---

## ❓ 常见问题

### Q: 为什么只支持这 3 个模型？

A: 经过 2025-10-08 的真实 API 测试，Factory AI 后端仅这 3 个 Claude 模型返回正常响应。其他模型（包括 GPT、Gemini、Grok 等）均返回错误或超时。

### Q: 推荐使用哪个模型？

A: **强烈推荐** `claude-sonnet-4-5-20250929`，它是最新版本，性能最强，适合几乎所有场景。

### Q: 支持流式响应吗？

A: 是的！所有 3 个模型都支持流式响应。只需设置 `stream: true` 参数。详见 [流式功能文档](STREAMING.md)。

### Q: 如何处理长文本？

A: 所有模型都支持 200K tokens 上下文，足以处理大多数长文本场景。

### Q: 其他模型什么时候支持？

A: 这取决于 Factory AI 官方。建议关注 Factory AI 的更新公告，或查看我们的 [SUPPORTED_MODELS.md](SUPPORTED_MODELS.md) 文档了解最新测试结果。

---

## 🔗 相关文档

- [支持的模型详细测试报告](SUPPORTED_MODELS.md) - 完整的测试结果和分析
- [快速开始](QUICK_START.md) - 5分钟快速上手
- [流式功能](STREAMING.md) - 流式响应使用指南
- [API Key 代理](API-KEY-PROXY.md) - API Key 管理说明
- [OpenAI 兼容模式](README-OpenAI.md) - OpenAI SDK 使用指南

---

## 📈 测试记录

**测试日期**: 2025-10-08  
**测试方法**: 真实 API 调用  
**总测试数**: 27 个模型  

### ✅ 可用模型（3个）

| 模型 ID | 测试状态 | 响应 | 备注 |
|---------|---------|------|------|
| `claude-sonnet-4-5-20250929` | ✅ 通过 | 正常 | ⭐ 推荐 |
| `claude-sonnet-4-20250514` | ✅ 通过 | 正常 | 稳定 |
| `claude-3-7-sonnet-20250219` | ✅ 通过 | 正常 | 可靠 |

### ❌ 不可用模型（24个）

以下模型经测试均不可用：

**Claude 其他版本（5个）**:
- `claude-3-5-sonnet-20241022`, `claude-3-5-sonnet-20250219`
- `claude-sonnet-4-1-20250514`
- `claude-3-5-haiku-20241022`, `claude-3-haiku-20240307`
- 原因: HTTP 400 - Unsupported OpenAI model ID

**OpenAI/GPT 系列（10个）**:
- `gpt-5-*`, `o1-*`, `o3-*`, `o4-*`, `gpt-4o-*`
- 原因: 超时（30秒无响应）- Responses API 端点问题

**Gemini 系列（4个）**:
- `gemini-2.5-*`, `gemini-2.0-*`, `gemini-exp-*`
- 原因: HTTP 405 - Method not allowed

**Grok 系列（3个）**:
- `grok-4`, `grok-beta`, `grok-vision-beta`
- 原因: HTTP 405/400

**其他（2个）**:
- `glm-4.6`, 等
- 原因: HTTP 400

详细测试报告请查看: [SUPPORTED_MODELS.md](SUPPORTED_MODELS.md)

---

**最后更新**: 2025-10-08  
**支持的模型数**: 3 个（Claude 系列）  
**测试验证**: ✅ 完成