<template>
  <view class="page stack">
    <up-search v-model="keyword" placeholder="搜索姓名" @search="load" @custom="load" />
    <up-tabs :list="tabs" :current="current" @change="changeTab" />
    <EmptyState v-if="items.length === 0" text="暂无生死簿记录" />
    <up-cell-group v-else>
      <up-cell v-for="item in items" :key="item.id" :title="item.soulName" :label="`阳寿 ${item.actualLifespan} · ${item.fateLevel}`" is-link @click="openDetail(item.id)">
        <template #value>
          <StatusTag :value="item.riskFlag" />
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
import { getLifeBookRecords } from '../../api/lifeBook'
import { showError } from '../../utils/format'

const keyword = ref('')
const current = ref(0)
const tabs = [
  { name: '全部', riskFlag: 'all' },
  { name: '正常', riskFlag: 'normal' },
  { name: '预警', riskFlag: 'warning' },
  { name: '严重', riskFlag: 'critical' }
]
const items = ref([])

function changeTab(event) {
  current.value = event.index
  load()
}

function openDetail(id) {
  uni.navigateTo({ url: `/pages/life-book/detail?id=${id}` })
}

async function load() {
  try {
    items.value = await getLifeBookRecords({ keyword: keyword.value, riskFlag: tabs[current.value].riskFlag })
  } catch (error) {
    showError(error)
  }
}

onShow(load)
</script>

