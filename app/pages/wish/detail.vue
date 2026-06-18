<template>
  <view class="page stack detail-page" v-if="item">
    <view class="detail-hero">
      <view class="detail-hero__top">
        <view>
          <text class="detail-hero__eyebrow">WISH TICKET #{{ item.id }}</text>
          <text class="detail-hero__title">{{ item.title }}</text>
        </view>
        <StatusTag :value="item.status" />
      </view>
      <text class="detail-hero__desc">{{ item.requesterName }} 发起 {{ wishTypeText }} 愿望，{{ decisionTone }}</text>
      <view class="detail-hero__stats">
        <view class="detail-stat">
          <text class="detail-stat__label">香火</text>
          <text class="detail-stat__value">{{ item.incenseAmount }}</text>
        </view>
        <view class="detail-stat">
          <text class="detail-stat__label">功德</text>
          <text class="detail-stat__value">{{ item.meritScore }}</text>
        </view>
        <view class="detail-stat">
          <text class="detail-stat__label">承办</text>
          <text class="detail-stat__value">{{ item.assignedDeity || '待分派' }}</text>
        </view>
      </view>
    </view>

    <view class="detail-card">
      <text class="detail-card__title">愿望档案</text>
      <view class="detail-grid">
        <view class="detail-field">
          <text class="detail-field__label">许愿人</text>
          <text class="detail-field__value">{{ item.requesterName }}</text>
        </view>
        <view class="detail-field">
          <text class="detail-field__label">愿望类型</text>
          <text class="detail-field__value">{{ wishTypeText }}</text>
        </view>
        <view class="detail-field">
          <text class="detail-field__label">当前状态</text>
          <text class="detail-field__value">{{ statusText }}</text>
        </view>
        <view class="detail-field">
          <text class="detail-field__label">更新时间</text>
          <text class="detail-field__value">{{ displayTime(item.updatedAt) }}</text>
        </view>
        <view class="detail-field detail-field--wide">
          <text class="detail-field__label">处理结果</text>
          <text class="detail-field__value">{{ item.resultNote || '尚未填写结果，建议先评估香火、功德和合理性。' }}</text>
        </view>
      </view>
    </view>

    <view class="detail-card">
      <text class="detail-card__title">处理流程</text>
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
      <text class="detail-card__title">判定建议</text>
      <view class="detail-grid">
        <view class="detail-field">
          <text class="detail-field__label">优先级</text>
          <text class="detail-field__value">{{ priorityHint }}</text>
        </view>
        <view class="detail-field">
          <text class="detail-field__label">下一步</text>
          <text class="detail-field__value">{{ actionHint }}</text>
        </view>
        <view class="detail-field detail-field--wide">
          <text class="detail-field__label">系统吐槽</text>
          <text class="detail-field__value">{{ jokeHint }}</text>
        </view>
      </view>
    </view>

    <FixedActionBar v-if="['pending', 'routed'].includes(item.status)">
      <up-button v-if="item.status === 'pending'" type="primary" text="转派" @click="showRoute = true" />
      <up-button type="success" text="办结" @click="resolve" />
      <up-button type="error" text="驳回" @click="reject" />
    </FixedActionBar>

    <up-action-sheet :show="showRoute" :actions="routeActions" title="选择承办神仙" @close="showRoute = false" @select="route" />
  </view>

  <view v-else class="page">
    <EmptyState text="愿望工单加载中" />
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import EmptyState from '../../components/EmptyState.vue'
import FixedActionBar from '../../components/FixedActionBar.vue'
import StatusTag from '../../components/StatusTag.vue'
import { getWish, rejectWish, resolveWish, routeWish } from '../../api/wishes'
import { displayTime, showError, showSuccess } from '../../utils/format'
import { labelOf } from '../../utils/status'

const id = ref('')
const item = ref(null)
const showRoute = ref(false)
const routeActions = [
  { name: '财神', note: '转财运组评估，香火不足也要看一眼。' },
  { name: '月老', note: '转姻缘组复核，禁止强行复合。' },
  { name: '文昌帝君', note: '转考试组处理，建议同步刷题。' },
  { name: '太白金星', note: '转疑难杂症组，先给一条体面建议。' }
]

const wishTypeText = computed(() => {
  const labels = {
    wealth: '求财',
    study: '求学',
    work: '职场',
    love: '姻缘',
    health: '健康',
    career: '事业',
    luck: '运气'
  }
  return labels[item.value?.wishType] || item.value?.wishType || '-'
})

const statusText = computed(() => labelOf(item.value?.status))

const decisionTone = computed(() => {
  if (!item.value) return ''
  if (item.value.status === 'resolved') return '已给出神仙建议'
  if (item.value.status === 'rejected') return '已被天规挡回'
  if (item.value.status === 'routed') return `已转给 ${item.value.assignedDeity || '承办神仙'}`
  return '等待分派或现场驳回'
})

const processSteps = computed(() => [
  { index: '1', title: '收愿登记', desc: `${item.value?.requesterName || '-'} 提交愿望，香火 ${item.value?.incenseAmount || 0}，功德 ${item.value?.meritScore || 0}。` },
  { index: '2', title: '神仙分派', desc: item.value?.assignedDeity ? `当前承办：${item.value.assignedDeity}。` : '尚未分派，可按愿望类型选择承办神仙。' },
  { index: '3', title: '办结留痕', desc: ['resolved', 'rejected'].includes(item.value?.status) ? '工单已完结，结果会进入审计日志。' : '办结或驳回后会写入结果备注。' }
])

const priorityHint = computed(() => {
  if (!item.value) return '-'
  const score = Number(item.value.incenseAmount || 0) + Number(item.value.meritScore || 0)
  if (score >= 35) return '高优先级'
  if (score >= 15) return '正常处理'
  return '低优先级'
})

const actionHint = computed(() => {
  if (!item.value) return '-'
  if (item.value.status === 'pending') return '可转派、办结或驳回'
  if (item.value.status === 'routed') return '等待承办神仙给出结论'
  return '已归档'
})

const jokeHint = computed(() => {
  if (!item.value) return '-'
  if (item.value.meritScore < 0) return '功德为负还想许愿，建议先还愿再排队。'
  if (item.value.incenseAmount < 3) return '香火太少，系统建议先补票。'
  if (item.value.wishType === 'wealth') return '求财类愿望需确认是否包含“不劳而获”关键词。'
  if (item.value.wishType === 'love') return '姻缘类愿望禁止绑定前任回滚。'
  return '愿望看起来正常，但仍建议自力更生。'
})

onLoad(query => {
  id.value = query.id
})

async function load() {
  if (!id.value) return
  try {
    item.value = await getWish(id.value)
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

async function route(action) {
  showRoute.value = false
  await run(() => routeWish(id.value, { assignedDeity: action.name, note: action.note }), '已转派')
}

const resolve = () => run(() => resolveWish(id.value, { note: '建议先自力更生' }), '已办结')
const reject = () => run(() => rejectWish(id.value, { note: '愿望不符合天规' }), '已驳回')

onShow(load)
</script>

<style scoped>
.detail-page {
  padding-bottom: 140rpx;
}
</style>
