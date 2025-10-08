# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- OpenAI 兼容模式支持
- 自动格式转换（OpenAI ↔ Anthropic）
- 完整的 system 消息处理
- 详细的请求/响应日志
- 健康检查端点
- GitHub Actions CI/CD 工作流
- 完整的文档和贡献指南

## [1.0.0] - 2025-10-08

### Added
- 🎉 首次发布
- Anthropic API 原生代理支持
- OpenAI 兼容接口
- Bedrock API 代理支持
- 自动认证头转换
- Factory Droid system prompt 注入
- 环境变量配置支持
- 启动脚本和示例代码
- MIT License

### Features
- ⚡ 极致性能（< 10ms 启动，~11MB 内存）
- 🔄 智能格式转换
- 🔐 多种认证方式支持
- 📊 完整的日志记录
- 🏥 健康检查端点

### Documentation
- 完整的 README.md
- OpenAI 兼容模式文档
- 贡献指南
- API 使用示例
- 部署指南

---

## 版本说明

### 版本号规则

- **主版本号 (Major)**: 不兼容的 API 修改
- **次版本号 (Minor)**: 向下兼容的功能性新增
- **修订号 (Patch)**: 向下兼容的问题修正

### 更新日志分类

- `Added`: 新增功能
- `Changed`: 功能变更
- `Deprecated`: 即将废弃的功能
- `Removed`: 已移除的功能
- `Fixed`: Bug 修复
- `Security`: 安全性修复

---

[Unreleased]: https://github.com/yourusername/factory-proxy/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/yourusername/factory-proxy/releases/tag/v1.0.0