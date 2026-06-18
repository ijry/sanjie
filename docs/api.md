# Sanjie API

所有接口以 `/api` 为前缀，返回统一结构：

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

本地演示身份通过请求头切换：

```text
X-User-ID: 1
```

辅助文件：

- `docs/api.http`：可执行 HTTP 调试集合，适用于 VS Code REST Client 或 JetBrains HTTP Client。
- `docs/openapi.yaml`：OpenAPI 3.0 基础契约，覆盖当前所有后端路由。

## 接口分组

- `GET /api/health`：健康检查。
- `GET /api/users/me`：当前用户。
- `POST /api/users/switch`：演示身份切换。
- `GET /api/dashboard`：工作台指标。
- `GET /api/alerts`：告警列表。
- `POST /api/alerts/{id}/handle`：处理告警。
- `GET /api/souls`、`GET /api/souls/{id}`：魂魄查询。
- `GET /api/life-book`、`GET /api/life-book/{id}`：生死簿查询。
- `POST /api/life-book/{id}/change-lifespan`：发起寿命变更审批。
- `POST /api/life-book/{id}/freeze`：冻结生死簿记录。
- `GET /api/life-book/{id}/audit-logs`：生死簿审计日志。
- `GET /api/capture-tasks`、`GET /api/capture-tasks/{id}`：勾魂任务。
- `POST /api/capture-tasks/{id}/start`：开始勾魂。
- `POST /api/capture-tasks/{id}/complete`：完成勾魂。
- `POST /api/capture-tasks/{id}/exception`：上报异常。
- `POST /api/capture-tasks/{id}/cancel`：撤销任务。
- `GET /api/reincarnations`、`GET /api/reincarnations/{id}`：投胎队列。
- `POST /api/reincarnations/{id}/approve`：通过投胎。
- `POST /api/reincarnations/{id}/reject`：驳回投胎。
- `POST /api/reincarnations/{id}/review`：转复核。
- `POST /api/reincarnations/{id}/assign-quota`：分配胎位。
- `POST /api/reincarnations/{id}/soup`：发放孟婆汤。
- `POST /api/reincarnations/{id}/complete`：完成投胎。
- `GET /api/approvals`、`GET /api/approvals/{id}`：审批中心。
- `POST /api/approvals/{id}/approve`：同意审批。
- `POST /api/approvals/{id}/reject`：驳回审批。
- `POST /api/approvals/{id}/transfer`：转交审批。
- `GET /api/hell/floors`：地狱楼层。
- `GET /api/hell/sentences`、`GET /api/hell/sentences/{id}`：刑期记录。
- `POST /api/hell/sentences/{id}/review`：发起刑期复核。
- `POST /api/hell/floors/{id}/dispatch`：跨层分流。
- `GET /api/wishes`、`GET /api/wishes/{id}`：愿望工单。
- `POST /api/wishes/{id}/route`：转派愿望。
- `POST /api/wishes/{id}/resolve`：办结愿望。
- `POST /api/wishes/{id}/reject`：驳回愿望。
- `GET /api/soup/inventory`：孟婆汤库存。
- `POST /api/soup/inventory/{id}/adjust`：调整库存。
- `GET /api/soup/records`：孟婆汤发放记录。
- `GET /api/audit-logs`：审计日志。
