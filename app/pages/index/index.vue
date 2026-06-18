<template>
  <view class="page stack">
    <view class="hero">
      <text class="hero__title">三界工作台</text>
      <text class="hero__sub">今日公务，按天条流程处理</text>
    </view>

    <view class="metric-grid">
      <MetricCard label="待勾魂" :value="metrics.pendingCaptureCount || 0" />
      <MetricCard label="待审批" :value="metrics.pendingApprovalCount || 0" />
      <MetricCard label="投胎积压" :value="metrics.reincarnationQueueCount || 0" />
      <MetricCard label="高危寿命" :value="metrics.highRiskLifeBookCount || 0" />
    </view>

    <view class="panel">
      <SectionTitle title="快捷入口" />
      <up-grid :border="false" col="4">
        <up-grid-item v-for="item in shortcuts" :key="item.url" @click="go(item.url)">
          <text class="shortcut">{{ item.text }}</text>
        </up-grid-item>
      </up-grid>
    </view>

    <view class="panel">
      <SectionTitle title="异常预警" />
      <EmptyState v-if="alerts.length === 0" text="暂无预警" />
      <up-cell-group v-else :border="false">
        <up-cell v-for="item in alerts" :key="item.id" :title="item.title" :label="item.content" clickable @click="markAlert(item)">
          <template #value>
            <StatusTag :value="item.level" />
          </template>
        </up-cell>
      </up-cell-group>
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
  { text: '生死簿', url: '/pages/life-book/search' },
  { text: '地狱', url: '/pages/hell/index' },
  { text: '愿望', url: '/pages/wish/index' },
  { text: '孟婆汤', url: '/pages/soup/index' },
  { text: '审计', url: '/pages/audit/index' }
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
.hero {
  padding: 34rpx 28rpx;
  color: #fffaf0;
  background: linear-gradient(135deg, #25211b, #7d1d14);
  border-radius: 8rpx;
}

.hero__title {
  display: block;
  font-size: 42rpx;
  font-weight: 800;
}

.hero__sub {
  display: block;
  margin-top: 10rpx;
  color: #e5d8bd;
  font-size: 24rpx;
}

.metric-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 18rpx;
}

.shortcut {
  padding: 20rpx 0;
  color: #a93226;
  font-size: 24rpx;
}
</style>

