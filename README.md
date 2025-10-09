# Factory Go API

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

**高性能 Factory AI API 网关 | 多模型支持 | OpenAI 兼容格式**

</div>

---

## 📖 简介

Factory Go API 是一个高性能的 AI 模型网关，提供 OpenAI 兼容的统一接口。支持 Claude、GPT 等多个 AI 模型家族。

### 核心特性

- ⚡ **极致性能** - Go 原生实现，启动 < 10ms，内存占用 ~11MB
- 🎯 **多模型支持** - 统一接口访问 Claude、GPT 等 5 个模型
- 🔄 **OpenAI 兼容** - 使用标准 OpenAI SDK 即可调用所有模型
- 🌊 **双模式响应** - 完整支持流式（SSE）和非流式响应
- 🔐 **安全认证** - Bearer Token 认证保护
- 📊 **详细日志** - 完整的请求/响应日志

## 🎯 支持的模型

**5 个 AI 模型**，覆盖两大家族，**每个模型都支持流式和非流式** = **10 种配置全部可用**：

### Claude 系列 (Anthropic)
| 模型 ID | 描述 | 流式 | 非流式 |
|---------|------|------|--------|
| `claude-opus-4-1-20250805` | 🧠 Claude Opus 4.1 - 最强推理 | ✅ | ✅ |
| `claude-sonnet-4-20250514` | ⚡ Claude Sonnet 4 - Extended Thinking 中等 | ✅ | ✅ |
| `claude-sonnet-4-5-20250929` | ⭐ Claude Sonnet 4.5 - Extended Thinking 高级（推荐） | ✅ | ✅ |

### GPT 系列 (OpenAI)
| 模型 ID | 描述 | 流式 | 非流式 |
|---------|------|------|--------|
| `gpt-5-2025-08-07` | 🚀 GPT-5 - 最新旗舰，Extended Thinking | ✅ | ✅ |
| `gpt-5-codex` | 💻 GPT-5 Codex - 代码生成专家 | ✅ | ✅ |

## 🚀 快速开始

### 安装

```bash
# 克隆仓库
git clone https://github.com/libaxuan/factory-go-api.git
cd factory-go-api

# 配置环境变量
cp .env.example .env
# 编辑 .env，设置 FACTORY_API_KEY
```

### 启动

**macOS/Linux:**
```bash
./start.sh
# 服务运行在 http://localhost:8003
```

**Windows:**
```cmd
start.bat
# 服务运行在 http://localhost:8003
```

### 测试

```bash
# 快速健康检查
curl http://localhost:8003/health

# 查看支持的模型
curl http://localhost:8003/v1/models \
  -H "Authorization: Bearer YOUR_API_KEY"

# 运行完整测试（测试所有 10 种配置）
./restart_and_test.sh  # 重启服务并测试
# 或直接测试
./test_models.sh  # macOS/Linux
test_models.bat   # Windows

# 测试单个模型（非流式）
curl -X POST http://localhost:8003/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -d '{
    "model": "gpt-5-2025-08-07",
    "messages": [{"role": "user", "content": "Hello!"}],
    "stream": false
  }'

# 测试流式响应
curl -X POST http://localhost:8003/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -d '{
    "model": "claude-sonnet-4-5-20250929",
    "messages": [{"role": "user", "content": "Tell me a story"}],
    "stream": true
  }'
```

## 💻 代码示例

### Python

```python
from openai import OpenAI

client = OpenAI(
    base_url="http://localhost:8003/v1",
    api_key="YOUR_FACTORY_API_KEY"
)

# 非流式
response = client.chat.completions.create(
    model="claude-sonnet-4-5-20250929",
    messages=[{"role": "user", "content": "Hello!"}]
)
print(response.choices[0].message.content)

# 流式
stream = client.chat.completions.create(
    model="gpt-5-2025-08-07",
    messages=[{"role": "user", "content": "Tell me a story"}],
    stream=True
)
for chunk in stream:
    if chunk.choices[0].delta.content:
        print(chunk.choices[0].delta.content, end="")
```

### JavaScript/TypeScript

```javascript
import OpenAI from 'openai';

const client = new OpenAI({
  baseURL: 'http://localhost:8003/v1',
  apiKey: 'YOUR_FACTORY_API_KEY'
});

// 非流式
const response = await client.chat.completions.create({
  model: 'claude-sonnet-4-5-20250929',
  messages: [{ role: 'user', content: 'Hello!' }]
});
console.log(response.choices[0].message.content);

// 流式
const stream = await client.chat.completions.create({
  model: 'gpt-5-2025-08-07',
  messages: [{ role: 'user', content: 'Tell me a story' }],
  stream: true
});
for await (const chunk of stream) {
  process.stdout.write(chunk.choices[0]?.delta?.content || '');
}
```

### Go

```go
package main

import (
    "context"
    "fmt"
    "github.com/sashabaranov/go-openai"
)

func main() {
    config := openai.DefaultConfig("YOUR_FACTORY_API_KEY")
    config.BaseURL = "http://localhost:8003/v1"
    client := openai.NewClientWithConfig(config)

    resp, err := client.CreateChatCompletion(
        context.Background(),
        openai.ChatCompletionRequest{
            Model: "claude-sonnet-4-5-20250929",
            Messages: []openai.ChatCompletionMessage{
                {
                    Role:    openai.ChatMessageRoleUser,
                    Content: "Hello!",
                },
            },
        },
    )
    if err != nil {
        panic(err)
    }
    fmt.Println(resp.Choices[0].Message.Content)
}
```

## ⚙️ 配置

### 环境变量

在 `.env` 文件中配置：

```bash
# 必需
FACTORY_API_KEY=your_factory_api_key

# 可选
PORT=8003
CONFIG_PATH=config.json
```

### 模型配置

编辑 `config.json` 添加或修改模型：

```json
{
  "port": 8003,
  "endpoints": [
    {
      "name": "anthropic",
      "base_url": "https://app.factory.ai/api/llm/a/v1/messages"
    },
    {
      "name": "openai",
      "base_url": "https://app.factory.ai/api/llm/o/v1/responses"
    }
  ],
  "models": [
    {
      "name": "Claude Sonnet 4.5",
      "id": "claude-sonnet-4-5-20250929",
      "type": "anthropic",
      "reasoning": "high"
    },
    {
      "name": "GPT-5",
      "id": "gpt-5-2025-08-07",
      "type": "openai",
      "reasoning": "high"
    }
  ]
}
```

## 🔌 API 端点

| 端点 | 方法 | 描述 |
|------|------|------|
| `/health` | GET | 健康检查 |
| `/v1/models` | GET | 模型列表 |
| `/v1/chat/completions` | POST | 聊天补全（OpenAI 兼容） |
| `/docs` | GET | API 文档页面 |

## 📊 性能

| 指标 | 数值 |
|------|------|
| 启动时间 | < 10ms |
| 内存占用 | ~11MB |
| 二进制大小 | ~8MB |
| 并发性能 | 优秀 |

## 🔧 开发

```bash
# 安装依赖
go mod tidy

# 开发模式运行
go run main_multimodel.go

# 构建
go build -o factory-api main_multimodel.go

# 格式化代码
gofmt -w .

# 代码检查
go vet ./...

# 运行完整模型测试
chmod +x restart_and_test.sh
./restart_and_test.sh  # 重启服务并测试所有 7 个模型配置
```

## 🚢 部署

### Docker

```bash
# 构建镜像
docker build -t factory-api .

# 运行容器
docker run -d \
  -p 8003:8003 \
  -e FACTORY_API_KEY=your_key \
  --name factory-api \
  factory-api
```

### systemd

创建 `/etc/systemd/system/factory-api.service`:

```ini
[Unit]
Description=Factory API Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/factory-api

Environment="FACTORY_API_KEY=your_key"
ExecStart=/opt/factory-api/factory-api
Restart=always

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl enable factory-api
sudo systemctl start factory-api
```

## 🔐 安全建议

1. **保护 API Key** - 使用环境变量，不要硬编码
2. **使用 HTTPS** - 生产环境配置反向代理（Nginx/Caddy）
3. **限流保护** - 在代理层面配置限流规则
4. **日志管理** - API Key 已脱敏显示

## 🆘 故障排除

### 端口被占用

```bash
# 查看端口占用
lsof -i :8003

# 使用其他端口
PORT=9000 ./start.sh
```

### 403 Forbidden

确保：
- 使用正确的 Factory API Key
- 请求包含 Authorization 头
- 环境变量正确配置

### 连接超时

```bash
# 检查目标服务
curl -I https://app.factory.ai

# 检查防火墙和网络
```

## 📝 更新日志

查看 [CHANGELOG.md](CHANGELOG.md) 了解版本历史和更新内容。

## 🤝 贡献

欢迎贡献！请遵循以下步骤：

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing`)
3. 提交更改 (`git commit -m 'feat: add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing`)
5. 创建 Pull Request

提交信息规范：
- `feat:` 新功能
- `fix:` Bug 修复
- `docs:` 文档更新
- `refactor:` 代码重构
- `test:` 测试相关

## 📄 许可证

本项目采用 [MIT License](LICENSE) 开源。

## 🔗 相关链接

- [Factory AI 官网](https://factory.ai)
- [OpenAI API 文档](https://platform.openai.com/docs)
- [问题反馈](https://github.com/libaxuan/factory-go-api/issues)

---

<div align="center">

**Made with ❤️ by Factory Go API Team**

如果这个项目对你有帮助，请给个 ⭐ Star！

</div>