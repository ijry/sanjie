# 接口文档

Sanjie 后端接口以 `/api` 为前缀，返回 JSON。

## 响应结构

成功：

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

失败：

```json
{
  "code": 40001,
  "message": "invalid request",
  "data": null
}
```

## 演示身份

本地 demo 使用请求头切换身份：

```text
X-User-ID: 1
```

默认用户：

| ID | 角色 |
|----|------|
| 1 | 玉帝 |
| 2 | 阎王 |
| 3 | 判官 |
| 4 | 黑无常 |
| 5 | 孟婆 |

## Mock 演示接口

以下接口只在 H5 mock 模式下使用，用于文档站内嵌演示的本地会话管理；真实 Go + SQLite 后端不需要这些接口。

```text
GET  /api/mock/session
POST /api/mock/reset
```

`GET /api/mock/session` 返回当前本地 mock 会话的开始时间、审计日志数量、审批单数量、发汤记录数量和已处理预警数量。`POST /api/mock/reset` 会恢复前端 mock 初始数据，并切回玉帝身份。

## 调试与契约

仓库内提供两份开发辅助文件：

- `docs/api.http`：可执行 HTTP 调试集合，可在 VS Code REST Client 或 JetBrains HTTP Client 中直接调用本地接口。
- `docs/openapi.yaml`：OpenAPI 3.0 基础契约，覆盖当前所有后端路由、路径参数和常用请求体。

默认本地接口地址：

```text
http://localhost:8080
```

## 接口分组

### 健康检查

```text
GET /api/health
```

### 用户

```text
GET  /api/users/me
POST /api/users/switch
```

### 工作台与告警

```text
GET  /api/dashboard
GET  /api/alerts
POST /api/alerts/{id}/handle
```

### 魂魄与生死簿

```text
GET  /api/souls
GET  /api/souls/{id}
GET  /api/life-book
GET  /api/life-book/{id}
POST /api/life-book/{id}/change-lifespan
POST /api/life-book/{id}/freeze
GET  /api/life-book/{id}/audit-logs
```

### 勾魂任务

```text
GET  /api/capture-tasks
GET  /api/capture-tasks/{id}
POST /api/capture-tasks/{id}/start
POST /api/capture-tasks/{id}/complete
POST /api/capture-tasks/{id}/exception
POST /api/capture-tasks/{id}/cancel
```

### 投胎队列

```text
GET  /api/reincarnations
GET  /api/reincarnations/{id}
POST /api/reincarnations/{id}/approve
POST /api/reincarnations/{id}/reject
POST /api/reincarnations/{id}/review
POST /api/reincarnations/{id}/assign-quota
POST /api/reincarnations/{id}/soup
POST /api/reincarnations/{id}/complete
```

### 审批中心

```text
GET  /api/approvals
GET  /api/approvals/{id}
POST /api/approvals/{id}/approve
POST /api/approvals/{id}/reject
POST /api/approvals/{id}/transfer
```

### 十八层地狱

```text
GET  /api/hell/floors
GET  /api/hell/sentences
GET  /api/hell/sentences/{id}
POST /api/hell/sentences/{id}/review
POST /api/hell/floors/{id}/dispatch
```

### 愿望工单

```text
GET  /api/wishes
GET  /api/wishes/{id}
POST /api/wishes/{id}/route
POST /api/wishes/{id}/resolve
POST /api/wishes/{id}/reject
```

### 孟婆汤

```text
GET  /api/soup/inventory
POST /api/soup/inventory/{id}/adjust
GET  /api/soup/records
```

### 审计日志

```text
GET /api/audit-logs
```
