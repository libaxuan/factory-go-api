
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
    api_key="YOUR_FACTORY_API_KEY",
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
  -H "Authorization: Bearer YOUR_FACTORY_API_KEY" \
  -d '{
    "model": "claude-sonnet-4-5-20250929",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'

# GPT-5 Mini
curl -X POST http://localhost:8003/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_FACTORY_API_KEY" \
  -d '{
    "model": "gpt-5-mini-2025-08-07",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'

# Grok 4
curl -X POST http://localhost:8003/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_FACTORY_API_KEY" \
  -d '{
    "model": "grok-4",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

### Node.js

```javascript
import OpenAI from 'openai';

const client = new OpenAI({
  apiKey: process.env.FACTORY_API_KEY,
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

## ⚙️ 模型参数

所有模型支持以下通用参数：

```python
response = client.chat.completions.create(
    model="claude-sonnet-4-5-20250929",
    messages=[
        {"role": "system", "content": "You are a helpful assistant."},
        {"role": "user", "content": "Hello!"}
    ],
    max_tokens=1000,        # 最大输出长度
    temperature=0.7,        # 随机性 (0-2)
    top_p=1.0,             # 核采样
    n=1,                   # 生成数量
    stream=False,          # 流式输出
    stop=None              # 停止词
)
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
# 推荐: claude-opus-4-1 (100K+ tokens)
model = "claude-opus-4-1-20250805"
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
4. **上下文长度**: 不同模型支持的最大 token 数不同

---

## 🔗 相关文档

- [快速开始](QUICK_START.md) - 5分钟快速上手
- [README](README.md) - 项目主文档
- [OpenAI 兼容模式](README-OpenAI.md) - OpenAI SDK 使用指南
- [贡献指南](CONTRIBUTING.md) - 如何贡献代码

---

**最后更新**: 2025-01-08  
**支持的模型数**: 25+