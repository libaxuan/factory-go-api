
# Factory Proxy - 快速开始指南

## 🚀 5分钟快速上手

### 步骤 1: 获取代码

```bash
git clone https://github.com/libaxuan/factory-proxy.git
cd factory-proxy/factory-go-api
```

### 步骤 2: 构建项目

```bash
# 使用 Makefile（推荐）
make build-openai

# 或者直接使用 Go
go build -o factory-proxy-openai main-openai.go
```

### 步骤 3: 启动服务

```bash
# 默认端口 8000
./factory-proxy-openai

# 或指定端口
PORT=8003 ./factory-proxy-openai
```

### 步骤 4: 测试服务

**健康检查：**
```bash
curl http://localhost:8003/v1/health
```

**发送请求：**
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

## 🐍 Python 示例

### 安装 OpenAI SDK

```bash
pip install openai
```

### 使用代码

```python
from openai import OpenAI

# 初始化客户端
client = OpenAI(
    api_key="YOUR_PROXY_API_KEY",  # 使用代理 Key
    base_url="http://localhost:8003/v1"
)

# 发送请求
response = client.chat.completions.create(
    model="claude-sonnet-4-5-20250929",
    messages=[
        {"role": "system", "content": "You are a helpful assistant."},
        {"role": "user", "content": "What is 2+2?"}
    ],
    max_tokens=100
)

# 打印响应
print(response.choices[0].message.content)
```

## 📦 Node.js 示例

### 安装依赖

```bash
npm install openai
```

### 使用代码

```javascript
import OpenAI from 'openai';

const client = new OpenAI({
  apiKey: process.env.PROXY_API_KEY,  // 使用代理 Key
  baseURL: 'http://localhost:8003/v1'
});

async function main() {
  const response = await client.chat.completions.create({
    model: 'claude-sonnet-4-5-20250929',
    messages: [
      { role: 'system', content: 'You are a helpful assistant.' },
      { role: 'user', content: 'What is 2+2?' }
    ],
    max_tokens: 100
  });

  console.log(response.choices[0].message.content);
}

main();
```

## 🐳 Docker 快速启动

### 使用 Docker Compose（推荐）

```bash
# 启动 OpenAI 兼容模式
docker-compose up -d factory-proxy-openai

# 查看日志
docker-compose logs -f factory-proxy-openai

# 停止服务
docker-compose down
```

### 使用 Docker

```bash
# 构建镜像
docker build -t factory-proxy --target openai .

# 运行容器
docker run -d \
  -p 8003:8003 \
  -e PORT=8003 \
  --name factory-proxy \
  factory-proxy

# 查看日志
docker logs -f factory-proxy
```

## 🛠️ 使用 Makefile

Factory Proxy 提供了方便的 Makefile 命令：

```bash
# 查看所有命令
make help

# 构建 OpenAI 兼容版本
make build-openai

# 运行（自动构建）
make run-openai

# 开发模式（无需构建）
make dev-openai

# 格式化代码
make fmt

# 运行测试
make test

# 清理构建文件
make clean
```

## 🔧 常用配置

### 环境变量

创建 `.env` 文件：

```bash
# 服务器端口
PORT=8003

# 目标 API URL（可选，已有默认值）
ANTHROPIC_TARGET_URL=https://your-endpoint.com
```

### 使用 .env 文件

```bash
# 加载环境变量
export $(cat .env | xargs)

# 启动服务
./factory-proxy-openai
```

## 📊 两种模式对比

### OpenAI 兼容模式 ⭐ 推荐

**优点：**
- ✅ 兼容所有 OpenAI SDK
- ✅ 无需修改现有代码
- ✅ 标准 API 格式

**端点：**
- `POST /v1/chat/completions`
- `GET /v1/health`

**认证：**
```bash
Authorization: Bearer YOUR_PROXY_API_KEY
```

### Anthropic 原生模式

**优点：**
- ✅ 直接使用原生格式
- ✅ 无格式转换开销

**端点：**
- `POST /anthropic/v1/messages`
- `GET /health`

**认证：**
```bash
x-api-key: YOUR_PROXY_API_KEY
```

## 🆘 常见问题

### 1. 端口被占用

```bash
# 查看端口占用
lsof -i :8003

# 使用其他端口
PORT=9000 ./factory-proxy-openai
```

### 2. 连接被拒绝

检查：
- ✅ 服务是否正在运行
- ✅ 端口号是否正确
- ✅ 防火墙设置

### 3. 认证失败

确保：
- ✅ API Key 正确
- ✅ 使用正确的认证头格式
- ✅ API Key 有效且未过期

### 4. 构建失败

```bash
# 清理并重新构建
make clean
go mod tidy
make build-openai
```

## 📚 更多资源

- [完整文档](README.md) - 详细的使用说明
- [OpenAI 兼容文档](README-OpenAI.md) - OpenAI 模式详细说明
- [项目结构](PROJECT_STRUCTURE.md) - 代码结构说明
- [贡献指南](CONTRIBUTING.md) - 如何贡献代码
- [更新日志](CHANGELOG.md) - 版本更新记录

## 💡 提示

1. **推荐使用 OpenAI 兼容模式** - 更容易集成到现有项目
2. **使用 Makefile** - 简化构建和运行流程
3. **查看日志** - 遇到问题时先查看服务器日志
4. **健康检查** - 部署后先测试健康检查端点

## 🎯 下一步

现在你已经成功运行了 Factory Proxy！

建议接下来：
1. 阅读 [README.md](README.md) 了解更多功能

2. 尝试使用 Python/Node.js SDK
3. 部署到生产环境
4. 查看 [GitHub Issues](https://github.com/libaxuan/factory-proxy/issues) 参与讨论

---

**开始你的 Factory Proxy 之旅！** 🚀