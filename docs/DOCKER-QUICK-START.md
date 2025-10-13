
# 🚀 Docker 快速部署指南

5 分钟快速在任何支持 Docker 的环境中部署 Factory Go API。

---

## ⚡ 一键部署

```bash
# 1. 克隆项目
git clone https://github.com/libaxuan/factory-go-api.git
cd factory-go-api

# 2. 配置环境变量
cp .env.example .env
nano .env  # 填写 FACTORY_API_KEY 和 PROXY_API_KEY

# 3. 启动服务
docker-compose up -d factory-proxy-openai

# 4. 验证部署
curl http://localhost:8003/health
```

✅ **完成！服务已在 http://localhost:8003 运行**

---

## 📋 前置要求

- Docker 20.10+
- Docker Compose 2.0+
- 1GB+ 可用内存
- 500MB+ 可用磁盘空间

### 快速安装 Docker

**Ubuntu/Debian:**
```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
newgrp docker
```

**其他系统:** 访问 [Docker 官网](https://docs.docker.com/get-docker/)

---

## 🎯 三种部署方式

### 方式 1: Docker Compose（推荐）⭐

```bash
# 启动
docker-compose up -d factory-proxy-openai

# 查看日志
docker-compose logs -f factory-proxy-openai

# 停止
docker-compose down
```

**优点:**
- 一键启动/停止
- 配置管理简单
- 自动重启
- 适合生产环境

### 方式 2: Docker CLI

```bash
# 构建
docker build --target openai -t factory-api .

# 运行
docker run -d \
  --name factory-api \
  -p 8003:8003 \
  --env-file .env \
  --restart unless-stopped \
  factory-api

# 查看日志
docker logs -f factory-api
```

**优点:**
- 更灵活的控制
- 适合自定义配置

### 方式 3: 预构建镜像（即将支持）

```bash
# 拉取镜像
docker pull ghcr.io/libaxuan/factory-go-api:latest

# 运行
docker run -d \
  --name factory-api \
  -p 8003:8003 \
  -e FACTORY_API_KEY=your_key \
  -e PROXY_API_KEY=your_proxy_key \
  ghcr.io/libaxuan/factory-go-api:latest
```

**优点:**
- 无需构建
- 节省时间

---

## ⚙️ 环境变量配置

### 必需配置

```bash
# .env 文件
FACTORY_API_KEY=fk-xxxxxxxxxx    # 从 factory.ai 获取
PROXY_API_KEY=your-secure-key    # 自定义安全密钥
```

### 可选配置

```bash
PORT=8003                         # 服务端口（默认 8003）
CONFIG_PATH=config.json           # 配置文件路径
```

### ⚠️ 修改 .env 后如何生效

**重要提示：修改 .env 文件后必须重启容器才能生效！**

```bash
# 方式1: 重启容器（推荐 - 快速）
docker-compose restart factory-proxy-openai

# 方式2: 停止并重新启动（完全重启）
docker-compose down
docker-compose up -d factory-proxy-openai

# 验证新配置
docker-compose logs -f factory-proxy-openai
```

**为什么需要重启？**
- Docker 容器在启动时读取环境变量
- 修改 .env 文件不会自动传递到运行中的容器
- 必须重启容器才能加载新的环境变量

---

## ✅ 验证部署

### 1. 健康检查

```bash
curl http://localhost:8003/health
```

**预期输出:**
```json
{
  "status": "healthy",
  "timestamp": "2025-10-13T09:00:00Z",
  "uptime": 123.456
}
```

### 2. 查看模型列表

```bash
curl http://localhost:8003/v1/models \
  -H "Authorization: Bearer YOUR_PROXY_API_KEY"
```

**预期输出:** 返回 5 个可用模型

### 3. 测试聊天接口

```bash
curl -X POST http://localhost:8003/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_PROXY_API_KEY" \
  -d '{
    "model": "claude-sonnet-4-5-20250929",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

---

## 📊 常用命令

### 服务管理

```bash
# 启动
docker-compose up -d factory-proxy-openai

# 停止
docker-compose stop factory-proxy-openai

# 重启
docker-compose restart factory-proxy-openai

# 删除
docker-compose down
```

### 日志查看

```bash
# 实时日志
docker-compose logs -f factory-proxy-openai

# 最近 100 行
docker-compose logs --tail=100 factory-proxy-openai
```

### 状态检查

```bash
# 容器状态
docker-compose ps

# 资源使用
docker stats factory-proxy-openai

# 容器详情
docker inspect factory-proxy-openai
```

---

## 🔧 故障排查

### 问题 1: 端口被占用

```bash
# 查看占用进程
sudo lsof -i :8003

# 修改端口（编辑 .env）
PORT=9000

# 重启服务
docker-compose down
docker-compose up -d factory-proxy-openai
```

### 问题 2: 容器无法启动

```bash
# 查看详细日志
docker-compose logs factory-proxy-openai

# 检查环境变量
cat .env

# 重新构建
docker-compose up -d --build factory-proxy-openai
```

### 问题 3: API 返回 401

```bash
# 确认使用正确的 Key
cat .env | grep PROXY_API_KEY

# 注意：使用 PROXY_API_KEY（不是 FACTORY_API_KEY）
curl http://localhost:8003/v1/models \
  -H "Authorization: Bearer YOUR_PROXY_API_KEY"
```

---

## 🎯 下一步

部署成功后，你可以:

### 1. 配置反向代理

使用 Nginx 或 Caddy 配置 HTTPS:

```nginx
server {
    listen 443 ssl;
    server_name api.yourdomain.com;
    
    location / {
        proxy_pass http://localhost:8003;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

### 2. 集成到应用

使用 OpenAI SDK:

```python
from openai import OpenAI

client = OpenAI(
    base_url="http://localhost:8003/v1",
    api_key="YOUR_PROXY_API_KEY"
)

response = client.chat.completions.create(
    model="claude-sonnet-4-5-20250929",
    
messages=[{"role": "user", "content": "Hello!"}]
)
print(response.choices[0].message.content)
```

### 3. 监控和日志

设置日志收集:

```bash
# 导出日志
docker logs factory-proxy-openai > /tmp/api.log

# 持续监控
docker stats factory-proxy-openai
```

### 4. 生产环境部署

查看完整部署指南:
- [Ubuntu Docker 部署](UBUNTU-DOCKER-DEPLOYMENT.md)
- [通用部署指南](DEPLOYMENT.md)

---

## 📚 更多资源

- **完整文档**: [README.md](../README.md)
- **Ubuntu 部署**: [UBUNTU-DOCKER-DEPLOYMENT.md](UBUNTU-DOCKER-DEPLOYMENT.md)
- **API 文档**: http://localhost:8003/docs (服务启动后访问)
- **项目主页**: https://github.com/libaxuan/factory-go-api
- **问题反馈**: https://github.com/libaxuan/factory-go-api/issues

---

## 🎉 部署检查清单

- [ ] Docker 已安装并运行
- [ ] 项目已克隆到本地
- [ ] .env 文件已配置
- [ ] 容器已成功启动
- [ ] 健康检查通过
- [ ] 可以访问模型列表
- [ ] 聊天接口测试成功

**全部完成？恭喜！🎊 你已成功部署 Factory Go API！**

---

<div align="center">

**5 分钟快速部署完成！**

有问题？查看 [故障排查文档](UBUNTU-DOCKER-DEPLOYMENT.md#故障排查)

**Made with ❤️ by Factory Go API Team**

</div>