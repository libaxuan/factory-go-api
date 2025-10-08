
# Factory Proxy

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

**高性能 Factory AI API 代理服务器 | 支持 OpenAI 兼容格式**

[English](README.md) | [简体中文](README.md)

</div>

---

## 📖 简介

Factory Proxy 是一个用 Go 语言编写的高性能代理服务器，专为 Factory AI API 设计。它提供两种工作模式：

1. **Anthropic 原生模式** - 直接代理 Factory AI 的原生 API
2. **OpenAI 兼容模式** ⭐ - 将 Factory AI 转换为标准 OpenAI API 格式

使用 OpenAI 兼容模式，你可以：
- 🔄 无缝迁移现有 OpenAI 项目
- 📦 使用标准 OpenAI SDK（Python、Node.js 等）
- 🚀 零代码改动，只需修改 `base_url`

## ✨ 特性

### 核心功能
- ⚡ **极致性能** - Go 原生实现，启动 < 10ms，内存占用 ~11MB
- 🔄 **格式转换** - 自动转换 OpenAI ↔ Anthropic 格式
- 🔐 **API Key 代理** - 双 Key 机制保护源头 API Key 🆕
- 🔐 **认证处理** - 支持 Bearer Token 和 API Key 认证
- 🎯 **智能路由** - 自动注入 Factory Droid system prompt
- 📊 **详细日志** - 完整的请求/响应日志记录
- 🏥 **健康检查** - 内置健康检查端点

### 支持的 API
- ✅ Anthropic Claude API（原生格式）
- ✅ OpenAI Chat Completions API（兼容格式）
- ✅ AWS Bedrock API（原生格式）

## 🚀 快速开始

### 安装

**前置要求**: Go 1.21 或更高版本

```bash
# 克隆仓库
git clone https://github.com/yourusername/factory-proxy.git
cd factory-proxy/factory-go-api

# 配置环境变量
cp .env.example .env
# 编辑 .env 文件:
# - FACTORY_API_KEY: 从 https://app.factory.ai/settings/api-keys 获取
# - PROXY_API_KEY: 自定义的安全字符串

# 编译
go build -o factory-proxy main.go              # Anthropic 原生模式
go build -o factory-proxy-openai main-openai.go  # OpenAI 兼容模式
```

### 使用 OpenAI 兼容模式 ⭐ 推荐

#### 1. 启动服务器

```bash
PORT=8003 ./factory-proxy-openai
```

输出：
```
🚀 Factory OpenAI-Compatible Proxy 启动中...
✅ 服务器已启动，监听于 http://localhost:8003
📋 OpenAI兼容接口:
   - POST /v1/chat/completions
   - GET /v1/health
```

#### 2. 使用 Python OpenAI SDK

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
        {"role": "user", "content": "Hello!"}
    ],
    max_tokens=100
)

print(response.choices[0].message.content)
```

#### 3. 使用 curl

```bash
curl -X POST http://localhost:8003/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_PROXY_API_KEY" \
  -d '{
    "model": "claude-sonnet-4-5-20250929",
    "messages": [{"role": "user", "content": "Hello!"}],
    "max_tokens": 100
  }'
```

响应（标准 OpenAI 格式）：
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
      "content": "Hello! How can I help you today?"
    },
    "finish_reason": "stop"
  }],
  "usage": {
    "prompt_tokens": 10,
    "completion_tokens": 8,
    "total_tokens": 18
  }
}
```

### 使用 Anthropic 原生模式

```bash
# 启动服务器
PORT=8001 ./factory-proxy

# 调用 API
curl -X POST http://localhost:8001/anthropic/v1/messages \
  -H "Content-Type: application/json" \
  -H "x-api-key: YOUR_PROXY_API_KEY" \
  -d '{
    "model": "claude-sonnet-4-5-20250929",
    "messages": [{"role": "user", "content": "Hello"}],
    "max_tokens": 100
  }'
```

## 📚 文档

- [🔐 API Key 代理功能](API-KEY-PROXY.md) - 双 Key 机制保护源头 API Key 🆕
- [OpenAI 兼容模式完整文档](README-OpenAI.md) - 详细的 OpenAI 兼容接口说明
- [快速开始指南](QUICK_START.md) - 5分钟快速上手
- [支持的模型列表](MODELS.md) - 25+ 模型完整列表 ⭐
- [项目结构说明](PROJECT_STRUCTURE.md) - 代码架构详解
- [贡献指南](CONTRIBUTING.md) - 如何参与项目开发
- [更新日志](CHANGELOG.md) - 版本更新记录
- [许可证](LICENSE) - MIT License

## 🎯 支持的模型

支持 **25+ 种模型**，包括：

### Claude 系列
- `claude-sonnet-4-5-20250929` ⭐ 推荐
- `claude-opus-4-1-20250805` - 最强推理
- `claude-3-7-sonnet-20250219`
- `claude-3-5-sonnet-20241022`
- `claude-3-5-haiku-20241022`

### GPT 系列
- `gpt-5-2025-08-07` - GPT-5 标准版
- `gpt-5-codex` - 代码专用
- `gpt-4o`, `gpt-4.1`
- `o1`, `o3`, `o4-mini`

### 其他模型
- `gemini-2.5-pro`, `gemini-2.5-flash` - Google Gemini
- `grok-4` - xAI Grok
- `glm-4.6` - 智谱 GLM

**查看完整列表**: [MODELS.md](MODELS.md) 📋

## ⚙️ 配置

### 环境变量

```bash
# 必需配置
export FACTORY_API_KEY="your_real_factory_api_key"  # 源头 Factory API Key (从 https://app.factory.ai/settings/api-keys 获取)
export PROXY_API_KEY="your_custom_proxy_key"        # 对外代理 Key (自定义)

# 可选配置
export PORT=8003  # 服务器端口（默认：8000）
export ANTHROPIC_TARGET_URL="https://your-endpoint.com"  # 已预配置
```




### 使用 .env 文件

复制 `.env.example` 并修改：

```bash
cp .env.example .env
# 编辑 .env 文件
```

## 📊 性能对比

| 指标 | Go 版本 | Deno 版本 |
|------|---------|-----------|
| **启动时间** | ⚡ < 10ms | 🐢 ~500ms |
| **内存占用** | 📉 ~11MB | 📈 ~50MB |
| **二进制大小** | 📦 ~8MB | ❌ 需要运行时 |
| **并发性能** | ⚡ 优秀 | ✅ 良好 |
| **部署复杂度** | ✅ 单文件 | ⚠️ 需要 Deno 环境 |

## 🔧 开发

### 项目结构

```
factory-go/
├── main.go              # Anthropic 原生模式
├── main-openai.go       # OpenAI 兼容模式 ⭐
├── go.mod & go.sum      # Go 依赖
├── README.md            # 主文档
├── README-OpenAI.md     # OpenAI 模式详细文档
├── CONTRIBUTING.md      # 贡献指南
├── LICENSE              # MIT 许可证
├── .gitignore           # Git 忽略文件
├── .env.example         # 环境变量示例
├── start.sh             # 启动脚本
└── test_openai_sdk.py   # Python 测试脚本
```

### 本地开发

```bash
# 克隆项目
git clone https://github.com/yourusername/factory-proxy.git
cd factory-proxy/factory-go

# 安装依赖
go mod tidy

# 运行（开发模式）
go run main-openai.go

# 构建
go build -o factory-proxy-openai main-openai.go

# 测试
go test -v ./...
```

### 代码格式化

```bash
# 格式化代码
gofmt -w .

# 检查代码
go vet ./...
```

## 🚢 部署

### 本地部署

```bash
# 使用启动脚本
./start.sh

# 或手动启动
PORT=8003 ./factory-proxy-openai
```

### Docker 部署

```bash
# 构建镜像
docker build -t factory-proxy .

# 运行容器
docker run -d \
  -p 8003:8003 \
  -e PORT=8003 \
  --name factory-proxy \
  factory-proxy
```

### 生产部署（systemd）

创建服务文件 `/etc/systemd/system/factory-proxy.service`：

```ini
[Unit]
Description=Factory Proxy Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/factory-proxy
Environment="PORT=8003"
ExecStart=/opt/factory-proxy/factory-proxy-openai
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl enable factory-proxy
sudo systemctl start factory-proxy
sudo systemctl status factory-proxy
```

## 🔍 API 端点

### OpenAI 兼容模式

| 端点 | 方法 | 描述 |
|------|------|------|
| `/v1/chat/completions` | POST | OpenAI 兼容的对话接口 |
| `/v1/health` | GET | 健康检查 |
| `/health` | GET | 健康检查（别名） |

### Anthropic 原生模式

| 端点 | 方法 | 描述 |
|------|------|------|
| `/anthropic/*` | POST | Anthropic API 代理 |
| `/openai/*` | POST | OpenAI API 代理 |
| `/bedrock/*` | POST | Bedrock API 代理 |
| `/health` | GET | 健康检查 |

## 📝 示例代码

### Node.js

```javascript
import OpenAI from 'openai';

const client = new OpenAI({
  apiKey: process.env.PROXY_API_KEY,  // 使用代理 Key
  baseURL: 'http://localhost:8003/v1'
});

const response = await client.chat.completions.create({
  model: 'claude-sonnet-4-5-20250929',
  messages: [
    { role: 'system', content: 'You are a helpful assistant.' },
    { role: 'user', content: 'Hello!' }
  ],
  max_tokens: 100
});

console.log(response.choices[0].message.content);
```

### Python

```python
from openai import OpenAI
import os

client = OpenAI(
    api_key=os.getenv("PROXY_API_KEY"),  # 使用代理 Key
    base_url="http://localhost:8003/v1"
)

response = client.chat.completions.create(
    model="claude-sonnet-4-5-20250929",
    messages=[
        {"role": "system", "content": "You are a helpful assistant."},
        {"role": "user", "content": "Hello!"}
    ],
    max_tokens=100
)

print(response.choices[0].message.content)
```

## 🔐 安全建议

1. **使用 API Key 代理** 🆕
   ```bash
   # 配置双 Key 机制
   export FACTORY_API_KEY="your_factory_key"  # 服务器端使用 (从 https://app.factory.ai/settings/api-keys 获取)
   export PROXY_API_KEY="your_proxy_key"      # 客户端使用 (自定义)
   
   # 客户端永远不会接触到源头的 Factory API Key
   ```

2. **保护 API Key**
   ```bash
   # 使用环境变量或 .env 文件
   # 不要在代码中硬编码 API Key
   # 不要提交 .env 文件到 Git（已在 .gitignore 中）
   ```

3. **使用 HTTPS**
   - 生产环境请使用反向代理（Nginx/Caddy）配置 HTTPS

4. **限流保护**
   - 建议在反向代理层面配置限流规则

5. **日志管理**
   - 日志中不包含敏感信息（API Key 已脱敏，只显示前 8 位）

6. **定期轮换 Key**
   - 可以独立轮换 PROXY_API_KEY 而不影响上游连接

## 🆘 故障排除

### 端口被占用

```bash
# 查看端口占用
lsof -i :8003

# 或使用其他端口
PORT=9000 ./factory-proxy-openai
```

### 403 Forbidden 错误

确保：
1. ✅ 使用正确的 Factory API Key
2. ✅ 请求包含正确的认证头
3. ✅ 服务器已正确配置环境变量

### 连接超时

```bash
# 检查目标服务是否可访问
curl -I https://your-target-endpoint.com

# 检查防火墙规则
# 检查网络连接
```

## 🤝 贡献

我们欢迎各种形式的贡献！请查看 [贡献指南](CONTRIBUTING.md)。

### 贡献者

感谢所有贡献者！

<a href="https://github.com/yourusername/factory-proxy/graphs/contributors">
  
  <img src="https://contrib.rocks/image?repo=yourusername/factory-proxy" />
</a>

## 📄 许可证

本项目采用 [MIT License](LICENSE) 开源。

## 🔗 相关链接

- [Factory AI 官网](https://factory.ai)
- [OpenAI API 文档](https://platform.openai.com/docs)
- [Anthropic API 文档](https://docs.anthropic.com)

## ⭐ Star History

如果这个项目对你有帮助，请给我们一个 Star！

[![Star History Chart](https://api.star-history.com/svg?repos=yourusername/factory-proxy&type=Date)](https://star-history.com/#yourusername/factory-proxy&Date)

## 📮 联系方式

- Issues: [GitHub Issues](https://github.com/yourusername/factory-proxy/issues)
- Discussions: [GitHub Discussions](https://github.com/yourusername/factory-proxy/discussions)

---

<div align="center">

**Made with ❤️ by the Factory Proxy Team**

[⬆ 回到顶部](#factory-proxy)

</div>