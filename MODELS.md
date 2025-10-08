
# 支持的模型列表

Factory Proxy API 支持以下所有模型。您可以在请求中使用任何这些模型 ID。

## 📋 完整模型列表

### 🤖 Claude 系列 (Anthropic)

#### Claude 3.x
- `claude-3-5-sonnet-20241022` - Claude 3.5 Sonnet (2024-10)
- `claude-3-5-haiku-20241022` - Claude 3.5 Haiku (2024-10)
- `anthropic.claude-3-haiku-20240307-v1:0` - Claude 3 Haiku (Bedrock)

#### Claude 4.x ⭐ 最新
- `claude-3-7-sonnet-20250219` - Claude 3.7 Sonnet
- `claude-sonnet-4-20250514` - Claude Sonnet 4
- `claude-sonnet-4-5-20250929` - Claude Sonnet 4.5 ⭐ 推荐
- `claude-opus-4-1-20250805` - Claude Opus 4.1 (最强推理)

### 🔷 Gemini 系列 (Google)
- `gemini-2.5-flash` - Gemini 2.5 Flash (快速)
- `gemini-2.5-pro` - Gemini 2.5 Pro (强大)

### 🟢 OpenAI 系列

#### O 系列 (推理模型)
- `o1` - O1 (推理优化)
- `o3` - O3 (增强推理)
- `o4-mini` - O4 Mini (轻量推理)
- `o4-mini-alpha-2025-07-11` - O4 Mini Alpha

#### GPT-4 系列
- `gpt-4o` - GPT-4 Optimized
- `gpt-4.1` - GPT-4.1

#### GPT-5 系列 ⭐ 最新
- `gpt-5-2025-08-07` - GPT-5 标准版
- `gpt-5-mini-2025-08-07` - GPT-5 Mini
- `gpt-5-nano-2025-08-07` - GPT-5 Nano (超轻量)
- `gpt-5-codex` - GPT-5 Codex (代码专用)
- `gpt-5-reasoning-alpha-2025-07-17` - GPT-5 推理 Alpha

#### 特殊模型
- `nectarine-alpha-2025-07-24` - Nectarine Alpha
- `nectarine-alpha-new-reasoning-effort-2025-07-25` - Nectarine 推理增强版

### 🦅 Grok 系列 (xAI)
- `grok-4` - Grok 4 (最新)

### 🇨🇳 GLM 系列 (智谱)
- `glm-4.6` - GLM 4.6

---

## 🎯 推荐模型

### 最佳平衡 (性能/成本)
- `claude-sonnet-4-5-20250929` - 最新 Claude，性能强大
- `gpt-5-mini-2025-08-07` - 轻量级 GPT-5
- `gemini-2.5-flash` - 快速响应

### 最强推理能力
- `claude-opus-4-1-20250805` - Claude 最强推理
- `o3` - OpenAI 推理优化
- `gpt-5-reasoning-alpha-2025-07-17` - GPT-5 推理版

### 代码专用
- `gpt-5-codex` - 专为代码设计
- `claude-sonnet-4-5-20250929` - 全能代码助手

### 成本优化
- `claude-3-5-haiku-20241022` - 快速且经济
- `gpt-5-nano-2025-08-07` - 超轻量级
- `gemini-2.5-flash` - 高性价比

---

## 📖 使用示例

### Python (OpenAI SDK)

```python
from openai import OpenAI

client = OpenAI(
    api_key="YOUR_PROXY_API_KEY",  # 使用代理 Key
    base_url="http://localhost:8003/v1"
)

# 使用 Claude Sonnet 4.5 (推荐)
response = client.chat.completions.create(
    model="claude-sonnet-4-5-20250929",
    messages=[
        {"role": "user", "content": "Hello!"}
    ]
)

# 使用 GPT-5
response = client.chat.completions.create(
    model="gpt-5-2025-08-07",
    messages=[
        {"role": "user", "content": "Hello!"}
    ]
)

# 使用 Gemini
response = client.chat.completions.create(
    model="gemini-2.5-pro",
    messages=[
        {"role": "user", "content": "Hello!"}
    ]
)
```

### cURL

```bash
# Claude Sonnet 4.5
curl -X POST http://localhost:8003/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_PROXY_API_KEY" \
  -d '{
    "model": "claude-sonnet-4-5-20250929",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'

# GPT-5 Mini
curl -X POST http://localhost:8003/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_PROXY_API_KEY" \
  -d '{
    "model": "gpt-5-mini-2025-08-07",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'

# Grok 4
curl -X POST http://localhost:8003/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_PROXY_API_KEY" \
  -d '{
    "model": "grok-4",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

### Node.js

```javascript
import OpenAI from 'openai';

const client = new OpenAI({
  apiKey: process.env.PROXY_API_KEY,  // 使用代理 Key
  baseURL: 'http://localhost:8003/v1'
});

// 使用不同模型
const models = [
  'claude-sonnet-4-5-20250929',
  'gpt-5-2025-08-07',
  'gemini-2.5-pro',
  'grok-4'
];

for (const model of models) {
  const response = await client.chat.completions.create({
    model: model,
    messages: [{ role: 'user', 
content: "Hello!" }]
  });
  console.log(`${model}: ${response.choices[0].message.content}`);
}
```

---

## 🔄 模型切换

所有模型使用相同的 API 格式，只需更改 `model` 参数即可切换：

```python
# 快速切换模型
models_to_try = [
    "claude-sonnet-4-5-20250929",  # 最佳全能
    "gpt-5-2025-08-07",            # GPT-5 标准
    "gemini-2.5-pro",              # Gemini Pro
    "grok-4"                       # Grok 最新
]

for model in models_to_try:
    response = client.chat.completions.create(
        model=model,
        messages=[{"role": "user", "content": "你好"}],
        max_tokens=100
    )
    print(f"{model}: {response.choices[0].message.content}")
```

---

## 📊 模型对比

| 模型系列 | 最新版本 | 特点 | 适用场景 |
|---------|---------|------|---------|
| Claude 4.x | `claude-sonnet-4-5-20250929` | 强大推理、长文本 | 全能型任务 |
| GPT-5 | `gpt-5-2025-08-07` | 最新 OpenAI、多模态 | 通用对话 |
| Gemini 2.5 | `gemini-2.5-pro` | 快速、高效 | 实时应用 |
| Grok | `grok-4` | 实时信息、幽默 | 信息检索 |
| GLM | `glm-4.6` | 中文优化 | 中文场景 |

---

## 📏 模型上下文限制

不同模型支持的最大上下文 token 数量不同，这会影响您可以处理的文本长度和对话历史：

| 模型名称 | 最大上下文 Tokens | 输入限制 | 输出限制 | 说明 |
|---------|------------------|---------|---------|------|
| **Claude 4.5 Sonnet**<br/>`claude-sonnet-4-5-20250929` | 1,000,000 | ~900K | ~100K | 通过 beta header 支持 1M，标准为 200K；适用于复杂代理和编码任务 ⭐ |
| **Claude 4 Opus**<br/>`claude-opus-4-1-20250805` | 200,000 | ~168K | ~32K | 标准窗口，专注于长运行任务 |
| **GPT-5**<br/>`gpt-5-2025-08-07` | 400,000 | ~272K | ~128K | API 总和；ChatGPT Plus 为 32K，Pro 为 128K |
| **GPT-5 Codex**<br/>`gpt-5-codex` | 192,000 | ~160K | ~32K | 优化编码，窗口较小以提升效率 |
| **Gemini 2.5 Pro**<br/>`gemini-2.5-pro` | 1,000,000<br/>(即将 2M) | ~983K | ~65K | Vertex AI 支持 1,048,576 输入；多模态强 |
| **Claude 3.5 Haiku**<br/>`claude-3-5-haiku-20241022` | 200,000 | ~168K | ~32K | 快速且经济的选择 |
| **Gemini 2.5 Flash**<br/>`gemini-2.5-flash` | 1,000,000 | ~983K | ~65K | 高性价比，快速响应 |

### 💡 上下文使用建议

1. **长文档处理**: 使用 `claude-sonnet-4-5-20250929` 或 `gemini-2.5-pro` (1M tokens)
2. **代码分析**: `gpt-5-codex` (192K) 或 `claude-sonnet-4-5` (1M)
3. **快速对话**: `claude-3-5-haiku` 或 `gemini-2.5-flash` (200K/1M)
4. **成本优化**: 根据实际需求选择合适的上下文长度

### ⚠️ 注意事项

- **实际可用长度**: 取决于平台、订阅计划和 beta 功能配置
- **输入+输出总和**: 大多数模型的上下文限制是输入和输出 tokens 的总和
- **成本考虑**: 更大的上下文窗口通常意味着更高的成本
- **性能影响**: 极长的上下文可能影响响应速度

---

## ⚙️ 模型参数

所有模型支持以下通用参数：

```python
response = client.chat.completions.create(
    model="claude-sonnet-4-5-20250929",
    messages=[
        {"role": "system", "content": "You are a helpful assistant."},
        {"role": "user", "content": "Hello!"}
    ],
    max_tokens=4096,       # 最大输出长度 (Claude 4.5 最大 ~100K)
    temperature=0.7,       # 随机性 (0-2)
    top_p=1.0,            # 核采样
    n=1,                  # 生成数量
    stream=False,         # 流式输出
    stop=None             # 停止词
)
```

### 📊 推荐的 max_tokens 设置

根据不同场景和模型选择合适的 `max_tokens`：

```python
# 简短对话 (所有模型)
max_tokens = 1024  # ~1K tokens

# 标准响应 (推荐)
max_tokens = 4096  # ~4K tokens

# 长文本生成 (Claude/Gemini)
max_tokens = 8192  # ~8K tokens

# 超长输出 (Claude 4.5 / Gemini 2.5)
max_tokens = 16384  # ~16K tokens

# 极限输出 (仅 Claude 4.5 Sonnet/Opus, Gemini 2.5)
max_tokens = 32768  # ~32K tokens (需考虑成本)
```

---

## 🎯 场景推荐

### 代码生成
```python
# 推荐: gpt-5-codex 或 claude-sonnet-4-5
model = "gpt-5-codex"
```

### 长文本分析
```python
# 推荐: claude-sonnet-4-5 (1M tokens) 或 gemini-2.5-pro (1M tokens)
model = "claude-sonnet-4-5-20250929"  # 最佳选择，支持 1M 上下文
# 或
model = "gemini-2.5-pro"              # 多模态，1M 上下文
# 或
model = "claude-opus-4-1-20250805"    # 200K 上下文
```

### 快速对话
```python
# 推荐: claude-haiku 或 gemini-flash
model = "claude-3-5-haiku-20241022"
```

### 复杂推理
```python
# 推荐: o3 或 claude-opus-4-1
model = "o3"
```

### 中文场景
```python
# 推荐: glm-4.6
model = "glm-4.6"
```

---

## 📝 注意事项

1. **模型可用性**: 某些模型可能需要特定的 API 权限
2. **定价差异**: 不同模型价格不同，请参考 Factory AI 官方定价
3. **速率限制**: 根据您的订阅计划可能有不同的速率限制
4. **上下文长度**: 
   - Claude 4.5 Sonnet: 最大 1M tokens (需 beta header)
   - Gemini 2.5 Pro: 最大 1M tokens (即将支持 2M)
   - GPT-5: 最大 400K tokens (272K 输入 + 128K 输出)
   - GPT-5 Codex: 最大 192K tokens
   - Claude 4 Opus: 最大 200K tokens
5. **输出限制**: 注意不同模型的最大输出 token 数限制
6. **成本优化**: 使用合适的 `max_tokens` 参数控制成本

---

## 🔗 相关文档

- [快速开始](QUICK_START.md) - 5分钟快速上手
- [README](README.md) - 项目主文档
- [OpenAI 兼容模式](README-OpenAI.md) - OpenAI SDK 使用指南
- [贡献指南](CONTRIBUTING.md) - 如何贡献代码

---

**最后更新**: 2025-10-08  
**支持的模型数**: 25+