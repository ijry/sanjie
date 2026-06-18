<template>
  <view class="page stack detail-page" v-if="task">
    <view class="detail-hero">
      <view class="detail-hero__top">
        <view>
          <text class="detail-hero__eyebrow">CAPTURE TASK #{{ task.id }}</text>
          <text class="detail-hero__title">{{ task.soulName }}</text>
        </view>
        <StatusTag :value="task.status" />
      </view>
      <text class="detail-hero__desc">{{ task.location }}</text>
      <view class="detail-hero__stats">
        <view class="detail-stat">
          <text class="detail-stat__label">执行人</text>
          <text class="detail-stat__value">{{ task.assigneeName || '待派单' }}</text>
        </view>
        <view class="detail-stat">
          <text class="detail-stat__label">计划</text>
          <text class="detail-stat__value">{{ shortTime(task.scheduledTime) }}</text>
        </view>
        <view class="detail-stat">
          <text class="detail-stat__label">实际</text>
          <text class="detail-stat__value">{{ shortTime(task.actualTime) }}</text>
        </view>
      </view>
    </view>

    <view class="detail-card">
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

    <view class="detail-card">
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

    <view class="detail-card">
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
</style>
