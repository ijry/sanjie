<template>
  <view class="page stack">
    <view class="panel stack">
      <SectionTitle title="楼层负载" />
      <view v-for="floor in floors" :key="floor.id" class="floor">
        <view>
          <text class="floor__name">{{ floor.floorNo }} 层 · {{ floor.name }}</text>
          <text class="floor__meta">{{ floor.occupancy }}/{{ floor.capacity }} · {{ floor.equipmentStatus }}</text>
        </view>
        <StatusTag :value="floor.loadLevel" />
        <up-button v-if="floor.loadLevel === 'critical'" size="mini" type="warning" text="分流" @click="dispatch(floor)" />
      </view>
    </view>

    <view class="panel">
      <SectionTitle title="刑期记录" />
      <EmptyState v-if="sentences.length === 0" text="暂无刑期记录" />
      <up-cell-group v-else :border="false">
        <up-cell v-for="item in sentences" :key="item.id" :title="item.soulName" :label="`${item.floorName} · ${item.crimeType}`" is-link @click="openSentence(item.id)">
          <template #value>
            <StatusTag :value="item.reviewStatus" />
          </template>
        </up-cell>
      </up-cell-group>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import SectionTitle from '../../components/SectionTitle.vue'
import StatusTag from '../../components/StatusTag.vue'
import EmptyState from '../../components/EmptyState.vue'
import { dispatchHellFloor, getHellFloors, getHellSentences } from '../../api/hell'
import { showError, showSuccess } from '../../utils/format'

const floors = ref([])
const sentences = ref([])

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
</style>

