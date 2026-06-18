<template>
  <view class="page stack detail-page" v-if="item">
    <view class="detail-hero">
      <view class="detail-hero__top">
        <view>
          <text class="detail-hero__eyebrow">REINCARNATION #{{ item.id }}</text>
          <text class="detail-hero__title">{{ item.soulName }}</text>
        </view>
        <StatusTag :value="item.status" />
      </view>
      <text class="detail-hero__desc">目标 {{ item.targetRealm }} · {{ item.reviewNote || '等待按流程安排轮回。' }}</text>
      <view class="detail-hero__stats">
        <view class="detail-stat">
          <text class="detail-stat__label">优先级</text>
          <text class="detail-stat__value">{{ item.priority }}</text>
        </view>
        <view class="detail-stat">
          <text class="detail-stat__label">胎位</text>
          <text class="detail-stat__value">{{ item.quotaType || '未分配' }}</text>
        </view>
        <view class="detail-stat">
          <text class="detail-stat__label">孟婆汤</text>
          <text class="detail-stat__value">{{ item.soupStatus }}</text>
        </view>
      </view>
    </view>

    <view class="detail-card">
      <text class="detail-card__title">轮回档案</text>
      <view class="detail-grid">
        <view class="detail-field">
          <text class="detail-field__label">魂魄 ID</text>
          <text class="detail-field__value">#{{ item.soulId }}</text>
        </view>
        <view class="detail-field">
          <text class="detail-field__label">目标道别</text>
          <text class="detail-field__value">{{ item.targetRealm }}</text>
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
          <text class="detail-field__label">复核备注</text>
          <text class="detail-field__value">{{ item.reviewNote || '暂无复核备注。' }}</text>
        </view>
      </view>
    </view>

    <view class="detail-card">
      <text class="detail-card__title">轮回流程</text>
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
      <text class="detail-card__title">下一步建议</text>
      <view class="detail-grid">
        <view class="detail-field">
          <text class="detail-field__label">当前动作</text>
          <text class="detail-field__value">{{ actionHint }}</text>
        </view>
        <view class="detail-field">
          <text class="detail-field__label">汤剂策略</text>
          <text class="detail-field__value">{{ soupHint }}</text>
        </view>
      </view>
    </view>

    <FixedActionBar>
      <up-button v-if="['queued', 'reviewing'].includes(item.status)" type="success" text="通过" @click="approve" />
      <up-button v-if="['queued', 'reviewing', 'approved'].includes(item.status)" type="error" text="驳回" @click="reject" />
      <up-button v-if="['queued', 'approved'].includes(item.status)" type="warning" text="复核" @click="review" />
      <up-button v-if="['approved', 'pending_soup'].includes(item.status)" type="primary" text="分胎位" @click="showQuota = true" />
      <up-button v-if="item.status === 'pending_soup' && item.soupStatus === 'pending'" type="primary" text="发汤" @click="issueSoup" />
      <up-button v-if="item.status === 'pending_soup' && item.soupStatus === 'issued'" type="success" text="完成" @click="complete" />
    </FixedActionBar>

    <up-action-sheet :show="showQuota" :actions="quotaActions" title="选择胎位" @close="showQuota = false" @select="assignQuota" />
  </view>

  <view v-else class="page">
    <EmptyState text="投胎单加载中" />
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import EmptyState from '../../components/EmptyState.vue'
import FixedActionBar from '../../components/FixedActionBar.vue'
import StatusTag from '../../components/StatusTag.vue'
import { approveReincarnation, assignReincarnationQuota, completeReincarnation, getReincarnation, issueReincarnationSoup, rejectReincarnation, reviewReincarnation } from '../../api/reincarnation'
import { getSoupInventory } from '../../api/soup'
import { displayTime, showError, showSuccess } from '../../utils/format'

const id = ref('')
const item = ref(null)
const showQuota = ref(false)
const quotaActions = ['富贵胎', '普通胎', '清贫胎', '天才胎', '困难模式胎'].map(name => ({ name }))

const processSteps = computed(() => [
  { index: '1', title: '排队复核', desc: ['queued', 'reviewing'].includes(item.value?.status) ? '当前需要判官确认是否允许进入轮回。' : '复核阶段已处理。' },
  { index: '2', title: '分配胎位', desc: item.value?.quotaType ? `已分配 ${item.value.quotaType}。` : '通过后需要分配胎位资源。' },
  { index: '3', title: '孟婆汤与归档', desc: item.value?.soupStatus === 'issued' ? '孟婆汤已发放，可完成轮回归档。' : '胎位确认后发放孟婆汤。' }
])

const actionHint = computed(() => {
  if (!item.value) return '-'
  if (['queued', 'reviewing'].includes(item.value.status)) return '审核通过或驳回'
  if (item.value.status === 'approved') return '分配胎位'
  if (item.value.status === 'pending_soup' && item.value.soupStatus === 'pending') return '发放孟婆汤'
  if (item.value.status === 'pending_soup' && item.value.soupStatus === 'issued') return '完成轮回'
  return '已归档'
})

const soupHint = computed(() => {
  if (!item.value) return '-'
  if (item.value.soupStatus === 'issued') return '已喝汤'
  if (item.value.status === 'pending_soup') return '待发汤'
  if (item.value.soupStatus === 'failed') return '需重新排队'
  return '暂不需要'
})

onLoad(query => {
  id.value = query.id
})

async function load() {
  if (!id.value) return
  try {
    item.value = await getReincarnation(id.value)
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

const approve = () => run(() => approveReincarnation(id.value, { note: '移动端通过' }), '已通过')
const reject = () => run(() => rejectReincarnation(id.value, { note: '移动端驳回' }), '已驳回')
const review = () => run(() => reviewReincarnation(id.value, { note: '进入复核' }), '已复核')
const complete = () => run(() => completeReincarnation(id.value, { note: '已入轮回' }), '已完成')

async function assignQuota(action) {
  showQuota.value = false
  await run(() => assignReincarnationQuota(id.value, { quotaType: action.name }), '已分配')
}

async function issueSoup() {
  try {
    const inventory = await getSoupInventory()
    const first = inventory[0]
    if (!first) throw new Error('暂无孟婆汤库存')
    await issueReincarnationSoup(id.value, { inventoryId: first.id, dose: 1, memoryAfter: 0 })
    showSuccess('已发汤')
    await load()
  } catch (error) {
    showError(error)
  }
}

onShow(load)
</script>

<style scoped>
.detail-page {
  padding-bottom: 170rpx;
}
</style>
