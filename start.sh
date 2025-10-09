#!/bin/bash

# Factory Go API - 多模型支持启动脚本

echo "🚀 Factory Go API - 多模型支持"
echo "=================================="

# 加载 .env 文件（如果存在）
if [ -f .env ]; then
    echo "📄 加载 .env 配置文件..."
    export $(cat .env | grep -v '^#' | grep -v '^$' | xargs)
    echo "✅ 环境变量已加载"
else
    echo "⚠️  未找到 .env 文件，将使用环境变量或默认值"
    echo "   提示: 复制 .env.example 为 .env 并配置 FACTORY_API_KEY"
fi

# 检查必需的环境变量
if [ -z "$FACTORY_API_KEY" ]; then
    echo "❌ 错误: 未设置 FACTORY_API_KEY 环境变量"
    echo "   请在 .env 文件中设置或通过环境变量设置"
    exit 1
fi
echo "✅ FACTORY_API_KEY 已配置（源头 Key）"

if [ -n "$PROXY_API_KEY" ]; then
    echo "✅ PROXY_API_KEY 已配置（对外代理 Key）"
else
    echo "⚠️  未设置 PROXY_API_KEY，将使用直连模式"
fi

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo "❌ Go 未安装，请先安装 Go: https://golang.org/dl/"
    exit 1
fi

echo "✅ Go 版本: $(go version)"

# 检查并关闭占用 8003 端口的进程
echo "🔍 检查端口 8003..."
PORT_PID=$(lsof -ti:8003)
if [ -n "$PORT_PID" ]; then
    echo "⚠️  端口 8003 已被进程 $PORT_PID 占用"
    echo "🔪 自动终止旧进程..."
    kill -9 $PORT_PID 2>/dev/null
    sleep 1
    echo "✅ 旧进程已终止"
else
    echo "✅ 端口 8003 可用"
fi

# 安装依赖
echo "📦 安装依赖..."
go mod tidy

# 构建多模型版本
echo "🔨 构建多模型支持版本..."
go build -o factory-api main_multimodel.go

if [ $? -eq 0 ]; then
    echo "✅ 构建成功！"
    
    echo ""
    echo "📍 服务信息:"
    echo "   地址: http://localhost:8003"
    echo "   文档: http://localhost:8003/docs"
    echo "   配置: config.json"
    echo ""
    echo "💡 快速测试:"
    echo "   curl http://localhost:8003/health"
    if [ -n "$PROXY_API_KEY" ]; then
        echo "   curl http://localhost:8003/v1/models -H \"Authorization: Bearer $PROXY_API_KEY\""
    else
        echo "   curl http://localhost:8003/v1/models -H \"Authorization: Bearer $FACTORY_API_KEY\""
    fi
    echo ""
    echo "📖 完整文档: cat README.md"
    echo "=================================="
    echo ""
    
    # 启动服务
    ./factory-api
    
else
    echo "❌ 构建失败"
    exit 1
fi