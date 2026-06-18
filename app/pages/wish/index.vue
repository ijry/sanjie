<template>
  <view class="page stack dense-page wish-list-page">
    <view class="wish-board">
      <view class="wish-board__ribbon">愿</view>
      <text class="wish-board__eyebrow">WISH TICKETS</text>
      <text class="wish-board__title">愿望工单</text>
      <text class="wish-board__desc">凡人许愿先排队，香火、功德、合理性一起看。</text>
      <view class="wish-board__stats">
        <view class="wish-board__stat">
          <text>{{ allItems.length }}</text>
          <text>工单</text>
        </view>
        <view class="wish-board__stat">
          <text>{{ countByStatus('pending') }}</text>
          <text>待处理</text>
        </view>
        <view class="wish-board__stat">
          <text>{{ incenseTotal }}</text>
          <text>总香火</text>
        </view>
      </view>
    </view>

    <view class="compact-tabs">
      <up-tabs :list="tabs" :current="current" @change="changeTab" />
    </view>

    <EmptyState v-if="items.length === 0" text="暂无愿望工单" />

    <view v-else class="wish-ticket-list">
      <view v-for="item in items" :key="item.id" class="wish-ticket-card" @click="openDetail(item.id)">
        <view class="wish-ticket-card__type">{{ wishTypeLabel(item.wishType) }}</view>
        <view class="wish-ticket-card__head">
          <view class="wish-ticket-card__title">
            <text>{{ item.title }}</text>
            <text>{{ item.requesterName }} · {{ item.wishType }}</text>
          </view>
          <StatusTag :value="item.status" />
        </view>
        <view class="wish-ticket-card__scale">
          <view>
            <text class="wish-ticket-card__amount">{{ item.incenseAmount }}</text>
            <text class="wish-ticket-card__label">香火</text>
          </view>
          <view class="wish-ticket-card__beam">
            <text :style="{ width: wishWeight(item) }"></text>
          </view>
          <view>
            <text class="wish-ticket-card__amount">{{ item.meritScore }}</text>
            <text class="wish-ticket-card__label">功德</text>
          </view>
        </view>
        <view class="wish-ticket-card__footer">
          <text>{{ item.assignedDeity || '待分配' }}</text>
          <text>{{ displayTime(item.updatedAt) }}</text>
        </view>
        <view v-if="item.resultNote" class="wish-ticket-card__note">{{ item.resultNote }}</view>
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

function wishTypeLabel(type) {
  const labels = {
    wealth: '财',
    study: '学',
    work: '工',
    love: '缘',
    health: '康',
    career: '业',
    luck: '运'
  }
  return labels[type] || '愿'
}

function wishWeight(item) {
  const total = Math.max(1, Number(item.incenseAmount || 0) + Math.max(0, Number(item.meritScore || 0)))
  return `${Math.min(92, Math.max(8, Math.round((Number(item.incenseAmount || 0) / total) * 100)))}%`
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

<style scoped>
.wish-list-page {
  background:
    radial-gradient(circle at 88% 0%, rgba(234, 179, 8, 0.2), transparent 30%),
    linear-gradient(180deg, #fff7ed 0%, #f7f1e3 34%, #f7f1e3 100%);
}

.wish-board {
  position: relative;
  overflow: hidden;
  padding: 26rpx;
  border: 1rpx solid rgba(217, 119, 6, 0.22);
  border-radius: 34rpx 34rpx 12rpx 12rpx;
  background:
    radial-gradient(circle at 12% 18%, rgba(251, 191, 36, 0.2), transparent 20%),
    linear-gradient(135deg, #fffbeb, #fff7ed 58%, #fef3c7);
  box-shadow: 0 18rpx 42rpx rgba(217, 119, 6, 0.12);
}

.wish-board__ribbon {
  position: absolute;
  right: -38rpx;
  top: 30rpx;
  width: 162rpx;
  padding: 8rpx 0;
  color: #fffaf0;
  background: #d97706;
  font-size: 28rpx;
  font-weight: 900;
  text-align: center;
  transform: rotate(34deg);
}

.wish-board__eyebrow,
.wish-board__title,
.wish-board__desc {
  display: block;
}

.wish-board__eyebrow {
  color: #d97706;
  font-size: 20rpx;
  font-weight: 900;
  letter-spacing: 2rpx;
}

.wish-board__title {
  margin-top: 6rpx;
  color: #422006;
  font-size: 44rpx;
  font-weight: 900;
}

.wish-board__desc {
  width: 78%;
  margin-top: 8rpx;
  color: #78350f;
  font-size: 24rpx;
  line-height: 1.42;
}

.wish-board__stats {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10rpx;
  margin-top: 22rpx;
}

.wish-board__stat {
  padding: 14rpx 12rpx;
  border: 1rpx solid rgba(217, 119, 6, 0.18);
  border-radius: 16rpx;
  background: rgba(255, 250, 240, 0.68);
}

.wish-board__stat text:first-child {
  display: block;
  color: #92400e;
  font-size: 30rpx;
  font-weight: 900;
}

.wish-board__stat text:last-child {
  display: block;
  margin-top: 4rpx;
  color: #a16207;
  font-size: 20rpx;
}

.wish-ticket-list {
  display: flex;
  flex-direction: column;
  gap: 14rpx;
}

.wish-ticket-card {
  position: relative;
  overflow: hidden;
  padding: 20rpx;
  border: 1rpx solid rgba(217, 119, 6, 0.18);
  border-radius: 24rpx 24rpx 12rpx 12rpx;
  background: linear-gradient(135deg, #fffbeb, #fffaf0);
  box-shadow: 0 10rpx 28rpx rgba(217, 119, 6, 0.08);
}

.wish-ticket-card__type {
  position: absolute;
  right: 18rpx;
  bottom: 18rpx;
  width: 78rpx;
  height: 78rpx;
  color: rgba(217, 119, 6, 0.2);
  border: 5rpx solid rgba(217, 119, 6, 0.18);
  border-radius: 50%;
  font-size: 40rpx;
  font-weight: 900;
  line-height: 78rpx;
  text-align: center;
  transform: rotate(-12deg);
}

.wish-ticket-card__head {
  position: relative;
  z-index: 1;
  display: flex;
  justify-content: space-between;
  gap: 12rpx;
}

.wish-ticket-card__title text:first-child {
  display: block;
  color: #422006;
  font-size: 32rpx;
  font-weight: 900;
}

.wish-ticket-card__title text:last-child {
  display: block;
  margin-top: 6rpx;
  color: #a16207;
  font-size: 22rpx;
}

.wish-ticket-card__scale {
  position: relative;
  z-index: 1;
  display: grid;
  grid-template-columns: 96rpx minmax(0, 1fr) 96rpx;
  gap: 12rpx;
  align-items: center;
  margin-top: 18rpx;
}

.wish-ticket-card__amount,
.wish-ticket-card__label {
  display: block;
  text-align: center;
}

.wish-ticket-card__amount {
  color: #92400e;
  font-size: 30rpx;
  font-weight: 900;
}

.wish-ticket-card__label {
  color: #a16207;
  font-size: 20rpx;
}

.wish-ticket-card__beam {
  height: 16rpx;
  overflow: hidden;
  border-radius: 999rpx;
  background: rgba(217, 119, 6, 0.16);
}

.wish-ticket-card__beam text {
  display: block;
  height: 100%;
  border-radius: 999rpx;
  background: linear-gradient(90deg, #f59e0b, #d97706);
}

.wish-ticket-card__footer {
  position: relative;
  z-index: 1;
  display: flex;
  justify-content: space-between;
  gap: 12rpx;
  margin-top: 14rpx;
  padding: 12rpx;
  border: 1rpx dashed rgba(217, 119, 6, 0.24);
  border-radius: 14rpx;
  color: #78350f;
  font-size: 22rpx;
  font-weight: 800;
}

.wish-ticket-card__note {
  position: relative;
  z-index: 1;
  margin-top: 12rpx;
  padding: 12rpx;
  color: #78350f;
  border-radius: 14rpx;
  background: rgba(254, 243, 199, 0.62);
  font-size: 23rpx;
  line-height: 1.42;
}
</style>
