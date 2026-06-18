<template>
  <view class="page detail-page">
    <view class="panel stack" v-if="item">
      <view class="detail-title">
        <text>{{ item.soulName }}</text>
        <StatusTag :value="item.riskFlag" />
      </view>
      <up-cell title="原定阳寿" :value="String(item.expectedLifespan)" />
      <up-cell title="当前阳寿" :value="String(item.actualLifespan)" />
      <up-cell title="命格" :value="item.fateLevel" />
      <up-cell title="劫数" :value="String(item.calamityCount)" />
      <up-cell title="死因" :value="item.deathReason" />
      <up-cell title="冻结" :value="item.locked ? '是' : '否'" />
    </view>

    <view class="panel stack">
      <SectionTitle title="审计记录" />
      <EmptyState v-if="logs.length === 0" text="暂无审计记录" />
      <up-cell v-for="log in logs" :key="log.id" :title="log.action" :label="log.note || log.createdAt" />
    </view>

    <FixedActionBar v-if="item">
      <up-button v-if="!item.locked" type="warning" text="冻结" @click="freeze" />
      <up-button type="primary" text="申请改寿" @click="changeLife" />
    </FixedActionBar>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import FixedActionBar from '../../components/FixedActionBar.vue'
import SectionTitle from '../../components/SectionTitle.vue'
import StatusTag from '../../components/StatusTag.vue'
import EmptyState from '../../components/EmptyState.vue'
import { changeLifespan, freezeLifeBookRecord, getLifeBookAuditLogs, getLifeBookRecord } from '../../api/lifeBook'
import { showError, showSuccess } from '../../utils/format'

const id = ref('')
const item = ref(null)
const logs = ref([])

onLoad(query => {
  id.value = query.id
})

async function load() {
  if (!id.value) return
  try {
    item.value = await getLifeBookRecord(id.value)
    logs.value = await getLifeBookAuditLogs(id.value)
  } catch (error) {
    showError(error)
  }
}

async function freeze() {
  try {
    await freezeLifeBookRecord(id.value)
    showSuccess('已冻结')
    await load()
  } catch (error) {
    showError(error)
  }
}

function changeLife() {
  uni.showModal({
    title: '寿命变更',
    content: '将申请把当前阳寿调整为 88 岁',
    success: async res => {
      if (!res.confirm) return
      try {
        await changeLifespan(id.value, { newLifespan: 88, reason: '移动端发起寿命变更' })
        showSuccess('已发起审批')
        await load()
      } catch (error) {
        showError(error)
      }
    }
  })
}

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
  font-size: 38rpx;
  font-weight: 800;
}
</style>

