<template>
  <view class="page stack dense-page hell-list-page">
    <view class="hell-ops-board">
      <view class="hell-ops-board__flame"></view>
      <text class="hell-ops-board__eyebrow">HELL OPERATIONS</text>
      <text class="hell-ops-board__title">十八层地狱</text>
      <text class="hell-ops-board__desc">楼层容量、设备状态和刑期复核统一调度，高负载楼层可一键分流。</text>
      <view class="hell-ops-board__stats">
        <view class="hell-ops-board__stat">
          <text>{{ floors.length }}</text>
          <text>楼层</text>
        </view>
        <view class="hell-ops-board__stat">
          <text>{{ criticalCount }}</text>
          <text>严重</text>
        </view>
        <view class="hell-ops-board__stat">
          <text>{{ totalOccupancy }}</text>
          <text>在押</text>
        </view>
      </view>
    </view>

    <view class="hell-panel stack">
      <SectionTitle title="楼层负载" />
      <view v-for="floor in floors" :key="floor.id" class="floor-monitor">
        <view>
          <text class="floor-monitor__name">{{ floor.floorNo }} 层 · {{ floor.name }}</text>
          <text class="floor-monitor__meta">{{ floor.occupancy }}/{{ floor.capacity }} · {{ floor.equipmentStatus }} · {{ percent(floor) }}%</text>
          <view class="floor-monitor__bar">
            <view class="floor-monitor__bar-inner" :style="{ width: `${percent(floor)}%` }"></view>
          </view>
        </view>
        <view class="floor-monitor__ops">
          <StatusTag :value="floor.loadLevel" />
          <up-button v-if="floor.loadLevel === 'critical'" size="mini" type="warning" text="分流" @click="dispatch(floor)" />
        </view>
      </view>
    </view>

    <view class="hell-panel">
      <SectionTitle title="刑期记录" />
      <EmptyState v-if="sentences.length === 0" text="暂无刑期记录" />
      <view v-else class="sentence-console-list">
        <view v-for="item in sentences" :key="item.id" class="sentence-console-card" @click="openSentence(item.id)">
          <view class="sentence-console-card__head">
            <view class="sentence-console-card__title">
              <text>{{ item.soulName }}</text>
              <text>{{ item.floorName }} · {{ item.crimeType }}</text>
            </view>
            <StatusTag :value="item.reviewStatus" />
          </view>
          <view class="sentence-console-card__gauges">
            <view>
              <text class="sentence-console-card__value">{{ item.sentenceDays }}</text>
              <text class="sentence-console-card__label">刑期 / 天</text>
            </view>
            <view>
              <text class="sentence-console-card__value">{{ item.painLevel }}/10</text>
              <text class="sentence-console-card__label">痛感</text>
            </view>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import SectionTitle from '../../components/SectionTitle.vue'
import StatusTag from '../../components/StatusTag.vue'
import EmptyState from '../../components/EmptyState.vue'
import { dispatchHellFloor, getHellFloors, getHellSentences } from '../../api/hell'
import { showError, showSuccess } from '../../utils/format'

const floors = ref([])
const sentences = ref([])
const criticalCount = computed(() => floors.value.filter(item => item.loadLevel === 'critical').length)
const totalOccupancy = computed(() => floors.value.reduce((total, item) => total + Number(item.occupancy || 0), 0))

async function load() {
  try {
    floors.value = await getHellFloors()
    sentences.value = await getHellSentences()
  } catch (error) {
    showError(error)
  }
}

function openSentence(id) {
  uni.navigateTo({ url: `/pages/hell/sentence-detail?id=${id}` })
}

async function dispatch(floor) {
  const target = floors.value.find(item => item.id !== floor.id)
  if (!target) return
  try {
    await dispatchHellFloor(floor.id, { targetFloorId: target.id, amount: 10 })
    showSuccess('已分流')
    await load()
  } catch (error) {
    showError(error)
  }
}

function percent(floor) {
  if (!floor.capacity) return 0
  return Math.min(100, Math.round((Number(floor.occupancy || 0) / Number(floor.capacity)) * 100))
}

onShow(load)
</script>

<style scoped>
.hell-list-page {
  background:
    radial-gradient(circle at 82% 0%, rgba(239, 68, 68, 0.2), transparent 26%),
    linear-gradient(180deg, #0f0a08 0%, #2b0f0b 30%, #f7f1e3 30%, #f7f1e3 100%);
}

.hell-ops-board {
  position: relative;
  overflow: hidden;
  padding: 26rpx;
  color: #fff7ed;
  border: 1rpx solid rgba(248, 113, 113, 0.28);
  border-radius: 28rpx;
  background:
    linear-gradient(90deg, rgba(248, 113, 113, 0.08) 1rpx, transparent 1rpx),
    linear-gradient(0deg, rgba(248, 113, 113, 0.08) 1rpx, transparent 1rpx),
    linear-gradient(135deg, #120806, #450a0a 56%, #111827);
  background-size: 36rpx 36rpx;
  box-shadow: 0 22rpx 46rpx rgba(69, 10, 10, 0.26);
}

.hell-ops-board__flame {
  position: absolute;
  right: -50rpx;
  top: -70rpx;
  width: 230rpx;
  height: 230rpx;
  border-radius: 50%;
  background: rgba(249, 115, 22, 0.26);
}

.hell-ops-board__eyebrow,
.hell-ops-board__title,
.hell-ops-board__desc {
  position: relative;
  z-index: 1;
  display: block;
}

.hell-ops-board__eyebrow {
  color: #fb923c;
  font-size: 20rpx;
  font-weight: 900;
  letter-spacing: 2rpx;
}

.hell-ops-board__title {
  margin-top: 6rpx;
  font-size: 44rpx;
  font-weight: 900;
}

.hell-ops-board__desc {
  width: 76%;
  margin-top: 8rpx;
  color: rgba(255, 247, 237, 0.72);
  font-size: 24rpx;
  line-height: 1.42;
}

.hell-ops-board__stats {
  position: relative;
  z-index: 1;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10rpx;
  margin-top: 22rpx;
}

.hell-ops-board__stat {
  padding: 14rpx 12rpx;
  border: 1rpx solid rgba(251, 146, 60, 0.24);
  border-radius: 14rpx;
  background: rgba(15, 23, 42, 0.52);
}

.hell-ops-board__stat text:first-child {
  display: block;
  color: #ffedd5;
  font-size: 30rpx;
  font-weight: 900;
}

.hell-ops-board__stat text:last-child {
  display: block;
  margin-top: 4rpx;
  color: rgba(255, 237, 213, 0.64);
  font-size: 20rpx;
}

.hell-panel {
  padding: 22rpx;
  border: 1rpx solid rgba(127, 29, 29, 0.18);
  border-radius: 20rpx;
  background: linear-gradient(180deg, #fffaf0, #fff7ed);
}

.floor-monitor {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16rpx;
  padding: 18rpx 0;
  border-bottom: 1rpx solid rgba(127, 29, 29, 0.16);
}

.floor-monitor__name {
  display: block;
  color: #431407;
  font-weight: 900;
}

.floor-monitor__meta {
  display: block;
  margin-top: 6rpx;
  color: #7f1d1d;
  font-size: 24rpx;
}

.floor-monitor__bar {
  width: 300rpx;
  height: 10rpx;
  margin-top: 10rpx;
  overflow: hidden;
  border-radius: 999rpx;
  background: rgba(229, 216, 189, 0.75);
}

.floor-monitor__bar-inner {
  height: 100%;
  border-radius: 999rpx;
  background: linear-gradient(90deg, #fb923c, #a93226);
}

.floor-monitor__ops {
  display: flex;
  flex-direction: column;
  gap: 10rpx;
  align-items: flex-end;
}

.sentence-console-list {
  display: flex;
  flex-direction: column;
  gap: 12rpx;
  margin-top: 14rpx;
}

.sentence-console-card {
  padding: 18rpx;
  color: #fff7ed;
  border-radius: 18rpx;
  background: linear-gradient(135deg, #1f0f0f, #7f1d1d);
  box-shadow: 0 10rpx 26rpx rgba(127, 29, 29, 0.16);
}

.sentence-console-card__head {
  display: flex;
  justify-content: space-between;
  gap: 12rpx;
}

.sentence-console-card__title text:first-child {
  display: block;
  font-size: 30rpx;
  font-weight: 900;
}

.sentence-console-card__title text:last-child {
  display: block;
  margin-top: 6rpx;
  color: rgba(255, 237, 213, 0.68);
  font-size: 22rpx;
}

.sentence-console-card__gauges {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10rpx;
  margin-top: 14rpx;
}

.sentence-console-card__gauges view {
  padding: 12rpx;
  border: 1rpx solid rgba(251, 146, 60, 0.22);
  border-radius: 14rpx;
  background: rgba(15, 23, 42, 0.38);
}

.sentence-console-card__value,
.sentence-console-card__label {
  display: block;
}

.sentence-console-card__value {
  font-size: 28rpx;
  font-weight: 900;
}

.sentence-console-card__label {
  margin-top: 4rpx;
  color: rgba(255, 237, 213, 0.66);
  font-size: 20rpx;
}
</style>
