# 架构说明

Sanjie 使用前后端分离的轻量架构，重点是本地可运行、逻辑可验证、文档可发布。

## 总体结构

```text
sanjie/
├─ app/          # uni-app + uview-plus 移动端
├─ server/       # Go + SQLite API
├─ docs-site/    # VitePress 文档站
├─ docs/         # 规格、计划和工程文档
└─ assets/       # 展示图和静态资源
```

## 移动端

移动端使用 `uni-app + Vue 3 + uview-plus`。页面以移动办公为主，不复刻 PC 后台的大表格体验。

核心 tab 建议为：

```text
工作台 / 勾魂 / 投胎 / 审批 / 我的
```

二级入口包括生死簿、十八层地狱、愿望工单、孟婆汤和审计日志。

## 后端

后端使用 Go 提供 JSON REST API，默认使用 SQLite 单文件数据库。模块边界按业务域拆分：

- dashboard
- users
- souls
- life-book
- capture
- reincarnation
- approvals
- hell
- wishes
- soup
- audit

## API 响应格式

所有接口使用统一响应结构：

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

错误响应：

```json
{
  "code": 40001,
  "message": "invalid status transition",
  "data": null
}
```

## 数据库

SQLite 负责保存演示数据和业务状态。核心表包括：

- `users`
- `souls`
- `life_book_records`
- `capture_tasks`
- `reincarnation_orders`
- `approvals`
- `hell_floors`
- `hell_sentences`
- `wish_tickets`
- `soup_inventory`
- `soup_records`
- `alerts`
- `audit_logs`

## 文档站

`docs-site` 使用 VitePress。GitHub Pages 发布路径为 `/sanjie/`，构建产物不提交。

## Docker 部署

Docker Compose 使用三个服务：`nginx` 对外暴露入口，`server` 运行 Go API，`app_h5` 提供 H5 静态资源。Nginx 将 `/api/` 转发到后端，其余路径转发到 H5，因此容器部署时移动端可以使用同源接口。
