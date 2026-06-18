<template>
  <view class="page stack detail-page wish-detail" v-if="item">
    <view class="wish-paper">
      <view class="wish-paper__ribbon">{{ wishTypeText }}</view>
      <view class="wish-paper__top">
        <view>
          <text class="wish-paper__eyebrow">INCENSE REQUEST #{{ item.id }}</text>
          <text class="wish-paper__title">{{ item.title }}</text>
        </view>
        <StatusTag :value="item.status" />
      </view>
      <text class="wish-paper__desc">{{ item.requesterName }} 发起 {{ wishTypeText }} 愿望，{{ decisionTone }}</text>
      <view class="wish-scale">
        <view class="wish-scale__item">
          <text class="wish-scale__value">{{ item.incenseAmount }}</text>
          <text class="wish-scale__label">香火</text>
        </view>
        <view class="wish-scale__beam">
          <text :style="{ width: incenseWeight }"></text>
        </view>
        <view class="wish-scale__item">
          <text class="wish-scale__value">{{ item.meritScore }}</text>
          <text class="wish-scale__label">功德</text>
        </view>
      </view>
      <view class="wish-assignee">
        <text class="wish-assignee__label">承办神仙</text>
        <text class="wish-assignee__value">{{ item.assignedDeity || '待分派' }}</text>
      </view>
    </view>

    <view class="wish-verdict">
      <text class="wish-verdict__title">{{ priorityHint }}</text>
      <text class="wish-verdict__desc">{{ jokeHint }}</text>
    </view>

    <view class="detail-card wish-card">
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

    <view class="detail-card wish-card">
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

    <view class="detail-card wish-card">
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

const incenseWeight = computed(() => {
  const total = Math.max(1, Number(item.value?.incenseAmount || 0) + Math.max(0, Number(item.value?.meritScore || 0)))
  return `${Math.min(92, Math.max(8, Math.round((Number(item.value?.incenseAmount || 0) / total) * 100)))}%`
})

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

.wish-detail {
  background:
    radial-gradient(circle at 88% 0%, rgba(234, 179, 8, 0.2), transparent 30%),
    linear-gradient(180deg, #fff7ed 0%, #f7f1e3 48%, #f7f1e3 100%);
}

.wish-paper {
  position: relative;
  overflow: hidden;
  padding: 30rpx;
  border: 1rpx solid rgba(217, 119, 6, 0.22);
  border-radius: 34rpx 34rpx 12rpx 12rpx;
  background:
    radial-gradient(circle at 12% 18%, rgba(251, 191, 36, 0.2), transparent 20%),
    linear-gradient(135deg, #fffbeb, #fff7ed 58%, #fef3c7);
  box-shadow: 0 18rpx 42rpx rgba(217, 119, 6, 0.12);
}

.wish-paper::after {
  content: "";
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  height: 16rpx;
  background:
    linear-gradient(135deg, transparent 8rpx, rgba(217, 119, 6, 0.18) 0) 0 0 / 24rpx 16rpx repeat-x,
    linear-gradient(225deg, transparent 8rpx, rgba(217, 119, 6, 0.18) 0) 0 0 / 24rpx 16rpx repeat-x;
}

.wish-paper__ribbon {
  position: absolute;
  right: -46rpx;
  top: 32rpx;
  width: 190rpx;
  padding: 8rpx 0;
  color: #fffaf0;
  background: #d97706;
  font-size: 22rpx;
  font-weight: 900;
  text-align: center;
  transform: rotate(34deg);
}

.wish-paper__top {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16rpx;
}

.wish-paper__eyebrow {
  display: block;
  color: #d97706;
  font-size: 20rpx;
  font-weight: 900;
  letter-spacing: 2rpx;
}

.wish-paper__title {
  display: block;
  margin-top: 8rpx;
  color: #422006;
  font-size: 50rpx;
  font-weight: 900;
  line-height: 1.13;
}

.wish-paper__desc {
  display: block;
  margin-top: 18rpx;
  color: #78350f;
  font-size: 25rpx;
  line-height: 1.5;
}

.wish-scale {
  display: grid;
  grid-template-columns: 116rpx minmax(0, 1fr) 116rpx;
  gap: 14rpx;
  align-items: center;
  margin-top: 26rpx;
}

.wish-scale__item {
  padding: 14rpx 10rpx;
  border: 1rpx solid rgba(217, 119, 6, 0.18);
  border-radius: 18rpx;
  background: rgba(255, 250, 240, 0.74);
  text-align: center;
}

.wish-scale__value {
  display: block;
  color: #92400e;
  font-size: 34rpx;
  font-weight: 900;
}

.wish-scale__label {
  display: block;
  margin-top: 4rpx;
  color: #a16207;
  font-size: 20rpx;
}

.wish-scale__beam {
  height: 18rpx;
  overflow: hidden;
  border-radius: 999rpx;
  background: rgba(217, 119, 6, 0.16);
}

.wish-scale__beam text {
  display: block;
  height: 100%;
  border-radius: 999rpx;
  background: linear-gradient(90deg, #f59e0b, #d97706);
}

.wish-assignee {
  display: flex;
  justify-content: space-between;
  gap: 12rpx;
  margin-top: 20rpx;
  padding: 16rpx;
  border: 1rpx dashed rgba(217, 119, 6, 0.3);
  border-radius: 16rpx;
  background: rgba(255, 250, 240, 0.64);
}

.wish-assignee__label {
  color: #92400e;
  font-size: 23rpx;
  font-weight: 800;
}

.wish-assignee__value {
  overflow: hidden;
  color: #422006;
  font-size: 24rpx;
  font-weight: 900;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.wish-verdict {
  padding: 20rpx;
  border: 1rpx solid rgba(217, 119, 6, 0.2);
  border-radius: 22rpx;
  background: linear-gradient(135deg, #422006, #92400e);
  color: #fffaf0;
}

.wish-verdict__title {
  display: block;
  font-size: 30rpx;
  font-weight: 900;
}

.wish-verdict__desc {
  display: block;
  margin-top: 8rpx;
  color: rgba(255, 250, 240, 0.76);
  font-size: 24rpx;
  line-height: 1.45;
}

.wish-card {
  border-color: rgba(217, 119, 6, 0.18);
  background: rgba(255, 250, 240, 0.94);
}
</style>
