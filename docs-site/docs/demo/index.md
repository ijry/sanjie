# 在线演示

Sanjie 文档站内置一个纯前端 H5 演示版本。它和 Lyshop 一样随 GitHub Pages 一起发布，不需要启动 Go 后端。

<div class="demo-layout">
  <div class="phone-shell">
    <div class="phone-speaker"></div>
    <iframe
      class="phone-frame"
      src="/sanjie/h5-demo/index.html#/pages/index/index"
      title="Sanjie H5 Mock Demo"
      loading="lazy"
    ></iframe>
  </div>
  <div class="demo-panel">
    <p class="demo-kicker">Mock H5</p>
    <h2>三界移动办公演示</h2>
    <p>内置完整 mock 数据，覆盖三界总控工作台、勾魂任务、投胎队列、审批中心、生死簿、十八层地狱、愿望工单、汤坊库存台、天眼审计时间线和多角色身份牌。核心页面都按业务拆成夜巡、票据、案牍、命簿、控制台、签文、汤坊和令牌等差异化视觉，适合在文档站内嵌窗口直接演示。</p>
    <p>
      <a class="demo-button" href="/sanjie/h5-demo/index.html#/pages/index/index" target="_blank" rel="noreferrer">全屏打开 H5 演示</a>
    </p>
  </div>
</div>

## Mock 数据范围

- 多角色：玉帝、阎王、判官、黑无常、孟婆、天兵统领、城隍、审计神仙。
- 工作台：三界总控大屏、流转图、天条预警和推荐演示路线都可跳转到对应模块。
- 勾魂管理：待执行、执行中、异常、已完成、已撤销任务均有多条样例，列表和详情都是夜巡定位任务台。
- 投胎管理：排队、复核、已通过、待喝汤、已完成、已驳回流程均有多条样例，列表和详情都是奈何桥排队票据。
- 生死簿：正常、预警、严重风险和冻结记录均有样例，列表和详情都是命簿卷宗，展示阳寿偏差、处置流程和审计记录。
- 十八层地狱：1 至 18 层容量、负载、设备状态和刑期记录完整展示，楼层、刑期列表和详情都是地狱控制台。
- 愿望工单：求财、求学、姻缘、职场、健康、事业、运气等类型齐全，列表和详情都是许愿签文并支持转派、办结和驳回。
- 其他模块：汤坊库存支持补货/扣减，审计日志按目标类型过滤，多角色身份页展示岗位权限和推荐入口。

mock 状态保存在前端内存中，页面刷新后恢复初始数据。按钮操作会即时修改当前会话数据，并补充审计记录或审批记录，适合演示流程和截图。Go + SQLite 真实后端的种子数据也保持同等密度，方便从 mock 演示切到真实 API 联调。

## 推荐体验路径

1. 从三界总控工作台查看流转图和天条预警，进入勾魂任务。
2. 进入审批中心，处理寿命变更、投胎插队或地狱复核审批。
3. 打开生死簿详情，冻结一条高危记录或发起改寿申请。
4. 打开十八层地狱，查看高负载楼层并进入刑期详情发起复核。
5. 打开愿望工单详情，把愿望转派给对应神仙，再办结或驳回。
6. 切换到孟婆身份进入汤坊库存台，补一批低库存汤剂。
7. 最后进入天眼审计日志，按类型过滤并确认敏感操作已经留痕。

## 本地运行

开发模式：

```bash
npm run app:install
npm run dev:h5:mock
```

构建普通 mock H5：

```bash
npm run app:build:mock
```

构建文档站内嵌演示：

```bash
npm run app:build:demo
npm run docs:build-with-demo
```

`app:build:demo` 会把 mock H5 输出到：

```text
docs-site/docs/public/h5-demo/
```

该目录是构建产物，不提交到 Git。GitHub Pages 发布时会在 workflow 中重新生成。

## 真实后端演示

真实后端演示使用 Go + SQLite 和种子数据：

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

<style>
.demo-layout {
  display: grid;
  grid-template-columns: minmax(280px, 390px) minmax(0, 1fr);
  gap: 32px;
  align-items: center;
  margin: 28px 0 36px;
}

.phone-shell {
  position: relative;
  width: min(100%, 390px);
  height: 780px;
  padding: 42px 16px 18px;
  border: 1px solid rgba(92, 64, 34, 0.2);
  border-radius: 38px;
  background:
    radial-gradient(circle at 50% 0%, rgba(210, 169, 93, 0.28), transparent 32%),
    linear-gradient(145deg, #2d261d, #0f0d0a);
  box-shadow: 0 24px 70px rgba(19, 13, 5, 0.28);
}

.phone-speaker {
  position: absolute;
  top: 20px;
  left: 50%;
  width: 74px;
  height: 6px;
  border-radius: 999px;
  background: rgba(255, 250, 238, 0.22);
  transform: translateX(-50%);
}

.phone-frame {
  width: 100%;
  height: 100%;
  border: 0;
  border-radius: 24px;
  background: #f8f1e4;
}

.demo-panel {
  padding: 30px;
  border: 1px solid rgba(151, 108, 50, 0.22);
  border-radius: 20px;
  background:
    linear-gradient(135deg, rgba(169, 50, 38, 0.08), rgba(214, 166, 80, 0.1)),
    var(--vp-c-bg-soft);
}

.demo-kicker {
  margin: 0 0 8px;
  color: #a93226;
  font-size: 13px;
  font-weight: 800;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.demo-panel h2 {
  margin: 0 0 12px;
  border: 0;
  padding: 0;
}

.demo-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  margin-top: 12px;
  padding: 10px 18px;
  border-radius: 999px;
  color: #fffaf0;
  background: linear-gradient(135deg, #a93226, #6f1b12);
  font-weight: 800;
  text-decoration: none;
}

@media (max-width: 780px) {
  .demo-layout {
    grid-template-columns: 1fr;
  }

  .phone-shell {
    height: 680px;
    margin: 0 auto;
  }
}
</style>
