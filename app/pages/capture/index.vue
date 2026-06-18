<template>
  <view class="page stack">
    <up-tabs :list="tabs" :current="current" @change="changeTab" />
    <EmptyState v-if="items.length === 0" text="暂无勾魂任务" />
    <up-cell-group v-else>
      <up-cell v-for="item in items" :key="item.id" :title="item.soulName" :label="`${item.location} · ${displayTime(item.scheduledTime)}`" is-link @click="openDetail(item.id)">
        <template #value>
          <StatusTag :value="item.status" />
        </template>
      </up-cell>
    </up-cell-group>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import EmptyState from '../../components/EmptyState.vue'
import StatusTag from '../../components/StatusTag.vue'
import { getCaptureTasks } from '../../api/capture'
import { displayTime, showError } from '../../utils/format'

const tabs = [
  { name: '全部', status: '' },
  { name: '待执行', status: 'pending' },
  { name: '执行中', status: 'running' },
  { name: '异常', status: 'exception' },
  { name: '已完成', status: 'completed' }
]
const current = ref(0)
const items = ref([])

function changeTab(event) {
  current.value = event.index
  load()
}

function openDetail(id) {
  uni.navigateTo({ url: `/pages/capture/detail?id=${id}` })
}

async function load() {
  try {
    items.value = await getCaptureTasks({ status: tabs[current.value].status })
  } catch (error) {
    showError(error)
  }
}

onShow(load)
</script>

