<template>
  <view class="page stack dense-page">
    <view class="module-hero">
      <text class="module-hero__eyebrow">SANJIE CONSOLE</text>
      <text class="module-hero__title">三界工作台</text>
      <text class="module-hero__desc">今日公务，按天条流程处理；点击指标或快捷入口可直接进入业务页。</text>
      <view class="module-hero__stats">
        <view class="module-hero__stat" @click="go('/pages/capture/index')">
          <text>{{ metrics.pendingCaptureCount || 0 }}</text>
          <text>待勾魂</text>
        </view>
        <view class="module-hero__stat" @click="go('/pages/approval/index')">
          <text>{{ metrics.pendingApprovalCount || 0 }}</text>
          <text>待审批</text>
        </view>
        <view class="module-hero__stat" @click="go('/pages/reincarnation/index')">
          <text>{{ metrics.reincarnationQueueCount || 0 }}</text>
          <text>投胎积压</text>
        </view>
      </view>
    </view>

    <view class="metric-grid">
      <view @click="go('/pages/capture/index')"><MetricCard label="待勾魂" :value="metrics.pendingCaptureCount || 0" /></view>
      <view @click="go('/pages/approval/index')"><MetricCard label="待审批" :value="metrics.pendingApprovalCount || 0" /></view>
      <view @click="go('/pages/reincarnation/index')"><MetricCard label="投胎积压" :value="metrics.reincarnationQueueCount || 0" /></view>
      <view @click="go('/pages/life-book/search')"><MetricCard label="高危寿命" :value="metrics.highRiskLifeBookCount || 0" /></view>
    </view>

    <view class="detail-card">
      <SectionTitle title="快捷入口" />
      <view class="shortcut-grid">
        <view v-for="item in shortcuts" :key="item.url" class="shortcut-card" @click="go(item.url)">
          <text class="shortcut-card__name">{{ item.text }}</text>
          <text class="shortcut-card__desc">{{ item.desc }}</text>
        </view>
      </view>
    </view>

    <view class="detail-card">
      <SectionTitle title="异常预警" />
      <EmptyState v-if="alerts.length === 0" text="暂无预警" />
      <view v-else class="mini-list">
        <view v-for="item in alerts" :key="item.id" class="mini-row" @click="markAlert(item)">
          <view class="mini-row__main">
            <text class="mini-row__title">{{ item.title }}</text>
            <text class="mini-row__meta">{{ item.content }}</text>
          </view>
          <view>
            <StatusTag :value="item.level" />
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import MetricCard from '../../components/MetricCard.vue'
import StatusTag from '../../components/StatusTag.vue'
import EmptyState from '../../components/EmptyState.vue'
import SectionTitle from '../../components/SectionTitle.vue'
import { getAlerts, getDashboard, handleAlert } from '../../api/dashboard'
import { showError, showSuccess } from '../../utils/format'

const metrics = ref({})
const alerts = ref([])
const shortcuts = [
  { text: '生死簿', desc: '查阳寿风险', url: '/pages/life-book/search' },
  { text: '地狱', desc: '看楼层负载', url: '/pages/hell/index' },
  { text: '愿望', desc: '处理凡人诉求', url: '/pages/wish/index' },
  { text: '孟婆汤', desc: '库存与发放', url: '/pages/soup/index' },
  { text: '审计', desc: '追踪高危操作', url: '/pages/audit/index' }
]

function go(url) {
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

<style scoped lang="scss">
.metric-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 18rpx;
}

.shortcut-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12rpx;
}

.shortcut-card {
  padding: 18rpx;
  border: 1rpx solid rgba(229, 216, 189, 0.72);
  border-radius: 16rpx;
  background: rgba(247, 241, 227, 0.72);
}

.shortcut-card__name {
  display: block;
  color: #a93226;
  font-size: 28rpx;
  font-weight: 900;
}

.shortcut-card__desc {
  display: block;
  margin-top: 6rpx;
  color: #756f63;
  font-size: 22rpx;
}
</style>
