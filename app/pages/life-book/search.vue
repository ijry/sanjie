<template>
  <view class="page stack dense-page life-list-page">
    <view class="life-index-board">
      <view class="life-index-board__seal">簿</view>
      <text class="life-index-board__eyebrow">LIFE BOOK</text>
      <text class="life-index-board__title">生死簿查询</text>
      <text class="life-index-board__desc">阳寿、命格、风险和冻结状态一屏核对，批量划名重点盯防。</text>
      <view class="life-index-board__stats">
        <view class="life-index-board__stat">
          <text>{{ allItems.length }}</text>
          <text>记录</text>
        </view>
        <view class="life-index-board__stat">
          <text>{{ countByRisk('critical') }}</text>
          <text>严重</text>
        </view>
        <view class="life-index-board__stat">
          <text>{{ lockedCount }}</text>
          <text>冻结</text>
        </view>
      </view>
    </view>

    <up-search v-model="keyword" placeholder="搜索姓名" @search="load" @custom="load" />

    <view class="compact-tabs">
      <up-tabs :list="tabs" :current="current" @change="changeTab" />
    </view>

    <EmptyState v-if="items.length === 0" text="暂无生死簿记录" />

    <view v-else class="life-record-list">
      <view v-for="item in items" :key="item.id" class="life-record-card" @click="openDetail(item.id)">
        <view class="life-record-card__seal" :class="{ 'life-record-card__seal--locked': item.locked }">{{ item.locked ? '封' : '命' }}</view>
        <view class="life-record-card__head">
          <view class="life-record-card__title">
            <text>{{ item.soulName }}</text>
            <text>{{ item.fateLevel }} · {{ item.locked ? '已冻结' : '可流转' }}</text>
          </view>
          <StatusTag :value="item.riskFlag" />
        </view>
        <view class="life-record-card__ruler">
          <view class="life-record-card__expected" :style="{ width: lifespanPercent(item.expectedLifespan) }"></view>
          <view class="life-record-card__actual" :style="{ width: lifespanPercent(item.actualLifespan) }"></view>
        </view>
        <view class="life-record-card__grid">
          <view>
            <text class="life-record-card__label">原定阳寿</text>
            <text class="life-record-card__value">{{ item.expectedLifespan }}</text>
          </view>
          <view>
            <text class="life-record-card__label">当前阳寿</text>
            <text class="life-record-card__value">{{ item.actualLifespan }}</text>
          </view>
          <view>
            <text class="life-record-card__label">劫数</text>
            <text class="life-record-card__value">{{ item.calamityCount }}</text>
          </view>
          <view>
            <text class="life-record-card__label">更新</text>
            <text class="life-record-card__value">{{ displayTime(item.updatedAt) }}</text>
          </view>
        </view>
        <view class="life-record-card__note">{{ item.deathReason }}</view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import EmptyState from '../../components/EmptyState.vue'
import StatusTag from '../../components/StatusTag.vue'
import { getLifeBookRecords } from '../../api/lifeBook'
import { displayTime, showError } from '../../utils/format'

const keyword = ref('')
const current = ref(0)
const tabs = [
  { name: '全部', riskFlag: 'all' },
  { name: '正常', riskFlag: 'normal' },
  { name: '预警', riskFlag: 'warning' },
  { name: '严重', riskFlag: 'critical' }
]
const allItems = ref([])
const items = computed(() => allItems.value)
const lockedCount = computed(() => allItems.value.filter(item => item.locked).length)

function changeTab(event) {
  current.value = event.index
  load()
}

function openDetail(id) {
  uni.navigateTo({ url: `/pages/life-book/detail?id=${id}` })
}

function countByRisk(riskFlag) {
  return allItems.value.filter(item => item.riskFlag === riskFlag).length
}

function lifespanPercent(value) {
  return `${Math.min(100, Math.round((Number(value || 0) / 120) * 100))}%`
}

async function load() {
  try {
    allItems.value = await getLifeBookRecords({ keyword: keyword.value, riskFlag: tabs[current.value].riskFlag })
  } catch (error) {
    showError(error)
  }
}

onShow(load)
</script>

<style scoped>
.life-list-page {
  background:
    radial-gradient(circle at 10% 0%, rgba(180, 83, 9, 0.12), transparent 25%),
    linear-gradient(180deg, #fef3c7 0%, #f7f1e3 34%, #f7f1e3 100%);
}

.life-index-board {
  position: relative;
  overflow: hidden;
  padding: 26rpx;
  border: 2rpx solid rgba(120, 53, 15, 0.18);
  border-radius: 10rpx 30rpx 10rpx 30rpx;
  background:
    linear-gradient(90deg, rgba(120, 53, 15, 0.05) 1rpx, transparent 1rpx),
    linear-gradient(0deg, rgba(120, 53, 15, 0.05) 1rpx, transparent 1rpx),
    linear-gradient(135deg, #fff7ed, #fef3c7);
  background-size: 34rpx 34rpx;
  box-shadow: 0 18rpx 38rpx rgba(120, 53, 15, 0.12);
}

.life-index-board__seal {
  position: absolute;
  right: 24rpx;
  top: 24rpx;
  width: 118rpx;
  height: 118rpx;
  color: rgba(169, 50, 38, 0.22);
  border: 8rpx double rgba(169, 50, 38, 0.2);
  border-radius: 50%;
  font-size: 60rpx;
  font-weight: 900;
  line-height: 118rpx;
  text-align: center;
  transform: rotate(13deg);
}

.life-index-board__eyebrow,
.life-index-board__title,
.life-index-board__desc {
  position: relative;
  z-index: 1;
  display: block;
}

.life-index-board__eyebrow {
  color: #92400e;
  font-size: 20rpx;
  font-weight: 900;
  letter-spacing: 2rpx;
}

.life-index-board__title {
  margin-top: 6rpx;
  color: #3f2a12;
  font-size: 44rpx;
  font-weight: 900;
}

.life-index-board__desc {
  width: 76%;
  margin-top: 8rpx;
  color: #7c2d12;
  font-size: 24rpx;
  line-height: 1.44;
}

.life-index-board__stats {
  position: relative;
  z-index: 1;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10rpx;
  margin-top: 22rpx;
}

.life-index-board__stat {
  padding: 14rpx 12rpx;
  border: 1rpx solid rgba(120, 53, 15, 0.16);
  background: rgba(255, 250, 240, 0.68);
}

.life-index-board__stat text:first-child {
  display: block;
  color: #3f2a12;
  font-size: 30rpx;
  font-weight: 900;
}

.life-index-board__stat text:last-child {
  display: block;
  margin-top: 4rpx;
  color: #92400e;
  font-size: 20rpx;
}

.life-record-list {
  display: flex;
  flex-direction: column;
  gap: 14rpx;
}

.life-record-card {
  position: relative;
  overflow: hidden;
  padding: 20rpx;
  border: 1rpx solid rgba(120, 53, 15, 0.18);
  border-radius: 12rpx 24rpx 12rpx 24rpx;
  background:
    linear-gradient(90deg, rgba(120, 53, 15, 0.04) 1rpx, transparent 1rpx),
    #fffaf0;
  background-size: 30rpx 30rpx;
  box-shadow: 0 10rpx 28rpx rgba(120, 53, 15, 0.08);
}

.life-record-card__seal {
  position: absolute;
  right: 18rpx;
  bottom: 18rpx;
  width: 86rpx;
  height: 86rpx;
  color: rgba(180, 83, 9, 0.2);
  border: 5rpx solid rgba(180, 83, 9, 0.18);
  border-radius: 50%;
  font-size: 42rpx;
  font-weight: 900;
  line-height: 86rpx;
  text-align: center;
  transform: rotate(12deg);
}

.life-record-card__seal--locked {
  color: rgba(169, 50, 38, 0.24);
  border-color: rgba(169, 50, 38, 0.22);
}

.life-record-card__head {
  position: relative;
  z-index: 1;
  display: flex;
  justify-content: space-between;
  gap: 12rpx;
}

.life-record-card__title text:first-child {
  display: block;
  color: #3f2a12;
  font-size: 32rpx;
  font-weight: 900;
}

.life-record-card__title text:last-child {
  display: block;
  margin-top: 6rpx;
  color: #92400e;
  font-size: 22rpx;
}

.life-record-card__ruler {
  position: relative;
  z-index: 1;
  height: 24rpx;
  margin-top: 16rpx;
  overflow: hidden;
  border: 1rpx solid rgba(120, 53, 15, 0.16);
  border-radius: 999rpx;
  background: rgba(255, 250, 240, 0.8);
}

.life-record-card__expected,
.life-record-card__actual {
  position: absolute;
  left: 0;
  top: 0;
  height: 100%;
}

.life-record-card__expected {
  background: rgba(180, 83, 9, 0.24);
}

.life-record-card__actual {
  height: 10rpx;
  top: 7rpx;
  background: #a93226;
}

.life-record-card__grid {
  position: relative;
  z-index: 1;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10rpx;
  margin-top: 14rpx;
}

.life-record-card__grid view {
  min-width: 0;
  padding: 12rpx;
  border-radius: 14rpx;
  background: rgba(255, 250, 240, 0.72);
}

.life-record-card__label,
.life-record-card__value {
  display: block;
}

.life-record-card__label {
  color: #92400e;
  font-size: 20rpx;
}

.life-record-card__value {
  overflow: hidden;
  margin-top: 4rpx;
  color: #3f2a12;
  font-size: 24rpx;
  font-weight: 900;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.life-record-card__note {
  position: relative;
  z-index: 1;
  margin-top: 12rpx;
  padding: 12rpx;
  color: #7c2d12;
  border-radius: 14rpx;
  background: rgba(254, 243, 199, 0.62);
  font-size: 23rpx;
  line-height: 1.42;
}
</style>
