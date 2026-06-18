# 快速开始

Sanjie 采用轻量本地开发方式：Go 后端运行 API，SQLite 保存演示数据，uni-app 移动端调用本地接口，VitePress 负责文档站。

## 后端

```bash
npm run dev:server
```

默认配置参考根目录 `config.example.yaml`：

```yaml
server:
  addr: ":8080"

database:
  driver: "sqlite"
  dsn: "sanjie.db"

demo:
  seed: true
```

后端会自动读取根目录 `config.yaml`。也可以通过环境变量覆盖：

```bash
SANJIE_CONFIG=config.yaml
SANJIE_SERVER_ADDR=:8080
SANJIE_DATABASE_DSN=sanjie.db
SANJIE_DEMO_SEED=true
```

## 移动端

```bash
npm run app:install
npm run dev:h5
```

移动端默认对接 `http://localhost:8080`。接口封装位于 `app/api/`，页面使用 uview-plus 组件优先实现。

不启动后端的 mock 演示：

```bash
npm run dev:h5:mock
```

构建 mock H5：

```bash
npm run app:build:mock
```

## 文档站

```bash
npm run docs:install
npm run dev:docs
```

构建文档：

```bash
npm run docs:build
```

构建产物位于 `docs-site/docs/.vitepress/dist`，该目录不提交到仓库。

## 完整验证

```bash
npm run verify
```

该命令依次执行后端测试、后端构建、后端 HTTP 烟测、移动端 H5 构建、移动端 H5 mock 构建和文档站构建。

只验证后端真实启动和基础接口：

```bash
npm run server:smoke
```

该命令会使用随机本地端口和临时 SQLite 数据库，不会污染工作区数据库。

## Docker 本地部署

```bash
docker compose up --build
```

默认访问：

```text
http://localhost:8088
```

Docker 方式使用 Nginx 同源代理：`/api/` 转发到 Go 后端，其余路径转发到 H5 静态站点。
