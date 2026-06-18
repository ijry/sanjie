# 贡献指南

Sanjie 是玩笑性质的开源 demo，但代码、接口和文档按可验证项目维护。

## 环境要求

- Go 1.25+
- Node.js 20+
- npm

首次安装依赖：

```bash
npm run app:install
npm run docs:install
```

## 常用命令

```bash
npm run dev:server
npm run dev:h5
npm run dev:docs
npm run server:smoke
npm run verify
```

`npm run server:smoke` 会启动临时后端并验证基础接口。`npm run verify` 会执行后端测试、后端构建、后端 HTTP 烟测、H5 构建和文档站构建。

## 提交约定

- 仓库使用 `.editorconfig` 固定基础缩进、换行和编码规则。
- 后端代码提交前运行 `gofmt`。
- 移动端优先使用 uview-plus 组件。
- 功能变更同步 README 或 docs-site。
- 接口变更同步 `docs/api.md` 和 `docs-site/docs/api/index.md`。
- 不提交运行期数据库、日志、依赖目录、构建产物和本地配置。

## PR 内容

建议每个 PR 说明：

- 变更目的。
- 影响模块。
- 验证命令和结果。
- 是否涉及接口、配置或部署变化。

根目录还提供：

- `CONTRIBUTING.md`：协作细则。
- `SECURITY.md`：安全边界和问题报告方式。
- `CHANGELOG.md`：版本变更记录。
