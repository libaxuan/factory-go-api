# 🌊 流式响应功能文档

Factory Proxy API 完整支持 Server-Sent Events (SSE) 流式响应，让您可以实时接收 AI 生成的内容。

---

## 📖 目录

- [什么是流式响应](#什么是流式响应)
- [为什么使用流式响应](#为什么使用流式响应)
- [快速开始](#快速开始)
- [详细示例](#详细示例)
- [SSE 格式说明](#sse-格式说明)
- [错误处理](#错误处理)
- [最佳实践](#最佳实践)
- [性能考虑](#性能考虑)
- [故障排除](#故障排除)

---

## 什么是流式响应

流式响应（Streaming）允许服务器在生成内容的同时实时发送给客户端，而不是等待完整响应生成后再一次性返回。

### 对比

| 特性 | 非流式 | 流式 |
|------|--------|------|
| **响应方式** | 一次性返回完整内容 | 实时逐块返回 |
| **用户体验** | 等待时间长 | 即时反馈 |
| **适用场景** | 短文本、批处理 | 长文本、交互式对话 |
| **实现复杂度** | 简单 | 稍复杂 |

---

## 为什么使用流式响应

✅ **更好的用户体验**
- 用户立即看到响应开始
- 减少感知延迟
- 类似 ChatGPT 的打字机效果

✅ **适合长文本生成**
- 代码生成
- 文章写作
- 长篇对话

✅ **更好的交互性**
- 用户可以提前看到部分结果
- 可以提前中断不需要的生成

---

## 快速开始

### Python

```python
from openai import OpenAI

client = OpenAI(
    api_key="YOUR_PROXY_API_KEY",
    base_url="http://localhost:8003/v1"
)

# 启用流式
stream = client.chat.completions.create(
    model="claude-sonnet-4-5-20250929",
    messages=[{"role": "user", "content": "写一个排序算法"}],
    max_tokens=1000,
    stream=True  # 🔑 关键参数
)

# 逐块处理响应
for chunk in stream:
    if chunk.choices[0].delta.content:
        print(chunk.choices[0].delta.content, end="", flush=True)
```

### Node.js

```javascript
import OpenAI from 'openai';

const client = new OpenAI({
  apiKey: process.env.PROXY_API_KEY,
  baseURL: 'http://localhost:8003/v1'
});

const stream = await client.chat.completions.create({
  model: 'claude-sonnet-4-5-20250929',
  messages: [{ role: 'user', content: '写一个排序算法' }],
  max_tokens: 1000,
  stream: true
});

for await (const chunk of stream) {
  if (chunk.choices[0]?.delta?.content) {
    process.stdout.write(chunk.choices[0].delta.content);
  }
}
```

### cURL

```bash
curl -N -X POST http://localhost:8003/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_PROXY_API_KEY" \
  -d '{
    "model": "claude-sonnet-4-5-20250929",
    "messages": [{"role": "user", "content": "写一个排序算法"}],
    "max_tokens": 1000,
    "stream": true
  }'
```

**注意**: `-N` 或 `--no-buffer` 参数很重要，它禁用缓冲以实时显示流式内容。

---

## 详细示例

### 示例 1: 基础流式对话

```python
from openai import OpenAI

client = OpenAI(
    api_key="YOUR_PROXY_API_KEY",
    base_url="http://localhost:8003/v1"
)

stream = client.chat.completions.create(
    model="claude-sonnet-4-5-20250929",
    messages=[
        {"role": "user", "content": "介绍一下 Python 的主要特点"}
    ],
    stream=True
)

print("AI: ", end="", flush=True)
for chunk in stream:
    if chunk.choices[0].delta.content:
        print(chunk.choices[0].delta.content, end="", flush=True)
print()  # 换行
```

### 示例 2: 收集完整响应

```python
stream = client.chat.completions.create(
    model="claude-sonnet-4-5-20250929",
    messages=[{"role": "user", "content": "生成一个 TODO 列表"}],
    stream=True
)

full_response = ""
for chunk in stream:
    if chunk.choices[0].delta.content:
        content = chunk.choices[0].delta.content
        full_response += content
        print(content, end="", flush=True)

print(f"\n\n完整响应长度: {len(full_response)} 字符")
```

### 示例 3: 带系统提示的流式

```python
stream = client.chat.completions.create(
    model="claude-sonnet-4-5-20250929",
    messages=[
        {"role": "system", "content": "你是一个专业的 Python 程序员"},
        {"role": "user", "content": "写一个快速排序的实现"}
    ],
    temperature=0.7,
    max_tokens=2000,
    stream=True
)

for chunk in stream:
    if chunk.choices[0].delta.content:
        print(chunk.choices[0].delta.content, end="", flush=True)
```

---

## SSE 格式说明

流式响应使用 **Server-Sent Events (SSE)** 格式。

### 响应头

```
Content-Type: text/event-stream
Cache-Control: no-cache
Connection: keep-alive
```

### 数据格式

每个事件的格式：

```
data: {"id":"chatcmpl-xxx","object":"chat.completion.chunk","created":1234567890,"model":"claude-sonnet-4-5-20250929","choices":[{"index":0,"delta":{"role":"assistant"},"finish_reason":null}]}

data: {"id":"chatcmpl-xxx","object":"chat.completion.chunk","created":1234567890,"model":"claude-sonnet-4-5-20250929","choices":[{"index":0,"delta":{"content":"Hello"},"finish_reason":null}]}

data: {"id":"chatcmpl-xxx","object":"chat.completion.chunk","created":1234567890,"model":"claude-sonnet-4-5-20250929","choices":[{"index":0,"delta":{},"finish_reason":"stop"}]}

data: [DONE]
```

---

## 错误处理

### 处理连接错误

```python
from openai import OpenAI, OpenAIError

try:
    