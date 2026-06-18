<template>
  <view class="page stack detail-page approval-detail" v-if="item">
    <view class="approval-dossier">
      <view class="approval-dossier__stamp">{{ stampText }}</view>
      <view class="approval-dossier__top">
        <view>
          <text class="approval-dossier__eyebrow">JUDGE DOSSIER #{{ item.id }}</text>
          <text class="approval-dossier__title">{{ item.title }}</text>
        </view>
        <StatusTag :value="item.status" />
      </view>
      <text class="approval-dossier__reason">{{ item.reason || '无申请说明。' }}</text>
      <view class="dossier-meta">
        <view class="dossier-meta__item">
          <text class="dossier-meta__label">案由</text>
          <text class="dossier-meta__value">{{ typeLabel(item.type) }}</text>
        </view>
        <view class="dossier-meta__item">
          <text class="dossier-meta__label">目标卷宗</text>
          <text class="dossier-meta__value">#{{ item.targetId }}</text>
        </view>
        <view class="dossier-meta__item">
          <text class="dossier-meta__label">当前动作</text>
          <text class="dossier-meta__value">{{ actionHint }}</text>
        </view>
      </view>
    </view>

    <view class="approval-chain">
      <view class="approval-chain__node">
        <text>申请</text>
        <text>用户 {{ item.applicantId || '-' }}</text>
      </view>
      <view class="approval-chain__arrow"></view>
      <view class="approval-chain__node approval-chain__node--focus">
        <text>判官</text>
        <text>{{ item.approverId ? `用户 ${item.approverId}` : '待审' }}</text>
      </view>
      <view class="approval-chain__arrow"></view>
      <view class="approval-chain__node">
        <text>归档</text>
        <text>{{ ['approved', 'rejected'].includes(item.status) ? '已留痕' : '待完成' }}</text>
      </view>
    </view>

    <view class="detail-card approval-card">
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

    <view class="detail-card approval-card">
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

    <view class="risk-seal-card">
      <view class="risk-seal-card__mark">审</view>
      <view class="risk-seal-card__main">
        <text class="risk-seal-card__title">风险提示</text>
        <text class="risk-seal-card__desc">{{ riskHint }}</text>
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

const stampText = computed(() => {
  if (!item.value) return '待审'
  if (item.value.status === 'approved') return '准'
  if (item.value.status === 'rejected') return '驳'
  if (item.value.status === 'transferred') return '转'
  return '审'
})

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

.approval-detail {
  background:
    radial-gradient(circle at 0% 0%, rgba(185, 28, 28, 0.1), transparent 26%),
    linear-gradient(180deg, #fff7ed 0%, #f7f1e3 52%, #f7f1e3 100%);
}

.approval-dossier {
  position: relative;
  overflow: hidden;
  padding: 30rpx;
  border: 2rpx solid rgba(120, 53, 15, 0.2);
  border-radius: 18rpx;
  background:
    linear-gradient(90deg, rgba(120, 53, 15, 0.05) 1rpx, transparent 1rpx),
    linear-gradient(0deg, rgba(120, 53, 15, 0.05) 1rpx, transparent 1rpx),
    #fff7ed;
  background-size: 28rpx 28rpx;
  box-shadow: 0 16rpx 40rpx rgba(120, 53, 15, 0.12);
}

.approval-dossier__stamp {
  position: absolute;
  right: 28rpx;
  bottom: 22rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 140rpx;
  height: 140rpx;
  color: rgba(185, 28, 28, 0.2);
  border: 8rpx solid rgba(185, 28, 28, 0.18);
  border-radius: 24rpx;
  font-size: 74rpx;
  font-weight: 900;
  transform: rotate(-14deg);
}

.approval-dossier__top {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16rpx;
}

.approval-dossier__eyebrow {
  display: block;
  color: #b45309;
  font-size: 20rpx;
  font-weight: 900;
  letter-spacing: 2rpx;
}

.approval-dossier__title {
  display: block;
  margin-top: 8rpx;
  color: #431407;
  font-size: 46rpx;
  font-weight: 900;
  line-height: 1.18;
}

.approval-dossier__reason {
  position: relative;
  z-index: 1;
  display: block;
  margin-top: 18rpx;
  color: #7c2d12;
  font-size: 25rpx;
  line-height: 1.55;
}

.dossier-meta {
  position: relative;
  z-index: 1;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10rpx;
  margin-top: 24rpx;
}

.dossier-meta__item {
  min-width: 0;
  padding: 14rpx 12rpx;
  border-bottom: 4rpx solid #b45309;
  background: rgba(255, 250, 240, 0.68);
}

.dossier-meta__label {
  display: block;
  color: #9a3412;
  font-size: 20rpx;
}

.dossier-meta__value {
  display: block;
  overflow: hidden;
  margin-top: 8rpx;
  color: #431407;
  font-size: 25rpx;
  font-weight: 900;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.approval-chain {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 34rpx minmax(0, 1fr) 34rpx minmax(0, 1fr);
  align-items: center;
  padding: 18rpx;
  border: 1rpx solid rgba(180, 83, 9, 0.18);
  border-radius: 20rpx;
  background: rgba(255, 250, 240, 0.9);
}

.approval-chain__node {
  min-width: 0;
  padding: 14rpx 12rpx;
  border: 1rpx dashed rgba(180, 83, 9, 0.22);
  border-radius: 14rpx;
  background: #fffbeb;
}

.approval-chain__node--focus {
  color: #fffaf0;
  border-color: transparent;
  background: linear-gradient(135deg, #7f1d1d, #b45309);
}

.approval-chain__node text {
  display: block;
  overflow: hidden;
  font-size: 21rpx;
  font-weight: 800;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.approval-chain__node text:first-child {
  margin-bottom: 6rpx;
  font-size: 25rpx;
  font-weight: 900;
}

.approval-chain__arrow {
  height: 2rpx;
  background: linear-gradient(90deg, #b45309, rgba(180, 83, 9, 0.2));
}

.approval-card {
  border-color: rgba(180, 83, 9, 0.18);
  background: #fffaf0;
}

.risk-seal-card {
  display: grid;
  grid-template-columns: 92rpx minmax(0, 1fr);
  gap: 18rpx;
  align-items: center;
  padding: 22rpx;
  border: 1rpx solid rgba(185, 28, 28, 0.22);
  border-radius: 22rpx;
  background: linear-gradient(135deg, #7f1d1d, #b91c1c);
  color: #fffaf0;
}

.risk-seal-card__mark {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 92rpx;
  height: 92rpx;
  border: 4rpx solid rgba(255, 250, 240, 0.56);
  border-radius: 20rpx;
  font-size: 50rpx;
  font-weight: 900;
  transform: rotate(-10deg);
}

.risk-seal-card__title {
  display: block;
  font-size: 29rpx;
  font-weight: 900;
}

.risk-seal-card__desc {
  display: block;
  margin-top: 8rpx;
  color: rgba(255, 250, 240, 0.78);
  font-size: 24rpx;
  line-height: 1.45;
}
</style>
