
# 🚀 启动脚本使用指南

Factory Proxy API 提供了便捷的启动脚本 `start.sh`，支持两种运行模式。

## 📖 使用方法

### OpenAI 兼容模式 ⭐ 推荐

```bash
# 默认启动 OpenAI 兼容模式
./start.sh

# 或明确指定
./start.sh openai
```

**输出示例**:
```
🚀 Factory Proxy API - 快速启动
==================================
✅ Go 版本: go version go1.25.0 darwin/arm64
📦 安装依赖...
🔨 构建 OpenAI 兼容模式... ⭐
✅ 构建成功！

📍 启动信息:
   模式: OpenAI 兼容模式
   端口: 8003
   服务: http://localhost:8003

📋 API 端点 (OpenAI 兼容模式) ⭐:
   - Chat:      http://localhost:8003/v1/chat/completions
   - 健康检查:  http://localhost:8003/v1/health

🔑 认证方式:
   Authorization: Bearer YOUR_PROXY_API_KEY

💡 快速测试:
   curl http://localhost:8003/v1/health

⏳ 启动服务器...
==================================
```

### Anthropic 原生模式

```bash
./start.sh anthropic
```

**输出示例**:
```
🚀 Factory Proxy API - 快速启动
==================================
✅ Go 版本: go version go1.25.0 darwin/arm64
📦 安装依赖...
🔨 构建 Anthropic 原生模式...
✅ 构建成功！

📍 启动信息:
   模式: Anthropic 原生模式
   端口: 8000
   服务: http://localhost:8000

📋 API 端点 (Anthropic 原生模式):
   - Anthropic: http://localhost:8000/anthropic/v1/messages
   - OpenAI:    http://localhost:8000/openai/v1/chat/completions
   - Bedrock:   http://localhost:8000/bedrock/v1/messages
   - 健康检查:  http://localhost:8000/health

🔑 认证方式:
   x-api-key: YOUR_FACTORY_API_KEY

⏳ 启动服务器...
==================================
```

## ⚙️ 环境变量

### 自定义端口

```bash
# OpenAI 模式使用自定义端口
PORT=9000 ./start.sh

# Anthropic 模式使用自定义端口
PORT=9001 ./start.sh anthropic
```

### 使用 .env 文件

```bash
# 创建 .env 文件
cp .env.example .env

# 编辑配置
vim .env

# 加载并启动
source .env && ./start.sh
```

## 📝 快速测试

### OpenAI 兼容模式测试

```bash
# 1. 启动服务
./start.sh

# 2. 在另一个终端测试健康检查
curl http://localhost:8003/v1/health

# 3. 测试 Chat Completions
curl -X POST http://localhost:8003/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_PROXY_API_KEY" \
  -d '{
    "model": "claude-sonnet-4-5-20250929",
    "messages": [{"role": "user", "content": "Hello!"}],
    "max_tokens": 100
  }'
```

### Anthropic 原生模式测试

```bash
# 1. 启动服务
./start.sh anthropic

# 2. 测试健康检查
curl http://localhost:8000/health

# 3. 测试 Anthropic API
curl -X POST http://localhost:8000/anthropic/v1/messages \
  -H "Content-Type: application/json" \
  -H "x-api-key: YOUR_FACTORY_API_KEY" \
  -d '{
    "model": "claude-sonnet-4-5-20250929",
    "messages": [{"role": "user", "content": "Hello!"}],
    "max_tokens": 100
  }'
```

## 🔍 常见问题

### Q: 如何停止服务？
A: 按 `Ctrl+C` 即可优雅停止服务器。

### Q: 端口被占用怎么办？
```bash
# 查看端口占用
lsof -i :8003

# 使用其他端口
PORT=9000 ./start.sh
```

### Q: 权限被拒绝？
```bash
# 添加执行权限
chmod +x start.sh
```

### Q: Go 未安装？
访问 https://golang.org/dl/ 下载安装 Go 1.21+

### Q: 如何切换模式？
```bash
# OpenAI 模式 (默认，推荐)
./start.sh

# Anthropic 原生模式
./start.sh anthropic
```

## 📊 模式对比

| 特性 | OpenAI 模式 ⭐ | Anthropic 模式 |
|------|---------------|----------------|
| **默认端口** | 8003 | 8000 |
| **API 格式** | OpenAI 标准 | Anthropic 原生 |
| **认证方式** | Bearer Token | x-api-key |
| **SDK 支持** | ✅ OpenAI SDK | ❌ 需自定义 |
| **易用性** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |
| **推荐场景** | 通用应用 | Factory 原生 |

## 🎯 推荐使用

**大多数情况下，我们推荐使用 OpenAI 兼容模式**:

```bash
./start.sh
```

优势：
- ✅ 标准 OpenAI API 格式
- ✅ 支持所有 OpenAI SDK
- ✅ 零学习成本
- ✅ 社区支持完善

## 🔗 相关文档

- 
[快速开始](QUICK_START.md) - 5分钟快速上手
- [支持的模型](MODELS.md) - 25+ 模型列表
- [完整文档](README.md) - 项目主文档
- [OpenAI 模式详解](README-OpenAI.md) - OpenAI 兼容接口

---

**推荐**: 使用 `./start.sh` 快速启动 OpenAI 兼容模式！ 🚀