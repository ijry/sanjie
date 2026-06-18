<template>
  <view class="page stack">
    <up-tabs :list="tabs" :current="current" @change="changeTab" />
    <EmptyState v-if="items.length === 0" text="暂无投胎单" />
    <up-cell-group v-else>
      <up-cell v-for="item in items" :key="item.id" :title="item.soulName" :label="`${item.targetRealm} · 优先级 ${item.priority} · ${item.quotaType || '未分配'}`" is-link @click="openDetail(item.id)">
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
import { getReincarnations } from '../../api/reincarnation'
import { showError } from '../../utils/format'

const tabs = [
  { name: '全部', status: '' },
  { name: '排队', status: 'queued' },
  { name: '复核', status: 'reviewing' },
  { name: '待喝汤', status: 'pending_soup' },
  { name: '完成', status: 'completed' },
  { name: '驳回', status: 'rejected' }
]
const current = ref(0)
const items = ref([])

function changeTab(event) {
  current.value = event.index
  load()
}

function openDetail(id) {
  uni.navigateTo({ url: `/pages/reincarnation/detail?id=${id}` })
}

async function load() {
  try {
    items.value = await getReincarnations({ status: tabs[current.value].status })
  } catch (error) {
    showError(error)
  }
}

onShow(load)
</script>

