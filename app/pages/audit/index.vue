<template>
  <view class="page stack dense-page">
    <view class="module-hero">
      <text class="module-hero__eyebrow">AUDIT TRAIL</text>
      <text class="module-hero__title">审计日志</text>
      <text class="module-hero__desc">所有敏感操作都留痕，方便追责批量划名和关系户投胎。</text>
      <view class="module-hero__stats">
        <view class="module-hero__stat">
          <text>{{ items.length }}</text>
          <text>日志</text>
        </view>
        <view class="module-hero__stat">
          <text>{{ approvalCount }}</text>
          <text>审批</text>
        </view>
        <view class="module-hero__stat">
          <text>{{ captureCount }}</text>
          <text>勾魂</text>
        </view>
      </view>
    </view>

    <EmptyState v-if="items.length === 0" text="暂无审计日志" />

    <view v-else class="card-list">
      <view v-for="item in items" :key="item.id" class="work-card">
        <view class="work-card__head">
          <view class="work-card__title">
            <text>{{ item.action }}</text>
            <text class="work-card__subtitle">{{ item.targetType }} #{{ item.targetId }}</text>
          </view>
          <text class="audit-actor">用户 {{ item.actorId || '-' }}</text>
        </view>
        <view class="work-card__note">{{ item.note || '无备注' }}</view>
        <view class="work-card__meta">{{ displayTime(item.createdAt) }}</view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import EmptyState from '../../components/EmptyState.vue'
import { getAuditLogs } from '../../api/audit'
import { displayTime, showError } from '../../utils/format'

const items = ref([])
const approvalCount = computed(() => items.value.filter(item => item.targetType === 'approval').length)
const captureCount = computed(() => items.value.filter(item => item.targetType === 'capture_task').length)

async function load() {
  try {
    items.value = await getAuditLogs()
  } catch (error) {
    showError(error)
  }
}

onShow(load)
</script>

<style scoped>
.audit-actor {
  flex-shrink: 0;
  color: #a93226;
  font-size: 22rpx;
  font-weight: 800;
}
</style>
