<template>
  <view class="page detail-page">
    <view class="panel stack" v-if="task">
      <view class="detail-title">
        <text>{{ task.soulName }}</text>
        <StatusTag :value="task.status" />
      </view>
      <up-cell title="地点" :value="task.location" />
      <up-cell title="执行人" :value="task.assigneeName || '-'" />
      <up-cell title="计划时间" :value="displayTime(task.scheduledTime)" />
      <up-cell title="实际时间" :value="displayTime(task.actualTime)" />
      <up-cell title="异常类型" :value="task.exceptionType || '-'" />
      <up-cell title="异常说明" :value="task.exceptionNote || '-'" />
    </view>

    <FixedActionBar v-if="task">
      <up-button v-if="task.status === 'pending'" type="primary" text="开始" @click="start" />
      <up-button v-if="task.status === 'running'" type="success" text="完成" @click="complete" />
      <up-button v-if="task.status === 'running'" type="warning" text="异常" @click="showException = true" />
      <up-button v-if="['pending', 'running', 'exception'].includes(task.status)" type="error" text="撤销" @click="cancel" />
    </FixedActionBar>

    <up-action-sheet :show="showException" :actions="exceptionActions" title="选择异常类型" @close="showException = false" @select="submitException" />
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
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

.detail-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 38rpx;
  font-weight: 800;
}
</style>

