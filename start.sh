#!/bin/bash

# Factory Go Proxy 快速启动脚本

echo "🚀 Factory Go Proxy - 快速启动"
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

# 构建项目
echo "🔨 构建可执行文件..."
go build -o factory-proxy main.go

if [ $? -eq 0 ]; then
    echo "✅ 构建成功！"
    
    # 设置环境变量（可选）
    export PORT=${PORT:-8000}
    
    echo ""
    echo "📍 启动信息:"
    echo "   端口: $PORT"
    echo "   服务: http://localhost:$PORT"
    echo ""
    echo "📋 API 端点:"
    echo "   - Anthropic: http://localhost:$PORT/anthropic/*"
    echo "   - OpenAI:   http://localhost:$PORT/openai/*" 
    echo "   - Bedrock:  http://localhost:$PORT/bedrock/*"
    echo "   - 健康检查: http://localhost:$PORT/health"
    echo ""
    echo "⏳ 启动服务器..."
    
    # 启动服务
    ./factory-proxy
    
else
    echo "❌ 构建失败"
    exit 1
fi