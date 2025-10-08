
# Contributing to Factory Proxy

感谢你对 Factory Proxy 项目的关注！我们欢迎各种形式的贡献。

## 🤝 如何贡献

### 报告问题 (Bug Reports)

如果你发现了 bug，请创建一个 issue 并包含以下信息：

1. **清晰的标题**：简要描述问题
2. **详细描述**：
   - 你期望的行为
   - 实际发生的行为
   - 复现步骤
3. **环境信息**：
   - Go 版本
   - 操作系统
   - 相关配置
4. **日志或错误信息**

### 功能请求 (Feature Requests)

我们欢迎新功能的建议！请创建 issue 并说明：

1. **功能描述**：你希望添加什么功能
2. **使用场景**：为什么需要这个功能
3. **实现建议**：（可选）你的实现想法

### 提交代码 (Pull Requests)

1. **Fork 项目**
   ```bash
   git clone https://github.com/YOUR_USERNAME/factory-proxy.git
   cd factory-proxy/factory-go
   ```

2. **创建分支**
   ```bash
   git checkout -b feature/your-feature-name
   # 或
   git checkout -b fix/your-bug-fix
   ```

3. **编写代码**
   - 遵循现有的代码风格
   - 添加必要的注释
   - 确保代码通过测试

4. **测试你的更改**
   ```bash
   # 编译测试
   go build -o factory-proxy main.go
   go build -o factory-proxy-openai main-openai.go
   
   # 运行测试
   go test -v ./...
   
   # 手动测试
   ./factory-proxy
   curl http://localhost:8000/health
   ```

5. **提交更改**
   ```bash
   git add .
   git commit -m "feat: 添加新功能描述"
   # 或
   git commit -m "fix: 修复bug描述"
   ```

6. **推送到 GitHub**
   ```bash
   git push origin feature/your-feature-name
   ```

7. **创建 Pull Request**
   - 访问你的 fork 仓库
   - 点击 "New Pull Request"
   - 填写 PR 描述，说明你的更改

## 📝 提交信息规范

我们使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

- `feat: 添加新功能`
- `fix: 修复bug`
- `docs: 文档更新`
- `style: 代码格式（不影响功能）`
- `refactor: 重构代码`
- `perf: 性能优化`
- `test: 添加测试`
- `chore: 构建/工具链更新`

示例：
```
feat: 添加流式响应支持

- 实现了 Server-Sent Events
- 添加了流式测试用例
- 更新了文档

Closes #123
```

## 🎨 代码风格

### Go 代码规范

1. **使用 `gofmt` 格式化代码**
   ```bash
   gofmt -w .
   ```

2. **遵循 Go 命名约定**
   - 包名：小写，简短
   - 函数名：驼峰命名，导出函数首字母大写
   - 变量名：驼峰命名

3. **添加注释**
   ```go
   // 导出函数必须有注释
   // ProxyHandler 处理代理请求
   func ProxyHandler(targetURL string) http.HandlerFunc {
       // 实现...
   }
   ```

4. **错误处理**
   ```go
   // 明确处理每个错误
   if err != nil {
       log.Printf("错误: %v", err)
       return err
   }
   ```

## 🧪 测试

### 运行测试

```bash
# 运行所有测试
go test -v ./...

# 运行特定测试
go test -v -run TestProxyHandler

# 测试覆盖率
go test -cover ./...
```

### 添加测试

为新功能添加测试：

```go
func TestNewFeature(t *testing.T) {
    // 准备
    expected := "expected value"
    
    // 执行
    result := YourNewFunction()
    
    // 验证
    if result != expected {
        t.Errorf("期望 %s, 得到 %s", expected, result)
    }
}
```

## 📚 文档

### 更新文档

如果你的更改影响了用户使用方式，请更新相应文档：

- `README.md` - 主要使用文档
- `README-OpenAI.md` - OpenAI 兼容接口文档
- 代码注释

### 文档风格

- 使用清晰、简洁的语言
- 提供代码示例
- 包含使用场景说明

## 🔍 代码审查

所有的 Pull Request 都会经过代码审查。审查者会关注：

1. **代码质量**
   - 是否遵循项目规范
   - 是否有充分的错误处理
   - 是否有必要的注释

2. **功能完整性**
   - 是否解决了问题
   - 是否有测试覆盖
   - 是否有文档更新

3. **向后兼容性**
   - 是否破坏了现有 API
   - 是否需要版本升级

## 📦 发布流程

项目维护者负责发布新版本：

1. 更新 `CHANGELOG.md`
2. 更新版本号
3. 创建 Git tag
4. 发布 GitHub Release
5. 更新文档

## 💬 交流

- **Issues**: 用于 bug 报告和功能请求
- **Discussions**: 用于一般性讨论和问题
- **Pull Requests**: 用于代码贡献

## 🙏 感谢

感谢所有贡献者让这个项目变得更好！

## 📄 许可证



本项目采用 [MIT License](LICENSE)，贡献的代码将在相同许可下发布。