# 部署说明

当前阶段支持本地开发、Docker Compose 单机部署和 GitHub Pages 文档发布。

## 本地后端

```bash
npm run dev:server
```

默认监听：

```text
http://localhost:8080
```

默认数据库：

```text
sanjie.db
```

配置读取顺序：

```text
默认值 -> config.yaml -> SANJIE_* 环境变量
```

支持的环境变量：

```text
SANJIE_CONFIG
SANJIE_SERVER_ADDR
SANJIE_DATABASE_DSN
SANJIE_DEMO_SEED
```

后端 HTTP 烟测：

```bash
npm run server:smoke
```

该命令会使用随机本地端口和临时 SQLite 数据库启动后端，并验证 `/api/health` 与 `/api/dashboard`。

## 本地移动端

```bash
npm run app:install
npm run dev:h5
```

移动端本地调试时连接 `http://localhost:8080`。

H5 构建时可通过 `VITE_SANJIE_API_BASE_URL` 指定接口地址。Docker 构建默认设置为空字符串，让 H5 走同源 `/api/`。

不启动后端的 mock H5：

```bash
npm run dev:h5:mock
npm run app:build:mock
```

mock 模式通过 `VITE_SANJIE_MOCK=true` 启用，数据来自 `app/mock/`。

## Docker Compose

```bash
docker compose up --build
```

默认访问地址：

```text
http://localhost:8088
```

服务组成：

- `nginx`：对外暴露 H5 和 `/api/`。
- `server`：Go + SQLite API。
- `app_h5`：uni-app H5 静态站点。

SQLite 数据保存在 Docker volume `sanjie-data`。如需改端口：

```bash
SANJIE_HTTP_PORT=8090 docker compose up --build
```

## 文档站

```bash
npm run docs:install
npm run docs:build
```

GitHub Pages workflow 会发布：

```text
docs-site/docs/.vitepress/dist
```

发布路径：

```text
https://ijry.github.io/sanjie/
```

## 后续部署方向

- 后端可打包为单个 Go 二进制。
- SQLite 适合 demo、单机版和轻量部署。
- 移动端 H5 可作为静态资源发布。
- App 和小程序端按 uni-app 标准流程构建。

## CI 验证

GitHub Actions 的 `CI` workflow 会分别执行：

- `server` 目录下的 `go test ./...` 和 `go build ./...`。
- 根目录 `npm run server:smoke` 的后端 HTTP 烟测。
- `app` 目录下的 `npm install --legacy-peer-deps` 和 `npm run build:h5`。
- `app` 目录下的 `npm run build:h5:mock`。
- `docs-site` 目录下的 `npm install` 和 `npm run docs:build`。
