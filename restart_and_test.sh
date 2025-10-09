#!/bin/bash

# 停止旧服务
echo "🛑 停止旧服务..."
pkill -f factory-api 2>/dev/null
sleep 2

# 加载环境变量
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | grep -v '^$' | xargs)
fi

# 重新编译
echo "🔨 重新编译..."
go build -o factory-api main_multimodel.go

if [ $? -ne 0 ]; then
    echo "❌ 编译失败"
    exit 1
fi

# 后台启动服务
echo "🚀 启动服务..."
./factory-api > /dev/null 2>&1 &
SERVER_PID=$!
echo "   PID: $SERVER_PID"

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 5

# 检查服务是否运行
if ! curl -s http://localhost:8003/health > /dev/null; then
    echo "❌ 服务启动失败"
    exit 1
fi

echo "✅ 服务已启动"
echo ""

# 运行测试
echo "🧪 开始测试..."
echo "=================================="
./test_models.sh

# 保存测试结果
TEST_RESULT=$?

echo ""
echo "=================================="
if [ $TEST_RESULT -eq 0 ]; then
    echo "✅ 所有测试通过！"
else
    echo "❌ 部分测试失败"
fi

exit $TEST_RESULT