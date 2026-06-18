# Contributing to Sanjie

感谢你愿意参与 Sanjie。这个项目是玩笑性质的三界移动办公系统，但工程约束按认真项目处理。

## 开发环境

需要：

- Go 1.25+
- Node.js 20+
- npm

首次准备：

```bash
npm run app:install
npm run docs:install
```

常用命令：

```bash
npm run dev:server
npm run dev:h5
npm run dev:docs
npm run server:smoke
npm run verify
```

## 提交前检查

提交 PR 前至少运行：

```bash
npm run verify
```

该命令会执行：

- 后端测试：`go test ./...`
- 后端构建：`go build ./...`
- 后端烟测：临时启动 API 并验证 `/api/health`、`/api/dashboard`
- 移动端 H5 构建：`npm run build:h5`
- 文档站构建：`npm run docs:build`

## 代码约定

- 仓库使用 `.editorconfig` 固定基础缩进、换行和编码规则。
- 后端使用 Go 标准格式化，提交前运行 `gofmt`。
- 移动端优先使用 uview-plus 组件。
- 新增或修改接口时，同步更新 `docs/api.md` 和 `docs-site/docs/api/index.md`。
- 新增或修改用户可见功能时，同步更新 README 或 docs-site 对应页面。
- 不提交 `node_modules`、构建产物、SQLite 运行库、日志文件和本地配置。

## 业务设定

Sanjie 可以保持幽默，但状态流转、接口语义和审计记录要可验证。新增梗时优先复用现有模块：

- 勾魂：`capture`
- 投胎：`reincarnation`
- 生死簿：`life-book`
- 十八层地狱：`hell`
- 愿望工单：`wishes`
- 孟婆汤：`soup`
- 审计：`audit`

## PR 建议

- 一个 PR 聚焦一个主题。
- 附上验证命令和结果。
- 如果修改部署、配置或 API，说明兼容影响。
- 如果目录已经初始化为 git 仓库，提交信息使用中文，格式参考 `AGENTS.md`。
