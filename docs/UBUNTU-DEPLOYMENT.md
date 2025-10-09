# 🚀 Ubuntu 服务器部署完整指南

本文档提供在 Ubuntu 服务器上部署 Factory Go API 的详细步骤，适用于直接运行方式（无需 Docker）。

---

## 📋 目录

- [环境要求](#环境要求)
- [第一步：安装 Go 环境](#第一步安装-go-环境)
- [第二步：上传项目到服务器](#第二步上传项目到服务器)
- [第三步：配置环境变量](#第三步配置环境变量)
- [第四步：构建项目](#第四步构建项目)
- [第五步：测试运行](#第五步测试运行)
- [第六步：配置 systemd 服务](#第六步配置-systemd-服务)
- [第七步：配置防火墙](#第七步配置防火墙)
- [第八步：配置 Nginx 反向代理](#第八步配置-nginx-反向代理可选)
- [日常运维](#日常运维)
- [故障排查](#故障排查)
- [附录](#附录)

---

## 环境要求

- **操作系统**: Ubuntu 18.04+ / Ubuntu 20.04+ / Ubuntu 22.04+
- **Go 版本**: 1.21 或更高
- **内存**: 最低 512MB，推荐 1GB+
- **磁盘**: 最低 100MB 可用空间
- **网络**: 需要访问外网（下载依赖和调用 Factory API）

---

## 第一步：安装 Go 环境

### 1.1 检查是否已安装 Go

```bash
go version
```

如果显示版本号（如 `go version go1.22.0 linux/amd64`），则已安装，可跳过此步骤。

### 1.2 下载并安装 Go

```bash
# 下载 Go 1.22.0（推荐最新稳定版）
cd /tmp
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz

# 删除旧版本（如果存在）
sudo rm -rf /usr/local/go

# 解压到 /usr/local
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz

# 清理安装包
rm go1.22.0.linux-amd64.tar.gz
```

### 1.3 配置环境变量

```bash
# 添加到 .bashrc
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc

# 立即生效
source ~/.bashrc

# 验证安装
go version
```

**预期输出：**
```
go version go1.22.0 linux/amd64
```

---

## 第二步：上传项目到服务器

### 方式 A：使用 Git Clone（推荐）

```bash
# 进入部署目录（推荐使用 /home 或 /opt）
cd /home

# 克隆项目
git clone https://github.com/your-username/factory-go-api.git

# 进入项目目录
cd factory-go-api

# 查看文件
ls -la
```

### 方式 B：使用 SCP 上传

```bash
# 在本地电脑执行（替换为你的服务器信息）
scp -r /path/to/local/factory-go-api root@your-server-ip:/home/

# 然后在服务器上
cd /home/factory-go-api
```

### 方式 C：使用 SFTP 工具

使用 FileZilla、WinSCP 等 SFTP 工具上传整个项目文件夹到 `/home/factory-go-api`

---

## 第三步：配置环境变量

### 3.1 复制环境变量模板

```bash
cd /home/factory-go-api
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

# 对外代理 API Key - 自定义一个复杂字符串
PROXY_API_KEY=your_custom_secure_key_here

# 服务监听端口（默认 8003）
PORT=8003
```

**重要提示：**
- ✅ 不要在值两边加引号
- ✅ 等号前后不要有空格
- ✅ 每行结尾不要有多余空格

### 3.4 保存文件

- 按 `Ctrl + X`
- 按 `Y` 确认
- 按 `Enter` 保存

---

## 第四步：构建项目

### 4.1 下载依赖

```bash
cd /home/factory-go-api
go mod download
go mod tidy
```

### 4.2 构建可执行文件

```bash
# 构建 OpenAI 兼容模式
go build -ldflags="-s -w" -o factory-proxy-openai main_multimodel.go

# 赋予执行权限
chmod +x factory-proxy-openai

# 验证构建成功
ls -lh factory-proxy-openai
```

**预期输出：**
```
-rwxr-xr-x 1 root root 8.5M Oct  9 18:00 factory-proxy-openai
```

---

## 第五步：测试运行

### 5.1 前台测试运行

```bash
cd /home/factory-go-api
./factory-proxy-openai
```

**预期输出：**
```
2025/10/09 18:00:00 🔐 代理模式: 已启用
2025/10/09 18:00:00 📖 加载配置文件: config.json
2025/10/09 18:00:00 ✅ 配置加载成功
2025/10/09 18:00:00 🚀 服务启动于 http://localhost:8003
```

### 5.2 新开一个终端测试 API

```bash
# 健康检查
curl http://localhost:8003/health

# 查看模型列表
curl http://localhost:8003/v1/models \
  -H "Authorization: Bearer YOUR_PROXY_API_KEY"
```

### 5.3 停止测试

在第一个终端按 `Ctrl + C` 停止服务。

---

## 第六步：配置 systemd 服务

### 6.1 创建服务文件

```bash
sudo nano /etc/systemd/system/factory-proxy.service
```

### 6.2 填写配置

```ini
[Unit]
Description=Factory Proxy OpenAI Compatible API
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/home/factory-go-api
EnvironmentFile=/home/factory-go-api/.env
ExecStart=/home/factory-go-api/factory-proxy-openai
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

保存：`Ctrl + X` → `Y` → `Enter`

### 6.3 启动服务

```bash
# 重载配置
sudo systemctl daemon-reload

# 启用开机自启
sudo systemctl enable factory-proxy

# 启动服务
sudo systemctl start factory-proxy

# 查看状态
sudo systemctl status factory-proxy
```

**成功标志：** 看到 `Active: active (running)` 和绿色圆点

### 6.4 查看日志

```bash
# 实时日志
sudo journalctl -u factory-proxy -f

# 最近50行
sudo journalctl -u factory-proxy -n 50
```

---

## 第七步：配置防火墙

```bash
# 检查状态
sudo ufw status

# 允许 8003 端口
sudo ufw allow 8003/tcp

# 启用防火墙（如果未启用）
sudo ufw enable

# 验证规则
sudo ufw status numbered
```

**云服务器额外步骤：**
登录云服务商控制台（阿里云/腾讯云等），在安全组中添加入站规则：TCP 8003

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
sudo nano /etc/nginx/sites-available/factory-proxy
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

    access_log /var/log/nginx/factory-proxy-access.log;
    error_log /var/log/nginx/factory-proxy-error.log;

    location / {
        proxy_pass http://factory_proxy;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
}
```

### 8.4 启用配置

```bash
# 创建软链接
sudo ln -s /etc/nginx/sites-available/factory-proxy /etc/nginx/sites-enabled/

# 测试配置
sudo nginx 
-t

# 重载配置
sudo systemctl reload nginx
```

### 8.5 配置 HTTPS（可选）

```bash
# 安装 Certbot
sudo apt install certbot python3-certbot-nginx -y

# 自动配置 SSL
sudo certbot --nginx -d api.yourdomain.com

# 测试自动续期
sudo certbot renew --dry-run
```

---

## 日常运维

### systemd 服务管理

```bash
# 查看服务状态
sudo systemctl status factory-proxy

# 启动服务
sudo systemctl start factory-proxy

# 停止服务
sudo systemctl stop factory-proxy

# 重启服务
sudo systemctl restart factory-proxy

# 查看是否开机自启
sudo systemctl is-enabled factory-proxy
```

### 日志管理

```bash
# 实时查看日志
sudo journalctl -u factory-proxy -f

# 查看最近 100 行
sudo journalctl -u factory-proxy -n 100

# 查看今天的日志
sudo journalctl -u factory-proxy --since today

# 查看特定时间段
sudo journalctl -u factory-proxy --since "2025-10-09 10:00:00" --until "2025-10-09 12:00:00"

# 清理旧日志（保留最近 7 天）
sudo journalctl --vacuum-time=7d
```

### 修改环境变量

```bash
# 1. 编辑 .env 文件
cd /home/factory-go-api
nano .env

# 2. 修改后保存

# 3. 重启服务使配置生效
sudo systemctl restart factory-proxy

# 4. 查看服务状态
sudo systemctl status factory-proxy

# 5. 查看日志确认新配置生效
sudo journalctl -u factory-proxy -n 20
```

### 代码更新

```bash
# 1. 进入项目目录
cd /home/factory-go-api

# 2. 拉取最新代码
git pull

# 3. 重新构建
go build -ldflags="-s -w" -o factory-proxy-openai main_multimodel.go

# 4. 重启服务
sudo systemctl restart factory-proxy

# 5. 验证服务
curl http://localhost:8003/health
```

### 性能监控

```bash
# 查看服务资源占用
sudo systemctl status factory-proxy

# 查看进程详细信息
ps aux | grep factory-proxy

# 查看端口监听
sudo netstat -tulnp | grep 8003

# 或使用 ss
sudo ss -tulnp | grep 8003
```

---

## 故障排查

### 问题 1：服务启动失败

**症状：** `sudo systemctl status factory-proxy` 显示 `failed` 或 `inactive`

**排查步骤：**

```bash
# 1. 查看详细日志
sudo journalctl -u factory-proxy -n 100 --no-pager

# 2. 常见错误及解决方案：

# 错误：FACTORY_API_KEY 未配置
# 解决：检查 .env 文件是否存在且格式正确
cat /home/factory-go-api/.env

# 错误：端口被占用 (address already in use)
# 解决：查找并停止占用端口的进程
sudo lsof -i :8003
sudo kill -9 <PID>

# 错误：文件不存在 (no such file)
# 解决：检查路径和文件是否存在
ls -la /home/factory-go-api/factory-proxy-openai

# 3. 手动测试运行
cd /home/factory-go-api
./factory-proxy-openai
```

### 问题 2：API 返回 401 未授权

**症状：** curl 请求返回 `{"error": {"message": "Invalid API key", ...}}`

**解决方案：**

```bash
# 1. 确认使用的是 PROXY_API_KEY（不是 FACTORY_API_KEY）
cat /home/factory-go-api/.env | grep PROXY_API_KEY

# 2. 测试正确的 Key
curl http://localhost:8003/v1/models \
  -H "Authorization: Bearer YOUR_PROXY_API_KEY_FROM_ENV"

# 3. 查看日志中的验证失败信息
sudo journalctl -u factory-proxy -n 50 | grep "验证失败"
```

### 问题 3：健康检查返回 404

**症状：** `curl http://localhost:8003/v1/health` 返回 `{"error": "Not found"}`

**解决方案：**

```bash
# 正确的健康检查路径是 /health（不是 /v1/health）
curl http://localhost:8003/health

# API 端点路径：
# ✅ /health
# ✅ /v1/models
# ✅ /v1/chat/completions
# ✅ /docs
```

### 问题 4：无法从外网访问

**排查步骤：**

```bash
# 1. 检查服务是否运行
sudo systemctl status factory-proxy

# 2. 检查本地端口监听
sudo netstat -tulnp | grep 8003

# 3. 检查防火墙规则
sudo ufw status

# 4. 检查云服务商安全组
# 登录控制台查看是否开放了 8003 端口

# 5. 测试从服务器本地访问
curl http://localhost:8003/health

# 6. 测试从公网访问（替换为实际 IP）
curl http://YOUR_SERVER_IP:8003/health
```

### 问题 5：修改 .env 后不生效

**解决方案：**

```bash
# systemd 服务需要重启才能重新加载环境变量
sudo systemctl restart factory-proxy

# 查看日志确认新配置
sudo journalctl -u factory-proxy -n 20

# 确认 systemd 配置中有 EnvironmentFile
cat /etc/systemd/system/factory-proxy.service | grep EnvironmentFile
```

### 问题 6：构建失败

**常见错误：**

```bash
# 错误：依赖下载失败
# 解决：使用国内镜像
go env -w GOPROXY=https://goproxy.cn,direct
go mod download

# 错误：磁盘空间不足
# 解决：清理空间
df -h
sudo apt clean
sudo journalctl --vacuum-time=3d

# 错误：Go 版本过低
# 解决：升级 Go 到 1.21+
go version
# 
参考第一步重新安装
```

### 快速诊断命令

```bash
# 一键诊断脚本
cat > /tmp/diagnose.sh << 'EOF'
#!/bin/bash
echo "=== Factory Proxy 诊断信息 ==="
echo ""
echo "1. Go 版本:"
go version
echo ""
echo "2. 服务状态:"
sudo systemctl status factory-proxy --no-pager
echo ""
echo "3. 端口监听:"
sudo netstat -tulnp | grep 8003
echo ""
echo "4. 最近日志:"
sudo journalctl -u factory-proxy -n 20 --no-pager
echo ""
echo "5. 环境文件:"
ls -la /home/factory-go-api/.env
echo ""
echo "6. 可执行文件:"
ls -lh /home/factory-go-api/factory-proxy-openai
echo ""
echo "7. 防火墙状态:"
sudo ufw status
echo ""
echo "8. 健康检查:"
curl -s http://localhost:8003/health
echo ""
EOF

chmod +x /tmp/diagnose.sh
/tmp/diagnose.sh
```

---

## 附录

### A. 完整的 .env 示例

```bash
# Factory API Key（必需）
FACTORY_API_KEY=fk-your-factory-api-key-here

# 对外代理 API Key（必需）
PROXY_API_KEY=your-custom-proxy-key-here

# 服务端口（可选，默认 8003）
PORT=8003
```

### B. 完整的 systemd 服务文件

```ini
[Unit]
Description=Factory Proxy OpenAI Compatible API
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/home/factory-go-api
EnvironmentFile=/home/factory-go-api/.env
ExecStart=/home/factory-go-api/factory-proxy-openai
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
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
  -H "Authorization: Bearer YOUR_KEY"

# 非流式对话
curl -X POST http://localhost:8003/v1/chat/completions \
  -H "Authorization: Bearer YOUR_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "claude-sonnet-4-5-20250929",
    "messages": [{"role": "user", "content": "Hello"}],
    "stream": false
  }'

# 流式对话
curl -X POST http://localhost:8003/v1/chat/completions \
  -H "Authorization: Bearer YOUR_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-5-2025-08-07",
    "messages": [{"role": "user", "content": "Tell me a story"}],
    "stream": true
  }'
```

### E. 常用运维命令速查

```bash
# 服务管理
sudo systemctl start factory-proxy      # 启动
sudo systemctl stop factory-proxy       # 停止
sudo systemctl restart factory-proxy    # 重启
sudo systemctl status factory-proxy     # 状态

# 日志查看
sudo journalctl -u factory-proxy -f     # 实时日志
sudo journalctl -u factory-proxy -n 50  # 最近50行

# 端口检查
sudo lsof -i :8003                      # 查看端口占用
sudo netstat -tulnp | grep 8003         # 查看端口监听

# 进程管理
ps aux | grep factory-proxy             # 查看进程
sudo pkill factory-proxy-openai         # 停止进程
```

### F. 支持的模型列表

| 模型 ID | 类型 | 说明 |
|---------|------|------|
| `claude-opus-4-1-20250805` | Anthropic | Claude Opus 4.1 |
| `claude-sonnet-4-20250514` | Anthropic | Claude Sonnet 4 |
| `claude-sonnet-4-5-20250929` | Anthropic | Claude Sonnet 4.5（推荐） |
| `gpt-5-2025-08-07` | OpenAI | GPT-5 |
| `gpt-5-codex` | OpenAI | GPT-5 Codex |

---

## 📞 获取帮助

- **项目主页**: https://github.com/your-username/factory-go-api
- **问题反馈**: https://github.com/your-username/factory-go-api/issues
- **Factory AI 官网**: https://factory.ai

---

## 🎉 部署成功检查清单

- [ ] Go 环境已安装（1.21+）
- [ ] 项目代码已上传到服务器
- [ ] .env 文件已正确配置
- [ ] 项目已成功构建
- [ ] 前台测试运行成功
- [ ] systemd 服务已配置并启动
- [ ] 服务状态显示 `active (running)`
- [ ] 健康检查返回 `{"status":"healthy"}`
- [ ] 防火墙已开放 8003 端口
- [ ] 可以从外网访问（如需要）
- [ ] Nginx 反向代理已配置（如需要）

---

**恭喜！🎊 你已成功在 Ubuntu 服务器上部署 Factory Go API！**

如有问题，请参考「故障排查」章节或提交 Issue。