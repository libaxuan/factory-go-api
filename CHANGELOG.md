# 更新日志

本文档记录 Factory Go API 项目的所有重要变更。

遵循 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/) 规范，
版本号遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [2.0.1] - 2025-10-10

### 🔄 变更

- **Claude Opus 4.1 Extended Thinking 对齐** - 将 [`claude-opus-4-1-20250805`](config.json:16) 的 `reasoning` 配置从 `"off"` 改为 `"high"`，与 Claude Sonnet 4.5 对齐
  - 启用 Extended Thinking 能力（budget_tokens: 24576）
  - 自动调整 max_tokens 逻辑：当 max_tokens ≤ 24576 时，自动增加到 28576
  - 确保 Claude Opus 4.1 也具备深度推理能力

### 🐛 修复

- **修复 Docker 构建失败** - 解决 `docker compose up -d` 时找不到 `main-openai.go` 文件的问题
  - 修正 [`Dockerfile`](Dockerfile:19) 复制文件路径，使用实际存在的 [`main_multimodel.go`](main_multimodel.go:1)
  - 添加必要的配置文件和目录复制：`config/`, `transformers/`, `config.json`, `docs.html`
  - 统一两个运行阶段的二进制文件为 `factory-proxy`
  - 更新 [`docker-compose.yml`](docker-compose.yml:1) 环境变量配置，使用 `FACTORY_API_KEY` 和 `PROXY_API_KEY`

### 📚 文档

- **测试脚本对齐** - 更新 [`test_models.sh`](test_models.sh:128) 和 [`test_models.bat`](test_models.bat:67)
  - 将 Claude Opus 4.1 测试的 max_tokens 从 100 改为 30000，与其他 Extended Thinking 模型一致

## [2.0.0] - 2025-10-09

### ✨ 新增

- **多模型支持架构** - 统一接口访问 Claude、GPT 等多个模型家族
- **5 个 AI 模型完整支持** - 所有模型均支持流式和非流式响应
  - Claude Opus 4.1, Sonnet 4, Sonnet 4.5（Anthropic）
  - GPT-5, GPT-5 Codex（OpenAI）
- **配置文件系统** - 通过 `config.json` 灵活配置模型和端点
- **API 文档页面** - 访问 `/docs` 查看交互式文档
- **模型列表端点** - `GET /v1/models` 获取所有可用模型
- **Extended Thinking 支持** - Claude Sonnet 4/4.5 和 GPT-5 系列的深度推理能力
- **双 Key 认证机制** - 支持 PROXY_API_KEY 和 FACTORY_API_KEY 分离
- **完整测试套件** - 测试所有 5 个模型的流式和非流式，共 10 种配置

### 🔄 变更

- **端口改回 8003** - 与原始配置保持一致
- **简化启动脚本** - `./start.sh` 直接启动多模型模式，无需参数
- **统一 API 格式** - 所有模型使用 OpenAI 兼容格式
- **优化项目结构** - 模块化设计，更易维护

### 🗑️ 移除

- 移除旧的单模型二进制文件
- 移除过时的测试报告文档
- 移除冗余的临时文档

### 🐛 修复

- 修复 GPT 系列非流式响应内容提取（`output[type=message].content[].text`）
- 修复 Claude Extended Thinking 响应解析（`content[type=text]`）
- 修复流式响应的内容处理和事件类型转换
- 修复不同模型类型的路由和转换器选择
- **修复 GPT-5 流式响应问题** - 禁用 Extended Thinking 模式的 reasoning 配置，避免只输出推理过程而无实际答案
- **修复 Factory OpenAI 事件格式** - 支持 GPT 的 `response.output_text.delta` 等新事件类型
- **优化 Token 分配** - 为启用推理的模型自动增加 max_output_tokens，确保有足够空间输出答案
- **修复 Claude Extended Thinking Token 限制** - 自动调整 max_tokens 大于 thinking.budget_tokens
- **修复 Windows 终端中文乱码** - 优化编码设置，移除 emoji 使用文本标记替代（`[OK]`, `[FAIL]`, `[WARN]`）
- **修复编译冲突** - 将旧的 main.go 和 main-openai.go 移至 legacy/ 目录，解决 startTime 重复声明问题

### 📚 文档

- 精简文档至 2 个：README.md（主文档）和 CHANGELOG.md（本文件）
- 更新 README.md 包含完整使用指南和模型支持表格
- 添加多语言代码示例（Python、JavaScript、Go）
- 完善配置说明和故障排除指南
- 更新 docs.html 突出所有模型均支持流式/非流式（10/10 配置）
- 清理多余测试脚本，保留核心启动和测试脚本
- 完善测试脚本覆盖所有 10 种模型配置（5 模型 × 2 模式）
- **Windows 脚本优化** - 移除 emoji 字符，使用 `[OK]`/`[FAIL]`/`[WARN]` 等文本标记，确保在任何终端编码下正常显示
- **测试用例优化** - 使用数学问题 "123 + 456" 替代简单问候，更易验证响应完整性

## [1.0.0] - 2025-10-08

### ✨ 新增

- 🎉 首次发布
- Anthropic Claude API 原生代理支持
- OpenAI 兼容接口
- 自动认证头转换
- Factory Droid system prompt 注入
- 环境变量配置支持
- 健康检查端点 `/health`
- 流式和非流式响应支持

### 📊 性能

- ⚡ 启动时间 < 10ms
- 📉 内存占用 ~11MB
- 📦 二进制大小 ~8MB
- ✅ 支持高并发请求

### 📚 文档

- 完整的 README.md
- OpenAI 兼容模式文档
- 贡献指南
- API 使用示例

---

## 版本说明

### 版本号规则

- **主版本号 (Major)**: 不兼容的 API 修改
- **次版本号 (Minor)**: 向下兼容的功能性新增
- **修订号 (Patch)**: 向下兼容的问题修正

### 更新类型

- `✨ 新增` (Added): 新功能
- `🔄 变更` (Changed): 功能变更
- `🗑️ 移除` (Removed): 移除的功能
- `🐛 修复` (Fixed): Bug 修复
- `🔐 安全` (Security): 安全性修复
- `📚 文档` (Documentation): 文档更新

---

[2.0.1]: https://github.com/yourusername/factory-go-api/compare/v2.0.0...v2.0.1
[2.0.0]: https://github.com/yourusername/factory-go-api/compare/v1.0.0...v2.0.0
[1.0.0]: https://github.com/yourusername/factory-go-api/releases/tag/v1.0.0