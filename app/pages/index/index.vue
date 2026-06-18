<template>
  <view class="page stack dense-page command-page">
    <view class="command-hero">
      <view class="command-orbit">
        <view class="command-orbit__ring"></view>
        <view class="command-orbit__core">
          <text>{{ metrics.unhandledAlertCount || 0 }}</text>
          <text>警</text>
        </view>
      </view>
      <view class="command-hero__main">
        <text class="command-hero__eyebrow">SANJIE WAR ROOM</text>
        <text class="command-hero__title">三界总控大屏</text>
        <text class="command-hero__desc">把勾魂、投胎、审批、命簿、地狱和汤坊压到一张移动端值班图里。</text>
      </view>
      <view class="command-hero__stats">
        <view v-for="item in commandStats" :key="item.label" class="command-stat" @click="go(item.url)">
          <text class="command-stat__value">{{ item.value }}</text>
          <text class="command-stat__label">{{ item.label }}</text>
        </view>
      </view>
    </view>

    <view class="realm-map">
      <view class="realm-map__head">
        <view>
          <text class="realm-map__kicker">REALM ROUTER</text>
          <text class="realm-map__title">今日流转图</text>
        </view>
        <text class="realm-map__hint">点击进入模块</text>
      </view>
      <view class="realm-lanes">
        <view
          v-for="lane in realmLanes"
          :key="lane.name"
          class="realm-lane"
          :class="`realm-lane--${lane.tone}`"
          @click="go(lane.url)"
        >
          <view class="realm-lane__signal">
            <text>{{ lane.signal }}</text>
          </view>
          <view class="realm-lane__main">
            <text class="realm-lane__name">{{ lane.name }}</text>
            <text class="realm-lane__desc">{{ lane.desc }}</text>
          </view>
          <text class="realm-lane__value">{{ lane.value }}</text>
        </view>
      </view>
    </view>

    <view class="demo-route">
      <view class="demo-route__title">
        <text>推荐演示路线</text>
        <text>从告警到审计闭环</text>
      </view>
      <view class="demo-route__steps">
        <view v-for="step in demoSteps" :key="step.title" class="demo-step" @click="go(step.url)">
          <text class="demo-step__index">{{ step.index }}</text>
          <view>
            <text class="demo-step__title">{{ step.title }}</text>
            <text class="demo-step__desc">{{ step.desc }}</text>
          </view>
        </view>
      </view>
    </view>

    <view class="alert-board">
      <view class="alert-board__head">
        <view>
          <text class="alert-board__kicker">HEAVENLY ALERTS</text>
          <text class="alert-board__title">天条预警</text>
        </view>
        <text class="alert-board__count">{{ unhandledCount }} 未处理</text>
      </view>

      <EmptyState v-if="alerts.length === 0" text="暂无预警" />
      <view v-else class="alert-stack">
        <view
          v-for="item in alerts"
          :key="item.id"
          class="alert-card"
          :class="{ 'alert-card--handled': item.handled }"
          @click="markAlert(item)"
        >
          <view class="alert-card__head">
            <text class="alert-card__title">{{ item.title }}</text>
            <StatusTag :value="item.level" />
          </view>
          <text class="alert-card__content">{{ item.content }}</text>
          <view class="alert-card__foot">
            <text>{{ item.handled ? '已归档' : '点击标记已处理' }}</text>
            <text>#{{ item.id }}</text>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import StatusTag from '../../components/StatusTag.vue'
import EmptyState from '../../components/EmptyState.vue'
import { getAlerts, getDashboard, handleAlert } from '../../api/dashboard'
import { showError, showSuccess } from '../../utils/format'

const metrics = ref({})
const alerts = ref([])

const commandStats = computed(() => [
  { label: '待勾魂', value: metrics.value.pendingCaptureCount || 0, url: '/pages/capture/index' },
  { label: '待审批', value: metrics.value.pendingApprovalCount || 0, url: '/pages/approval/index' },
  { label: '投胎积压', value: metrics.value.reincarnationQueueCount || 0, url: '/pages/reincarnation/index' },
  { label: '高危命簿', value: metrics.value.highRiskLifeBookCount || 0, url: '/pages/life-book/search' }
])

const unhandledCount = computed(() => alerts.value.filter(item => !item.handled).length)

const realmLanes = computed(() => [
  { name: '地府外勤', desc: '勾魂夜巡任务池', signal: '魂', value: metrics.value.pendingCaptureCount || 0, tone: 'blue', url: '/pages/capture/index' },
  { name: '奈何桥', desc: '投胎排队与胎位', signal: '桥', value: metrics.value.reincarnationQueueCount || 0, tone: 'teal', url: '/pages/reincarnation/index' },
  { name: '判官台', desc: '改寿插队审批', signal: '审', value: metrics.value.pendingApprovalCount || 0, tone: 'amber', url: '/pages/approval/index' },
  { name: '命簿阁', desc: '生死簿风险索引', signal: '命', value: metrics.value.highRiskLifeBookCount || 0, tone: 'paper', url: '/pages/life-book/search' },
  { name: '十八层', desc: '刑期和楼层负载', signal: '狱', value: '18', tone: 'red', url: '/pages/hell/index' },
  { name: '愿力司', desc: '凡人许愿工单', signal: '愿', value: '票', tone: 'gold', url: '/pages/wish/index' },
  { name: '汤坊', desc: '孟婆汤库存发放', signal: '汤', value: '炉', tone: 'violet', url: '/pages/soup/index' },
  { name: '天眼', desc: '敏感操作审计', signal: '眼', value: unhandledCount.value, tone: 'ink', url: '/pages/audit/index' }
])

const demoSteps = [
  { index: '01', title: '处理一条勾魂任务', desc: '进入夜巡任务详情，开始、完成或上报异常。', url: '/pages/capture/index' },
  { index: '02', title: '走一次判官审批', desc: '审批改寿、插队或刑期复核，观察状态变化。', url: '/pages/approval/index' },
  { index: '03', title: '冻结高危命簿', desc: '从生死簿发起改寿申请，确认必须留痕。', url: '/pages/life-book/search' },
  { index: '04', title: '回到天眼审计', desc: '检查刚才操作写入的审计记录。', url: '/pages/audit/index' }
]

function go(url) {
  if (['/pages/index/index', '/pages/capture/index', '/pages/reincarnation/index', '/pages/approval/index', '/pages/mine/index'].includes(url)) {
    uni.switchTab({ url })
    return
  }
  uni.navigateTo({ url })
}

async function load() {
  try {
    metrics.value = await getDashboard()
    alerts.value = await getAlerts()
  } catch (error) {
    showError(error)
  }
}

async function markAlert(item) {
  if (item.handled) return
  try {
    await handleAlert(item.id)
    showSuccess('已处理')
    await load()
  } catch (error) {
    showError(error)
  }
}

onShow(load)
</script>

<style scoped>
.command-page {
  background:
    radial-gradient(circle at 12% 0%, rgba(217, 119, 6, 0.22), transparent 26%),
    radial-gradient(circle at 92% 10%, rgba(14, 165, 233, 0.16), transparent 24%),
    linear-gradient(180deg, #18130d 0%, #2f1f12 30%, #f7f1e3 30%, #f7f1e3 100%);
}

.command-hero {
  display: grid;
  grid-template-columns: 148rpx minmax(0, 1fr);
  gap: 18rpx;
  padding: 26rpx;
  color: #fffaf0;
  border: 1rpx solid rgba(245, 158, 11, 0.28);
  border-radius: 30rpx;
  background:
    linear-gradient(90deg, rgba(245, 158, 11, 0.08) 1rpx, transparent 1rpx),
    linear-gradient(0deg, rgba(245, 158, 11, 0.08) 1rpx, transparent 1rpx),
    linear-gradient(135deg, #111827, #3d1f12 58%, #7f1d1d);
  background-size: 36rpx 36rpx;
  box-shadow: 0 20rpx 46rpx rgba(49, 28, 13, 0.26);
}

.command-orbit {
  position: relative;
  width: 138rpx;
  height: 138rpx;
}

.command-orbit__ring {
  position: absolute;
  inset: 0;
  border: 2rpx dashed rgba(251, 191, 36, 0.46);
  border-radius: 50%;
}

.command-orbit__ring::before,
.command-orbit__ring::after {
  content: "";
  position: absolute;
  inset: 22rpx;
  border: 1rpx solid rgba(125, 211, 252, 0.28);
  border-radius: 50%;
}

.command-orbit__ring::after {
  inset: 46rpx;
}

.command-orbit__core {
  position: absolute;
  inset: 26rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: linear-gradient(135deg, #f59e0b, #a93226);
}

.command-orbit__core text:first-child {
  font-size: 36rpx;
  font-weight: 900;
  line-height: 1;
}

.command-orbit__core text:last-child {
  margin-top: 4rpx;
  font-size: 20rpx;
  font-weight: 900;
}

.command-hero__main {
  min-width: 0;
}

.command-hero__eyebrow,
.command-hero__title,
.command-hero__desc {
  display: block;
}

.command-hero__eyebrow {
  color: #fbbf24;
  font-size: 20rpx;
  font-weight: 900;
  letter-spacing: 2rpx;
}

.command-hero__title {
  margin-top: 6rpx;
  font-size: 46rpx;
  font-weight: 900;
  line-height: 1.12;
}

.command-hero__desc {
  margin-top: 10rpx;
  color: rgba(255, 250, 240, 0.72);
  font-size: 24rpx;
  line-height: 1.45;
}

.command-hero__stats {
  grid-column: 1 / -1;
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 10rpx;
}

.command-stat {
  min-width: 0;
  padding: 14rpx 10rpx;
  border: 1rpx solid rgba(255, 250, 240, 0.14);
  border-radius: 16rpx;
  background: rgba(255, 250, 240, 0.08);
}

.command-stat__value,
.command-stat__label {
  display: block;
  text-align: center;
}

.command-stat__value {
  color: #fef3c7;
  font-size: 32rpx;
  font-weight: 900;
}

.command-stat__label {
  margin-top: 6rpx;
  color: rgba(255, 250, 240, 0.66);
  font-size: 20rpx;
}

.realm-map,
.demo-route,
.alert-board {
  padding: 22rpx;
  border: 1rpx solid rgba(120, 53, 15, 0.16);
  border-radius: 24rpx;
  background: rgba(255, 250, 240, 0.92);
  box-shadow: 0 12rpx 30rpx rgba(66, 41, 20, 0.08);
}

.realm-map__head,
.alert-board__head,
.demo-route__title {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16rpx;
}

.realm-map__kicker,
.alert-board__kicker {
  display: block;
  color: #a93226;
  font-size: 19rpx;
  font-weight: 900;
  letter-spacing: 2rpx;
}

.realm-map__title,
.alert-board__title,
.demo-route__title text:first-child {
  display: block;
  margin-top: 4rpx;
  color: #25211b;
  font-size: 32rpx;
  font-weight: 900;
}

.realm-map__hint,
.alert-board__count,
.demo-route__title text:last-child {
  color: #8c806d;
  font-size: 22rpx;
  font-weight: 800;
}

.realm-lanes {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12rpx;
  margin-top: 18rpx;
}

.realm-lane {
  display: grid;
  grid-template-columns: 54rpx minmax(0, 1fr) auto;
  gap: 12rpx;
  align-items: center;
  min-width: 0;
  padding: 16rpx;
  border-radius: 18rpx;
  background: #fff7ed;
}

.realm-lane__signal {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 54rpx;
  height: 54rpx;
  color: #fffaf0;
  border-radius: 16rpx;
  font-size: 24rpx;
  font-weight: 900;
}

.realm-lane__main {
  min-width: 0;
}

.realm-lane__name,
.realm-lane__desc,
.realm-lane__value {
  display: block;
}

.realm-lane__name {
  overflow: hidden;
  color: #25211b;
  font-size: 26rpx;
  font-weight: 900;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.realm-lane__desc {
  overflow: hidden;
  margin-top: 4rpx;
  color: #756f63;
  font-size: 20rpx;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.realm-lane__value {
  color: #25211b;
  font-size: 28rpx;
  font-weight: 900;
}

.realm-lane--blue .realm-lane__signal { background: #0284c7; }
.realm-lane--teal .realm-lane__signal { background: #0f766e; }
.realm-lane--amber .realm-lane__signal { background: #b45309; }
.realm-lane--paper .realm-lane__signal { background: #92400e; }
.realm-lane--red .realm-lane__signal { background: #991b1b; }
.realm-lane--gold .realm-lane__signal { background: #d97706; }
.realm-lane--violet .realm-lane__signal { background: #7c3aed; }
.realm-lane--ink .realm-lane__signal { background: #1f2937; }

.demo-route {
  background:
    linear-gradient(90deg, rgba(169, 50, 38, 0.08) 1rpx, transparent 1rpx),
    linear-gradient(135deg, #fffaf0, #fef3c7);
  background-size: 28rpx 28rpx;
}

.demo-route__steps {
  display: flex;
  flex-direction: column;
  gap: 10rpx;
  margin-top: 16rpx;
}

.demo-step {
  display: grid;
  grid-template-columns: 58rpx minmax(0, 1fr);
  gap: 12rpx;
  padding: 14rpx;
  border: 1rpx dashed rgba(169, 50, 38, 0.2);
  border-radius: 16rpx;
  background: rgba(255, 250, 240, 0.7);
}

.demo-step__index {
  color: #a93226;
  font-size: 24rpx;
  font-weight: 900;
}

.demo-step__title,
.demo-step__desc {
  display: block;
}

.demo-step__title {
  color: #25211b;
  font-size: 25rpx;
  font-weight: 900;
}

.demo-step__desc {
  margin-top: 4rpx;
  color: #756f63;
  font-size: 21rpx;
  line-height: 1.38;
}

.alert-stack {
  display: flex;
  flex-direction: column;
  gap: 12rpx;
  margin-top: 16rpx;
}

.alert-card {
  padding: 18rpx;
  border-left: 8rpx solid #a93226;
  border-radius: 18rpx;
  background: linear-gradient(135deg, #fff7ed, #fffaf0);
}

.alert-card--handled {
  opacity: 0.62;
  border-left-color: #94a3b8;
}

.alert-card__head,
.alert-card__foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12rpx;
}

.alert-card__title {
  overflow: hidden;
  color: #25211b;
  font-size: 28rpx;
  font-weight: 900;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.alert-card__content {
  display: block;
  margin-top: 10rpx;
  color: #756f63;
  font-size: 23rpx;
  line-height: 1.45;
}

.alert-card__foot {
  margin-top: 12rpx;
  color: #9a3412;
  font-size: 20rpx;
  font-weight: 800;
}
</style>
