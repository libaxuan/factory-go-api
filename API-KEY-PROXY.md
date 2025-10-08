# 🔐 API Key 代理功能说明

## 概述

API Key 代理功能允许您保护源头的 Factory API Key，通过设置自定义的对外 API Key，实现 Key 的隔离和保护。

## 工作原理

```
客户端请求 (使用 PROXY_API_KEY)
    ↓
代理服务器验证 PROXY_API_KEY
    ↓
转发请求到 Factory API (使用 FACTORY_API_KEY)
    ↓
返回响应给客户端
```

### 安全优势

1. **源头 Key 保护**: Factory API Key 只存储在服务器端，不会暴露给客户端
2. **灵活的访问控制**: 可以随时更换对外的 PROXY_API_KEY，无需修改源头 Key
3. **多用户支持**: 未来可扩展为支持多个 PROXY_API_KEY，分配给不同用户
4. **审计日志**: 所有请求都经过代理，便于记录和审计

## 配置方法

### 1. 创建 .env 文件

```bash
cp .env.example .env
```

### 2. 编辑 .env 文件

```bash
# 源头 Factory API Key（不会暴露）
FACTORY_API_KEY=your_real_factory_api_key_here

# 对外代理 Key（客户端使用）
PROXY_API_KEY=your_custom_proxy_key_here

# 可选配置
PORT=8003
```

**重要提示**:
- `FACTORY_API_KEY`: 您从 Factory AI 获取的真实 API Key
- `PROXY_API_KEY`: 您自定义的任意字符串，建议使用复杂的随机字符串

### 3. 生成安全的 PROXY_API_KEY

推荐使用以下方法生成安全的代理 Key：

```bash
# 方法1: 使用 openssl
openssl rand -hex 32

# 方法2: 使用 uuidgen
uuidgen

# 方法3: 使用 Python
python3 -c "import secrets; print(secrets.token_urlsafe(32))"
```

## 使用方法

### 启动服务

```bash
./start.sh
```

服务启动后会显示：
```
🔐 API Key 代理: 已启用
   - 对外 Key: abc12345***
   - 源头 Key: def67890***
```

### 客户端调用

使用 **PROXY_API_KEY** 而不是 Factory API Key：

#### cURL

```bash
curl -X POST http://localhost:8003/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_PROXY_API_KEY" \
  -d '{
    "model": "claude-sonnet-4-5-20250929",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

#### Python (OpenAI SDK)

```python
from openai import OpenAI

client = OpenAI(
    base_url="http://localhost:8003/v1",
    api_key="YOUR_PROXY_API_KEY"  # 使用代理 Key
)

response = client.chat.completions.create(
    model="claude-sonnet-4-5-20250929",
    messages=[{"role": "user", "content": "Hello!"}]
)
print(response.choices[0].message.content)
```

#### Node.js

```javascript
const OpenAI = require('openai');

const client = new OpenAI({
    baseURL: 'http://localhost:8003/v1',
    apiKey: 'YOUR_PROXY_API_KEY'  // 使用代理 Key
});

const response = await client.chat.completions.create({
    model: 'claude-sonnet-4-5-20250929',
    messages: [{ role: 'user', content: 'Hello!' }]
});
console.log(response.choices[0].message.content);
```

## 错误处理

### 缺少环境变量

如果未设置必需的环境变量，服务启动时会报错：

```
❌ 错误: 必须设置 FACTORY_API_KEY 环境变量
❌ 错误: 必须设置 PROXY_API_KEY 环境变量
```

**解决方法**: 确保 .env 文件中包含这两个配置项。

### API Key 验证失败

如果客户端提供的 Key 不匹配，会返回：

```json
{
  "error": {
    "message": "Invalid API key",
    "type": "authentication_error"
  }
}
```

**解决方法**: 检查客户端使用的是否为正确的 PROXY_API_KEY。

## 安全建议

1. **保护 .env 文件**
   - 已添加到 .gitignore，不会被提交到代码库
   - 在服务器上设置适当的文件权限: `chmod 600 .env`

2. **定期轮换 Key**
   - 建议定期更换 PROXY_API_KEY
   - 可以在不影响客户端的情况下更换 FACTORY_API_KEY

3. **使用 HTTPS**
   - 生产环境应使用 HTTPS 加密传输
   - 可通过 Nginx 反向代理实现

4. **日志安全**
   - 日志中只显示 Key 的前 8 位
   - 完整的 Key 不会出现在日志中

## Docker 部署

使用 Docker 时，通过环境变量传递 Key：

```bash
docker run -d \
  -p 8003:8003 \
  -e FACTORY_API_KEY=your_factory_key \
  -e PROXY_API_KEY=your_proxy_key \
  