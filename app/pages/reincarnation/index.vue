<template>
  <view class="page stack dense-page reincarnation-list-page">
    <view class="queue-board">
      <view>
        <text class="queue-board__eyebrow">REINCARNATION QUEUE</text>
        <text class="queue-board__title">投胎队列</text>
        <text class="queue-board__desc">胎位、优先级、孟婆汤状态集中处理，避免奈何桥堵车。</text>
      </view>
      <view class="queue-board__stats">
        <view class="queue-board__stat">
          <text>{{ allItems.length }}</text>
          <text>投胎单</text>
        </view>
        <view class="queue-board__stat">
          <text>{{ countByStatus('queued') + countByStatus('reviewing') }}</text>
          <text>待复核</text>
        </view>
        <view class="queue-board__stat">
          <text>{{ countByStatus('pending_soup') }}</text>
          <text>待喝汤</text>
        </view>
      </view>
    </view>

    <view class="compact-tabs">
      <up-tabs :list="tabs" :current="current" @change="changeTab" />
    </view>

    <EmptyState v-if="items.length === 0" text="暂无投胎单" />

    <view v-else class="ticket-list">
      <view v-for="item in items" :key="item.id" class="queue-ticket" @click="openDetail(item.id)">
        <view class="queue-ticket__edge queue-ticket__edge--left"></view>
        <view class="queue-ticket__edge queue-ticket__edge--right"></view>
        <view class="queue-ticket__head">
          <view class="queue-ticket__title">
            <text>{{ item.soulName }}</text>
            <text>目标 {{ item.targetRealm }} · 投胎单 #{{ item.id }}</text>
          </view>
          <StatusTag :value="item.status" />
        </view>
        <view class="queue-ticket__route">
          <text>鬼门关</text>
          <text>奈何桥</text>
          <text>{{ item.targetRealm }}</text>
        </view>
        <view class="queue-ticket__grid">
          <view>
            <text class="queue-ticket__label">优先级</text>
            <text class="queue-ticket__value">{{ item.priority }}</text>
          </view>
          <view>
            <text class="queue-ticket__label">胎位</text>
            <text class="queue-ticket__value">{{ item.quotaType || '未分配' }}</text>
          </view>
          <view>
            <text class="queue-ticket__label">孟婆汤</text>
            <text class="queue-ticket__value">{{ item.soupStatus }}</text>
          </view>
          <view>
            <text class="queue-ticket__label">更新时间</text>
            <text class="queue-ticket__value">{{ displayTime(item.updatedAt) }}</text>
          </view>
        </view>
        <view v-if="item.reviewNote" class="queue-ticket__note">{{ item.reviewNote }}</view>
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

<style scoped>
.reincarnation-list-page {
  background:
    radial-gradient(circle at 80% 0%, rgba(20, 184, 166, 0.16), transparent 28%),
    linear-gradient(180deg, #ecfeff 0%, #f7f1e3 30%, #f7f1e3 100%);
}

.queue-board {
  padding: 24rpx;
  border: 2rpx solid rgba(15, 118, 110, 0.18);
  border-radius: 30rpx;
  background:
    linear-gradient(90deg, rgba(15, 118, 110, 0.08) 0 2rpx, transparent 2rpx 24rpx),
    linear-gradient(135deg, #f0fdfa, #fff7ed);
  box-shadow: 0 18rpx 42rpx rgba(15, 118, 110, 0.12);
}

.queue-board__eyebrow,
.queue-board__title,
.queue-board__desc {
  display: block;
}

.queue-board__eyebrow {
  color: #0f766e;
  font-size: 20rpx;
  font-weight: 900;
  letter-spacing: 2rpx;
}

.queue-board__title {
  margin-top: 6rpx;
  color: #134e4a;
  font-size: 44rpx;
  font-weight: 900;
}

.queue-board__desc {
  margin-top: 8rpx;
  color: #475569;
  font-size: 24rpx;
  line-height: 1.42;
}

.queue-board__stats {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10rpx;
  margin-top: 20rpx;
}

.queue-board__stat {
  padding: 14rpx 12rpx;
  border: 1rpx dashed rgba(15, 118, 110, 0.24);
  border-radius: 16rpx;
  background: rgba(255, 250, 240, 0.68);
}

.queue-board__stat text:first-child {
  display: block;
  color: #134e4a;
  font-size: 30rpx;
  font-weight: 900;
}

.queue-board__stat text:last-child {
  display: block;
  margin-top: 4rpx;
  color: #64748b;
  font-size: 20rpx;
}

.ticket-list {
  display: flex;
  flex-direction: column;
  gap: 14rpx;
}

.queue-ticket {
  position: relative;
  overflow: hidden;
  padding: 20rpx;
  border: 1rpx solid rgba(15, 118, 110, 0.18);
  border-radius: 22rpx;
  background: linear-gradient(135deg, #f0fdfa, #fffaf0);
  box-shadow: 0 10rpx 28rpx rgba(15, 118, 110, 0.08);
}

.queue-ticket__edge {
  position: absolute;
  top: 50%;
  width: 34rpx;
  height: 34rpx;
  border-radius: 50%;
  background: #f7f1e3;
  transform: translateY(-50%);
}

.queue-ticket__edge--left {
  left: -17rpx;
}

.queue-ticket__edge--right {
  right: -17rpx;
}

.queue-ticket__head {
  display: flex;
  justify-content: space-between;
  gap: 12rpx;
}

.queue-ticket__title {
  min-width: 0;
}

.queue-ticket__title text:first-child {
  display: block;
  color: #134e4a;
  font-size: 32rpx;
  font-weight: 900;
}

.queue-ticket__title text:last-child {
  display: block;
  margin-top: 6rpx;
  color: #64748b;
  font-size: 22rpx;
}

.queue-ticket__route {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8rpx;
  margin-top: 14rpx;
}

.queue-ticket__route text {
  padding: 10rpx 8rpx;
  color: #0f766e;
  border: 1rpx dashed rgba(15, 118, 110, 0.24);
  border-radius: 999rpx;
  font-size: 21rpx;
  font-weight: 900;
  text-align: center;
}

.queue-ticket__grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10rpx;
  margin-top: 14rpx;
}

.queue-ticket__grid view {
  min-width: 0;
  padding: 12rpx;
  border-radius: 14rpx;
  background: rgba(255, 250, 240, 0.78);
}

.queue-ticket__label,
.queue-ticket__value {
  display: block;
}

.queue-ticket__label {
  color: #64748b;
  font-size: 20rpx;
}

.queue-ticket__value {
  overflow: hidden;
  margin-top: 4rpx;
  color: #134e4a;
  font-size: 24rpx;
  font-weight: 900;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.queue-ticket__note {
  margin-top: 12rpx;
  padding: 12rpx;
  color: #0f766e;
  border-radius: 14rpx;
  background: rgba(20, 184, 166, 0.1);
  font-size: 23rpx;
  line-height: 1.42;
}
</style>
