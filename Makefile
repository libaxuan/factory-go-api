.PHONY: all build build-openai clean test run run-openai help install dev fmt lint

# 默认目标
all: build build-openai

# 构建 Anthropic 原生模式
build:
	@echo "🔨 构建 Anthropic 原生模式..."
	go build -ldflags="-s -w" -o factory-proxy main.go
	@echo "✅ 构建完成: factory-proxy"

# 构建 OpenAI 兼容模式
build-openai:
	@echo "🔨 构建 OpenAI 兼容模式..."
	go build -ldflags="-s -w" -o factory-proxy-openai main-openai.go
	@echo "✅ 构建完成: factory-proxy-openai"

# 安装依赖
install:
	@echo "📦 安装依赖..."
	go mod download
	go mod verify
	@echo "✅ 依赖安装完成"

# 运行 Anthropic 原生模式
run: build
	@echo "🚀 启动 Anthropic 原生模式..."
	./factory-proxy

# 运行 OpenAI 兼容模式
run-openai: build-openai
	@echo "🚀 启动 OpenAI 兼容模式..."
	./factory-proxy-openai

# 开发模式（不构建，直接运行）
dev:
	@echo "🔧 开发模式 - Anthropic 原生..."
	go run main.go

# 开发模式 - OpenAI 兼容
dev-openai:
	@echo "🔧 开发模式 - OpenAI 兼容..."
	go run main-openai.go

# 运行测试
test:
	@echo "🧪 运行测试..."
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	@echo "✅ 测试完成"

# 代码格式化
fmt:
	@echo "🎨 格式化代码..."
	go fmt ./...
	gofmt -w .
	@echo "✅ 格式化完成"

# 代码检查
lint:
	@echo "🔍 检查代码..."
	go vet ./...
	@echo "✅ 检查完成"

# 清理构建文件
clean:
	@echo "🧹 清理构建文件..."
	rm -f factory-proxy factory-proxy-openai
	rm -f *.log *.pid
	rm -f coverage.txt
	go clean -cache
	@echo "✅ 清理完成"

# 构建所有平台的二进制文件
build-all:
	@echo "🔨 构建所有平台..."
	@mkdir -p dist
	# Linux AMD64
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/factory-proxy-linux-amd64 main.go
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/factory-proxy-openai-linux-amd64 main-openai.go
	# Linux ARM64
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dist/factory-proxy-linux-arm64 main.go
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dist/factory-proxy-openai-linux-arm64 main-openai.go
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/factory-proxy-darwin-amd64 main.go
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/factory-proxy-openai-darwin-amd64 main-openai.go
	# macOS ARM64
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/factory-proxy-darwin-arm64 main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/factory-proxy-openai-darwin-arm64 main-openai.go
	# Windows AMD64
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/factory-proxy-windows-amd64.exe main.go
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/factory-proxy-openai-windows-amd64.exe main-openai.go
	@echo "✅ 所有平台构建完成，文件位于 dist/ 目录"

# 显示帮助信息
help:
	@echo "Factory Proxy - Makefile 命令"
	@echo ""
	@echo "使用方法: make [目标]"
	@echo ""
	@echo "可用目标:"
	@echo "  all          - 构建所有版本（默认）"
	@echo "  build        - 构建 Anthropic 原生模式"
	@echo "  build-openai - 构建 OpenAI 兼容模式"
	@echo "  build-all    - 构建所有平台的二进制文件"
	@echo "  install      - 安装 Go 依赖"
	@echo "  run          - 构建并运行 Anthropic 原生模式"
	@echo "  run-openai   - 构建并运行 OpenAI 兼容模式"
	@echo "  dev          - 开发模式（Anthropic，不构建）"
	@echo "  dev-openai   - 开发模式（OpenAI，不构建）"
	@echo "  test         - 运行测试"
	@echo "  fmt          - 格式化代码"
	@echo "  lint         - 代码检查"
	@echo "  clean        - 清理构建文件"
	@echo "  help         - 显示此帮助信息"
	@echo ""
	@echo "示例:"
	@echo "  make build-openai    # 构建 OpenAI 兼容版本"
	@echo "  make run-openai      # 构建并运行 OpenAI 版本"
	@echo "  make dev-openai      # 开发模式运行 OpenAI 版本"
	@echo "  make test            # 运行测试"
	@echo "  make clean           # 清理所有构建文件"