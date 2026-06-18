<template>
  <view class="page stack detail-page reincarnation-detail" v-if="item">
    <view class="reincarnation-ticket">
      <view class="ticket-edge ticket-edge--left"></view>
      <view class="ticket-edge ticket-edge--right"></view>
      <view class="reincarnation-ticket__top">
        <view>
          <text class="reincarnation-ticket__eyebrow">NAIHE BRIDGE QUEUE #{{ item.id }}</text>
          <text class="reincarnation-ticket__title">{{ item.soulName }}</text>
        </view>
        <StatusTag :value="item.status" />
      </view>
      <view class="ticket-route">
        <text>鬼门关</text>
        <text>奈何桥</text>
        <text>{{ item.targetRealm }}</text>
      </view>
      <text class="reincarnation-ticket__desc">目标 {{ item.targetRealm }} · {{ item.reviewNote || '等待按流程安排轮回。' }}</text>
      <view class="ticket-stub">
        <view class="ticket-stub__qr">
          <text></text><text></text><text></text><text></text>
        </view>
        <view class="ticket-stub__main">
          <view class="ticket-stub__row">
            <text class="ticket-stub__label">优先级</text>
            <text class="ticket-stub__value">{{ item.priority }}</text>
          </view>
          <view class="ticket-stub__row">
            <text class="ticket-stub__label">胎位</text>
            <text class="ticket-stub__value">{{ item.quotaType || '未分配' }}</text>
          </view>
          <view class="ticket-stub__row">
            <text class="ticket-stub__label">孟婆汤</text>
            <text class="ticket-stub__value">{{ item.soupStatus }}</text>
          </view>
        </view>
      </view>
    </view>

    <view class="bridge-progress">
      <view v-for="step in bridgeSteps" :key="step.title" class="bridge-step" :class="{ 'bridge-step--active': step.active }">
        <text class="bridge-step__marker">{{ step.index }}</text>
        <text class="bridge-step__title">{{ step.title }}</text>
        <text class="bridge-step__desc">{{ step.desc }}</text>
      </view>
    </view>

    <view class="detail-card reincarnation-card">
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

    <view class="detail-card reincarnation-card">
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

    <view class="detail-card reincarnation-card">
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

const bridgeSteps = computed(() => [
  { index: '一', title: '复核', desc: ['queued', 'reviewing'].includes(item.value?.status) ? '排队中' : '已处理', active: ['queued', 'reviewing'].includes(item.value?.status) },
  { index: '二', title: '胎位', desc: item.value?.quotaType || '待分配', active: item.value?.status === 'approved' },
  { index: '三', title: '孟婆汤', desc: item.value?.soupStatus === 'issued' ? '已发放' : '待发汤', active: item.value?.status === 'pending_soup' },
  { index: '四', title: '入轮回', desc: item.value?.status === 'completed' ? '已归档' : '未完成', active: item.value?.status === 'completed' }
])

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

.reincarnation-detail {
  background:
    radial-gradient(circle at 78% 4%, rgba(20, 184, 166, 0.16), transparent 28%),
    linear-gradient(180deg, #ecfeff 0%, #f7f1e3 40%, #f7f1e3 100%);
}

.reincarnation-ticket {
  position: relative;
  overflow: hidden;
  padding: 28rpx;
  border: 2rpx solid rgba(15, 118, 110, 0.2);
  border-radius: 28rpx;
  background:
    linear-gradient(90deg, rgba(15, 118, 110, 0.09) 0 2rpx, transparent 2rpx 22rpx),
    linear-gradient(135deg, #f0fdfa, #fff7ed 56%, #fef3c7);
  box-shadow: 0 18rpx 42rpx rgba(15, 118, 110, 0.14);
}

.ticket-edge {
  position: absolute;
  top: 50%;
  width: 44rpx;
  height: 44rpx;
  border-radius: 50%;
  background: #f7f1e3;
  transform: translateY(-50%);
}

.ticket-edge--left {
  left: -22rpx;
}

.ticket-edge--right {
  right: -22rpx;
}

.reincarnation-ticket__top {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16rpx;
}

.reincarnation-ticket__eyebrow {
  display: block;
  color: #0f766e;
  font-size: 20rpx;
  font-weight: 900;
  letter-spacing: 2rpx;
}

.reincarnation-ticket__title {
  display: block;
  margin-top: 8rpx;
  color: #134e4a;
  font-size: 48rpx;
  font-weight: 900;
  line-height: 1.15;
}

.ticket-route {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8rpx;
  margin-top: 22rpx;
}

.ticket-route text {
  padding: 12rpx;
  color: #0f766e;
  border: 1rpx dashed rgba(15, 118, 110, 0.28);
  border-radius: 999rpx;
  background: rgba(240, 253, 250, 0.78);
  font-size: 22rpx;
  font-weight: 900;
  text-align: center;
}

.reincarnation-ticket__desc {
  display: block;
  margin-top: 18rpx;
  color: #475569;
  font-size: 25rpx;
  line-height: 1.45;
}

.ticket-stub {
  display: grid;
  grid-template-columns: 112rpx minmax(0, 1fr);
  gap: 18rpx;
  margin-top: 24rpx;
  padding-top: 20rpx;
  border-top: 2rpx dashed rgba(15, 118, 110, 0.2);
}

.ticket-stub__qr {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8rpx;
  padding: 10rpx;
  border-radius: 18rpx;
  background: #134e4a;
}

.ticket-stub__qr text {
  border-radius: 8rpx;
  background: #ccfbf1;
}

.ticket-stub__main {
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.ticket-stub__row {
  display: flex;
  justify-content: space-between;
  gap: 12rpx;
  padding: 8rpx 0;
  border-bottom: 1rpx solid rgba(15, 118, 110, 0.12);
}

.ticket-stub__label {
  color: #64748b;
  font-size: 22rpx;
}

.ticket-stub__value {
  overflow: hidden;
  color: #134e4a;
  font-size: 24rpx;
  font-weight: 900;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.bridge-progress {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 10rpx;
}

.bridge-step {
  min-width: 0;
  padding: 16rpx 12rpx;
  border: 1rpx solid rgba(15, 118, 110, 0.16);
  border-radius: 18rpx;
  background: rgba(255, 250, 240, 0.76);
}

.bridge-step--active {
  color: #fffaf0;
  background: linear-gradient(135deg, #0f766e, #14b8a6);
}

.bridge-step__marker {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 38rpx;
  height: 38rpx;
  border-radius: 50%;
  background: rgba(15, 118, 110, 0.12);
  font-size: 20rpx;
  font-weight: 900;
}

.bridge-step--active .bridge-step__marker {
  background: rgba(255, 250, 240, 0.22);
}

.bridge-step__title {
  display: block;
  margin-top: 10rpx;
  font-size: 24rpx;
  font-weight: 900;
}

.bridge-step__desc {
  display: block;
  overflow: hidden;
  margin-top: 4rpx;
  color: #64748b;
  font-size: 20rpx;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.bridge-step--active .bridge-step__desc {
  color: rgba(255, 250, 240, 0.76);
}

.reincarnation-card {
  border-color: rgba(15, 118, 110, 0.18);
  background: rgba(255, 250, 240, 0.92);
}
</style>
