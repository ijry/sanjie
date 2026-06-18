<template>
  <view class="page stack dense-page approval-list-page">
    <view class="dossier-board">
      <view class="dossier-board__stamp">审</view>
      <text class="dossier-board__eyebrow">APPROVAL CENTER</text>
      <text class="dossier-board__title">审批中心</text>
      <text class="dossier-board__desc">改寿、插队、刑期复核都进审批池，关系户也要留痕。</text>
      <view class="dossier-board__stats">
        <view class="dossier-board__stat">
          <text>{{ allItems.length }}</text>
          <text>总审批</text>
        </view>
        <view class="dossier-board__stat">
          <text>{{ countByStatus('pending') }}</text>
          <text>待审</text>
        </view>
        <view class="dossier-board__stat">
          <text>{{ countByStatus('transferred') }}</text>
          <text>转交</text>
        </view>
      </view>
    </view>

    <view class="compact-tabs">
      <up-tabs :list="tabs" :current="current" @change="changeTab" />
    </view>

    <EmptyState v-if="items.length === 0" text="暂无审批" />

    <view v-else class="dossier-list">
      <view v-for="item in items" :key="item.id" class="dossier-card" @click="openDetail(item.id)">
        <view class="dossier-card__clip"></view>
        <view class="dossier-card__head">
          <view class="dossier-card__title">
            <text>{{ item.title }}</text>
            <text>{{ typeLabel(item.type) }} · 目标 #{{ item.targetId }}</text>
          </view>
          <StatusTag :value="item.status" />
        </view>
        <view class="dossier-card__reason">{{ item.reason || '无申请说明' }}</view>
        <view class="dossier-card__flow">
          <view>
            <text class="dossier-card__label">申请人</text>
            <text class="dossier-card__value">用户 {{ item.applicantId || '-' }}</text>
          </view>
          <view>
            <text class="dossier-card__label">审批人</text>
            <text class="dossier-card__value">用户 {{ item.approverId || '-' }}</text>
          </view>
        </view>
        <view v-if="item.resultNote" class="dossier-card__note">{{ item.resultNote }}</view>
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

function typeLabel(type) {
  const labels = {
    lifespan_change: '寿命变更',
    hell_review: '刑期复核',
    reincarnation_jump: '投胎插队',
    capture_exception: '勾魂异常',
    soup_extra_dose: '孟婆汤加量'
  }
  return labels[type] || type || '-'
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

<style scoped>
.approval-list-page {
  background:
    radial-gradient(circle at 0% 0%, rgba(185, 28, 28, 0.1), transparent 26%),
    linear-gradient(180deg, #fff7ed 0%, #f7f1e3 32%, #f7f1e3 100%);
}

.dossier-board {
  position: relative;
  overflow: hidden;
  padding: 26rpx;
  border: 2rpx solid rgba(120, 53, 15, 0.18);
  border-radius: 18rpx;
  background:
    linear-gradient(90deg, rgba(120, 53, 15, 0.05) 1rpx, transparent 1rpx),
    linear-gradient(0deg, rgba(120, 53, 15, 0.05) 1rpx, transparent 1rpx),
    #fff7ed;
  background-size: 30rpx 30rpx;
  box-shadow: 0 18rpx 38rpx rgba(120, 53, 15, 0.12);
}

.dossier-board__stamp {
  position: absolute;
  right: 28rpx;
  top: 26rpx;
  width: 116rpx;
  height: 116rpx;
  color: rgba(185, 28, 28, 0.2);
  border: 8rpx solid rgba(185, 28, 28, 0.16);
  border-radius: 24rpx;
  font-size: 62rpx;
  font-weight: 900;
  line-height: 116rpx;
  text-align: center;
  transform: rotate(-14deg);
}

.dossier-board__eyebrow,
.dossier-board__title,
.dossier-board__desc {
  position: relative;
  z-index: 1;
  display: block;
}

.dossier-board__eyebrow {
  color: #b45309;
  font-size: 20rpx;
  font-weight: 900;
  letter-spacing: 2rpx;
}

.dossier-board__title {
  margin-top: 6rpx;
  color: #431407;
  font-size: 44rpx;
  font-weight: 900;
}

.dossier-board__desc {
  width: 76%;
  margin-top: 8rpx;
  color: #7c2d12;
  font-size: 24rpx;
  line-height: 1.44;
}

.dossier-board__stats {
  position: relative;
  z-index: 1;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10rpx;
  margin-top: 22rpx;
}

.dossier-board__stat {
  padding: 14rpx 12rpx;
  border-bottom: 4rpx solid #b45309;
  background: rgba(255, 250, 240, 0.72);
}

.dossier-board__stat text:first-child {
  display: block;
  color: #431407;
  font-size: 30rpx;
  font-weight: 900;
}

.dossier-board__stat text:last-child {
  display: block;
  margin-top: 4rpx;
  color: #9a3412;
  font-size: 20rpx;
}

.dossier-list {
  display: flex;
  flex-direction: column;
  gap: 14rpx;
}

.dossier-card {
  position: relative;
  padding: 22rpx;
  border: 1rpx solid rgba(180, 83, 9, 0.18);
  border-radius: 16rpx;
  background: #fffaf0;
  box-shadow: 0 10rpx 28rpx rgba(120, 53, 15, 0.08);
}

.dossier-card__clip {
  position: absolute;
  left: 24rpx;
  top: -6rpx;
  width: 74rpx;
  height: 18rpx;
  border-radius: 0 0 8rpx 8rpx;
  background: #b45309;
}

.dossier-card__head {
  display: flex;
  justify-content: space-between;
  gap: 12rpx;
}

.dossier-card__title {
  min-width: 0;
}

.dossier-card__title text:first-child {
  display: block;
  color: #431407;
  font-size: 32rpx;
  font-weight: 900;
}

.dossier-card__title text:last-child {
  display: block;
  margin-top: 6rpx;
  color: #9a3412;
  font-size: 22rpx;
}

.dossier-card__reason {
  margin-top: 14rpx;
  padding: 14rpx;
  color: #7c2d12;
  border: 1rpx dashed rgba(180, 83, 9, 0.22);
  border-radius: 14rpx;
  background: #fff7ed;
  font-size: 24rpx;
  line-height: 1.42;
}

.dossier-card__flow {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10rpx;
  margin-top: 12rpx;
}

.dossier-card__flow view {
  min-width: 0;
  padding: 12rpx;
  border-radius: 14rpx;
  background: #fffbeb;
}

.dossier-card__label,
.dossier-card__value {
  display: block;
}

.dossier-card__label {
  color: #9a3412;
  font-size: 20rpx;
}

.dossier-card__value {
  margin-top: 4rpx;
  color: #431407;
  font-size: 24rpx;
  font-weight: 900;
}

.dossier-card__note {
  margin-top: 12rpx;
  padding: 12rpx;
  color: #7f1d1d;
  border-radius: 14rpx;
  background: #fef2f2;
  font-size: 23rpx;
  line-height: 1.42;
}
</style>
