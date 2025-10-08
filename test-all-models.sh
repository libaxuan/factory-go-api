#!/bin/bash

# 测试所有模型的脚本
API_URL="http://localhost:8003/v1/chat/completions"
API_KEY="0000"

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "================================"
echo "测试 Factory AI 代理支持的模型"
echo "================================"
echo ""

# 存储结果
WORKING_MODELS=()
FAILED_MODELS=()

# 测试函数
test_model() {
    local model_name=$1
    local provider=$2
    
    echo -n "测试 $provider: $model_name ... "
    
    # 发送请求，设置15秒超时
    response=$(curl -s --max-time 15 -w "\n%{http_code}" -X POST "$API_URL" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $API_KEY" \
        -d "{
            \"model\": \"$model_name\",
            \"messages\": [{\"role\": \"user\", \"content\": \"Hi\"}]
        }" 2>&1)
    
    # 获取HTTP状态码（最后一行）
    http_code=$(echo "$response" | tail -n 1)
    # 获取响应体（除最后一行外）
    body=$(echo "$response" | sed '$d')
    
    # 检查是否成功
    if [ "$http_code" = "200" ]; then
        echo -e "${GREEN}✅ 成功${NC}"
        WORKING_MODELS+=("$provider: $model_name")
        return 0
    else
        # 提取错误信息
        error_msg=$(echo "$body" | grep -o '"message":"[^"]*"' | head -1 | cut -d'"' -f4)
        if [ -z "$error_msg" ]; then
            error_msg="HTTP $http_code"
        fi
        echo -e "${RED}❌ 失败${NC} - $error_msg"
        FAILED_MODELS+=("$provider: $model_name - $error_msg")
        return 1
    fi
}

echo "📋 开始测试..."
echo ""

# Claude 模型
echo "=== Anthropic (Claude) ==="
test_model "claude-3-5-sonnet-20241022" "Claude"
test_model "claude-3-5-sonnet-20250219" "Claude"
test_model "claude-3-7-sonnet-20250219" "Claude"
test_model "claude-sonnet-4-20250514" "Claude"
test_model "claude-sonnet-4-1-20250514" "Claude"
test_model "claude-sonnet-4-5-20250929" "Claude"
test_model "claude-3-5-haiku-20241022" "Claude"
test_model "claude-3-haiku-20240307" "Claude"
echo ""

# OpenAI 模型
echo "=== OpenAI (GPT) ==="
test_model "gpt-5-2025-08-07" "OpenAI"
test_model "gpt-5-mini-2025-08-07" "OpenAI"
test_model "gpt-5-nano-2025-08-07" "OpenAI"
test_model "gpt-5-codex" "OpenAI"
test_model "o1-2024-12-17" "OpenAI"
test_model "o1-mini-2024-09-12" "OpenAI"
test_model "o1-preview-2024-09-12" "OpenAI"
test_model "o3-mini" "OpenAI"
test_model "o4-mini" "OpenAI"
test_model "gpt-4.1-preview" "OpenAI"
test_model "gpt-4o-2024-11-20" "OpenAI"
test_model "gpt-4o-mini-2024-07-18" "OpenAI"
echo ""

# Google 模型
echo "=== Google (Gemini) ==="
test_model "gemini-2.5-flash" "Gemini"
test_model "gemini-2.5-pro" "Gemini"
test_model "gemini-2.0-flash-exp" "Gemini"
test_model "gemini-exp-1206" "Gemini"
echo ""

# xAI 模型
echo "=== xAI (Grok) ==="
test_model "grok-4" "Grok"
test_model "grok-beta" "Grok"
test_model "grok-vision-beta" "Grok"
echo ""

# 其他模型
echo "=== 其他模型 ==="
test_model "glm-4.6" "GLM"
echo ""

# 输出汇总
echo "================================"
echo "测试完成！汇总结果："
echo "================================"
echo ""

echo -e "${GREEN}✅ 工作正常的模型 (${#WORKING_MODELS[@]})：${NC}"
if [ ${#WORKING_MODELS[@]} -eq 0 ]; then
    echo "  无"
else
    for model in "${WORKING_MODELS[@]}"; do
        echo "  - $model"
    done
fi
echo ""

echo -e "${RED}❌ 失败的模型 (${#FAILED_MODELS[@]})：${NC}"
if [ ${#FAILED_MODELS[@]} -eq 0 ]; then
    echo "  无"
else
    for model in "${FAILED_MODELS[@]}"; do
        echo "  - $model"
    done
fi
echo ""

echo "================================"
echo "支持率: ${#WORKING_MODELS[@]}/$((${#WORKING_MODELS[@]} + ${#FAILED_MODELS[@]}))"
echo "================================"