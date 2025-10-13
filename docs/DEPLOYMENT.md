
# 🚀 部署指南

Factory Proxy API 提供多种部署方式，适用于不同的使用场景。

## 📋 目录

- [快速开始](#快速开始)
- [方式1：使用启动脚本](#方式1使用启动脚本)
- [方式2：使用 Makefile](#方式2使用-makefile)
- [方式3：使用 Docker](#方式3使用-docker)
- [方式4：使用 Docker Compose](#方式4使用-docker-compose)
- [方式5：直接运行二进制文件](#方式5直接运行二进制文件)
- [生产环境部署](#生产环境部署)

---

## 快速开始

**推荐使用 OpenAI 兼容模式** ⭐

```bash
# 方式1: 使用启动脚本 (最简单)
./start.sh

# 方式2: 使用 Makefile (推荐)
make start

# 方式3: 使用 Docker Compose (容器化)
docker-compose up -d factory-proxy-openai
```

---

## 方式1：使用启动脚本

### OpenAI 兼容模式 ⭐ 推荐

```bash
# 默认启动
./start.sh

# 自定义端口
PORT=9000 ./start.sh
```

### Anthropic 原生模式

```bash
./start.sh anthropic
```

### 特点
- ✅ 最简单的启动方式
- ✅ 自动安装依赖
- ✅ 自动构建
- ✅ 友好的输出信息

详细文档: [START.md](START.md)

---

## 方式2：使用 Makefile

### 快速命令

```bash
# 🌟 推荐：一键启动 (OpenAI 模式)
make start

# 或者
make run-openai

# Anthropic 原生模式
make run
```

### 所有可用命令

```bash
# 查看帮助
make help

# 构建
make build-openai    # 构建 OpenAI 版本
make build           # 构建 Anthropic 版本
make build-all       # 构建所有平台

# 运行
make start           # 快速启动 (推荐)
make run-openai      # OpenAI 模式
make run             # Anthropic 模式

# 开发
make dev-openai      # 开发模式 (OpenAI)
make dev             # 开发模式 (Anthropic)

# 工具
make test            # 运行测试
make fmt             # 格式化代码
make lint            # 代码检查
make clean           # 清理构建文件
```

### 特点
- ✅ 专业的构建工具
- ✅ 多种构建选项
- ✅ 支持多平台交叉编译
- ✅ 集成测试和代码检查

---

## 方式3：使用 Docker

### OpenAI 兼容模式 ⭐

```bash
# 构建镜像
docker build --target openai -t factory-proxy-openai .

# 运行容器
docker run -d \
  --name factory-proxy-openai \
  -p 8003:8003 \
  -e PORT=8003 \
  factory-proxy-openai

# 查看日志
docker logs -f factory-proxy-openai

# 测试
curl http://localhost:8003/v1/health
```

### Anthropic 原生模式

```bash
# 构建镜像
docker build --target anthropic -t factory-proxy-openai-anthropic .

# 运行容器
docker run -d \
  --name factory-proxy-openai-anthropic \
  -p 8001:8000 \
  -e PORT=8000 \
  factory-proxy-openai-anthropic

# 测试
curl http://localhost:8001/health
```

### 特点
- ✅ 容器化部署
- ✅ 隔离环境
- ✅ 多阶段构建（镜像更小）
- ✅ 内置健康检查

---

## 方式4：使用 Docker Compose

### 启动 OpenAI 兼容模式 ⭐

```bash
# 启动服务
docker-compose up -d factory-proxy-openai

# 查看状态
docker-compose ps

# 查看日志
docker-compose logs -f factory-proxy-openai

# 测试
curl http://localhost:8003/v1/health
```

### 启动 Anthropic 原生模式

```bash
docker-compose up -d factory-proxy-openai-anthropic

# 测试
curl http://localhost:8001/health
```

### 同时启动两个模式

```bash
# 启动所有服务
docker-compose up -d

# OpenAI 模式: http://localhost:8003
# Anthropic 模式: http://localhost:8001
```

### 管理命令

```bash
# 停止服务
docker-compose stop

# 重启服务
docker-compose restart factory-proxy-openai

# 停止并删除容器
docker-compose down

# 查看日志
docker-compose logs -f

# 重新构建并启动
docker-compose up -d --build
```

### 特点
- ✅ 最简单的容器化部署
- ✅ 一键启动多个服务
- ✅ 自动重启
- ✅ 健康检查

---

## 方式5：直接运行二进制文件

### 构建

```bash
# OpenAI 模式
go build -ldflags="-s -w" -o factory-proxy-openai main-openai.go

# Anthropic 模式
go build -ldflags="-s -w" -o factory-proxy-openai main.go
```

### 运行

```bash
# OpenAI 模式 (推荐)
PORT=8003 ./factory-proxy-openai

# Anthropic 模式
PORT=8000 ./factory-proxy-openai
```

### 特点
- ✅ 最轻量级
- ✅ 单文件部署
- ✅ 无依赖
- ✅ 启动最快

---

## 生产环境部署

### 1. 使用 systemd (推荐)

#### OpenAI 兼容模式

创建服务文件 `/etc/systemd/system/factory-proxy-openai.service`:

```ini
[Unit]
Description=Factory Proxy API - OpenAI Compatible Mode
After=network.target

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/opt/factory-proxy-openai
Environment="PORT=8003"
ExecStart=/opt/factory-proxy-openai/factory-proxy-openai
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

# 安全加固
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/factory-proxy-openai

[Install]
WantedBy=multi-user.target
```

启动服务:

```bash
# 重载配置
sudo systemctl daemon-reload

# 启用开机自启
sudo systemctl enable factory-proxy-openai

# 启动服务
sudo systemctl start factory-proxy-openai

# 查看状态
sudo systemctl status factory-proxy-openai

# 查看日志
sudo journalctl -u factory-proxy-openai -f
```

### 2. 使用 Nginx 反向代理

```nginx
# /etc/nginx/sites-available/factory-proxy-openai

upstream factory_proxy {
    server 127.0.0.1:8003;
    keepalive 32;
}

server {
    listen 
80;
    server_name api.yourdomain.com;

    # SSL 配置 (Let's Encrypt)
    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;

    # 安全头
    add_header X-Content-Type-Options nosniff;
    add_header X-Frame-Options DENY;
    add_header X-XSS-Protection "1; mode=block";

    # 日志
    access_log /var/log/nginx/factory-proxy-openai-access.log;
    error_log /var/log/nginx/factory-proxy-openai-error.log;

    # 代理配置
    location / {
        proxy_pass http://factory_proxy;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # 超时设置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # 健康检查端点 (不记录日志)
    location /v1/health {
        proxy_pass http://factory_proxy;
        access_log off;
    }
}

# HTTP 重定向到 HTTPS
server {
    listen 80;
    server_name api.yourdomain.com;
    return 301 https://$server_name$request_uri;
}
```

启用配置:

```bash
# 测试配置
sudo nginx -t

# 重载配置
sudo systemctl reload nginx
```

### 3. 使用 Docker Swarm (集群部署)

```yaml
# docker-stack.yml
version: '3.8'

services:
  factory-proxy-openai:
    image: 
factory-proxy-openai:latest
    deploy:
      replicas: 3
      restart_policy:
        condition: on-failure
        delay: 5s
        max_attempts: 3
      resources:
        limits:
          cpus: '0.5'
          memory: 256M
        reservations:
          cpus: '0.25'
          memory: 128M
    ports:
      - "8003:8003"
    environment:
      - PORT=8003
    networks:
      - factory-network
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8003/v1/health"]
      interval: 30s
      timeout: 3s
      retries: 3

networks:
  factory-network:
    driver: overlay
```

部署:

```bash
# 初始化 Swarm (如果还没有)
docker swarm init

# 部署服务
docker stack deploy -c docker-stack.yml factory-proxy-openai

# 查看服务
docker service ls

# 查看服务日志
docker service logs -f factory-proxy-openai_factory-proxy-openai

# 扩容
docker service scale factory-proxy-openai_factory-proxy-openai=5

# 更新服务
docker service update factory-proxy-openai_factory-proxy-openai

# 删除服务
docker stack rm factory-proxy-openai
```

### 4. 使用 Kubernetes

```yaml
# k8s-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: factory-proxy-openai
  labels:
    app: factory-proxy-openai
spec:
  replicas: 3
  selector:
    matchLabels:
      app: factory-proxy-openai
  template:
    metadata:
      labels:
        app: factory-proxy-openai
    spec:
      containers:
      - name: factory-proxy-openai
        image: factory-proxy-openai:latest
        ports:
        - containerPort: 8003
        env:
        - name: PORT
          value: "8003"
        resources:
          limits:
            cpu: "500m"
            memory: "256Mi"
          requests:
            cpu: "250m"
            memory: "128Mi"
        livenessProbe:
          httpGet:
            path: /v1/health
            port: 8003
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /v1/health
            port: 8003
          initialDelaySeconds: 5
          periodSeconds: 10
---
apiVersion: v1
kind: Service
metadata:
  name: factory-proxy-openai
spec:
  selector:
    app: factory-proxy-openai
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8003
  type: LoadBalancer
```

部署:

```bash
# 应用配置
kubectl apply -f k8s-deployment.yaml

# 查看状态
kubectl get pods
kubectl get services

# 查看日志
kubectl logs -f deployment/factory-proxy-openai

# 扩容
kubectl scale deployment/factory-proxy-openai --replicas=5

# 删除
kubectl delete -f k8s-deployment.yaml
```

---

## 📊 部署方式对比

| 方式 | 难度 | 性能 | 适用场景 | 推荐度 |
|------|------|------|----------|--------|
| **启动脚本** | ⭐ | ⭐⭐⭐⭐⭐ | 开发/测试 | ⭐⭐⭐⭐⭐ |
| **Makefile** | ⭐⭐ | ⭐⭐⭐⭐⭐ | 开发/测试 | ⭐⭐⭐⭐⭐ |
| **Docker** | ⭐⭐ | ⭐⭐⭐⭐ | 单机生产 | ⭐⭐⭐⭐ |
| **Docker Compose** | ⭐⭐ | ⭐⭐⭐⭐ | 单机生产 | ⭐⭐⭐⭐⭐ |
| **直接运行** | ⭐ | ⭐⭐⭐⭐⭐ | 生产环境 | ⭐⭐⭐ |
| **systemd** | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | 生产环境 | ⭐⭐⭐⭐⭐ |
| **Docker Swarm** | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | 集群部署 | ⭐⭐⭐⭐ |
| **Kubernetes** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | 大规模集群 | ⭐⭐⭐⭐⭐ |

---

## 🎯 推荐部署方案

### 开发环境
```bash
# 最简单
./start.sh

# 或使用 Makefile
make start
```

### 测试环境
```bash
# Docker Compose
docker-compose up -d factory-proxy-openai
```

### 单机生产环境
```bash
# systemd + Nginx
sudo systemctl start factory-proxy-openai
```

### 集群生产环境
```bash
# Kubernetes
kubectl apply -f k8s-deployment.yaml
```

---

## 🔐 安全建议

1. **使用 HTTPS**: 生产环境必须使用 HTTPS
2. **限流保护**: 配置 Nginx 或 API Gateway 限流
3. **防火墙**: 只开放必要的端口
4. **监控**: 配置日志和监控系统
5. **备份**: 定期备份配置文件
6. **更新**: 及时更新依赖和系统补丁

---

## 📝 环境变量

```bash
# 端口配置
PORT=8003

# Anthropic 目标 URL (可选)
ANTHROPIC_TARGET_URL=https://your-target-url.com

# 日志级别 (可选)
LOG_LEVEL=info
```

---

## 🔗 相关文档

- [快速开始](QUICK_START.md) - 5分钟快速上手
- [启动脚本使用](START.md) - start.sh 详细说明
- [支持的模型](SUPPORTED_MODELS.md) - 已测试可用的模型列表
- [完整文档](README.md) - 项目主文档

---

**推荐**: 开发时使用 `./start.sh`，生产环境使用 `systemd` 
+ Docker Compose！ 🚀