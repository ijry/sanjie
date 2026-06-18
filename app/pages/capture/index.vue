<template>
  <view class="page stack dense-page capture-list-page">
    <view class="capture-board">
      <view class="capture-board__radar">
        <view class="capture-board__ring"></view>
        <view class="capture-board__sweep"></view>
        <text>{{ countByStatus('running') }}</text>
      </view>
      <view class="capture-board__main">
        <text class="capture-board__eyebrow">CAPTURE DISPATCH</text>
        <text class="capture-board__title">勾魂任务</text>
        <text class="capture-board__desc">外勤鬼差任务池，按状态快速切换，异常任务优先复盘。</text>
      </view>
      <view class="capture-board__stats">
        <view class="capture-board__stat">
          <text>{{ allItems.length }}</text>
          <text>全部任务</text>
        </view>
        <view class="capture-board__stat">
          <text>{{ countByStatus('pending') }}</text>
          <text>待执行</text>
        </view>
        <view class="capture-board__stat">
          <text>{{ countByStatus('exception') }}</text>
          <text>异常</text>
        </view>
      </view>
    </view>

    <view class="compact-tabs">
      <up-tabs :list="tabs" :current="current" @change="changeTab" />
    </view>

    <EmptyState v-if="items.length === 0" text="暂无勾魂任务" />

    <view v-else class="capture-task-list">
      <view v-for="item in items" :key="item.id" class="capture-task-card" @click="openDetail(item.id)">
        <view class="capture-task-card__line"></view>
        <view class="capture-task-card__head">
          <view class="capture-task-card__title">
            <text>{{ item.soulName }}</text>
            <text>任务 #{{ item.id }} · {{ item.assigneeName || '待派单' }}</text>
          </view>
          <StatusTag :value="item.status" />
        </view>
        <view class="capture-task-card__location">{{ item.location }}</view>
        <view class="capture-task-card__timeline">
          <view>
            <text class="capture-task-card__label">计划窗口</text>
            <text class="capture-task-card__value">{{ displayTime(item.scheduledTime) }}</text>
          </view>
          <view>
            <text class="capture-task-card__label">回执时间</text>
            <text class="capture-task-card__value">{{ displayTime(item.actualTime) }}</text>
          </view>
        </view>
        <view v-if="item.exceptionNote" class="capture-task-card__alert">{{ item.exceptionNote }}</view>
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

<style scoped>
.capture-list-page {
  background:
    radial-gradient(circle at 18% 0%, rgba(56, 189, 248, 0.18), transparent 26%),
    linear-gradient(180deg, #0f172a 0%, #1e293b 26%, #f7f1e3 26%, #f7f1e3 100%);
}

.capture-board {
  display: grid;
  grid-template-columns: 142rpx minmax(0, 1fr);
  gap: 18rpx;
  padding: 24rpx;
  color: #f8fafc;
  border: 1rpx solid rgba(125, 211, 252, 0.26);
  border-radius: 28rpx;
  background:
    linear-gradient(90deg, rgba(125, 211, 252, 0.08) 1rpx, transparent 1rpx),
    linear-gradient(0deg, rgba(125, 211, 252, 0.08) 1rpx, transparent 1rpx),
    linear-gradient(135deg, #020617, #172554);
  background-size: 38rpx 38rpx;
  box-shadow: 0 18rpx 42rpx rgba(15, 23, 42, 0.24);
}

.capture-board__radar {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 132rpx;
  height: 132rpx;
  border: 1rpx solid rgba(125, 211, 252, 0.35);
  border-radius: 50%;
  background: rgba(14, 165, 233, 0.08);
}

.capture-board__ring,
.capture-board__sweep {
  position: absolute;
  inset: 24rpx;
  border: 1rpx dashed rgba(125, 211, 252, 0.38);
  border-radius: 50%;
}

.capture-board__sweep {
  inset: 0;
  border: 0;
  background: conic-gradient(from 20deg, rgba(56, 189, 248, 0.36), transparent 32%);
}

.capture-board__radar text {
  position: relative;
  z-index: 1;
  color: #bae6fd;
  font-size: 42rpx;
  font-weight: 900;
}

.capture-board__main {
  min-width: 0;
}

.capture-board__eyebrow,
.capture-board__title,
.capture-board__desc {
  display: block;
}

.capture-board__eyebrow {
  color: #7dd3fc;
  font-size: 20rpx;
  font-weight: 900;
  letter-spacing: 2rpx;
}

.capture-board__title {
  margin-top: 6rpx;
  font-size: 44rpx;
  font-weight: 900;
}

.capture-board__desc {
  margin-top: 8rpx;
  color: rgba(248, 250, 252, 0.72);
  font-size: 24rpx;
  line-height: 1.42;
}

.capture-board__stats {
  grid-column: 1 / -1;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10rpx;
}

.capture-board__stat {
  padding: 12rpx;
  border-left: 4rpx solid #38bdf8;
  background: rgba(15, 23, 42, 0.72);
}

.capture-board__stat text:first-child {
  display: block;
  font-size: 30rpx;
  font-weight: 900;
}

.capture-board__stat text:last-child {
  display: block;
  margin-top: 4rpx;
  color: rgba(226, 232, 240, 0.66);
  font-size: 20rpx;
}

.capture-task-list {
  display: flex;
  flex-direction: column;
  gap: 14rpx;
}

.capture-task-card {
  position: relative;
  overflow: hidden;
  padding: 20rpx;
  border: 1rpx solid rgba(56, 189, 248, 0.2);
  border-radius: 20rpx;
  background: linear-gradient(135deg, #ffffff, #f8fafc);
  box-shadow: 0 10rpx 28rpx rgba(15, 23, 42, 0.08);
}

.capture-task-card__line {
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 8rpx;
  background: linear-gradient(180deg, #38bdf8, #ef4444);
}

.capture-task-card__head {
  display: flex;
  justify-content: space-between;
  gap: 12rpx;
}

.capture-task-card__title {
  min-width: 0;
}

.capture-task-card__title text:first-child {
  display: block;
  color: #0f172a;
  font-size: 32rpx;
  font-weight: 900;
}

.capture-task-card__title text:last-child {
  display: block;
  margin-top: 6rpx;
  color: #64748b;
  font-size: 22rpx;
}

.capture-task-card__location {
  margin-top: 14rpx;
  padding: 14rpx;
  color: #0f172a;
  border: 1rpx dashed rgba(56, 189, 248, 0.28);
  border-radius: 14rpx;
  background: rgba(224, 242, 254, 0.42);
  font-size: 24rpx;
  font-weight: 800;
}

.capture-task-card__timeline {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10rpx;
  margin-top: 12rpx;
}

.capture-task-card__timeline view {
  min-width: 0;
  padding: 12rpx;
  border-radius: 14rpx;
  background: #f1f5f9;
}

.capture-task-card__label,
.capture-task-card__value {
  display: block;
}

.capture-task-card__label {
  color: #64748b;
  font-size: 20rpx;
}

.capture-task-card__value {
  overflow: hidden;
  margin-top: 4rpx;
  color: #0f172a;
  font-size: 23rpx;
  font-weight: 900;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.capture-task-card__alert {
  margin-top: 12rpx;
  padding: 12rpx;
  color: #7f1d1d;
  border-radius: 14rpx;
  background: #fef2f2;
  font-size: 23rpx;
  line-height: 1.42;
}
</style>
