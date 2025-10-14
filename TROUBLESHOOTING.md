# 故障排查指南

## Vercel 安全检查点问题

### 问题描述
当 Docker 容器调用接口时返回 Vercel 安全检查点页面（HTML），而不是预期的 JSON API 响应。

### 根本原因
Vercel 使用安全机制检测和拦截可疑请求。缺少正确的浏览器特征请求头会触发安全检查点。

### 解决方案

我们已经在 `transformers/request.go` 中添加了完整的浏览器特征请求头：

#### 1. 对于 Factory OpenAI 端点
```go
headers := map[string]string{
    "content-type":             "application/json",
    "authorization":            authHeader,
    "accept":                   "application/json, text/event-stream",
    "accept-encoding":          "gzip, deflate, br",
    "accept-language":          "en-US,en;q=0.9",
    "connection":               "keep-alive",
    "origin":                   "https://app.factory.ai",
    "referer":                  "https://app.factory.ai/",
    "sec-fetch-dest":           "empty",
    "sec-fetch-mode":           "cors",
    "sec-fetch-site":           "same-origin",
    // ... 其他头部
}
```

#### 2. 对于 Anthropic 端点
```go
headers := map[string]string{
    "content-type":      "application/json",
    "authorization":     authHeader,
    "anthropic-version": "2023-06-01",
    "accept-encoding":   "gzip, deflate, br",
    "accept-language":   "en-US,en;q=0.9",
    "connection":        "keep-alive",
    "origin":            "https://app.factory.ai",
    "referer":           "https://app.factory.ai/",
    "sec-fetch-dest":    "empty",
    "sec-fetch-mode":    "cors",
    "sec-fetch-site":    "same-origin",
}
```

### 关键修改

添加的关键请求头：
- `accept-encoding`: 表明支持压缩
- `accept-language`: 语言偏好
- `origin`: 请求来源
- `referer`: 引用页面
- `sec-fetch-*`: 浏览器安全策略头

### 重新部署步骤

1. **使用自动脚本（推荐）**：
   ```bash
   ./rebuild_docker.sh
   ```

2. **手动部署**：
   ```bash
   # 停止现有容器
   docker-compose down
   
   # 清理旧镜像
   docker rmi factory-proxy-openai factory-proxy-openai-anthropic
   
   # 重新构建（不使用缓存）
   docker-compose build --no-cache
   
   # 启动服务
   docker-compose up -d
   
   # 查看日志
   docker-compose logs -f
   ```

### 验证修复

1. **检查健康状态**：
   ```bash
   curl http://localhost:8003/v1/health
   ```

2. **测试 API 调用**：
   ```bash
   curl -X POST http://localhost:8003/v1/chat/completions \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer your_proxy_api_key" \
     -d '{
       "model": "gpt-5-codex",
       "messages": [{"role": "user", "content": "Hello"}],
       "stream": false
     }'
   ```

3. **检查响应类型**：
   - ✅ 正确：返回 JSON 格式的 API 响应
   - ❌ 错误：返回 HTML 格式的 Vercel 安全检查点页面

### 其他可能的问题

#### 问题 1: FACTORY_API_KEY 无效
**症状**：401 Unauthorized 错误

**解决方案**：
```bash
# 检查 .env 文件
cat .env

# 确保 FACTORY_API_KEY 正确设置
echo "FACTORY_API_KEY=your_actual_key" > .env
```

#### 问题 2: 端口冲突
**症状**：容器无法启动

**解决方案**：
```bash
# 检查端口占用
lsof -i :8003
lsof -i :8001

# 修改 docker-compose.yml 中的端口映射
```

#### 问题 3: 网络连接问题
**症状**：请求超时

**解决方案**：
```bash
# 从容器内测试网络
docker exec factory-proxy-openai wget -O- https://app.factory.ai

# 检查 DNS 解析
docker exec factory-proxy-openai nslookup app.factory.ai
```

### 日志分析

查看详细日志以诊断问题：

```bash
# 实时日志
docker-compose logs -f factory-proxy-openai

# 特定时间段的日志
docker-compose logs --since 10m factory-proxy-openai

# 搜索错误
docker-compose logs factory-proxy-openai | grep -i error
```

### 常见日志消息

1. **✅ 成功响应**：
   ```
   ✅ gpt-5-codex [openai] stream=false
   📥 Factory OpenAI 响应: 200
   ```

2. **❌ Vercel 拦截**：
   ```
   📥 Factory OpenAI 响应: 403
   错误: 解析响应失败: invalid character '<' looking for beginning of value
   ```

3. **❌ API Key 错误**：
   ```
   ❌ API Key 验证失败
   ```

### 联系支持

如果问题仍未解决：
1. 收集完整的日志：`docker-compose logs > debug.log`
2. 检查网络连接
3. 验证 API Key 的有效性
4. 查看 Factory.ai 的状态页面

### 更新历史

- **2025-01-14**: 添加完整的浏览器特征请求头以绕过 Vercel 安全检查点