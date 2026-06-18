<template>
  <view class="page stack dense-page">
    <view class="module-hero">
      <text class="module-hero__eyebrow">CAPTURE DISPATCH</text>
      <text class="module-hero__title">勾魂任务</text>
      <text class="module-hero__desc">外勤鬼差任务池，按状态快速切换，异常任务优先复盘。</text>
      <view class="module-hero__stats">
        <view class="module-hero__stat">
          <text>{{ allItems.length }}</text>
          <text>全部任务</text>
        </view>
        <view class="module-hero__stat">
          <text>{{ countByStatus('pending') }}</text>
          <text>待执行</text>
        </view>
        <view class="module-hero__stat">
          <text>{{ countByStatus('exception') }}</text>
          <text>异常</text>
        </view>
      </view>
    </view>

    <view class="compact-tabs">
      <up-tabs :list="tabs" :current="current" @change="changeTab" />
    </view>

    <EmptyState v-if="items.length === 0" text="暂无勾魂任务" />

    <view v-else class="card-list">
      <view v-for="item in items" :key="item.id" class="work-card" @click="openDetail(item.id)">
        <view class="work-card__head">
          <view class="work-card__title">
            <text>{{ item.soulName }}</text>
            <text class="work-card__subtitle">任务 #{{ item.id }} · {{ item.assigneeName || '待派单' }}</text>
          </view>
          <StatusTag :value="item.status" />
        </view>
        <view class="work-card__meta">{{ item.location }}</view>
        <view class="work-card__grid">
          <view class="work-card__field">
            <text class="work-card__field-label">计划</text>
            <text class="work-card__field-value">{{ displayTime(item.scheduledTime) }}</text>
          </view>
          <view class="work-card__field">
            <text class="work-card__field-label">实际</text>
            <text class="work-card__field-value">{{ displayTime(item.actualTime) }}</text>
          </view>
        </view>
        <view v-if="item.exceptionNote" class="work-card__note">{{ item.exceptionNote }}</view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import EmptyState from '../../components/EmptyState.vue'
import StatusTag from '../../components/StatusTag.vue'
import { getCaptureTasks } from '../../api/capture'
import { displayTime, showError } from '../../utils/format'

const tabs = [
  { name: '全部', status: '' },
  { name: '待执行', status: 'pending' },
  { name: '执行中', status: 'running' },
  { name: '异常', status: 'exception' },
  { name: '已完成', status: 'completed' },
  { name: '已撤销', status: 'canceled' }
]
const current = ref(0)
const allItems = ref([])
const items = computed(() => {
  const status = tabs[current.value].status
  return status ? allItems.value.filter(item => item.status === status) : allItems.value
})

function changeTab(event) {
  current.value = event.index
}

function openDetail(id) {
  uni.navigateTo({ url: `/pages/capture/detail?id=${id}` })
}

function countByStatus(status) {
  return allItems.value.filter(item => item.status === status).length
}

async function load() {
  try {
    allItems.value = await getCaptureTasks()
  } catch (error) {
    showError(error)
  }
}

onShow(load)
</script>
