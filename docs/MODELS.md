
# 支持的模型列表

## ⚠️ 重要说明

**经过实际测试验证，Factory AI 后端目前仅支持 Claude 系列模型。**

通过 OpenAI 兼容格式的代理，您可以使用标准的 OpenAI SDK 访问所有 Claude 模型。

## ✅ 支持的模型 (Claude 系列)

Factory Proxy API 通过 Factory AI 后端支持以下 Claude 模型：

### 🤖 Claude 4.x ⭐ 最新推荐

- **`claude-sonnet-4-5-20250929`** - Claude Sonnet 4.5 ⭐ **强烈推荐**
  - ✅ 已测试验证可用
  - 🚀 最新版本，性能最强
  - 📊 最大 1M tokens 上下文（需 beta header，标准 200K）
  - 💡 适合所有场景：对话、代码生成、长文本分析

- **`claude-opus-4-1-20250805`** - Claude Opus 4.1
  - ✅ 已测试验证可用
  - 🧠 最强推理能力
  - 📊 200K tokens 上下文
  - 💡 适合复杂任务、深度分析、长时间推理

- `claude-sonnet-4-20250514` - Claude Sonnet 4
- `claude-3-7-sonnet-20250219` - Claude 3.7 Sonnet

### 🤖 Claude 3.x - 稳定可用

- **`claude-3-5-sonnet-20241022`** - Claude 3.5 Sonnet (2024-10)
  - 💰 性价比高
  - 📊 200K tokens 上下文
  - 💡 适合通用对话和标准任务

- **`claude-3-5-haiku-20241022`** - Claude 3.5 Haiku (2024-10)
  - ⚡ 响应最快
  - 💰 成本最优
  - 📊 200K tokens 上下文
  - 💡 适合简单对话、快速响应

- `anthropic.claude-3-haiku-20240307-v1:0` - Claude 3 Haiku (Bedrock)

---

## 🎯 模型选择建议

### 🏆 最佳全能：Claude Sonnet 4.5
```python
model = "claude-sonnet-4-5-20250929"
```
- 最新最强，适合 90% 的场景
- 1M 超长上下文，处理大型文档
- 代码生成、对话、分析全能

### 🧠 最强推理：Claude Opus 4.1
```python
model = "claude-opus-4-1-20250805"
```
- 需要深度思考时的最佳选择
- 复杂问题、科研分析、策略规划

### ⚡ 快速经济：Claude 3.5 Haiku
```python
model = "claude-3-5-haiku-20241022"
```
- 追求速度和成本优化
- 简单对话、实时交互

---

## 📖 使用示例

### Python (OpenAI SDK)

```python
from openai import OpenAI

client = OpenAI(
    api_key="YOUR_PROXY_API_KEY",  # 使用您的代理 Key
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
    stream=True  # 启用流式
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

# 流式请求（注意使用 -N 参数）
curl -N -X POST http://localhost:8003/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_PROXY_API_KEY" \
  -d '{
    "model": "claude-sonnet-4-5-20250929",
    "messages": [{"role": "user", "content": "讲个笑话"}],
    "max_tokens": 1024,
    "stream": true
  }'
```

### Node.js

```javascript
import OpenAI from 'openai';

const client = new OpenAI({
  apiKey: process.env.PROXY_API_KEY,  // 使用代理 Key
  baseURL: 'http://localhost:8003/v1'
});

// 非流式
const response = await client.chat.completions.create({
  model: 'claude-sonnet-4-5-20250929',
  messages: [{ role: 'user', content: 'Hello!' }],
  stream: false
});

console.log(response.choices[0].message.content);

// 流式
const stream = await client.chat.completions.create({
  model: 'claude-sonnet-4-5-20250929',
  messages: [{ role: 'user', 
content: 'Hello!' }],
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

所有 Claude 模型使用相同的 API 格式，只需更改 `model` 参数即可切换：

```python
# 快速测试不同 Claude 模型
models_to_try = [
    "claude-sonnet-4-5-20250929",   # 最强全能
    "claude-opus-4-1-20250805",     # 最强推理
    "claude-3-5-haiku-20241022"     # 最快最省
]

for model in models_to_try:
    response = client.chat.completions.create(
        model=model,
        messages=[{"role": "user", "content": "你好，介绍一下你自己"}],
        max_tokens=200
    )
    print(f"\n{'='*50}")
    print(f"模型: {model}")
    print(f"{'='*50}")
    print(response.choices[0].message.content)
```

---

## 📊 模型对比

| 模型 | 版本 | 上下文 | 特点 | 适用场景 | 测试状态 |
|------|------|--------|------|----------|---------|
| **Claude Sonnet 4.5** | `claude-sonnet-4-5-20250929` | 1M | 最新最强、长文本 | 全能型任务 | ✅ 已验证 |
| **Claude Opus 4.1** | `claude-opus-4-1-20250805` | 200K | 最强推理 | 复杂分析 | ✅ 已验证 |
| **Claude 3.5 Sonnet** | `claude-3-5-sonnet-20241022` | 200K | 性价比高 | 通用对话 | ✅ 可用 |
| **Claude 3.5 Haiku** | `claude-3-5-haiku-20241022` | 200K | 快速经济 | 简单任务 | ✅ 可用 |

---

## 📏 上下文长度限制

| 模型名称 | 最大上下文 | 输入限制 | 输出限制 | 说明 |
|---------|-----------|---------|---------|------|
| **Claude 4.5 Sonnet**<br/>`claude-sonnet-4-5-20250929` | 1,000,000 | ~900K | ~100K | 通过 beta header 支持 1M，标准为 200K ⭐ |
| **Claude 4 Opus**<br/>`claude-opus-4-1-20250805` | 200,000 | ~168K | ~32K | 标准窗口，专注推理 |
| **Claude 3.5 Sonnet**<br/>`claude-3-5-sonnet-20241022` | 200,000 | ~168K | ~32K | 性价比之选 |
| **Claude 3.5 Haiku**<br/>`claude-3-5-haiku-20241022` | 200,000 | ~168K | ~32K | 快速且经济 |

### 💡 使用建议

1. **长文档处理**: 使用 `claude-sonnet-4-5-20250929` (1M tokens)
2. **代码分析**: `claude-sonnet-4-5-20250929` (1M tokens)
3. **快速对话**: `claude-3-5-haiku-20241022` (200K tokens)
4. **复杂推理**: `claude-opus-4-1-20250805` (200K tokens)

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

### 📊 推荐的 max_tokens 设置

```python
# 简短对话
max_tokens = 1024      # ~1K tokens

# 标准响应 (推荐)
max_tokens = 4096      # ~4K tokens

# 长文本生成
max_tokens = 8192      # ~8K tokens

# 超长输出 (Claude 4.5)
max_tokens = 16384     # ~16K tokens

# 极限输出 (仅 Claude 4.5)
max_tokens = 32768     # ~32K tokens (注意成本)
```

---

## 🎯 场景推荐

### 💻 代码生成与分析
```python
model = "claude-sonnet-4-5-20250929"  # 最佳选择
# 或
model = "claude-opus-4-1-20250805"    # 复杂重构
```

### 📚 长文档处理
```python
model = "claude-sonnet-4-5-20250929"  # 1M 上下文
```

### 💬 快速对话
```python
model = "claude-3-5-haiku-20241022"   # 速度优先
```

### 🧠 复杂推理
```python
model = "claude-opus-4-1-20250805"    # 推理能力最强
```

### 🌏 中文场景
```python
model = "claude-sonnet-4-5-20250929"  # 中英文全能
```

---

## 📝 注意事项

1. **模型可用性**: 目前仅支持 Claude 系列模型
2. **API 兼容性**: 完全兼容 OpenAI SDK 和 API 格式
3. **定价差异**: 不同模型价格不同，请参考 Factory AI 
官方定价
4. **速率限制**: 根据您的订阅计划可能有不同的速率限制
5. **上下文长度**: 
   - Claude 4.5 Sonnet: 最大 1M tokens (需 beta header)
   - 其他 Claude 模型: 标准 200K tokens
6. **输出限制**: 不同模型有不同的最大输出 token 限制
7. **成本优化**: 使用合适的 `max_tokens` 参数控制成本
8. **流式支持**: 所有模型都支持流式（`stream: true`）和非流式响应

---

## ❓ 常见问题

### Q: 为什么只支持 Claude 系列？

A: Factory AI 后端目前仅通过 Anthropic Messages API 提供服务，因此只支持 Claude 系列模型。虽然某些文档中提到了 GPT-5、Gemini 等模型，但实际测试这些模型会返回 "Unknown Anthropic model ID" 错误。

### Q: 如何选择合适的模型？

A: 
- **首选**: `claude-sonnet-4-5-20250929` - 适合 90% 的场景
- **深度思考**: `claude-opus-4-1-20250805` - 需要最强推理时
- **速度优先**: `claude-3-5-haiku-20241022` - 快速响应、低成本

### Q: 支持流式响应吗？

A: 是的！所有 Claude 模型都支持流式响应。只需设置 `stream: true` 参数即可。详见 [流式功能文档](STREAMING.md)。

### Q: 如何处理长文本？

A: 使用 `claude-sonnet-4-5-20250929`，它支持最大 1M tokens 的上下文窗口，可以处理超大型文档和代码库。

---

## 🔗 相关文档

- [快速开始](QUICK_START.md) - 5分钟快速上手
- [流式功能](STREAMING.md) - 流式响应使用指南
- [API Key 代理](API-KEY-PROXY.md) - API Key 管理说明
- [README](README.md) - 项目主文档
- [OpenAI 兼容模式](README-OpenAI.md) - OpenAI SDK 使用指南

---

## 📈 测试记录

**测试日期**: 2025-10-08

### ✅ 已验证可用的模型

| 模型 ID | 测试状态 | 响应时间 | 备注 |
|---------|---------|---------|------|
| `claude-sonnet-4-5-20250929` | ✅ 通过 | 正常 | 推荐使用 |
| `claude-opus-4-1-20250805` | ✅ 通过 | 正常 | 最强推理 |
| `claude-3-5-sonnet-20241022` | ✅ 可用 | 正常 | - |
| `claude-3-5-haiku-20241022` | ✅ 可用 | 快速 | - |

### ❌ 已测试不可用的模型

以下模型经测试不可用，返回 "Unknown Anthropic model ID" 错误：

- `gpt-5`, `gpt-5-2025-08-07`, `gpt-5-codex`
- `gemini-2.5-pro`, `gemini-2.5-flash`
- `grok-4`
- `o1`, `o3`, `o4-mini`
- `glm-4.6`

**结论**: Factory AI 后端仅支持 Claude 系列模型。

---

**最后更新**: 2025-10-08  
**支持的模型数**: 7+ (Claude 系列)  
**测试验证**: ✅ 完成