---
layout: home

hero:
  name: Sanjie
  text: 三界管理系统
  tagline: 一个认真到离谱的三界移动办公系统
  actions:
    - theme: brand
      text: 快速开始
      link: /README
    - theme: alt
      text: 功能介绍
      link: /guide/features
    - theme: alt
      text: Mock 演示
      link: /demo/
    - theme: alt
      text: 贡献指南
      link: /dev/contributing

features:
  - title: 勾魂任务
    details: 面向黑白无常和鬼差的移动执行流，支持开始、完成、异常上报和撤销。
  - title: 投胎队列
    details: 管理魂魄排队、胎位分配、孟婆汤状态和轮回完成。
  - title: 审批中心
    details: 统一处理寿命变更、投胎插队、地狱刑期复核等高危流程。
  - title: 生死簿
    details: 查询寿命记录、冻结高危记录，并发起可审计的寿命变更申请。
  - title: 十八层地狱
    details: 查看楼层容量、设备状态、刑期记录和复核状态。
  - title: 文档站
    details: 使用 VitePress 发布功能、接口、部署和开发路线文档，并提供 OpenAPI 与 HTTP 调试集合。
  - title: Mock 演示
    details: H5 支持纯前端 mock 模式，不启动后端也能走查核心业务页面。
---

## 项目定位

Sanjie 是一个开源 demo / 产品原型，用移动办公系统的方式包装神话世界里的荒诞业务。前端计划使用 `uni-app + uview-plus`，后端计划使用 `Go + SQLite`，适合本地运行、演示和二次开发。

## 当前状态

当前版本已完成 Go + SQLite API、uni-app H5 页面、H5 mock 演示、VitePress 文档站、Docker Compose 本地部署示例和开源协作模板。完整开发计划位于 `docs/superpowers/plans/2026-06-18-sanjie-mobile-full-implementation.md`。
