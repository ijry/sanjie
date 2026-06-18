<template>
  <view class="page stack">
    <up-tabs :list="tabs" :current="current" @change="changeTab" />
    <EmptyState v-if="items.length === 0" text="暂无愿望工单" />
    <up-cell-group v-else>
      <up-cell v-for="item in items" :key="item.id" :title="item.title" :label="`${item.requesterName} · 香火 ${item.incenseAmount}`" is-link @click="openDetail(item.id)">
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
import { getWishes } from '../../api/wishes'
import { showError } from '../../utils/format'

const tabs = [
  { name: '全部', status: '' },
  { name: '待处理', status: 'pending' },
  { name: '已转派', status: 'routed' },
  { name: '已办结', status: 'resolved' },
  { name: '已驳回', status: 'rejected' }
]
const current = ref(0)
const items = ref([])

function changeTab(event) {
  current.value = event.index
  load()
}

function openDetail(id) {
  uni.navigateTo({ url: `/pages/wish/detail?id=${id}` })
}

async function load() {
  try {
    items.value = await getWishes({ status: tabs[current.value].status })
  } catch (error) {
    showError(error)
  }
}

onShow(load)
</script>

