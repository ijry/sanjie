<template>
  <view class="page stack detail-page" v-if="item">
    <view class="detail-hero">
      <view class="detail-hero__top">
        <view>
          <text class="detail-hero__eyebrow">LIFE BOOK #{{ item.id }}</text>
          <text class="detail-hero__title">{{ item.soulName }}</text>
        </view>
        <StatusTag :value="item.riskFlag" />
      </view>
      <text class="detail-hero__desc">{{ lifeSummary }}</text>
      <view class="detail-hero__stats">
        <view class="detail-stat">
          <text class="detail-stat__label">原定</text>
          <text class="detail-stat__value">{{ item.expectedLifespan }} 岁</text>
        </view>
        <view class="detail-stat">
          <text class="detail-stat__label">当前</text>
          <text class="detail-stat__value">{{ item.actualLifespan }} 岁</text>
        </view>
        <view class="detail-stat">
          <text class="detail-stat__label">偏差</text>
          <text class="detail-stat__value">{{ lifespanDeltaText }}</text>
        </view>
      </view>
    </view>

    <view class="detail-card">
      <text class="detail-card__title">命簿档案</text>
      <view class="detail-grid">
        <view class="detail-field">
          <text class="detail-field__label">魂魄 ID</text>
          <text class="detail-field__value">#{{ item.soulId }}</text>
        </view>
        <view class="detail-field">
          <text class="detail-field__label">命格</text>
          <text class="detail-field__value">{{ item.fateLevel }}</text>
        </view>
        <view class="detail-field">
          <text class="detail-field__label">劫数</text>
          <text class="detail-field__value">{{ item.calamityCount }} 次</text>
        </view>
        <view class="detail-field">
          <text class="detail-field__label">冻结状态</text>
          <text class="detail-field__value">{{ item.locked ? '已冻结' : '可操作' }}</text>
        </view>
        <view class="detail-field">
          <text class="detail-field__label">风险等级</text>
          <text class="detail-field__value">{{ riskText }}</text>
        </view>
        <view class="detail-field">
          <text class="detail-field__label">更新时间</text>
          <text class="detail-field__value">{{ displayTime(item.updatedAt) }}</text>
        </view>
        <view class="detail-field detail-field--wide">
          <text class="detail-field__label">登记死因</text>
          <text class="detail-field__value">{{ item.deathReason }}</text>
        </view>
      </view>
    </view>

    <view class="detail-card">
      <text class="detail-card__title">处置流程</text>
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
        <view class="detail-field">
          <text class="detail-field__label">当前建议</text>
          <text class="detail-field__value">{{ actionHint }}</text>
        </view>
        <view class="detail-field">
          <text class="detail-field__label">锁定策略</text>
          <text class="detail-field__value">{{ lockHint }}</text>
        </view>
        <view class="detail-field detail-field--wide">
          <text class="detail-field__label">系统判断</text>
          <text class="detail-field__value">{{ riskHint }}</text>
        </view>
      </view>
    </view>

    <view class="detail-card">
      <text class="detail-card__title">审计记录</text>
      <EmptyState v-if="logs.length === 0" text="暂无审计记录" />
      <view v-else class="mini-list">
        <view v-for="log in logs" :key="log.id" class="mini-row">
          <view class="mini-row__main">
            <text class="mini-row__title">{{ log.action }}</text>
            <text class="mini-row__meta">{{ log.note || '无备注' }} · {{ displayTime(log.createdAt) }}</text>
          </view>
        </view>
      </view>
    </view>

    <FixedActionBar>
      <up-button v-if="!item.locked" type="warning" text="冻结" @click="freeze" />
      <up-button type="primary" text="申请改寿" @click="changeLife" />
    </FixedActionBar>
  </view>

  <view v-else class="page">
    <EmptyState text="生死簿加载中" />
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import FixedActionBar from '../../components/FixedActionBar.vue'
import StatusTag from '../../components/StatusTag.vue'
import EmptyState from '../../components/EmptyState.vue'
import { changeLifespan, freezeLifeBookRecord, getLifeBookAuditLogs, getLifeBookRecord } from '../../api/lifeBook'
import { displayTime, showError, showSuccess } from '../../utils/format'
import { labelOf } from '../../utils/status'

const id = ref('')
const item = ref(null)
const logs = ref([])

const lifespanDelta = computed(() => {
  if (!item.value) return 0
  return Number(item.value.actualLifespan || 0) - Number(item.value.expectedLifespan || 0)
})

const lifespanDeltaText = computed(() => {
  if (lifespanDelta.value === 0) return '持平'
  return `${lifespanDelta.value > 0 ? '+' : ''}${lifespanDelta.value} 年`
})

const riskText = computed(() => labelOf(item.value?.riskFlag))

const lifeSummary = computed(() => {
  if (!item.value) return ''
  return `${item.value.deathReason} · ${item.value.fateLevel} · ${item.value.locked ? '记录已冻结' : '记录可申请变更'}`
})

const processSteps = computed(() => [
  { index: '1', title: '核对命簿', desc: `比对原定 ${item.value?.expectedLifespan || '-'} 岁与当前 ${item.value?.actualLifespan || '-'} 岁，确认是否存在阳寿偏差。` },
  { index: '2', title: '风险处置', desc: item.value?.locked ? '当前记录已冻结，后续变更必须走审批链路。' : '如发现关系户续命、误划名等风险，可先冻结记录。' },
  { index: '3', title: '发起改寿', desc: '改寿不会直接修改生死簿，会生成审批单，审批通过后由后端写回阳寿。' }
])

const actionHint = computed(() => {
  if (!item.value) return '-'
  if (item.value.riskFlag === 'critical') return '优先冻结并交由判官复核'
  if (item.value.riskFlag === 'warning') return '建议核对死因和阳寿偏差'
  return '低风险，可按需发起改寿申请'
})

const lockHint = computed(() => {
  if (!item.value) return '-'
  return item.value.locked ? '已锁定，等待审批或审计处理' : '未锁定，可冻结或申请改寿'
})

const riskHint = computed(() => {
  if (!item.value) return '-'
  if (lifespanDelta.value > 20) return '当前阳寿明显高于原定，疑似香火充值或关系户续命。'
  if (lifespanDelta.value < -10) return '当前阳寿明显短于原定，需排查误勾魂或提前划名。'
  if (item.value.calamityCount >= 8) return '劫数较多，建议检查是否进入地狱或投胎复核链路。'
  return '阳寿偏差处于可解释范围，仍需保留审计记录。'
})

onLoad(query => {
  id.value = query.id
})

async function load() {
  if (!id.value) return
  try {
    item.value = await getLifeBookRecord(id.value)
    logs.value = await getLifeBookAuditLogs(id.value)
  } catch (error) {
    showError(error)
  }
}

async function freeze() {
  try {
    await freezeLifeBookRecord(id.value)
    showSuccess('已冻结')
    await load()
  } catch (error) {
    showError(error)
  }
}

function changeLife() {
  uni.showModal({
    title: '寿命变更',
    content: '将申请把当前阳寿调整为 88 岁',
    success: async res => {
      if (!res.confirm) return
      try {
        await changeLifespan(id.value, { newLifespan: 88, reason: '移动端发起寿命变更' })
        showSuccess('已发起审批')
        await load()
      } catch (error) {
        showError(error)
      }
    }
  })
}

onShow(load)
</script>

<style scoped>
.detail-page {
  padding-bottom: 140rpx;
}
</style>
