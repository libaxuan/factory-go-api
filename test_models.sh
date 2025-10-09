#!/bin/bash

# Factory Go API - 模型测试脚本 (macOS/Linux)

echo "🧪 Factory Go API - 模型测试"
echo "=================================="

# 加载 .env 文件
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | grep -v '^$' | xargs)
fi

# 获取 API Key
if [ -n "$PROXY_API_KEY" ]; then
    API_KEY="$PROXY_API_KEY"
    echo "🔑 使用 PROXY_API_KEY"
else
    API_KEY="$FACTORY_API_KEY"
    echo "🔑 使用 FACTORY_API_KEY"
fi

if [ -z "$API_KEY" ]; then
    echo "❌ 错误: 未设置 API Key"
    exit 1
fi

# 服务地址
BASE_URL="http://localhost:8003"

# 测试提示词（使用简单的数学问题验证完整输出）
TEST_PROMPT="What is 123 + 456? Just give me the number."

# Extended Thinking 模型的 max_tokens（需要 > 24576）
MAX_TOKENS_HIGH=30000
MAX_TOKENS_LOW=100

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试计数
TOTAL=0
SUCCESS=0
FAILED=0

echo ""
echo "📋 开始测试所有模型..."
echo "=================================="

# 测试函数
test_model() {
    local model_id=$1
    local model_name=$2
    local use_stream=$3
    local max_tokens=$4
    
    TOTAL=$((TOTAL + 1))
    
    echo ""
    echo "[$TOTAL] 测试: $model_name"
    echo "    模型: $model_id"
    echo "    流式: $use_stream | max_tokens: $max_tokens"
    
    # 发送请求
    response=$(curl -s -X POST "$BASE_URL/v1/chat/completions" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $API_KEY" \
        --max-time 30 \
        -d "{
            \"model\": \"$model_id\",
            \"messages\": [{\"role\": \"user\", \"content\": \"$TEST_PROMPT\"}],
            \"stream\": $use_stream,
            \"max_tokens\": $max_tokens
        }")
    
    # 检查响应
    if [ "$use_stream" = "true" ]; then
        # 流式响应检查
        if echo "$response" | grep -q "data:"; then
            echo -e "    状态: ${GREEN}✅ 成功（流式）${NC}"
            # 提取第一个有内容的 chunk
            content=$(echo "$response" | grep -o '"delta":{"content":"[^"]*"' | head -1 | sed 's/.*"content":"\([^"]*\)".*/\1/')
            if [ -n "$content" ]; then
                echo "    响应片段: $content"
            fi
            SUCCESS=$((SUCCESS + 1))
        elif echo "$response" | grep -q '"error"'; then
            error=$(echo "$response" | jq -r '.error.message' 2>/dev/null)
            echo -e "    状态: ${RED}❌ 失败${NC}"
            echo "    错误: $error"
            FAILED=$((FAILED + 1))
        else
            echo -e "    状态: ${YELLOW}⚠️  无数据${NC}"
            echo "    原始: $(echo "$response" | head -c 100)"
            FAILED=$((FAILED + 1))
        fi
    else
        # 非流式响应检查
        if echo "$response" | grep -q '"choices"'; then
            content=$(echo "$response" | jq -r '.choices[0].message.content' 2>/dev/null)
            if [ -n "$content" ] && [ "$content" != "null" ]; then
                echo -e "    状态: ${GREEN}✅ 成功${NC}"
                echo "    响应: $content"
                SUCCESS=$((SUCCESS + 1))
            else
                echo -e "    状态: ${YELLOW}⚠️  响应为空${NC}"
                FAILED=$((FAILED + 1))
            fi
        elif echo "$response" | grep -q '"error"'; then
            error=$(echo "$response" | jq -r '.error.message' 2>/dev/null)
            echo -e "    状态: ${RED}❌ 失败${NC}"
            echo "    错误: $error"
            FAILED=$((FAILED + 1))
        else
            echo -e "    状态: ${RED}❌ 超时或无响应${NC}"
            echo "    原始响应: $(echo "$response" | head -c 200)"
            FAILED=$((FAILED + 1))
        fi
    fi
}

# 测试所有模型（每个模型测试流式和非流式）
echo ""
echo "🤖 Claude 系列 (Anthropic)"
echo "--------------------------------"
test_model "claude-opus-4-1-20250805" "Claude Opus 4.1 (非流式)" "false" "$MAX_TOKENS_LOW"
test_model "claude-opus-4-1-20250805" "Claude Opus 4.1 (流式)" "true" "$MAX_TOKENS_LOW"
test_model "claude-sonnet-4-20250514" "Claude Sonnet 4 (非流式)" "false" "$MAX_TOKENS_HIGH"
test_model "claude-sonnet-4-20250514" "Claude Sonnet 4 (流式)" "true" "$MAX_TOKENS_HIGH"
test_model "claude-sonnet-4-5-20250929" "Claude Sonnet 4.5 (非流式)" "false" "$MAX_TOKENS_HIGH"
test_model "claude-sonnet-4-5-20250929" "Claude Sonnet 4.5 (流式)" "true" "$MAX_TOKENS_HIGH"

echo ""
echo "🚀 GPT 系列 (OpenAI)"
echo "--------------------------------"
test_model "gpt-5-2025-08-07" "GPT-5 (非流式)" "false" "$MAX_TOKENS_HIGH"
test_model "gpt-5-2025-08-07" "GPT-5 (流式)" "true" "$MAX_TOKENS_HIGH"
test_model "gpt-5-codex" "GPT-5 Codex (非流式)" "false" "$MAX_TOKENS_LOW"
test_model "gpt-5-codex" "GPT-5 Codex (流式)" "true" "$MAX_TOKENS_LOW"

# 汇总结果
echo ""
echo "=================================="
echo "📊 测试结果汇总"
echo "=================================="
echo "总测试数: $TOTAL"
echo -e "成功: ${GREEN}$SUCCESS${NC}"
echo -e "失败: ${RED}$FAILED${NC}"
if command -v bc &> /dev/null; then
    SUCCESS_RATE=$(echo "scale=1; ($SUCCESS/$TOTAL)*100" | bc)
    echo "成功率: ${SUCCESS_RATE}%"
else
    echo "成功率: $((SUCCESS*100/TOTAL))%"
fi
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 所有模型测试通过！${NC}"
    exit 0
else
    echo -e "${YELLOW}⚠️  部分模型测试失败，请检查日志${NC}"
    exit 1
fi