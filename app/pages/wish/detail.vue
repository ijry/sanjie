<template>
  <view class="page detail-page">
    <view class="panel stack" v-if="item">
      <view class="detail-title">
        <text>{{ item.title }}</text>
        <StatusTag :value="item.status" />
      </view>
      <up-cell title="许愿人" :value="item.requesterName" />
      <up-cell title="类型" :value="item.wishType" />
      <up-cell title="香火" :value="String(item.incenseAmount)" />
      <up-cell title="功德" :value="String(item.meritScore)" />
      <up-cell title="转派" :value="item.assignedDeity || '-'" />
      <up-cell title="结果" :value="item.resultNote || '-'" />
    </view>

    <FixedActionBar v-if="item && ['pending', 'routed'].includes(item.status)">
      <up-button v-if="item.status === 'pending'" type="primary" text="转财神" @click="route" />
      <up-button type="success" text="办结" @click="resolve" />
      <up-button type="error" text="驳回" @click="reject" />
    </FixedActionBar>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import FixedActionBar from '../../components/FixedActionBar.vue'
import StatusTag from '../../components/StatusTag.vue'
import { getWish, rejectWish, resolveWish, routeWish } from '../../api/wishes'
import { showError, showSuccess } from '../../utils/format'

const id = ref('')
const item = ref(null)

onLoad(query => {
  id.value = query.id
})

async function load() {
  if (!id.value) return
  try {
    item.value = await getWish(id.value)
  } catch (error) {
    showError(error)
  }
}

async function run(action, text) {
  try {
    await action()
    showSuccess(text)
    await load()
  } catch (error) {
    showError(error)
  }
}

const route = () => run(() => routeWish(id.value, { assignedDeity: '财神', note: '转财运组评估' }), '已转派')
const resolve = () => run(() => resolveWish(id.value, { note: '建议先自力更生' }), '已办结')
const reject = () => run(() => rejectWish(id.value, { note: '愿望不符合天规' }), '已驳回')

onShow(load)
</script>

<style scoped>
.detail-page {
  padding-bottom: 140rpx;
}

.detail-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 34rpx;
  font-weight: 800;
}
</style>

