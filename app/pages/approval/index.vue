<template>
  <view class="page stack dense-page">
    <view class="module-hero">
      <text class="module-hero__eyebrow">APPROVAL CENTER</text>
      <text class="module-hero__title">审批中心</text>
      <text class="module-hero__desc">改寿、插队、刑期复核都进审批池，关系户也要留痕。</text>
      <view class="module-hero__stats">
        <view class="module-hero__stat">
          <text>{{ allItems.length }}</text>
          <text>总审批</text>
        </view>
        <view class="module-hero__stat">
          <text>{{ countByStatus('pending') }}</text>
          <text>待审</text>
        </view>
        <view class="module-hero__stat">
          <text>{{ countByStatus('transferred') }}</text>
          <text>转交</text>
        </view>
      </view>
    </view>

    <view class="compact-tabs">
      <up-tabs :list="tabs" :current="current" @change="changeTab" />
    </view>

    <EmptyState v-if="items.length === 0" text="暂无审批" />

    <view v-else class="card-list">
      <view v-for="item in items" :key="item.id" class="work-card" @click="openDetail(item.id)">
        <view class="work-card__head">
          <view class="work-card__title">
            <text>{{ item.title }}</text>
            <text class="work-card__subtitle">{{ item.type }} · 目标 #{{ item.targetId }}</text>
          </view>
          <StatusTag :value="item.status" />
        </view>
        <view class="work-card__meta">{{ item.reason || '无申请说明' }}</view>
        <view class="work-card__grid">
          <view class="work-card__field">
            <text class="work-card__field-label">申请人</text>
            <text class="work-card__field-value">用户 {{ item.applicantId || '-' }}</text>
          </view>
          <view class="work-card__field">
            <text class="work-card__field-label">审批人</text>
            <text class="work-card__field-value">用户 {{ item.approverId || '-' }}</text>
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
import { getApprovals } from '../../api/approvals'
import { showError } from '../../utils/format'

const tabs = [
  { name: '待审', status: 'pending' },
  { name: '转交', status: 'transferred' },
  { name: '通过', status: 'approved' },
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
  uni.navigateTo({ url: `/pages/approval/detail?id=${id}` })
}

function countByStatus(status) {
  return allItems.value.filter(item => item.status === status).length
}

async function load() {
  try {
    allItems.value = await getApprovals()
  } catch (error) {
    showError(error)
  }
}

onShow(load)
</script>
