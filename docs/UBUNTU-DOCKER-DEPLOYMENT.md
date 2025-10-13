
# 🐳 Ubuntu 服务器 Docker 部署指南

本文档提供在 Ubuntu 服务器上使用 Docker 部署 Factory Go API 的详细步骤。

---

## 📋 目录

- [环境要求](#环境要求)
- [第一步：安装 Docker](#第一步安装-docker)
- [第二步：上传项目到服务器](#第二步上传项目到服务器)
- [第三步：配置环境变量](#第三步配置环境变量)
- [第四步：使用 Docker Compose 部署](#第四步使用-docker-compose-部署推荐)
- [第五步：使用 Docker 手动部署](#第五步使用-docker-手动部署可选)
- [第六步：验证部署](#第六步验证部署)
- [第七步：配置防火墙](#第七步配置防火墙)
- [第八步：配置 Nginx 反向代理](#第八步配置-nginx-反向代理可选)
- [日常运维](#日常运维)
- [故障排查](#故障排查)
- [附录](#附录)

---

## 环境要求

- **操作系统**: Ubuntu 18.04+ / Ubuntu 20.04+ / Ubuntu 22.04+
- **Docker**: 20.10+ 或更高
- **Docker Compose**: 2.0+ 或更高
- **内存**: 最低 512MB，推荐 1GB+
- **磁盘**: 最低 500MB 可用空间（包含镜像）
- **网络**: 需要访问外网（拉取镜像和调用 Factory API）

---

## 第一步：安装 Docker

### 1.1 检查是否已安装 Docker

```bash
docker --version
docker-compose --version
```

如果显示版本号，则已安装，可跳过此步骤。

### 1.2 安装 Docker（使用官方脚本 - 推荐）

```bash
# 下载并运行 Docker 官方安装脚本
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 将当前用户添加到 docker 组（避免每次使用 sudo）
sudo usermod -aG docker $USER

# 启动 Docker 服务
sudo systemctl start docker
sudo systemctl enable docker

# 注销并重新登录使用户组生效，或运行
newgrp docker

# 验证安装
docker --version
```

### 1.3 安装 Docker Compose

```bash
# Docker Compose v2 通常随 Docker Desktop 一起安装
# 验证是否已安装
docker-compose version

# 如果未安装，手动安装
sudo apt-get update
sudo apt-get install docker-compose-plugin

# 或使用独立版本
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 验证安装
docker-compose --version
```

**预期输出：**
```
Docker version 24.0.7, build afdd53b
Docker Compose version v2.23.0
```

### 1.4 配置 Docker（可选优化）

```bash
# 配置 Docker 镜像加速（中国大陆用户推荐）
sudo mkdir -p /etc/docker
sudo tee /etc/docker/daemon.json <<-'EOF'
{
  "registry-mirrors": [
    "https://docker.mirrors.ustc.edu.cn",
    "https://hub-mirror.c.163.com"
  ],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}
EOF

# 重启 Docker 服务
sudo systemctl daemon-reload
sudo systemctl restart docker

# 验证配置
docker info | grep -A 5 "Registry Mirrors"
```

---

## 第二步：上传项目到服务器

### 方式 A：使用 Git Clone（推荐）

```bash
# 进入部署目录
cd /opt

# 克隆项目
sudo git clone https://github.com/libaxuan/factory-go-api.git

# 进入项目目录
cd factory-go-api

# 查看文件
ls -la
```

### 方式 B：使用 SCP 上传

```bash
# 在本地电脑执行
scp -r /path/to/local/factory-go-api root@your-server-ip:/opt/

# 然后在服务器上
cd /opt/factory-go-api
```

### 方式 C：使用 SFTP 工具

使用 FileZilla、WinSCP 等工具上传整个项目到 `/opt/factory-go-api`

---

## 第三步：配置环境变量

### 3.1 复制环境变量模板

```bash
cd /opt/factory-go-api
cp .env.example .env
```

### 3.2 编辑环境变量

```bash
nano .env
```

### 3.3 填写配置信息

```bash
# Factory API Key - 从 https://app.factory.ai/settings/api-keys 获取
FACTORY_API_KEY=fk-xxxxxxxxxxxxxxxxxx

# 对外代理 API Key - 自定义一个安全的密钥
PROXY_API_KEY=your_custom_secure_key_here

# 服务端口（默认 8003）
PORT=8003

# 配置文件路径（默认 config.json）
CONFIG_PATH=config.json
```

**重要提示：**
- ✅ 不要在值两边加引号
- ✅ 等号前后不要有空格
- ✅ 每行结尾不要有多余空格
- ✅ `PROXY_API_KEY` 用于客户端调用本代理服务
- ✅ `FACTORY_API_KEY` 用于本服务调用 Factory 


### 3.4 保存文件

- 按 `Ctrl + X`
- 按 `Y` 确认
- 按 `Enter` 保存

### 3.5 ⚠️ 重要：修改 .env 后如何生效

**.env 文件修改后，必须重启容器才能生效！**

```bash
# 如果已经启动了容器，修改 .env 后需要重启
docker-compose restart factory-proxy-openai

# 或者停止后重新启动
docker-compose down
docker-compose up -d factory-proxy-openai

# 验证新配置是否生效
docker-compose logs -f factory-proxy-openai
```

**为什么需要重启？**
- Docker 容器在启动时读取环境变量
- 修改 .env 文件不会自动传递到运行中的容器
- 必须重启容器才能加载新的环境变量

---

## 第四步：使用 Docker Compose 部署（推荐）

### 4.1 查看 Docker Compose 配置

```bash
cd /opt/factory-go-api
cat docker-compose.yml
```

### 4.2 启动服务（OpenAI 兼容模式）

```bash
# 启动 OpenAI 兼容模式（推荐）
docker-compose up -d factory-proxy-openai

# 查看启动日志
docker-compose logs -f factory-proxy-openai
```

**预期输出：**
```
[+] Running 2/2
 ✔ Network factory-go-api_factory-network  Created
 ✔ Container factory-proxy-openai          Started
```

### 4.3 查看服务状态

```bash
# 查看容器状态
docker-compose ps

# 查看详细日志
docker-compose logs factory-proxy-openai

# 实时查看日志
docker-compose logs -f factory-proxy-openai
```

**成功标志：** 日志中显示
```
🔐 代理模式: 已启用
📖 加载配置文件: config.json
✅ 配置加载成功
📍 支持的模型 (5):
   • claude-opus-4-1-20250805 [anthropic]
   • claude-sonnet-4-20250514 [anthropic]
   • claude-sonnet-4-5-20250929 [anthropic]
   • gpt-5-2025-08-07 [openai]
   • gpt-5-codex [openai]
🚀 服务启动于 http://localhost:8003
📖 文档: http://localhost:8003/docs
```

### 4.4 Docker Compose 常用命令

```bash
# 启动服务
docker-compose up -d factory-proxy-openai

# 停止服务
docker-compose stop factory-proxy-openai

# 重启服务
docker-compose restart factory-proxy-openai

# 停止并删除容器
docker-compose down

# 查看日志
docker-compose logs -f factory-proxy-openai

# 重新构建并启动
docker-compose up -d --build factory-proxy-openai

# 查看服务状态
docker-compose ps

# 进入容器（调试用）
docker-compose exec factory-proxy-openai sh
```

---

## 第五步：使用 Docker 手动部署（可选）

如果不使用 Docker Compose，也可以手动使用 Docker 命令部署。

### 5.1 构建 Docker 镜像

```bash
cd /opt/factory-go-api

# 构建 OpenAI 兼容模式镜像
docker build --target openai -t factory-proxy-openai:latest .
```

**构建过程说明：**
- 使用多阶段构建，最终镜像约 20MB
- 基于 Alpine Linux，安全且轻量
- 自动包含所有依赖

### 5.2 运行容器

```bash
# 从 .env 文件加载环境变量并运行
docker run -d \
  --name factory-proxy-openai \
  --restart unless-stopped \
  -p 8003:8003 \
  --env-file .env \
  factory-proxy-openai:latest

# 或手动指定环境变量
docker run -d \
  --name factory-proxy-openai \
  --restart unless-stopped \
  -p 8003:8003 \
  -e PORT=8003 \
  -e FACTORY_API_KEY=fk-your-key \
  -e PROXY_API_KEY=your-proxy-key \
  factory-proxy-openai:latest
```

### 5.3 管理容器

```bash
# 查看运行状态
docker ps

# 查看所有容器（包括停止的）
docker ps -a

# 查看日志
docker logs -f factory-proxy-openai

# 停止容器
docker stop factory-proxy-openai

# 启动容器
docker start factory-proxy-openai

# 重启容器
docker restart factory-proxy-openai

# 删除容器
docker rm -f factory-proxy-openai

# 进入容器（调试用）
docker exec -it factory-proxy-openai sh
```

---

## 第六步：验证部署

### 6.1 健康检查

```bash
# 从服务器本地测试
curl http://localhost:8003/health

# 预期输出
# {"status":"healthy","timestamp":"2025-10-13T09:00:00Z","uptime":123.456}
```

### 6.2 查看服务信息

```bash
curl http://localhost:8003/

# 预期输出
# {
#   "service": "Factory Go API - Multi-Model Support",
#   "version": "2.0",
#   "endpoints": ["/health", "/v1/models", "/v1/chat/completions"]
# }
```

### 6.3 查看模型列表

```bash
# 替换 YOUR_PROXY_API_KEY 为 .env 中配置的 PROXY_API_KEY
curl http://localhost:8003/v1/models \
  -H "Authorization: Bearer YOUR_PROXY_API_KEY"

# 预期输出：返回 5 个模型的列表
```

### 6.4 测试聊天接口（非流式）

```bash
curl -X POST http://localhost:8003/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_PROXY_API_KEY" \
  -d '{
    "model": "claude-sonnet-4-5-20250929",
    "messages": [{"role": "user", "content": "你好"}],
    "stream": false
  }'
```

---

## 第七步：配置防火墙

```bash
# 检查防火墙状态
sudo ufw status

# 允许 8003 端口
sudo ufw allow 8003/tcp

# 如果使用 Nginx 反向代理，还需开放 80 和 443
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# 启用防火墙（如果未启用）
sudo ufw enable

# 验证规则
sudo ufw status numbered
```

**云服务器额外步骤：**

登录云服务商控制台（阿里云/腾讯云/AWS等），在安全组中添加入站规则：
- TCP 8003（API 服务）
- TCP 80（HTTP，如使用 Nginx）
- TCP 443（HTTPS，如使用 Nginx）

---

## 第八步：配置 Nginx 反向代理（可选）

### 8.1 安装 Nginx

```bash
sudo apt update
sudo apt install nginx -y
sudo systemctl start nginx
sudo systemctl enable nginx
```

### 8.2 创建配置文件

```bash
sudo nano /etc/nginx/sites-available/factory-proxy-openai
```

### 8.3 填写配置

```nginx
upstream factory_proxy {
    server 127.0.0.1:8003;
    keepalive 32;
}

server {
    listen 80;
    server_name api.yourdomain.com;

    access_log /var/log/nginx/factory-proxy-openai-access.log;
    error_log /var/log/nginx/factory-proxy-openai-error.log;

    # 客户端请求体大小限制
    client_max_body_size 10M;

    location / {
        proxy_pass http://factory_proxy;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # 超时设置（SSE 流式响应需要较长超时）
        proxy_connect_timeout 60s;
        proxy_send_timeout 300s;
        proxy_read_timeout 300s;
        
        # 禁用缓冲（支持流式响应）
        proxy_buffering off;
        proxy_cache off;
    }

    # 健康检查端点（不记录访问日志）
    location /health {
        proxy_pass http://factory_proxy;
        access_log off;
    }
}
```

### 8.4 启用配置

```bash
# 创建软链接
sudo ln -s /etc/nginx/sites-available/factory-proxy-openai /etc/nginx/sites-enabled/

# 测试配置
sudo nginx -t

# 重载配置
sudo systemctl reload nginx
```

### 8.5 配置 HTTPS（推荐）

```bash
# 安装 Certbot
sudo apt install certbot python3-certbot-nginx -y

# 自动配置 SSL（会自动修改 Nginx 配置）
sudo certbot --nginx -d api.yourdomain.com

# 测试自动续期
sudo certbot renew 
--dry-run

# 查看证书状态
sudo certbot certificates
```

---

## 日常运维

### Docker Compose 服务管理

```bash
# 查看服务状态
docker-compose ps

# 启动服务
docker-compose start factory-proxy-openai

# 停止服务
docker-compose stop factory-proxy-openai

# 重启服务
docker-compose restart factory-proxy-openai

# 查看实时日志
docker-compose logs -f factory-proxy-openai

# 查看最近 100 行日志
docker-compose logs --tail=100 factory-proxy-openai

# 停止并删除容器
docker-compose down

# 停止并删除容器及网络、卷
docker-compose down -v
```

### Docker 容器管理（手动部署）

```bash
# 查看容器状态
docker ps
docker ps -a  # 包含停止的容器

# 启动/停止/重启
docker start factory-proxy-openai
docker stop factory-proxy-openai
docker restart factory-proxy-openai

# 查看日志
docker logs -f factory-proxy-openai
docker logs --tail=100 factory-proxy-openai

# 查看容器资源使用
docker stats factory-proxy-openai

# 进入容器
docker exec -it factory-proxy-openai sh

# 删除容器
docker stop factory-proxy-openai
docker rm factory-proxy-openai
```

### 更新部署

```bash
# 方式1: 使用 Docker Compose（推荐）
cd /opt/factory-go-api

# 拉取最新代码
git pull

# 重新构建并启动
docker-compose up -d --build factory-proxy-openai

# 查看日志确认启动成功
docker-compose logs -f factory-proxy-openai
```

```bash
# 方式2: 手动构建
cd /opt/factory-go-api

# 拉取最新代码
git pull

# 停止旧容器
docker stop factory-proxy-openai
docker rm factory-proxy-openai

# 重新构建镜像
docker build --target openai -t factory-proxy-openai:latest .

# 启动新容器
docker run -d \
  --name factory-proxy-openai \
  --restart unless-stopped \
  -p 8003:8003 \
  --env-file .env \
  factory-proxy-openai:latest

# 查看日志
docker logs -f factory-proxy-openai
```

### 修改环境变量

```bash
# 1. 编辑 .env 文件
cd /opt/factory-go-api
nano .env

# 2. 修改后保存

# 3. 重启容器使配置生效
docker-compose restart factory-proxy-openai

# 或手动部署方式
docker stop factory-proxy-openai
docker rm factory-proxy-openai
docker run -d \
  --name factory-proxy-openai \
  --restart unless-stopped \
  -p 8003:8003 \
  --env-file .env \
  factory-proxy-openai:latest

# 4. 查看日志确认新配置生效
docker-compose logs -f factory-proxy-openai
```

### 清理 Docker 资源

```bash
# 清理停止的容器
docker container prune

# 清理未使用的镜像
docker image prune -a

# 清理未使用的卷
docker volume prune

# 清理所有未使用资源（慎用）
docker system prune -a

# 查看 Docker 磁盘使用
docker system df
```

### 监控和日志

```bash
# 实时查看容器资源使用
docker stats

# 查看特定容器资源
docker stats factory-proxy-openai

# 导出日志到文件
docker logs factory-proxy-openai > /tmp/factory-proxy-openai.log

# 查看容器详细信息
docker inspect factory-proxy-openai

# 查看容器进程
docker top factory-proxy-openai
```

---

## 故障排查

### 问题 1：容器无法启动

**症状：** `docker-compose ps` 显示容器状态为 `Exit` 或 `Restarting`

**排查步骤：**

```bash
# 1. 查看详细日志
docker-compose logs factory-proxy-openai

# 2. 常见错误及解决方案：

# 错误：端口被占用
# 解决：查找并停止占用端口的进程
sudo lsof -i :8003
sudo kill -9 <PID>
# 或修改 .env 中的 PORT

# 错误：环境变量未配置
# 解决：检查 .env 文件
cat .env
# 确保 FACTORY_API_KEY 和 PROXY_API_KEY 已设置

# 错误：权限问题
# 解决：检查文件权限
ls -la /opt/factory-go-api
sudo chown -R $USER:$USER /opt/factory-go-api
```

### 问题 2：镜像构建失败

**症状：** `docker build` 或 `docker-compose up --build` 失败

**解决方案：**

```bash
# 1. 检查网络连接
ping -c 3 registry-1.docker.io

# 2. 使用国内镜像加速
sudo nano /etc/docker/daemon.json
# 添加镜像源（见第一步的配置）

# 3. 清理构建缓存重试
docker builder prune -a
docker-compose build --no-cache factory-proxy-openai

# 4. 查看详细构建日志
docker-compose build --progress=plain factory-proxy-openai
```

### 问题 3：API 返回 401 未授权

**症状：** `curl` 请求返回 `{"error": {"message": "Invalid API key", ...}}`

**解决方案：**

```bash
# 1. 确认使用的是 PROXY_API_KEY（不是 FACTORY_API_KEY）
cat .env | grep PROXY_API_KEY

# 2. 检查容器环境变量
docker exec factory-proxy-openai env | grep API_KEY

# 3. 测试正确的 Key
curl http://localhost:8003/v1/models \
  -H "Authorization: Bearer $(grep PROXY_API_KEY .env | cut -d'=' -f2)"

# 4. 查看日志中的验证信息
docker-compose logs factory-proxy-openai | grep "验证"
```

### 问题 4：容器健康检查失败

**症状：** `docker ps` 显示容器状态为 `unhealthy`

**解决方案：**

```bash
# 1. 查看健康检查日志
docker inspect factory-proxy-openai | grep -A 10 Health

# 2. 手动测试健康检查端点
curl http://localhost:8003/health

# 3. 进入容器测试
docker exec -it factory-proxy-openai sh
wget --spider http://localhost:8003/health

# 4. 检查端口监听
docker exec factory-proxy-openai netstat -tlnp

# 5. 
如果健康检查正常但容器显示 unhealthy，可能是时间问题，等待片刻
```

### 问题 5：无法从外网访问

**排查步骤：**

```bash
# 1. 检查容器是否运行
docker ps | grep factory-proxy-openai

# 2. 检查端口映射
docker port factory-proxy-openai

# 3. 检查本地访问
curl http://localhost:8003/health

# 4. 检查防火墙
sudo ufw status
sudo ufw allow 8003/tcp

# 5. 检查云服务商安全组
# 登录控制台查看是否开放 8003 端口

# 6. 测试从外网访问（替换为实际 IP）
curl http://YOUR_SERVER_IP:8003/health
```

### 问题 6：Docker Compose version 警告

**症状：** 显示 `the attribute 'version' is obsolete`

**说明：** 这是警告信息，不影响功能。Docker Compose v2 不再需要 version 字段。

**解决方案（可选）：**

```bash
# 删除 docker-compose.yml 第一行的 version: '3.8'
sed -i '1d' docker-compose.yml
```

### 快速诊断脚本

```bash
# 创建诊断脚本
cat > /tmp/docker-diagnose.sh << 'EOF'
#!/bin/bash
echo "=== Docker 部署诊断信息 ==="
echo ""
echo "1. Docker 版本:"
docker --version
docker-compose --version
echo ""
echo "2. 容器状态:"
docker-compose ps
echo ""
echo "3. 容器日志（最近20行）:"
docker-compose logs --tail=20 factory-proxy-openai
echo ""
echo "4. 端口监听:"
sudo netstat -tlnp | grep 8003
echo ""
echo "5. 环境变量:"
ls -la /opt/factory-go-api/.env
echo ""
echo "6. 防火墙状态:"
sudo ufw status
echo ""
echo "7. 磁盘使用:"
docker system df
echo ""
echo "8. 健康检查:"
curl -s http://localhost:8003/health || echo "健康检查失败"
echo ""
EOF

chmod +x /tmp/docker-diagnose.sh
/tmp/docker-diagnose.sh
```

---

## 附录

### A. Docker Compose 完整配置示例

```yaml
version: '3.8'

services:
  factory-proxy-openai:
    build:
      context: .
      target: openai
    container_name: factory-proxy-openai
    ports:
      - "8003:8003"
    environment:
      - PORT=8003
      - FACTORY_API_KEY=${FACTORY_API_KEY}
      - PROXY_API_KEY=${PROXY_API_KEY}
      - CONFIG_PATH=config.json
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8003/health"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 5s
    networks:
      - factory-network

networks:
  factory-network:
    driver: bridge
```

### B. .env 文件完整示例

```bash
# Factory API Key（必需）
FACTORY_API_KEY=fk-your-factory-api-key-here

# 对外代理 API Key（必需）
PROXY_API_KEY=your-custom-proxy-key-here

# 服务端口（可选，默认 8003）
PORT=8003

# 配置文件路径（可选，默认 config.json）
CONFIG_PATH=config.json
```

### C. API 端点速查

| 端点 | 方法 | 认证 | 说明 |
|------|------|------|------|
| `/` | GET | ❌ | 服务信息 |
| `/health` | GET | ❌ | 健康检查 |
| `/v1/models` | GET | ✅ | 模型列表 |
| `/v1/chat/completions` | POST | ✅ | 聊天接口 |
| `/docs` | GET | ❌ | API 文档 |

### D. 测试命令速查

```bash
# 健康检查
curl http://localhost:8003/health

# 查看模型
curl http://localhost:8003/v1/models \
  -H "Authorization: Bearer YOUR_PROXY_API_KEY"

# 非流式对话
curl -X POST http://localhost:8003/v1/chat/completions \
  -H "Authorization: Bearer YOUR_PROXY_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "claude-sonnet-4-5-20250929",
    "messages": [{"role": "user", "content": "Hello"}],
    "stream": false
  }'

# 流式对话
curl -X POST http://localhost:8003/v1/chat/completions \
  -H "Authorization: Bearer YOUR_PROXY_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-5-2025-08-07",
    "messages": [{"role": "user", "content": "Tell me a story"}],
    "stream": true
  }'
```

### E. Docker 命令速查

```bash
# 服务管理
docker-compose up -d factory-proxy-openai      # 启动
docker-compose stop factory-proxy-openai       # 停止
docker-compose restart factory-proxy-openai    # 重启
docker-compose down                            # 停止并删除

# 日志查看
docker-compose logs -f factory-proxy-openai    # 实时日志
docker-compose logs --tail=50 factory-proxy-openai  # 最近50行

# 状态检查
docker-compose ps                              # 服务状态
docker stats factory-proxy-openai              # 资源使用

# 容器管理
docker exec -it factory-proxy-openai sh        # 进入容器
docker inspect factory-proxy-openai            # 详细信息
```

### F. 支持的模型列表

| 模型 ID | 类型 | 说明 |
|---------|------|------|
| `claude-opus-4-1-20250805` | Anthropic | Claude Opus 4.1 - 最强推理 |
| `claude-sonnet-4-20250514` | Anthropic | Claude Sonnet 4 - 中等性能 |
| `claude-sonnet-4-5-20250929` | Anthropic | Claude Sonnet 4.5（推荐） |
| `gpt-5-2025-08-07` | OpenAI | GPT-5 - 最新旗舰 |
| `gpt-5-codex` | OpenAI | 
GPT-5 Codex - 代码生成专家 |

### G. 性能优化建议

```bash
# 1. 限制容器资源使用
docker run -d \
  --name factory-proxy-openai \
  --restart unless-stopped \
  -p 8003:8003 \
  --memory="256m" \
  --cpus="0.5" \
  --env-file .env \
  factory-proxy-openai:latest

# 2. 使用 Docker Compose 限制资源
# 在 docker-compose.yml 中添加：
#   deploy:
#     resources:
#       limits:
#         cpus: '0.5'
#         memory: 256M

# 3. 定期清理 Docker 资源
docker system prune -a --volumes -f

# 4. 配置日志轮转
# 在 /etc/docker/daemon.json 中添加：
# {
#   "log-driver": "json-file",
#   "log-opts": {
#     "max-size": "10m",
#     "max-file": "3"
#   }
# }
```

---

## 🎉 部署成功检查清单

完成以下所有项目即表示部署成功：

- [ ] Docker 环境已安装（20.10+）
- [ ] Docker Compose 已安装（2.0+）
- [ ] 项目代码已上传到服务器
- [ ] .env 文件已正确配置
- [ ] 容器已成功构建并启动
- [ ] `docker-compose ps` 显示容器运行中
- [ ] 健康检查返回 `{"status":"healthy"}`
- [ ] 模型列表接口正常返回
- [ ] 防火墙已开放 8003 端口
- [ ] 可以从外网访问（如需要）
- [ ] Nginx 反向代理已配置（如需要）
- [ ] HTTPS 证书已配置（如需要）

---

## 📊 部署方式对比

| 特性 | Docker 部署 | 直接部署 |
|------|------------|---------|
| **安装难度** | ⭐⭐ | ⭐⭐⭐ |
| **启动速度** | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **资源占用** | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **隔离性** | ⭐⭐⭐⭐⭐ | ⭐⭐ |
| **可移植性** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |
| **维护难度** | ⭐⭐ | ⭐⭐⭐ |
| **更新便捷** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| **推荐度** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |

**推荐使用 Docker 部署的场景：**
- ✅ 需要快速部署和迁移
- ✅ 多环境部署（开发/测试/生产）
- ✅ 需要版本管理和回滚
- ✅ 团队协作开发
- ✅ 微服务架构

**推荐使用直接部署的场景：**
- ✅ 追求极致性能
- ✅ 资源受限环境
- ✅ 简单的单机部署
- ✅ 不需要容器化

---

## 🔗 相关文档

- [直接部署指南](UBUNTU-DEPLOYMENT.md) - 不使用 Docker 的部署方式
- [通用部署指南](DEPLOYMENT.md) - 多种部署方式对比
- [项目主文档](../README.md) - 项目完整文档
- [快速开始](QUICK_START.md) - 5分钟快速上手
- [Docker 官方文档](https://docs.docker.com/)
- [Docker Compose 文档](https://docs.docker.com/compose/)

---

## 📞 获取帮助

- **项目主页**: https://github.com/libaxuan/factory-go-api
- **问题反馈**: https://github.com/libaxuan/factory-go-api/issues
- **Factory AI 官网**: https://factory.ai
- **Docker 社区**: https://forums.docker.com/

---

## 🎯 下一步

部署成功后，你可以：

1. **集成到应用** - 使用 OpenAI SDK 调用服务
2. **配置监控** - 使用 Prometheus/Grafana 监控服务
3. **性能测试** - 使用压测工具验证性能
4. **扩展部署** - 使用 Docker Swarm 或 Kubernetes 集群部署
5. **自动化部署** - 配置 CI/CD 流程

---

<div align="center">

**恭喜！🎊 你已成功在 Ubuntu 服务器上使用 Docker 部署 Factory Go API！**

如有问题，请参考「故障排查」章节或提交 Issue。

---

**Made with ❤️ by Factory Go API Team**

</div>