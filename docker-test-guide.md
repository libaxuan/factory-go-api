# Docker 部署测试指南

## 当前状态
❌ Docker Desktop 未找到或未安装

## 安装步骤

### 1. 安装 Docker Desktop for Mac

请按照以下步骤安装 Docker Desktop：

1. 访问 Docker 官网下载页面：
   - Apple Silicon (M1/M2/M3): https://desktop.docker.com/mac/main/arm64/Docker.dmg
   - Intel 芯片: https://desktop.docker.com/mac/main/amd64/Docker.dmg

2. 下载后双击 `.dmg` 文件

3. 将 Docker 图标拖到 Applications 文件夹

4. 从 Applications 文件夹启动 Docker Desktop

5. 首次启动需要授权并等待 Docker 引擎启动完成

### 2. 验证 Docker 安装

安装完成后，在终端运行：

```bash
# 检查 Docker 版本
docker --version
docker-compose --version

# 检查 Docker 是否运行
docker ps
```

## 测试计划

安装完成后，将执行以下测试：

### ✅ 已完成
- [x] 检查 Docker 环境

### 📋 待执行
- [ ] 启动 Docker Desktop
- [ ] 验证 Docker 环境
- [ ] 测试 Docker 镜像构建（OpenAI 模式）
- [ ] 测试 Docker Compose 启动服务
- [ ] 验证服务健康状态
- [ ] 测试 API 接口
- [ ] 清理测试资源

## 快速测试命令

安装完成后，执行以下命令进行完整测试：

### 方式1: 使用 Docker Compose (推荐)

```bash
# 构建并启动 OpenAI 兼容模式
docker-compose up -d factory-proxy-openai

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f factory-proxy-openai

# 测试健康检查
curl http://localhost:8003/v1/health

# 测试 API 文档页面
curl http://localhost:8003/

# 停止服务
docker-compose down
```

### 方式2: 使用纯 Docker 命令

```bash
# 构建镜像
docker build --target openai -t factory-proxy-openai .

# 运行容器
docker run -d \
  --name factory-proxy-openai \
  -p 8003:8003 \
  -e PORT=8003 \
  factory-proxy-openai

# 查看容器状态
docker ps

# 查看日志
docker logs -f factory-proxy-openai

# 测试健康检查
curl http://localhost:8003/v1/health

# 停止并删除容器
docker stop factory-proxy-openai
docker rm factory-proxy-openai

# 删除镜像
docker rmi factory-proxy-openai
```

## 测试 Anthropic 原生模式

```bash
# 使用 Docker Compose
docker-compose up -d factory-proxy-openai-anthropic

# 测试
curl http://localhost:8001/health

# 停止
docker-compose down
```

## 同时测试两种模式

```bash
# 启动所有服务
docker-compose up -d

# 测试 OpenAI 模式
curl http://localhost:8003/v1/health

# 测试 Anthropic 模式
curl http://localhost:8001/health

# 查看所有容器
docker-compose ps

# 停止所有服务
docker-compose down
```

## 常见问题

### 1. Docker 命令未找到
确保 Docker Desktop 已启动，可以在 macOS 顶部菜单栏看到 Docker 图标

### 2. 端口被占用
如果端口 8003 或 8001 被占用，可以修改 `docker-compose.yml` 中的端口映射：
```yaml
ports:
  - "9003:8003"  # 将主机端口改为 9003
```

### 3. 构建失败
检查网络连接，确保可以访问 Go 模块代理和 Docker Hub

### 4. 容器启动失败
查看容器日志：
```bash
docker logs factory-proxy-openai
```

## 性能测试

部署成功后，可以进行性能测试：

```bash
# 使用 Apache Bench (需要安装)
ab -n 100 -c 10 http://localhost:8003/v1/health

# 或使用 curl 循环测试
for i in {1..10}; do
  curl -w "\nTime: %{time_total}s\n" http://localhost:8003/v1/health
done
```

## 下一步

请先安装 Docker Desktop，然后我可以继续帮您完成自动化测试。