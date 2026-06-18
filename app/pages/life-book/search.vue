<template>
  <view class="page stack dense-page">
    <view class="module-hero">
      <text class="module-hero__eyebrow">LIFE BOOK</text>
      <text class="module-hero__title">生死簿查询</text>
      <text class="module-hero__desc">阳寿、命格、风险和冻结状态一屏核对，批量划名重点盯防。</text>
      <view class="module-hero__stats">
        <view class="module-hero__stat">
          <text>{{ allItems.length }}</text>
          <text>记录</text>
        </view>
        <view class="module-hero__stat">
          <text>{{ countByRisk('critical') }}</text>
          <text>严重</text>
        </view>
        <view class="module-hero__stat">
          <text>{{ lockedCount }}</text>
          <text>冻结</text>
        </view>
      </view>
    </view>

    <up-search v-model="keyword" placeholder="搜索姓名" @search="load" @custom="load" />

    <view class="compact-tabs">
      <up-tabs :list="tabs" :current="current" @change="changeTab" />
    </view>

    <EmptyState v-if="items.length === 0" text="暂无生死簿记录" />

    <view v-else class="card-list">
      <view v-for="item in items" :key="item.id" class="work-card" @click="openDetail(item.id)">
        <view class="work-card__head">
          <view class="work-card__title">
            <text>{{ item.soulName }}</text>
            <text class="work-card__subtitle">{{ item.fateLevel }} · {{ item.locked ? '已冻结' : '可流转' }}</text>
          </view>
          <StatusTag :value="item.riskFlag" />
        </view>
        <view class="work-card__grid">
          <view class="work-card__field">
            <text class="work-card__field-label">原定阳寿</text>
            <text class="work-card__field-value">{{ item.expectedLifespan }}</text>
          </view>
          <view class="work-card__field">
            <text class="work-card__field-label">当前阳寿</text>
            <text class="work-card__field-value">{{ item.actualLifespan }}</text>
          </view>
          <view class="work-card__field">
            <text class="work-card__field-label">劫数</text>
            <text class="work-card__field-value">{{ item.calamityCount }}</text>
          </view>
          <view class="work-card__field">
            <text class="work-card__field-label">更新</text>
            <text class="work-card__field-value">{{ displayTime(item.updatedAt) }}</text>
          </view>
        </view>
        <view class="work-card__note">{{ item.deathReason }}</view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import EmptyState from '../../components/EmptyState.vue'
import StatusTag from '../../components/StatusTag.vue'
import { getLifeBookRecords } from '../../api/lifeBook'
import { displayTime, showError } from '../../utils/format'

const keyword = ref('')
const current = ref(0)
const tabs = [
  { name: '全部', riskFlag: 'all' },
  { name: '正常', riskFlag: 'normal' },
  { name: '预警', riskFlag: 'warning' },
  { name: '严重', riskFlag: 'critical' }
]
const allItems = ref([])
const items = computed(() => allItems.value)
const lockedCount = computed(() => allItems.value.filter(item => item.locked).length)

function changeTab(event) {
  current.value = event.index
  load()
}

function openDetail(id) {
  uni.navigateTo({ url: `/pages/life-book/detail?id=${id}` })
}

function countByRisk(riskFlag) {
  return allItems.value.filter(item => item.riskFlag === riskFlag).length
}

async function load() {
  try {
    allItems.value = await getLifeBookRecords({ keyword: keyword.value, riskFlag: tabs[current.value].riskFlag })
  } catch (error) {
    showError(error)
  }
}

onShow(load)
</script>
