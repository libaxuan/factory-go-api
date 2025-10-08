#!/bin/bash

# Factory Proxy API - 快速启动脚本
# 默认启动 OpenAI 兼容模式 (推荐)

echo "🚀 Factory Proxy API - 快速启动"
echo "=================================="

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo "❌ Go 未安装，请先安装 Go: https://golang.org/dl/"
    exit 1
fi

echo "✅ Go 版本: $(go version)"

# 安装依赖
echo "📦 安装依赖..."
go mod tidy

# 检查命令行参数，决定构建哪个版本
MODE=${1:-openai}  # 默认为 openai 模式

if [ "$MODE" = "anthropic" ]; then
    echo "🔨 构建 Anthropic 原生模式..."
    go build -o factory-proxy main.go
    BINARY="factory-proxy"
    API_MODE="Anthropic 原生模式"
else
    echo "🔨 构建 OpenAI 兼容模式... ⭐"
    go build -o factory-proxy-openai main-openai.go
    BINARY="factory-proxy-openai"
    API_MODE="OpenAI 兼容模式"
fi

if [ $? -eq 0 ]; then
    echo "✅ 构建成功！"
    
    # 设置环境变量（可选）
    export PORT=${PORT:-8003}
    
    echo ""
    echo "📍 启动信息:"
    echo "   模式: $API_MODE"
    echo "   端口: $PORT"
    echo "   服务: http://localhost:$PORT"
    echo ""
    
    if [ "$MODE" = "anthropic" ]; then
        echo "📋 API 端点 (Anthropic 原生模式):"
        echo "   - Anthropic: http://localhost:$PORT/anthropic/v1/messages"
        echo "   - OpenAI:    http://localhost:$PORT/openai/v1/chat/completions"
        echo "   - Bedrock:   http://localhost:$PORT/bedrock/v1/messages"
        echo "   - 健康检查:  http://localhost:$PORT/health"
        echo ""
        echo "🔑 认证方式:"
        echo "   x-api-key: YOUR_FACTORY_API_KEY"
    else
        echo "📋 API 端点 (OpenAI 兼容模式) ⭐:"
        echo "   - Chat:      http://localhost:$PORT/v1/chat/completions"
        echo "   - 健康检查:  http://localhost:$PORT/v1/health"
        echo ""
        echo "🔑 认证方式:"
        echo "   Authorization: Bearer YOUR_FACTORY_API_KEY"
        echo ""
        echo "💡 快速测试:"
        echo "   curl http://localhost:$PORT/v1/health"
    fi
    
    echo ""
    echo "⏳ 启动服务器..."
    echo "=================================="
    echo ""
    
    # 启动服务
    ./$BINARY
    
else
    echo "❌ 构建失败"
    exit 1
fi