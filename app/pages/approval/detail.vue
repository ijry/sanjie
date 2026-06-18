<template>
  <view class="page detail-page">
    <view class="panel stack" v-if="item">
      <view class="detail-title">
        <text>{{ item.title }}</text>
        <StatusTag :value="item.status" />
      </view>
      <up-cell title="类型" :value="item.type" />
      <up-cell title="目标 ID" :value="String(item.targetId)" />
      <up-cell title="申请人" :value="String(item.applicantId)" />
      <up-cell title="审批人" :value="String(item.approverId || '-')" />
      <up-cell title="理由" :value="item.reason || '-'" />
      <up-cell title="结果" :value="item.resultNote || '-'" />
    </view>

    <FixedActionBar v-if="item && ['pending', 'transferred'].includes(item.status)">
      <up-button type="success" text="同意" @click="approve" />
      <up-button type="error" text="驳回" @click="reject" />
      <up-button v-if="item.status === 'pending'" type="primary" text="转交" @click="transfer" />
    </FixedActionBar>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import FixedActionBar from '../../components/FixedActionBar.vue'
import StatusTag from '../../components/StatusTag.vue'
import { approveApproval, getApproval, rejectApproval, transferApproval } from '../../api/approvals'
import { showError, showSuccess } from '../../utils/format'

const id = ref('')
const item = ref(null)

onLoad(query => {
  id.value = query.id
})

async function load() {
  if (!id.value) return
  try {
    item.value = await getApproval(id.value)
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

const approve = () => run(() => approveApproval(id.value, { note: '移动端审批通过' }), '已同意')
const reject = () => run(() => rejectApproval(id.value, { note: '移动端审批驳回' }), '已驳回')
const transfer = () => run(() => transferApproval(id.value, { note: '转交上级复核' }), '已转交')

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

