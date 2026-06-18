<template>
  <view class="page stack detail-page" v-if="item">
    <view class="detail-hero">
      <view class="detail-hero__top">
        <view>
          <text class="detail-hero__eyebrow">HELL SENTENCE #{{ item.id }}</text>
          <text class="detail-hero__title">{{ item.soulName }}</text>
        </view>
        <StatusTag :value="item.reviewStatus" />
      </view>
      <text class="detail-hero__desc">{{ item.floorName }} · {{ item.crimeType }} · {{ reviewText }}</text>
      <view class="detail-hero__stats">
        <view class="detail-stat">
          <text class="detail-stat__label">刑期</text>
          <text class="detail-stat__value">{{ item.sentenceDays }} 天</text>
        </view>
        <view class="detail-stat">
          <text class="detail-stat__label">痛感</text>
          <text class="detail-stat__value">{{ item.painLevel }}/10</text>
        </view>
        <view class="detail-stat">
          <text class="detail-stat__label">设备</text>
          <text class="detail-stat__value">{{ item.equipmentId }}</text>
        </view>
      </view>
    </view>

    <view class="detail-card">
      <text class="detail-card__title">刑期档案</text>
      <view class="detail-grid">
        <view class="detail-field">
          <text class="detail-field__label">魂魄 ID</text>
          <text class="detail-field__value">#{{ item.soulId }}</text>
        </view>
        <view class="detail-field">
          <text class="detail-field__label">楼层 ID</text>
          <text class="detail-field__value">#{{ item.floorId }}</text>
        </view>
        <view class="detail-field">
          <text class="detail-field__label">执行楼层</text>
          <text class="detail-field__value">{{ item.floorName }}</text>
        </view>
        <view class="detail-field">
          <text class="detail-field__label">复核状态</text>
          <text class="detail-field__value">{{ reviewText }}</text>
        </view>
        <view class="detail-field">
          <text class="detail-field__label">创建时间</text>
          <text class="detail-field__value">{{ displayTime(item.createdAt) }}</text>
        </view>
        <view class="detail-field">
          <text class="detail-field__label">更新时间</text>
          <text class="detail-field__value">{{ displayTime(item.updatedAt) }}</text>
        </view>
        <view class="detail-field detail-field--wide">
          <text class="detail-field__label">登记罪行</text>
          <text class="detail-field__value">{{ item.crimeType }}</text>
        </view>
      </view>
    </view>

    <view class="detail-card">
      <text class="detail-card__title">执行流程</text>
      <view class="process-list">
        <view v-for="step in processSteps" :key="step.title" class="process-row">
          <text class="process-row__dot">{{ step.index }}</text>
          <view class="process-row__main">
            <text class="process-row__title">{{ step.title }}</text>
            <text class="process-row__desc">{{ step.desc }}</text>
          </view>
        </view>
      </view>
    </view>

    <view class="detail-card">
      <text class="detail-card__title">复核建议</text>
      <view class="detail-grid">
        <view class="detail-field">
          <text class="detail-field__label">量刑等级</text>
          <text class="detail-field__value">{{ sentenceLevel }}</text>
        </view>
        <view class="detail-field">
          <text class="detail-field__label">痛感风险</text>
          <text class="detail-field__value">{{ painHint }}</text>
        </view>
        <view class="detail-field detail-field--wide">
          <text class="detail-field__label">系统判断</text>
          <text class="detail-field__value">{{ reviewHint }}</text>
        </view>
      </view>
    </view>

    <FixedActionBar v-if="item.reviewStatus === 'none'">
      <up-button type="warning" text="发起复核" @click="review" />
    </FixedActionBar>
  </view>

  <view v-else class="page">
    <EmptyState text="刑期记录加载中" />
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import EmptyState from '../../components/EmptyState.vue'
import FixedActionBar from '../../components/FixedActionBar.vue'
import StatusTag from '../../components/StatusTag.vue'
import { getHellSentence, reviewHellSentence } from '../../api/hell'
import { displayTime, showError, showSuccess } from '../../utils/format'
import { labelOf } from '../../utils/status'

const id = ref('')
const item = ref(null)

const reviewText = computed(() => labelOf(item.value?.reviewStatus))

const sentenceLevel = computed(() => {
  const days = Number(item.value?.sentenceDays || 0)
  if (days >= 700) return '长期服刑'
  if (days >= 180) return '标准刑期'
  return '短期体验'
})

const painHint = computed(() => {
  const pain = Number(item.value?.painLevel || 0)
  if (pain >= 8) return '高痛感'
  if (pain >= 5) return '中痛感'
  return '低痛感'
})

const processSteps = computed(() => [
  { index: '1', title: '判罚登记', desc: `${item.value?.soulName || '-'} 因“${item.value?.crimeType || '-'}”进入 ${item.value?.floorName || '-'}。` },
  { index: '2', title: '设备执行', desc: `执行设备 ${item.value?.equipmentId || '-'}，痛感等级 ${item.value?.painLevel || '-'}。` },
  { index: '3', title: '刑期复核', desc: item.value?.reviewStatus === 'none' ? '当前未发起复核，可提交审批单给判官或阎王。' : '复核流程已进入审批中心，等待处理结果。' }
])

const reviewHint = computed(() => {
  if (!item.value) return '-'
  if (item.value.reviewStatus !== 'none') return '复核已发起，避免重复提交。'
  if (item.value.sentenceDays >= 365 || item.value.painLevel >= 8) return '刑期或痛感偏高，建议发起复核。'
  return '刑期处于常规范围，可继续执行。'
})

onLoad(query => {
  id.value = query.id
})

async function load() {
  if (!id.value) return
  try {
    item.value = await getHellSentence(id.value)
  } catch (error) {
    showError(error)
  }
}

async function review() {
  try {
    await reviewHellSentence(id.value, { note: '移动端发起刑期复核' })
    showSuccess('已发起')
    await load()
  } catch (error) {
    showError(error)
  }
}

onShow(load)
</script>

<style scoped>
.detail-page {
  padding-bottom: 140rpx;
}
</style>
