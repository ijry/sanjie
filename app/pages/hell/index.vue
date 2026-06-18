<template>
  <view class="page stack dense-page">
    <view class="module-hero">
      <text class="module-hero__eyebrow">HELL OPERATIONS</text>
      <text class="module-hero__title">十八层地狱</text>
      <text class="module-hero__desc">楼层容量、设备状态和刑期复核统一调度，高负载楼层可一键分流。</text>
      <view class="module-hero__stats">
        <view class="module-hero__stat">
          <text>{{ floors.length }}</text>
          <text>楼层</text>
        </view>
        <view class="module-hero__stat">
          <text>{{ criticalCount }}</text>
          <text>严重</text>
        </view>
        <view class="module-hero__stat">
          <text>{{ totalOccupancy }}</text>
          <text>在押</text>
        </view>
      </view>
    </view>

    <view class="detail-card stack">
      <SectionTitle title="楼层负载" />
      <view v-for="floor in floors" :key="floor.id" class="floor">
        <view>
          <text class="floor__name">{{ floor.floorNo }} 层 · {{ floor.name }}</text>
          <text class="floor__meta">{{ floor.occupancy }}/{{ floor.capacity }} · {{ floor.equipmentStatus }} · {{ percent(floor) }}%</text>
          <view class="floor__bar">
            <view class="floor__bar-inner" :style="{ width: `${percent(floor)}%` }"></view>
          </view>
        </view>
        <view class="floor__ops">
          <StatusTag :value="floor.loadLevel" />
          <up-button v-if="floor.loadLevel === 'critical'" size="mini" type="warning" text="分流" @click="dispatch(floor)" />
        </view>
      </view>
    </view>

    <view class="detail-card">
      <SectionTitle title="刑期记录" />
      <EmptyState v-if="sentences.length === 0" text="暂无刑期记录" />
      <view v-else class="card-list">
        <view v-for="item in sentences" :key="item.id" class="work-card" @click="openSentence(item.id)">
          <view class="work-card__head">
            <view class="work-card__title">
              <text>{{ item.soulName }}</text>
              <text class="work-card__subtitle">{{ item.floorName }} · {{ item.crimeType }}</text>
            </view>
            <StatusTag :value="item.reviewStatus" />
          </view>
          <view class="work-card__grid">
            <view class="work-card__field">
              <text class="work-card__field-label">刑期</text>
              <text class="work-card__field-value">{{ item.sentenceDays }} 天</text>
            </view>
            <view class="work-card__field">
              <text class="work-card__field-label">痛感</text>
              <text class="work-card__field-value">{{ item.painLevel }}/10</text>
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
.floor {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16rpx;
  padding: 18rpx 0;
  border-bottom: 1rpx solid #e5d8bd;
}

.floor__name {
  display: block;
  font-weight: 700;
}

.floor__meta {
  display: block;
  margin-top: 6rpx;
  color: #756f63;
  font-size: 24rpx;
}

.floor__bar {
  width: 300rpx;
  height: 10rpx;
  margin-top: 10rpx;
  overflow: hidden;
  border-radius: 999rpx;
  background: rgba(229, 216, 189, 0.75);
}

.floor__bar-inner {
  height: 100%;
  border-radius: 999rpx;
  background: linear-gradient(90deg, #b8860b, #a93226);
}

.floor__ops {
  display: flex;
  flex-direction: column;
  gap: 10rpx;
  align-items: flex-end;
}
</style>
