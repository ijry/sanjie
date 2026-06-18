<template>
  <view class="page detail-page">
    <view class="panel stack" v-if="item">
      <view class="detail-title">
        <text>{{ item.soulName }}</text>
        <StatusTag :value="item.reviewStatus" />
      </view>
      <up-cell title="楼层" :value="item.floorName" />
      <up-cell title="罪行" :value="item.crimeType" />
      <up-cell title="刑期" :value="`${item.sentenceDays} 天`" />
      <up-cell title="痛感" :value="String(item.painLevel)" />
      <up-cell title="设备" :value="item.equipmentId" />
    </view>
    <FixedActionBar v-if="item && item.reviewStatus === 'none'">
      <up-button type="warning" text="发起复核" @click="review" />
    </FixedActionBar>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import FixedActionBar from '../../components/FixedActionBar.vue'
import StatusTag from '../../components/StatusTag.vue'
import { getHellSentence, reviewHellSentence } from '../../api/hell'
import { showError, showSuccess } from '../../utils/format'

const id = ref('')
const item = ref(null)

onLoad(query => {
  id.value = query.id
})

async function load() {
  if (!id.value) return
  try {
    item.value = await getHellSentence(id.value)
  } catch (error) {
    showError(error)
  }
}

async function review() {
  try {
    await reviewHellSentence(id.value, { note: '移动端发起刑期复核' })
    showSuccess('已发起')
    await load()
  } catch (error) {
    showError(error)
  }
}

onShow(load)
</script>

<style scoped>
.detail-page {
  padding-bottom: 120rpx;
}

.detail-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 38rpx;
  font-weight: 800;
}
</style>

