#!/bin/bash

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}================================================${NC}"
echo -e "${GREEN}  Factory Go API - Docker 重新构建脚本${NC}"
echo -e "${GREEN}================================================${NC}"

# 检查 .env 文件
if [ ! -f .env ]; then
    echo -e "${RED}❌ 错误: .env 文件不存在${NC}"
    echo -e "${YELLOW}请创建 .env 文件并设置以下环境变量:${NC}"
    echo "FACTORY_API_KEY=your_factory_api_key"
    echo "PROXY_API_KEY=your_proxy_api_key"
    exit 1
fi

# 加载环境变量
source .env

# 检查必要的环境变量
if [ -z "$FACTORY_API_KEY" ]; then
    echo -e "${RED}❌ 错误: FACTORY_API_KEY 未设置${NC}"
    exit 1
fi

echo -e "${YELLOW}📦 步骤 1: 停止并删除现有容器...${NC}"
docker-compose down

echo -e "${YELLOW}📦 步骤 2: 删除旧镜像...${NC}"
docker rmi factory-proxy-openai 2>/dev/null || true
docker rmi factory-proxy-openai-anthropic 2>/dev/null || true

echo -e "${YELLOW}📦 步骤 3: 清理未使用的镜像...${NC}"
docker image prune -f

echo -e "${YELLOW}📦 步骤 4: 重新构建镜像...${NC}"
docker-compose build --no-cache

echo -e "${YELLOW}🚀 步骤 5: 启动容器...${NC}"
docker-compose up -d

echo -e "${YELLOW}⏳ 等待服务启动...${NC}"
sleep 5

echo -e "${GREEN}================================================${NC}"
echo -e "${GREEN}✅ 部署完成！${NC}"
echo -e "${GREEN}================================================${NC}"

echo ""
echo -e "${YELLOW}📊 容器状态:${NC}"
docker-compose ps

echo ""
echo -e "${YELLOW}📝 服务日志 (最近 20 行):${NC}"
echo -e "${GREEN}--- OpenAI 兼容模式 (端口 8003) ---${NC}"
docker-compose logs --tail=20 factory-proxy-openai

echo ""
echo -e "${GREEN}--- Anthropic 原生模式 (端口 8001) ---${NC}"
docker-compose logs --tail=20 factory-proxy-openai-anthropic

echo ""
echo -e "${YELLOW}🔗 可用端点:${NC}"
echo "  • OpenAI 兼容: http://localhost:8003/v1/chat/completions"
echo "  • Anthropic 原生: http://localhost:8001/v1/messages"
echo "  • 健康检查: http://localhost:8003/v1/health"
echo "  • 模型列表: http://localhost:8003/v1/models"
echo ""
echo -e "${YELLOW}💡 查看实时日志:${NC}"
echo "  docker-compose logs -f factory-proxy-openai"
echo ""