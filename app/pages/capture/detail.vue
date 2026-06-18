<template>
  <view class="page stack detail-page capture-detail" v-if="task">
    <view class="capture-hero">
      <view class="capture-hero__map">
        <view class="capture-hero__scan"></view>
        <view class="capture-hero__pin"></view>
      </view>
      <view class="capture-hero__top">
        <view>
          <text class="capture-hero__eyebrow">NIGHT PATROL / TASK #{{ task.id }}</text>
          <text class="capture-hero__title">{{ task.soulName }}</text>
        </view>
        <StatusTag :value="task.status" />
      </view>
      <text class="capture-hero__desc">{{ task.location }}</text>
      <view class="capture-coordinates">
        <text>魂魄 #{{ task.soulId }}</text>
        <text>{{ shortTime(task.scheduledTime) }}</text>
        <text>{{ task.exceptionType || 'NO ALERT' }}</text>
      </view>
      <view class="capture-hero__stats">
        <view class="capture-stat">
          <text class="capture-stat__label">夜巡鬼差</text>
          <text class="capture-stat__value">{{ task.assigneeName || '待派单' }}</text>
        </view>
        <view class="capture-stat">
          <text class="capture-stat__label">计划窗口</text>
          <text class="capture-stat__value">{{ shortTime(task.scheduledTime) }}</text>
        </view>
        <view class="capture-stat">
          <text class="capture-stat__label">回执时间</text>
          <text class="capture-stat__value">{{ shortTime(task.actualTime) }}</text>
        </view>
      </view>
    </view>

    <view class="capture-route">
      <view class="capture-route__node">
        <text class="capture-route__dot"></text>
        <text class="capture-route__title">定位</text>
        <text class="capture-route__desc">{{ task.location }}</text>
      </view>
      <view class="capture-route__line"></view>
      <view class="capture-route__node">
        <text class="capture-route__dot capture-route__dot--active"></text>
        <text class="capture-route__title">勾魂</text>
        <text class="capture-route__desc">{{ actionHint }}</text>
      </view>
      <view class="capture-route__line"></view>
      <view class="capture-route__node">
        <text class="capture-route__dot"></text>
        <text class="capture-route__title">归档</text>
        <text class="capture-route__desc">{{ task.actualTime ? displayTime(task.actualTime) : '等待回执' }}</text>
      </view>
    </view>

    <view class="detail-card capture-card">
      <text class="detail-card__title">任务档案</text>
      <view class="detail-grid">
        <view class="detail-field">
          <text class="detail-field__label">魂魄 ID</text>
          <text class="detail-field__value">#{{ task.soulId }}</text>
        </view>
        <view class="detail-field">
          <text class="detail-field__label">任务状态</text>
          <text class="detail-field__value">{{ task.status }}</text>
        </view>
        <view class="detail-field detail-field--wide">
          <text class="detail-field__label">勾魂地点</text>
          <text class="detail-field__value">{{ task.location }}</text>
        </view>
        <view class="detail-field">
          <text class="detail-field__label">计划时间</text>
          <text class="detail-field__value">{{ displayTime(task.scheduledTime) }}</text>
        </view>
        <view class="detail-field">
          <text class="detail-field__label">实际时间</text>
          <text class="detail-field__value">{{ displayTime(task.actualTime) }}</text>
        </view>
      </view>
    </view>

    <view class="detail-card capture-card">
      <text class="detail-card__title">执行闭环</text>
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

    <view class="detail-card capture-card">
      <text class="detail-card__title">异常与备注</text>
      <view class="detail-grid">
        <view class="detail-field">
          <text class="detail-field__label">异常类型</text>
          <text class="detail-field__value">{{ task.exceptionType || '无异常' }}</text>
        </view>
        <view class="detail-field">
          <text class="detail-field__label">可操作</text>
          <text class="detail-field__value">{{ actionHint }}</text>
        </view>
        <view class="detail-field detail-field--wide">
          <text class="detail-field__label">异常说明</text>
          <text class="detail-field__value">{{ task.exceptionNote || '当前任务未登记异常，按正常勾魂流程处理。' }}</text>
        </view>
      </view>
    </view>

    <FixedActionBar>
      <up-button v-if="task.status === 'pending'" type="primary" text="开始" @click="start" />
      <up-button v-if="task.status === 'running'" type="success" text="完成" @click="complete" />
      <up-button v-if="task.status === 'running'" type="warning" text="异常" @click="showException = true" />
      <up-button v-if="['pending', 'running', 'exception'].includes(task.status)" type="error" text="撤销" @click="cancel" />
    </FixedActionBar>

    <up-action-sheet :show="showException" :actions="exceptionActions" title="选择异常类型" @close="showException = false" @select="submitException" />
  </view>

  <view v-else class="page">
    <EmptyState text="勾魂任务加载中" />
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import EmptyState from '../../components/EmptyState.vue'
import StatusTag from '../../components/StatusTag.vue'
import FixedActionBar from '../../components/FixedActionBar.vue'
import { cancelCaptureTask, completeCaptureTask, getCaptureTask, reportCaptureException, startCaptureTask } from '../../api/capture'
import { displayTime, showError, showSuccess } from '../../utils/format'

const id = ref('')
const task = ref(null)
const showException = ref(false)
const exceptionActions = [
  { name: '目标抢救中', type: 'emergency_rescue' },
  { name: '勾错魂', type: 'wrong_soul' },
  { name: '魂魄逃逸', type: 'escaped' },
  { name: '查无此人', type: 'not_found' },
  { name: '疑似飞升', type: 'ascended' }
]

const processSteps = computed(() => [
  { index: '1', title: '接单定位', desc: `鬼差前往 ${task.value?.location || '目标地点'}，核对生死簿与魂魄编号。` },
  { index: '2', title: '执行勾魂', desc: task.value?.status === 'pending' ? '当前尚未开始，点击底部“开始”进入执行中。' : '任务已进入执行流，需记录实际结果。' },
  { index: '3', title: '闭环归档', desc: task.value?.status === 'completed' ? '任务已完成归档。' : '完成、异常或撤销后都会进入审计留痕。' }
])

const actionHint = computed(() => {
  if (!task.value) return '-'
  if (task.value.status === 'pending') return '可开始/撤销'
  if (task.value.status === 'running') return '可完成/上报异常/撤销'
  if (task.value.status === 'exception') return '可撤销或等待复盘'
  return '已归档'
})

onLoad(query => {
  id.value = query.id
})

async function load() {
  if (!id.value) return
  try {
    task.value = await getCaptureTask(id.value)
  } catch (error) {
    showError(error)
  }
}

async function run(action, successText) {
  try {
    await action()
    showSuccess(successText)
    await load()
  } catch (error) {
    showError(error)
  }
}

const start = () => run(() => startCaptureTask(id.value), '已开始')
const complete = () => run(() => completeCaptureTask(id.value, { note: '移动端确认完成' }), '已完成')
const cancel = () => run(() => cancelCaptureTask(id.value), '已撤销')

function shortTime(value) {
  const text = displayTime(value)
  if (text === '-') return '-'
  return text.slice(11, 16)
}

async function submitException(action) {
  showException.value = false
  await run(() => reportCaptureException(id.value, { type: action.type, note: action.name }), '已上报')
}

onShow(load)
</script>

<style scoped>
.detail-page {
  padding-bottom: 140rpx;
}

.capture-detail {
  background:
    radial-gradient(circle at 18% 0%, rgba(74, 111, 165, 0.18), transparent 28%),
    linear-gradient(180deg, #111827 0%, #1f2937 32%, #f7f1e3 32%, #f7f1e3 100%);
}

.capture-hero {
  position: relative;
  overflow: hidden;
  min-height: 390rpx;
  padding: 26rpx;
  color: #f8fafc;
  border: 1rpx solid rgba(148, 163, 184, 0.32);
  border-radius: 30rpx;
  background:
    linear-gradient(90deg, rgba(148, 163, 184, 0.08) 1rpx, transparent 1rpx),
    linear-gradient(0deg, rgba(148, 163, 184, 0.08) 1rpx, transparent 1rpx),
    radial-gradient(circle at 76% 18%, rgba(34, 197, 94, 0.18), transparent 26%),
    linear-gradient(135deg, #020617, #172554 56%, #111827);
  background-size: 42rpx 42rpx, 42rpx 42rpx, auto, auto;
  box-shadow: 0 22rpx 44rpx rgba(2, 6, 23, 0.28);
}

.capture-hero__map {
  position: absolute;
  right: -42rpx;
  top: -26rpx;
  width: 310rpx;
  height: 310rpx;
  border: 1rpx solid rgba(125, 211, 252, 0.28);
  border-radius: 50%;
}

.capture-hero__scan {
  position: absolute;
  inset: 42rpx;
  border: 1rpx dashed rgba(125, 211, 252, 0.34);
  border-radius: 50%;
}

.capture-hero__pin {
  position: absolute;
  left: 128rpx;
  top: 132rpx;
  width: 26rpx;
  height: 26rpx;
  border: 8rpx solid rgba(248, 113, 113, 0.9);
  border-radius: 50%;
  background: #fff7ed;
  box-shadow: 0 0 26rpx rgba(248, 113, 113, 0.7);
}

.capture-hero__top {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16rpx;
}

.capture-hero__eyebrow {
  display: block;
  color: #7dd3fc;
  font-size: 20rpx;
  font-weight: 900;
  letter-spacing: 2rpx;
}

.capture-hero__title {
  display: block;
  margin-top: 10rpx;
  font-size: 48rpx;
  font-weight: 900;
  line-height: 1.15;
}

.capture-hero__desc {
  position: relative;
  z-index: 1;
  display: block;
  width: 72%;
  margin-top: 18rpx;
  color: rgba(248, 250, 252, 0.76);
  font-size: 25rpx;
  line-height: 1.45;
}

.capture-coordinates {
  position: relative;
  z-index: 1;
  display: flex;
  flex-wrap: wrap;
  gap: 10rpx;
  margin-top: 22rpx;
}

.capture-coordinates text {
  padding: 8rpx 12rpx;
  color: #bae6fd;
  border: 1rpx solid rgba(125, 211, 252, 0.22);
  border-radius: 999rpx;
  background: rgba(14, 165, 233, 0.1);
  font-size: 20rpx;
  font-weight: 800;
}

.capture-hero__stats {
  position: relative;
  z-index: 1;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10rpx;
  margin-top: 26rpx;
}

.capture-stat {
  min-width: 0;
  padding: 14rpx 12rpx;
  border-left: 4rpx solid #38bdf8;
  background: rgba(15, 23, 42, 0.72);
}

.capture-stat__label {
  display: block;
  color: rgba(226, 232, 240, 0.62);
  font-size: 20rpx;
}

.capture-stat__value {
  display: block;
  overflow: hidden;
  margin-top: 8rpx;
  color: #f8fafc;
  font-size: 25rpx;
  font-weight: 900;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.capture-route {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 34rpx minmax(0, 1fr) 34rpx minmax(0, 1fr);
  align-items: stretch;
  padding: 20rpx;
  border: 1rpx solid rgba(56, 189, 248, 0.22);
  border-radius: 22rpx;
  background: #f8fafc;
  box-shadow: 0 12rpx 30rpx rgba(15, 23, 42, 0.1);
}

.capture-route__node {
  min-width: 0;
}

.capture-route__dot {
  display: block;
  width: 18rpx;
  height: 18rpx;
  border-radius: 50%;
  background: #94a3b8;
}

.capture-route__dot--active {
  background: #ef4444;
  box-shadow: 0 0 18rpx rgba(239, 68, 68, 0.7);
}

.capture-route__title {
  display: block;
  margin-top: 10rpx;
  color: #0f172a;
  font-size: 25rpx;
  font-weight: 900;
}

.capture-route__desc {
  display: block;
  overflow: hidden;
  margin-top: 6rpx;
  color: #475569;
  font-size: 21rpx;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.capture-route__line {
  align-self: center;
  height: 2rpx;
  background: linear-gradient(90deg, #94a3b8, #38bdf8);
}

.capture-card {
  border-color: rgba(56, 189, 248, 0.24);
  background: linear-gradient(180deg, #ffffff, #f8fafc);
}
</style>
