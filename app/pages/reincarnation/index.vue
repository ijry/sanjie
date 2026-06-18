<template>
  <view class="page stack dense-page">
    <view class="module-hero">
      <text class="module-hero__eyebrow">REINCARNATION QUEUE</text>
      <text class="module-hero__title">投胎队列</text>
      <text class="module-hero__desc">胎位、优先级、孟婆汤状态集中处理，避免奈何桥堵车。</text>
      <view class="module-hero__stats">
        <view class="module-hero__stat">
          <text>{{ allItems.length }}</text>
          <text>投胎单</text>
        </view>
        <view class="module-hero__stat">
          <text>{{ countByStatus('queued') + countByStatus('reviewing') }}</text>
          <text>待复核</text>
        </view>
        <view class="module-hero__stat">
          <text>{{ countByStatus('pending_soup') }}</text>
          <text>待喝汤</text>
        </view>
      </view>
    </view>

    <view class="compact-tabs">
      <up-tabs :list="tabs" :current="current" @change="changeTab" />
    </view>

    <EmptyState v-if="items.length === 0" text="暂无投胎单" />

    <view v-else class="card-list">
      <view v-for="item in items" :key="item.id" class="work-card" @click="openDetail(item.id)">
        <view class="work-card__head">
          <view class="work-card__title">
            <text>{{ item.soulName }}</text>
            <text class="work-card__subtitle">目标 {{ item.targetRealm }} · 投胎单 #{{ item.id }}</text>
          </view>
          <StatusTag :value="item.status" />
        </view>
        <view class="work-card__grid">
          <view class="work-card__field">
            <text class="work-card__field-label">优先级</text>
            <text class="work-card__field-value">{{ item.priority }}</text>
          </view>
          <view class="work-card__field">
            <text class="work-card__field-label">胎位</text>
            <text class="work-card__field-value">{{ item.quotaType || '未分配' }}</text>
          </view>
          <view class="work-card__field">
            <text class="work-card__field-label">孟婆汤</text>
            <text class="work-card__field-value">{{ item.soupStatus }}</text>
          </view>
          <view class="work-card__field">
            <text class="work-card__field-label">更新时间</text>
            <text class="work-card__field-value">{{ displayTime(item.updatedAt) }}</text>
          </view>
        </view>
        <view v-if="item.reviewNote" class="work-card__note">{{ item.reviewNote }}</view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import EmptyState from '../../components/EmptyState.vue'
import StatusTag from '../../components/StatusTag.vue'
import { getReincarnations } from '../../api/reincarnation'
import { displayTime, showError } from '../../utils/format'

const tabs = [
  { name: '全部', status: '' },
  { name: '排队', status: 'queued' },
  { name: '复核', status: 'reviewing' },
  { name: '待喝汤', status: 'pending_soup' },
  { name: '完成', status: 'completed' },
  { name: '驳回', status: 'rejected' }
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
  uni.navigateTo({ url: `/pages/reincarnation/detail?id=${id}` })
}

function countByStatus(status) {
  return allItems.value.filter(item => item.status === status).length
}

async function load() {
  try {
    allItems.value = await getReincarnations()
  } catch (error) {
    showError(error)
  }
}

onShow(load)
</script>
