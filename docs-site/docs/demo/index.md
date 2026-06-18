# 演示

Sanjie 支持两种演示方式：真实后端演示和纯前端 mock 演示。

## Mock H5 演示

mock 模式不需要启动 Go 后端，适合文档站、截图和快速产品走查。

```bash
npm run app:install
npm run dev:h5:mock
```

构建 mock H5：

```bash
npm run app:build:mock
```

mock 模式通过 `VITE_SANJIE_MOCK=true` 启用，数据来自 `app/mock/`。页面刷新后会恢复初始 mock 数据。

## 真实后端演示

真实后端演示会使用 Go + SQLite 和种子数据：

```bash
npm run dev:server
npm run dev:h5
```

也可以使用 Docker Compose：

```bash
docker compose up --build
```

默认访问：

```text
http://localhost:8088
```

## 适用场景

- mock 演示：无需后端、适合在线静态预览和 UI 流程展示。
- 真实后端演示：适合验证 API、SQLite 状态流转和审计日志。
