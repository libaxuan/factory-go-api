
# Factory Proxy - 项目结构

## 📁 目录结构

```
factory-go-api/
├── .github/
│   └── workflows/
│       └── ci.yml              # GitHub Actions CI/CD 配置
├── main.go                     # Anthropic 原生模式主程序
├── main-openai.go              # OpenAI 兼容模式主程序 ⭐
├── go.mod                      # Go 模块定义
├── go.sum                      # Go 依赖校验
├── Makefile                    # 构建和管理脚本
├── Dockerfile                  # Docker 多阶段构建配置
├── docker-compose.yml          # Docker Compose 编排
├── start.sh                    # 快速启动脚本
├── .env.example                # 环境变量示例
├── .gitignore                  # Git 忽略文件
├── LICENSE                     # MIT 许可证
├── README.md                   # 主文档 ⭐
├── README-OpenAI.md            # OpenAI 兼容模式详细文档
├── CONTRIBUTING.md             # 贡献指南
├── CHANGELOG.md                # 更新日志
├── PROJECT_STRUCTURE.md        # 本文件
└── test_openai_sdk.py          # Python SDK 测试脚本
```

## 🔑 核心文件说明

### 源代码文件

#### main.go
- **功能**: Anthropic API 原生格式代理
- **端点**: `/anthropic/*`, `/openai/*`, `/bedrock/*`, `/health`
- **认证**: `x-api-key` 头
- **特点**: 直接代理，无格式转换

#### main-openai.go ⭐
- **功能**: OpenAI 兼容格式代理
- **端点**: `/v1/chat/completions`, `/v1/health`, `/health`
- **认证**: `Authorization: Bearer` 头
- **特点**: 
  - 自动格式转换（OpenAI ↔ Anthropic）
  - 支持 system 消息
  - 自动注入 Factory Droid prompt
  - 完全兼容 OpenAI SDK

### 配置文件

#### go.mod & go.sum
- Go 模块依赖管理
- 依赖: `gorilla/mux` (路由)

#### .env.example
```bash
PORT=8000
ANTHROPIC_TARGET_URL=https://your-endpoint.com
```

#### Makefile
常用命令：
```bash
make build-openai   # 构建 OpenAI 版本
make run-openai     # 运行 OpenAI 版本
make dev-openai     # 开发模式
make test           # 运行测试
make clean          # 清理
make help           # 帮助
```

### 文档文件

#### README.md
- 项目主文档
- 快速开始指南
- 两种模式对比
- 完整的使用示例

#### README-OpenAI.md
- OpenAI 兼容模式详细文档
- API 端点说明
- 格式转换细节
- Python/Node.js SDK 示例
- 生产部署指南

#### CONTRIBUTING.md
- 贡献指南
- 代码规范
- 提交规范
- Pull Request 流程

#### CHANGELOG.md
- 版本更新记录
- 遵循 Keep a Changelog 规范
- 语义化版本控制

### 部署文件

#### Dockerfile
- 多阶段构建
- 两个 target: `anthropic` 和 `openai`
- 最小化镜像 (~20MB)
- 非 root 用户运行
- 内置健康检查

#### docker-compose.yml
- 同时运行两个服务
- Anthropic 模式: `localhost:8001`
- OpenAI 模式: `localhost:8003`
- 自动重启和健康检查

#### start.sh
```bash
./start.sh          # 一键启动
```

### CI/CD

#### .github/workflows/ci.yml
- 多平台构建（Linux, macOS, Windows）
- 多 Go 版本测试（1.21, 1.22）
- 自动化测试
- 代码覆盖率
- 自动发布（tag 触发）

### 测试文件

#### test_openai_sdk.py
- Python OpenAI SDK 兼容性测试
- 测试用例：
  1. 基础对话
  2. System 消息
  3. 多轮对话

## 🎯 使用流程

### 开发流程

```bash
# 1. 克隆项目
git clone https://github.com/yourusername/factory-proxy.git
cd factory-proxy/factory-go

# 2. 安装依赖
make install

# 3. 开发模式运行
make dev-openai

# 4. 测试
make test

# 5. 构建
make build-openai
```

### 部署流程

```bash
# 本地部署
make build-openai
./factory-proxy-openai

# Docker 部署
docker-compose up -d factory-proxy-openai

# 生产部署（systemd）
# 参考 README.md 的 systemd 配置
```

## 📊 代码统计

| 文件 | 行数 | 说明 |
|------|------|------|
| main.go | ~350 | Anthropic 原生模式 |
| main-openai.go | ~350 | OpenAI 兼容模式 |
| 总代码 | ~700 | 核心功能代码 |
| 文档 | ~1500 | 完整文档 |

## 🔄 工作流程

### 请求流程（OpenAI 模式）

```
客户端
  ↓ POST /v1/chat/completions
  ↓ Authorization: Bearer <factory-key>
  ↓ OpenAI 格式请求
  ↓
factory-proxy-openai
  ↓ 1. 提取 API Key
  ↓ 2. 转换为 Anthropic 格式
  ↓ 3. 注入 Factory Droid prompt
  ↓ 4. 
发送到目标 API
  ↓
Factory AI API
  ↓ 返回 Anthropic 格式
  ↓
factory-proxy-openai
  ↓ 5. 转换为 OpenAI 格式
  ↓ 6. 返回响应
  ↓
客户端
  ↓ 收到 OpenAI 格式响应
```

## 🔧 关键技术点

### 1. 格式转换

**OpenAI → Anthropic**
```go
// 提取 system 消息
// 合并 Factory Droid prompt
// 转换 messages 数组
```

**Anthropic → OpenAI**
```go
// 转换 content 数组
// 映射 finish_reason
// 转换 usage 统计
```

### 2. 认证处理

```go
// OpenAI 格式: Authorization: Bearer <key>
// 转换为: x-api-key: <key>
```

### 3. Factory Droid Prompt

自动注入的 system prompt：
```
You are Droid, a helpful assistant created by Factory.
```

## 🚀 性能优化

1. **编译优化**: `-ldflags="-s -w"` 减小二进制大小
2. **静态编译**: `CGO_ENABLED=0` 无外部依赖
3. **多阶段构建**: Docker 镜像 ~20MB
4. **并发处理**: Go 原生协程支持

## 📈 未来计划

- [ ] 流式响应支持（Server-Sent Events）
- [ ] 请求限流和缓存
- [ ] 更多模型支持
- [ ] 监控和指标收集
- [ ] 配置热重载
- [ ] gRPC 接口支持

## 🤝 维护指南

### 添加新功能

1. 修改相应的 `main.go` 或 `main-openai.go`
2. 添加测试用例
3. 更新文档
4. 提交 Pull Request

### 发布新版本

1. 更新 `CHANGELOG.md`
2. 创建 Git tag: `git tag v1.x.x`
3. 推送 tag: `git push origin v1.x.x`
4. GitHub Actions 自动构建和发布

---

**最后更新**: 2025-01-08