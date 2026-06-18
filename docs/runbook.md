# Sanjie Local Runbook

## Root Scripts

```powershell
npm run dev:server
npm run dev:h5
npm run dev:h5:mock
npm run dev:docs
npm run server:smoke
npm run verify
```

`server:smoke` 会用随机本地端口和临时 SQLite 数据库启动后端，验证 `/api/health` 和 `/api/dashboard`。`verify` 会依次执行后端测试、后端构建、后端 HTTP 烟测、移动端 H5 构建和文档站构建。首次运行移动端或文档站前先安装依赖：

```powershell
npm run app:install
npm run docs:install
```

## Backend

```powershell
npm run server:test
npm run server:smoke
npm run dev:server
```

Backend URL:

```text
http://localhost:8080
```

Health check:

```powershell
Invoke-RestMethod http://localhost:8080/api/health
```

Config resolution:

```text
config.yaml
SANJIE_CONFIG
SANJIE_SERVER_ADDR
SANJIE_DATABASE_DSN
SANJIE_DEMO_SEED
```

`config.yaml` can be copied from `config.example.yaml`. Environment variables override file values.

## Mobile H5

```powershell
npm run app:install
npm run dev:h5
```

Mock mode without backend:

```powershell
npm run dev:h5:mock
```

Build:

```powershell
npm run app:build
npm run app:build:mock
```

## Docs Site

```powershell
npm run docs:install
npm run docs:build
```

## Docker Compose

```powershell
docker compose up --build
```

Default URL:

```text
http://localhost:8088
```

Health check through Nginx:

```powershell
Invoke-RestMethod http://localhost:8088/api/health
```

Override host port:

```powershell
$env:SANJIE_HTTP_PORT=8090
docker compose up --build
```
