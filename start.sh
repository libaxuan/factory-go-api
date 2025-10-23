#!/bin/bash

# Factory Go API - 快速启动脚本
# 设置 VERBOSE=true 启用详细日志

# 加载 .env 文件（如果存在）
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | grep -v '^$' | xargs)
fi

# 检查必需的环境变量
if [ -z "$FACTORY_API_KEY" ]; then
    echo "❌ 错误: 未设置 FACTORY_API_KEY"
    echo "   请在 .env 文件中配置或设置环境变量"
    exit 1
fi

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo "❌ Go 未安装: https://golang.org/dl/"
    exit 1
fi

# 检查并关闭占用 8003 端口的进程
PORT_PID=$(lsof -ti:8003 2>/dev/null)
if [ -n "$PORT_PID" ]; then
    [ "$VERBOSE" = "true" ] && echo "🔪 终止旧进程 $PORT_PID..."
    kill -9 $PORT_PID 2>/dev/null
    sleep 1
fi

# 安装依赖和构建（静默模式）
if [ "$VERBOSE" = "true" ]; then
    echo "📦 安装依赖..."
    go mod tidy
    echo "🔨 构建服务..."
    go build -o factory-api main_multimodel.go
else
    go mod tidy > /dev/null 2>&1
    go build -o factory-api main_multimodel.go > /dev/null 2>&1
fi

if [ $? -eq 0 ]; then
    echo "🚀 Factory Go API 启动成功"
    echo "   地址: http://localhost:8003"
    [ "$VERBOSE" = "true" ] && echo "   文档: http://localhost:8003/docs"
    echo ""
    
    # 启动服务
    ./factory-api
else
    echo "❌ 构建失败"
    exit 1
fi