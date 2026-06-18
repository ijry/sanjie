<template>
  <view class="page detail-page">
    <view class="panel stack" v-if="item">
      <view class="detail-title">
        <text>{{ item.soulName }}</text>
        <StatusTag :value="item.status" />
      </view>
      <up-cell title="目标道别" :value="item.targetRealm" />
      <up-cell title="优先级" :value="String(item.priority)" />
      <up-cell title="胎位" :value="item.quotaType || '-'" />
      <up-cell title="孟婆汤" :value="item.soupStatus" />
      <up-cell title="复核备注" :value="item.reviewNote || '-'" />
    </view>

    <FixedActionBar v-if="item">
      <up-button v-if="['queued', 'reviewing'].includes(item.status)" type="success" text="通过" @click="approve" />
      <up-button v-if="['queued', 'reviewing', 'approved'].includes(item.status)" type="error" text="驳回" @click="reject" />
      <up-button v-if="['queued', 'approved'].includes(item.status)" type="warning" text="复核" @click="review" />
      <up-button v-if="['approved', 'pending_soup'].includes(item.status)" type="primary" text="分胎位" @click="showQuota = true" />
      <up-button v-if="item.status === 'pending_soup' && item.soupStatus === 'pending'" type="primary" text="发汤" @click="issueSoup" />
      <up-button v-if="item.status === 'pending_soup' && item.soupStatus === 'issued'" type="success" text="完成" @click="complete" />
    </FixedActionBar>

    <up-action-sheet :show="showQuota" :actions="quotaActions" title="选择胎位" @close="showQuota = false" @select="assignQuota" />
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import FixedActionBar from '../../components/FixedActionBar.vue'
import StatusTag from '../../components/StatusTag.vue'
import { approveReincarnation, assignReincarnationQuota, completeReincarnation, getReincarnation, issueReincarnationSoup, rejectReincarnation, reviewReincarnation } from '../../api/reincarnation'
import { getSoupInventory } from '../../api/soup'
import { showError, showSuccess } from '../../utils/format'

const id = ref('')
const item = ref(null)
const showQuota = ref(false)
const quotaActions = ['富贵胎', '普通胎', '清贫胎', '天才胎', '困难模式胎'].map(name => ({ name }))

onLoad(query => {
  id.value = query.id
})

async function load() {
  if (!id.value) return
  try {
    item.value = await getReincarnation(id.value)
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

const approve = () => run(() => approveReincarnation(id.value, { note: '移动端通过' }), '已通过')
const reject = () => run(() => rejectReincarnation(id.value, { note: '移动端驳回' }), '已驳回')
const review = () => run(() => reviewReincarnation(id.value, { note: '进入复核' }), '已复核')
const complete = () => run(() => completeReincarnation(id.value, { note: '已入轮回' }), '已完成')

async function assignQuota(action) {
  showQuota.value = false
  await run(() => assignReincarnationQuota(id.value, { quotaType: action.name }), '已分配')
}

async function issueSoup() {
  try {
    const inventory = await getSoupInventory()
    const first = inventory[0]
    if (!first) throw new Error('暂无孟婆汤库存')
    await issueReincarnationSoup(id.value, { inventoryId: first.id, dose: 1, memoryAfter: 0 })
    showSuccess('已发汤')
    await load()
  } catch (error) {
    showError(error)
  }
}

onShow(load)
</script>

<style scoped>
.detail-page {
  padding-bottom: 170rpx;
}

.detail-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 38rpx;
  font-weight: 800;
}
</style>

