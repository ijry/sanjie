<template>
  <view class="page stack">
    <EmptyState v-if="items.length === 0" text="暂无审计日志" />
    <up-cell-group v-else>
      <up-cell v-for="item in items" :key="item.id" :title="item.action" :label="`${item.targetType} #${item.targetId} · ${item.note || item.createdAt}`" :value="`用户 ${item.actorId}`" />
    </up-cell-group>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import EmptyState from '../../components/EmptyState.vue'
import { getAuditLogs } from '../../api/audit'
import { showError } from '../../utils/format'

const items = ref([])

async function load() {
  try {
    items.value = await getAuditLogs()
  } catch (error) {
    showError(error)
  }
}

onShow(load)
</script>

