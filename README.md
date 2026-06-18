# Sanjie 三界管理系统

玩笑性质的三界移动办公系统，基于 Go + SQLite + uni-app + uview-plus。

[![CI](https://github.com/ijry/sanjie/actions/workflows/ci.yml/badge.svg)](https://github.com/ijry/sanjie/actions/workflows/ci.yml)
[![Docs](https://github.com/ijry/sanjie/actions/workflows/docs.yml/badge.svg)](https://ijry.github.io/sanjie/)
[![License: MIT](https://img.shields.io/badge/License-MIT-red.svg)](LICENSE)

Sanjie 把“天兵管理、投胎管理、十八层地狱、勾魂管理、生死簿、孟婆汤、愿望工单”等神话梗做成一个认真到离谱的移动端管理系统。它定位为开源 demo / 产品原型，不宣称生产可用。

## 在线入口

- 项目文档：<https://ijry.github.io/sanjie/>
- 功能介绍：<https://ijry.github.io/sanjie/guide/features>
- H5 在线演示：<https://ijry.github.io/sanjie/h5-demo/index.html#/pages/index/index>
- 演示说明：<https://ijry.github.io/sanjie/demo/>
- 接口文档：<https://ijry.github.io/sanjie/api/>
- 开发路线：<https://ijry.github.io/sanjie/dev/roadmap>

## 特性

- **移动办公优先**：围绕黑白无常、判官、孟婆、阎王等角色的移动端任务流设计，工作台采用三界总控大屏，支持从指标、流转图和推荐路线直接进入业务闭环。
- **勾魂任务**：支持任务列表、详情、开始、完成、异常上报和撤销，列表与详情采用移动端高密度卡片展示。
- **投胎队列**：支持投胎审核、胎位分配、孟婆汤状态和轮回完成，详情页展示轮回流程和下一步建议。
- **审批中心**：统一处理寿命变更、投胎插队、勾魂异常、地狱刑期复核等高危操作，详情页展示审批路径和风险提示。
- **生死簿查询**：支持寿命记录检索、详情风险研判、冻结和发起寿命变更申请，操作会进入审计链路。
- **十八层地狱看板**：查看楼层容量、设备状态、刑期详情和复核流程，高负载楼层支持分流。
- **愿望工单**：把凡人许愿做成可转派、可办结、可驳回、可建议自力更生的工单。
- **孟婆汤管理**：管理库存、发放记录和记忆残留检测，汤坊库存台支持补货、扣减和低库存雷达提示。
- **审计日志**：敏感操作全量留痕，天眼时间线按审批、命簿、勾魂、地狱、愿望、汤坊分类追踪。
- **多角色演示**：我的页展示仙籍令牌、角色推荐入口和权限模拟，切换身份页提供完整身份牌。
- **文档站**：使用 VitePress 提供功能、接口、部署和开发文档。
- **完整演示数据**：H5 mock 和 Go + SQLite 种子数据均覆盖主要状态，工作台、列表页、详情页、汤坊、审计和身份页都采用差异化视觉，按钮操作可形成完整演示闭环。
- **可复用 H5 演示**：mock 操作会保存在本地存储，刷新页面不丢状态；我的页可一键重置为初始演示数据。

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go + SQLite |
| 移动端 | uni-app + Vue 3 + uview-plus |
| 文档站 | VitePress |
| CI/CD | GitHub Actions + GitHub Pages |
| 数据库 | SQLite 单文件默认运行 |

## 目录结构

```text
sanjie/
├── server/        # Go 后端 API
├── app/           # uni-app 移动端
├── docs-site/     # VitePress 文档站
├── docs/          # 设计规格、实施计划和工程文档
├── assets/        # 截图、展示图和品牌资源
├── scripts/       # 构建和辅助脚本
└── .github/       # CI 与文档发布 workflow
```

## 快速开始

根目录提供统一脚本，便于本地开发和 CI 复用：

```bash
npm run server:test
npm run server:smoke
npm run app:build
npm run docs:build
npm run verify
```

### 1. 后端

```bash
npm run dev:server
```

默认配置参考 `config.example.yaml`，数据库使用本地 `sanjie.db`。后端启动时会读取根目录 `config.yaml`，也支持环境变量覆盖：

```bash
SANJIE_SERVER_ADDR=:8080
SANJIE_DATABASE_DSN=sanjie.db
SANJIE_DEMO_SEED=true
```

### 2. 移动端

```bash
npm run app:install
npm run dev:h5
```

移动端使用 uni-app + uview-plus，第一目标是 H5，本地联调后端默认地址为 `http://localhost:8080`。

不启动后端的 mock 演示：

```bash
npm run dev:h5:mock
```

mock 演示会把当前会话状态保存在浏览器本地存储中，适合连续走查勾魂、审批、发汤、审计等流程；如需恢复初始数据，可在移动端“我的”页点击“重置演示”。

构建可嵌入文档站的 H5 mock 演示：

```bash
npm run app:build:demo
```

### 3. 文档站

```bash
npm run docs:install
npm run dev:docs
```

构建文档：

```bash
npm run docs:build
```

构建文档并生成内嵌 H5 演示：

```bash
npm run docs:build-with-demo
```

完整验证：

```bash
npm run verify
```

### 4. Docker 本地部署

```bash
docker compose up --build
```

默认访问地址：

```text
http://localhost:8088
```

Docker 方式会通过 Nginx 暴露同源 H5 和 `/api/`，SQLite 数据保存在 `sanjie-data` volume。

## 文档

- [快速开始](https://ijry.github.io/sanjie/README)
- [功能介绍](https://ijry.github.io/sanjie/guide/features)
- [架构说明](https://ijry.github.io/sanjie/guide/architecture)
- [接口文档](https://ijry.github.io/sanjie/api/)
- [H5 在线演示](https://ijry.github.io/sanjie/h5-demo/index.html#/pages/index/index)
- [演示说明](https://ijry.github.io/sanjie/demo/)
- [OpenAPI 契约](docs/openapi.yaml)
- [HTTP 调试集合](docs/api.http)
- [部署说明](https://ijry.github.io/sanjie/deploy/)
- [开发路线](https://ijry.github.io/sanjie/dev/roadmap)
- [贡献指南](CONTRIBUTING.md)
- [安全说明](SECURITY.md)
- [变更日志](CHANGELOG.md)

## License

[MIT](LICENSE)
