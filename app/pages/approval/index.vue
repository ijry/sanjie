<template>
  <view class="page stack">
    <up-tabs :list="tabs" :current="current" @change="changeTab" />
    <EmptyState v-if="items.length === 0" text="暂无审批" />
    <up-cell-group v-else>
      <up-cell v-for="item in items" :key="item.id" :title="item.title" :label="`${item.type} · ${item.reason}`" is-link @click="openDetail(item.id)">
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
import { getApprovals } from '../../api/approvals'
import { showError } from '../../utils/format'

const tabs = [
  { name: '待审', status: 'pending' },
  { name: '转交', status: 'transferred' },
  { name: '通过', status: 'approved' },
  { name: '驳回', status: 'rejected' }
]
const current = ref(0)
const items = ref([])

function changeTab(event) {
  current.value = event.index
  load()
}

function openDetail(id) {
  uni.navigateTo({ url: `/pages/approval/detail?id=${id}` })
}

async function load() {
  try {
    items.value = await getApprovals({ status: tabs[current.value].status })
  } catch (error) {
    showError(error)
  }
}

onShow(load)
</script>

