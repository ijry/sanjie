<template>
  <view class="page stack dense-page">
    <view class="module-hero">
      <text class="module-hero__eyebrow">WISH TICKETS</text>
      <text class="module-hero__title">愿望工单</text>
      <text class="module-hero__desc">凡人许愿先排队，香火、功德、合理性一起看。</text>
      <view class="module-hero__stats">
        <view class="module-hero__stat">
          <text>{{ allItems.length }}</text>
          <text>工单</text>
        </view>
        <view class="module-hero__stat">
          <text>{{ countByStatus('pending') }}</text>
          <text>待处理</text>
        </view>
        <view class="module-hero__stat">
          <text>{{ incenseTotal }}</text>
          <text>总香火</text>
        </view>
      </view>
    </view>

    <view class="compact-tabs">
      <up-tabs :list="tabs" :current="current" @change="changeTab" />
    </view>

    <EmptyState v-if="items.length === 0" text="暂无愿望工单" />

    <view v-else class="card-list">
      <view v-for="item in items" :key="item.id" class="work-card" @click="openDetail(item.id)">
        <view class="work-card__head">
          <view class="work-card__title">
            <text>{{ item.title }}</text>
            <text class="work-card__subtitle">{{ item.requesterName }} · {{ item.wishType }}</text>
          </view>
          <StatusTag :value="item.status" />
        </view>
        <view class="work-card__grid">
          <view class="work-card__field">
            <text class="work-card__field-label">香火</text>
            <text class="work-card__field-value">{{ item.incenseAmount }}</text>
          </view>
          <view class="work-card__field">
            <text class="work-card__field-label">功德</text>
            <text class="work-card__field-value">{{ item.meritScore }}</text>
          </view>
          <view class="work-card__field">
            <text class="work-card__field-label">转派</text>
            <text class="work-card__field-value">{{ item.assignedDeity || '待分配' }}</text>
          </view>
          <view class="work-card__field">
            <text class="work-card__field-label">更新时间</text>
            <text class="work-card__field-value">{{ displayTime(item.updatedAt) }}</text>
          </view>
        </view>
        <view v-if="item.resultNote" class="work-card__note">{{ item.resultNote }}</view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import EmptyState from '../../components/EmptyState.vue'
import StatusTag from '../../components/StatusTag.vue'
import { getWishes } from '../../api/wishes'
import { displayTime, showError } from '../../utils/format'

const tabs = [
  { name: '全部', status: '' },
  { name: '待处理', status: 'pending' },
  { name: '已转派', status: 'routed' },
  { name: '已办结', status: 'resolved' },
  { name: '已驳回', status: 'rejected' }
]
const current = ref(0)
const allItems = ref([])
const items = computed(() => {
  const status = tabs[current.value].status
  return status ? allItems.value.filter(item => item.status === status) : allItems.value
})
const incenseTotal = computed(() => allItems.value.reduce((total, item) => total + Number(item.incenseAmount || 0), 0))

function changeTab(event) {
  current.value = event.index
}

function openDetail(id) {
  uni.navigateTo({ url: `/pages/wish/detail?id=${id}` })
}

function countByStatus(status) {
  return allItems.value.filter(item => item.status === status).length
}

async function load() {
  try {
    allItems.value = await getWishes()
  } catch (error) {
    showError(error)
  }
}

onShow(load)
</script>
