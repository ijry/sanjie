<template>
  <view class="page stack detail-page" v-if="item">
    <view class="detail-hero">
      <view class="detail-hero__top">
        <view>
          <text class="detail-hero__eyebrow">APPROVAL #{{ item.id }}</text>
          <text class="detail-hero__title">{{ item.title }}</text>
        </view>
        <StatusTag :value="item.status" />
      </view>
      <text class="detail-hero__desc">{{ item.reason || '无申请说明。' }}</text>
      <view class="detail-hero__stats">
        <view class="detail-stat">
          <text class="detail-stat__label">类型</text>
          <text class="detail-stat__value">{{ typeLabel(item.type) }}</text>
        </view>
        <view class="detail-stat">
          <text class="detail-stat__label">目标</text>
          <text class="detail-stat__value">#{{ item.targetId }}</text>
        </view>
        <view class="detail-stat">
          <text class="detail-stat__label">动作</text>
          <text class="detail-stat__value">{{ actionHint }}</text>
        </view>
      </view>
    </view>

    <view class="detail-card">
      <text class="detail-card__title">审批档案</text>
      <view class="detail-grid">
        <view class="detail-field">
          <text class="detail-field__label">申请人</text>
          <text class="detail-field__value">用户 {{ item.applicantId || '-' }}</text>
        </view>
        <view class="detail-field">
          <text class="detail-field__label">审批人</text>
          <text class="detail-field__value">用户 {{ item.approverId || '-' }}</text>
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
          <text class="detail-field__label">申请理由</text>
          <text class="detail-field__value">{{ item.reason || '-' }}</text>
        </view>
        <view class="detail-field detail-field--wide">
          <text class="detail-field__label">处理结果</text>
          <text class="detail-field__value">{{ item.resultNote || '尚未填写处理结果。' }}</text>
        </view>
      </view>
    </view>

    <view class="detail-card">
      <text class="detail-card__title">审批路径</text>
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
      <text class="detail-card__title">风险提示</text>
      <view class="detail-grid">
        <view class="detail-field detail-field--wide">
          <text class="detail-field__label">系统判断</text>
          <text class="detail-field__value">{{ riskHint }}</text>
        </view>
      </view>
    </view>

    <FixedActionBar v-if="['pending', 'transferred'].includes(item.status)">
      <up-button type="success" text="同意" @click="approve" />
      <up-button type="error" text="驳回" @click="reject" />
      <up-button v-if="item.status === 'pending'" type="primary" text="转交" @click="transfer" />
    </FixedActionBar>
  </view>

  <view v-else class="page">
    <EmptyState text="审批单加载中" />
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import EmptyState from '../../components/EmptyState.vue'
import FixedActionBar from '../../components/FixedActionBar.vue'
import StatusTag from '../../components/StatusTag.vue'
import { approveApproval, getApproval, rejectApproval, transferApproval } from '../../api/approvals'
import { displayTime, showError, showSuccess } from '../../utils/format'

const id = ref('')
const item = ref(null)

const processSteps = computed(() => [
  { index: '1', title: '提交申请', desc: `申请人用户 ${item.value?.applicantId || '-'} 提交 ${typeLabel(item.value?.type)}。` },
  { index: '2', title: '判官复核', desc: item.value?.status === 'pending' ? '当前等待审批处理。' : '已进入处理链路。' },
  { index: '3', title: '归档留痕', desc: ['approved', 'rejected'].includes(item.value?.status) ? '审批已完结并写入审计日志。' : '同意、驳回或转交后都会留痕。' }
])

const actionHint = computed(() => {
  if (!item.value) return '-'
  if (item.value.status === 'pending') return '可同意/驳回/转交'
  if (item.value.status === 'transferred') return '等待上级处理'
  return '已完结'
})

const riskHint = computed(() => {
  const type = item.value?.type
  if (type === 'lifespan_change') return '改寿会影响生死簿，建议核对审计来源。'
  if (type === 'reincarnation_jump') return '投胎插队风险高，需确认功德依据。'
  if (type === 'capture_exception') return '勾魂异常需保留现场说明，避免错勾。'
  if (type === 'hell_review') return '刑期复核会影响地狱容量和楼层调度。'
  if (type === 'soup_extra_dose') return '孟婆汤加量会影响库存和记忆残留。'
  return '高危操作需留痕。'
})

onLoad(query => {
  id.value = query.id
})

async function load() {
  if (!id.value) return
  try {
    item.value = await getApproval(id.value)
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

const approve = () => run(() => approveApproval(id.value, { note: '移动端审批通过' }), '已同意')
const reject = () => run(() => rejectApproval(id.value, { note: '移动端审批驳回' }), '已驳回')
const transfer = () => run(() => transferApproval(id.value, { note: '转交上级复核' }), '已转交')

function typeLabel(type) {
  const labels = {
    lifespan_change: '寿命变更',
    hell_review: '刑期复核',
    reincarnation_jump: '投胎插队',
    capture_exception: '勾魂异常',
    soup_extra_dose: '孟婆汤加量'
  }
  return labels[type] || type || '-'
}

onShow(load)
</script>

<style scoped>
.detail-page {
  padding-bottom: 140rpx;
}
</style>
