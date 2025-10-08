
# Factory Proxy - OpenAI 兼容接口

这是一个将 Factory AI API 转换为 OpenAI 兼容格式的代理服务器。使用 Go 语言编写，性能优异，支持标准 OpenAI SDK。

## 🌟 特性

- ✅ **完全兼容 OpenAI API** - 使用标准 `/v1/chat/completions` 端点
- ✅ **自动格式转换** - OpenAI 格式 ↔ Anthropic 格式无缝转换
- ✅ **支持 system 消息** - 自动处理 OpenAI 的 system 角色
- ✅ **Factory 集成** - 自动添加 Factory Droid system prompt
- ✅ **高性能** - Go 原生实现，启动快，内存占用低
- ✅ **详细日志** - 完整的请求/响应日志记录

## 🚀 快速开始

### 1. 编译

```bash
cd factory-go
go build -o factory-proxy-openai main-openai.go
```

### 2. 启动服务器

```bash
PORT=8003 ./factory-proxy-openai
```

输出：
```
🚀 Factory OpenAI-Compatible Proxy 启动中...
📍 端口: 8003
➡️  目标: https://gibuoilncyzqebelqjqz.supabase.co/functions/v1/smooth-handler/https://app.factory.ai/api/llm/a/v1/messages
✅ 服务器已启动，监听于 http://localhost:8003
📋 OpenAI兼容接口:
   - POST /v1/chat/completions -> 需要 Authorization: Bearer <factory-api-key>
   - GET /health 或 /v1/health -> 健康检查
```

### 3. 使用 curl 测试

```bash
curl -X POST http://localhost:8003/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_PROXY_API_KEY" \
  -d '{
    "model": "claude-sonnet-4-5-20250929",
    "messages": [
      {"role": "user", "content": "你好，请用中文简短回复"}
    ],
    "max_tokens": 50
  }'
```

响应（OpenAI 格式）：
```json
{
  "id": "msg_01Hn4rd36DodnMx6Ggpnv1Q1",
  "object": "chat.completion",
  "created": 1759899206,
  "model": "claude-sonnet-4-5-20250929",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "你好！我是 Droid，很高兴为你服务。有什么我可以帮助你的吗？"
      },
      "finish_reason": "stop"
    }
  ],
  "usage": {
    "prompt_tokens": 34,
    "completion_tokens": 37,
    "total_tokens": 71
  }
}
```

## 📚 使用示例

### Python + OpenAI SDK

```python
from openai import OpenAI

client = OpenAI(
    api_key="YOUR_PROXY_API_KEY",  # 使用代理 Key，不是 Factory Key
    base_url="http://localhost:8003/v1"
)

response = client.chat.completions.create(
    model="claude-sonnet-4-5-20250929",
    messages=[
        {"role": "system", "content": "You are a helpful assistant."},
        {"role": "user", "content": "Hello, how are you?"}
    ],
    max_tokens=100
)

print(response.choices[0].message.content)
```

### Node.js + OpenAI SDK

```javascript
import OpenAI from 'openai';

const client = new OpenAI({
  apiKey: 'YOUR_PROXY_API_KEY',  // 使用代理 Key
  baseURL: 'http://localhost:8003/v1'
});

const response = await client.chat.completions.create({
  model: 'claude-sonnet-4-5-20250929',
  messages: [
    { role: 'system', content: 'You are a helpful assistant.' },
    { role: 'user', content: 'Hello, how are you?' }
  ],
  max_tokens: 100
});

console.log(response.choices[0].message.content);
```

### 使用 system 消息

```bash
curl -X POST http://localhost:8003/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_PROXY_API_KEY" \
  -d '{
    "model": "claude-sonnet-4-5-20250929",
    "messages": [
      {"role": "system", "content": "You are a helpful math tutor."},
      {"role": "user", "content": "What is 25 * 4?"}
    ],
    "max_tokens": 50,
    "temperature": 0.7
  }'
```

## 🔧 API 端点

### POST /v1/chat/completions

标准 OpenAI chat completions 端点。

**请求头：**
- `Content-Type: application/json`
- `Authorization: Bearer <proxy-api-key>`  （使用 PROXY_API_KEY，不是 FACTORY_API_KEY）

**请求体：**
```json
{
  "model": "claude-sonnet-4-5-20250929",
  "messages": [
    {"role": "system", "content": "System prompt (可选)"},
    {"role": "user", "content": "User message"},
    {"role": "assistant", "content": "Assistant message (可选)"}
  ],
  "max_tokens": 100,
  "temperature": 0.7,
  "stream": false
}
```

**响应：**
```json
{
  "id": "msg_xxx",
  "object": "chat.completion",
  "created": 1234567890,
  "model": "claude-sonnet-4-5-20250929",
  "choices": [{
    "index": 0,
    "message": {
      "role": "assistant",
      "content": "Response text"
    },
    "finish_reason": "stop"
  }],
  "usage": {
    "prompt_tokens": 10,
    "completion_tokens": 20,
    "total_tokens": 30
  }
}
```

### GET /health 或 /v1/health

健康检查端点。

**响应：**
```json
{
  "status": "healthy",
  "timestamp": "2025-10-08T04:53:50Z",
  "uptime": 32.73
}
```

## 🔄 格式转换说明

### OpenAI → Anthropic

1. **messages 转换**
   - `system` 角色消息 → 提取并放入 `system` 字段
   - 自动添加 Factory Droid system prompt
   - `user` 和 `assistant` 消息保持不变

2. **参数映射**
   - 
   - `max_tokens` → `max_tokens` (保持不变)
   - `temperature` → `temperature` (保持不变)
   - `stream` → `stream` (保持不变)

### Anthropic → OpenAI

1. **响应格式转换**
   - Anthropic 的 `content` 数组 → OpenAI 的 `choices[0].message.content`
   - `id` 保持不变
   - 添加 OpenAI 标准字段：`object`, `created`

2. **finish_reason 映射**
   - `end_turn` → `stop`
   - `max_tokens` → `length`
   - `stop_sequence` → `stop`

3. **usage 统计转换**
   - `input_tokens` → `prompt_tokens`
   - `output_tokens` → `completion_tokens`
   - 自动计算 `total_tokens`

## 🎯 支持的模型

所有 Factory 支持的 Claude 模型都可以使用，例如：

- `claude-sonnet-4-5-20250929`
- `claude-sonnet-3-5-20240620`
- `claude-opus-4-20250514`

## ⚙️ 配置

### 环境变量

- `PORT` - 服务器监听端口（默认：8000）
- `ANTHROPIC_TARGET_URL` - Anthropic API 目标 URL（已预配置）

### 示例

```bash
# 自定义端口
PORT=9000 ./factory-proxy-openai
```

## 📊 性能特点

| 指标 | 数值 |
|------|------|
| **启动时间** | < 10ms |
| **内存占用** | ~11MB |
| **响应延迟** | ~5-6s (取决于网络) |
| **并发支持** | 高并发 |
| **二进制大小** | ~6MB |

## 🔍 日志示例

```
2025/10/08 12:53:20 收到OpenAI格式请求: POST /v1/chat/completions
2025/10/08 12:53:20 API Key已获取: fk-nTguzhI...
2025/10/08 12:53:20 OpenAI请求: model=claude-sonnet-4-5-20250929, messages数量=1
2025/10/08 12:53:20 已转换为Anthropic格式，请求体大小: 231 bytes
2025/10/08 12:53:26 收到响应: 状态码 200
2025/10/08 12:53:26 [POST] /v1/chat/completions - 200 - 5.8s
```

## 🚢 生产部署

### 使用 systemd

创建 `/etc/systemd/system/factory-proxy.service`：

```ini
[Unit]
Description=Factory OpenAI Proxy
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/factory-proxy
Environment="PORT=8003"
ExecStart=/opt/factory-proxy/factory-proxy-openai
Restart=always

[Install]
WantedBy=multi-user.target
```

启动服务：
```bash
sudo systemctl enable factory-proxy
sudo systemctl start factory-proxy
```

## 📝 许可证

MIT License

---

**注意**: 请妥善保管你的 Factory API Key，不要在公共代码库中泄露。