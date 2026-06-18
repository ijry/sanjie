<script setup>
import { computed } from 'vue'
import { useData } from 'vitepress'

const { site } = useData()
const base = computed(() => site.value.base || '/')

const metrics = [
  { value: '10+', label: '三界业务模块' },
  { value: '18', label: '地狱楼层看板' },
  { value: 'H5', label: '内嵌 mock 演示' },
  { value: 'MIT', label: '开源协议' }
]

const modules = [
  { title: '勾魂调度', desc: '黑白无常外勤任务、异常上报、撤销和完成闭环。' },
  { title: '投胎管理', desc: '排队、复核、胎位分配、孟婆汤和轮回完成。' },
  { title: '生死簿', desc: '寿命查询、冻结、改寿申请和审计追踪。' },
  { title: '十八层地狱', desc: '楼层容量、设备状态、刑期记录和跨层分流。' },
  { title: '愿望工单', desc: '凡人许愿进入工单系统，支持转派、办结和驳回。' },
  { title: '审计与告警', desc: '高危操作留痕，重点盯防关系户投胎和批量划名。' }
]

const stack = ['Go', 'SQLite', 'uni-app', 'Vue 3', 'uview-plus', 'VitePress', 'Docker', 'GitHub Pages']

const flow = [
  { name: '魂魄入库', desc: '勾魂任务、生死簿记录和风险标记进入统一视图。' },
  { name: '判官审批', desc: '寿命、刑期、投胎插队等高危动作走审批链。' },
  { name: '轮回发车', desc: '胎位分配、孟婆汤发放和投胎完成形成闭环。' }
]

const links = [
  { title: 'H5 在线演示', desc: '打开完整 mock 演示', href: 'h5-demo/index.html#/pages/index/index' },
  { title: '快速开始', desc: '本地运行后端、H5 和文档站', href: 'README' },
  { title: '接口文档', desc: '查看 Go API 与 OpenAPI 契约', href: 'api/' },
  { title: '部署说明', desc: 'Docker Compose 与 Pages 发布', href: 'deploy/' }
]
</script>

<template>
  <main class="sanjie-home">
    <section class="sj-hero">
      <div class="sj-hero__glow sj-hero__glow--red"></div>
      <div class="sj-hero__glow sj-hero__glow--gold"></div>
      <div class="sj-container sj-hero__inner">
        <div class="sj-hero__copy">
          <p class="sj-eyebrow">Open Source Myth Office</p>
          <h1>Sanjie 三界管理系统</h1>
          <p class="sj-lead">
            把天兵巡逻、勾魂外勤、生死簿审批、十八层地狱负载和投胎队列做成一个认真到离谱的移动办公系统。
          </p>
          <div class="sj-actions">
            <a class="sj-button sj-button--primary" :href="`${base}h5-demo/index.html#/pages/index/index`" target="_blank" rel="noreferrer">打开 H5 演示</a>
            <a class="sj-button sj-button--ghost" :href="`${base}README`">快速开始</a>
            <a class="sj-button sj-button--ghost" href="https://github.com/ijry/sanjie" target="_blank" rel="noreferrer">GitHub</a>
          </div>
          <div class="sj-metrics" aria-label="项目指标">
            <div v-for="item in metrics" :key="item.label" class="sj-metric">
              <strong>{{ item.value }}</strong>
              <span>{{ item.label }}</span>
            </div>
          </div>
        </div>

        <div class="sj-hero__visual" aria-label="三界演示预览">
          <div class="sj-tablet">
            <div class="sj-tablet__bar">
              <span></span><span></span><span></span>
              <em>三界总览 / Dispatch Console</em>
            </div>
            <div class="sj-dashboard">
              <div class="sj-dashboard__header">
                <div>
                  <b>今日三界公务</b>
                  <small>按天条流程处理</small>
                </div>
                <span>实时</span>
              </div>
              <div class="sj-dashboard__grid">
                <div><strong>2</strong><span>待勾魂</span></div>
                <div><strong>3</strong><span>待审批</span></div>
                <div><strong>4</strong><span>投胎积压</span></div>
                <div><strong>3</strong><span>高危寿命</span></div>
              </div>
              <div class="sj-scroll">
                <p><b>生死簿批量划名风险</b><span>严重</span></p>
                <p><b>奈何桥排队拥堵</b><span>预警</span></p>
                <p><b>拔舌地狱容量过高</b><span>严重</span></p>
              </div>
            </div>
          </div>

          <div class="sj-phone">
            <iframe :src="`${base}h5-demo/index.html#/pages/index/index`" title="Sanjie H5 mock demo" loading="lazy"></iframe>
          </div>
        </div>
      </div>
    </section>

    <section class="sj-section sj-section--paper">
      <div class="sj-container sj-split">
        <div>
          <p class="sj-kicker">What It Solves</p>
          <h2>神话梗外壳，真实管理系统骨架</h2>
        </div>
        <p>
          Sanjie 是一个开源 demo / 产品原型。它不追求生产级严肃业务，却把移动端任务流、审批流、审计日志、接口契约和部署链路完整串起来，适合学习、展示和二次开发。
        </p>
      </div>
    </section>

    <section class="sj-section">
      <div class="sj-container">
        <div class="sj-heading">
          <p class="sj-kicker">Modules</p>
          <h2>三界模块一屏看懂</h2>
          <p>每个模块都能在 H5 mock 演示里直接走查，刷新后恢复初始数据。</p>
        </div>
        <div class="sj-module-grid">
          <article v-for="item in modules" :key="item.title" class="sj-card">
            <h3>{{ item.title }}</h3>
            <p>{{ item.desc }}</p>
          </article>
        </div>
      </div>
    </section>

    <section class="sj-section sj-section--dark">
      <div class="sj-container sj-architecture">
        <div class="sj-architecture__copy">
          <p class="sj-kicker sj-kicker--dark">Architecture</p>
          <h2>Go + SQLite 后端，uni-app 移动端，VitePress 文档站</h2>
          <p>
            后端提供可验证 API 和 SQLite 种子数据；移动端优先 H5，使用 uview-plus 组件；文档站发布在线说明和内嵌 H5 mock 演示。
          </p>
          <div class="sj-stack">
            <span v-for="item in stack" :key="item">{{ item }}</span>
          </div>
        </div>
        <div class="sj-map" aria-label="系统架构示意">
          <div class="sj-map__row">
            <span>uni-app H5</span>
            <span>VitePress Docs</span>
            <span>Docker Nginx</span>
          </div>
          <div class="sj-map__core">
            <strong>Go API</strong>
            <div>
              <span>capture</span><span>life-book</span><span>approval</span><span>hell</span><span>reincarnation</span><span>audit</span>
            </div>
          </div>
          <div class="sj-map__row">
            <span>SQLite</span>
            <span>Mock State</span>
            <span>OpenAPI</span>
          </div>
        </div>
      </div>
    </section>

    <section class="sj-section">
      <div class="sj-container">
        <div class="sj-heading">
          <p class="sj-kicker">Workflow</p>
          <h2>从勾魂到轮回，流程闭环</h2>
        </div>
        <div class="sj-flow">
          <article v-for="(item, index) in flow" :key="item.name" class="sj-flow__item">
            <span>{{ index + 1 }}</span>
            <h3>{{ item.name }}</h3>
            <p>{{ item.desc }}</p>
          </article>
        </div>
      </div>
    </section>

    <section class="sj-section sj-section--paper">
      <div class="sj-container">
        <div class="sj-heading">
          <p class="sj-kicker">Start Here</p>
          <h2>常用入口</h2>
        </div>
        <div class="sj-link-grid">
          <a v-for="item in links" :key="item.title" class="sj-link-card" :href="`${base}${item.href}`">
            <strong>{{ item.title }}</strong>
            <span>{{ item.desc }}</span>
          </a>
        </div>
      </div>
    </section>

    <section class="sj-final">
      <div class="sj-container sj-final__inner">
        <div>
          <p class="sj-kicker sj-kicker--dark">Ready</p>
          <h2>直接打开演示，或者从源码开始二开</h2>
          <p>文档、H5 mock 演示、Go 后端、SQLite 数据和 CI 发布链路已经打通。</p>
        </div>
        <div class="sj-actions">
          <a class="sj-button sj-button--light" :href="`${base}h5-demo/index.html#/pages/index/index`" target="_blank" rel="noreferrer">查看演示</a>
          <a class="sj-button sj-button--ghost-light" :href="`${base}deploy/`">部署说明</a>
        </div>
      </div>
    </section>
  </main>
</template>
