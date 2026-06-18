<template>
  <view class="page stack detail-page hell-detail" v-if="item">
    <view class="hell-console">
      <view class="hell-console__glow"></view>
      <view class="hell-console__top">
        <view>
          <text class="hell-console__eyebrow">INFERNO CONTROL / SENTENCE #{{ item.id }}</text>
          <text class="hell-console__title">{{ item.soulName }}</text>
        </view>
        <StatusTag :value="item.reviewStatus" />
      </view>
      <text class="hell-console__desc">{{ item.floorName }} · {{ item.crimeType }} · {{ reviewText }}</text>
      <view class="hell-gauges">
        <view class="hell-gauge">
          <text class="hell-gauge__value">{{ item.sentenceDays }}</text>
          <text class="hell-gauge__label">刑期 / 天</text>
        </view>
        <view class="hell-gauge hell-gauge--hot">
          <text class="hell-gauge__value">{{ item.painLevel }}</text>
          <text class="hell-gauge__label">痛感 / 10</text>
        </view>
        <view class="hell-gauge">
          <text class="hell-gauge__value">{{ item.equipmentId }}</text>
          <text class="hell-gauge__label">设备编号</text>
        </view>
      </view>
    </view>

    <view class="hell-warning-panel">
      <view class="hell-warning-panel__bars">
        <text></text><text></text><text></text>
      </view>
      <view>
        <text class="hell-warning-panel__title">{{ sentenceLevel }} · {{ painHint }}</text>
        <text class="hell-warning-panel__desc">{{ reviewHint }}</text>
      </view>
    </view>

    <view class="detail-card hell-card">
      <text class="detail-card__title">刑期档案</text>
      <view class="detail-grid">
        <view class="detail-field">
          <text class="detail-field__label">魂魄 ID</text>
          <text class="detail-field__value">#{{ item.soulId }}</text>
        </view>
        <view class="detail-field">
          <text class="detail-field__label">楼层 ID</text>
          <text class="detail-field__value">#{{ item.floorId }}</text>
        </view>
        <view class="detail-field">
          <text class="detail-field__label">执行楼层</text>
          <text class="detail-field__value">{{ item.floorName }}</text>
        </view>
        <view class="detail-field">
          <text class="detail-field__label">复核状态</text>
          <text class="detail-field__value">{{ reviewText }}</text>
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
          <text class="detail-field__label">登记罪行</text>
          <text class="detail-field__value">{{ item.crimeType }}</text>
        </view>
      </view>
    </view>

    <view class="detail-card hell-card">
      <text class="detail-card__title">执行流程</text>
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

    <view class="detail-card hell-card">
      <text class="detail-card__title">复核建议</text>
      <view class="detail-grid">
        <view class="detail-field">
          <text class="detail-field__label">量刑等级</text>
          <text class="detail-field__value">{{ sentenceLevel }}</text>
        </view>
        <view class="detail-field">
          <text class="detail-field__label">痛感风险</text>
          <text class="detail-field__value">{{ painHint }}</text>
        </view>
        <view class="detail-field detail-field--wide">
          <text class="detail-field__label">系统判断</text>
          <text class="detail-field__value">{{ reviewHint }}</text>
        </view>
      </view>
    </view>

    <FixedActionBar v-if="item.reviewStatus === 'none'">
      <up-button type="warning" text="发起复核" @click="review" />
    </FixedActionBar>
  </view>

  <view v-else class="page">
    <EmptyState text="刑期记录加载中" />
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import EmptyState from '../../components/EmptyState.vue'
import FixedActionBar from '../../components/FixedActionBar.vue'
import StatusTag from '../../components/StatusTag.vue'
import { getHellSentence, reviewHellSentence } from '../../api/hell'
import { displayTime, showError, showSuccess } from '../../utils/format'
import { labelOf } from '../../utils/status'

const id = ref('')
const item = ref(null)

const reviewText = computed(() => labelOf(item.value?.reviewStatus))

const sentenceLevel = computed(() => {
  const days = Number(item.value?.sentenceDays || 0)
  if (days >= 700) return '长期服刑'
  if (days >= 180) return '标准刑期'
  return '短期体验'
})

const painHint = computed(() => {
  const pain = Number(item.value?.painLevel || 0)
  if (pain >= 8) return '高痛感'
  if (pain >= 5) return '中痛感'
  return '低痛感'
})

const processSteps = computed(() => [
  { index: '1', title: '判罚登记', desc: `${item.value?.soulName || '-'} 因“${item.value?.crimeType || '-'}”进入 ${item.value?.floorName || '-'}。` },
  { index: '2', title: '设备执行', desc: `执行设备 ${item.value?.equipmentId || '-'}，痛感等级 ${item.value?.painLevel || '-'}。` },
  { index: '3', title: '刑期复核', desc: item.value?.reviewStatus === 'none' ? '当前未发起复核，可提交审批单给判官或阎王。' : '复核流程已进入审批中心，等待处理结果。' }
])

const reviewHint = computed(() => {
  if (!item.value) return '-'
  if (item.value.reviewStatus !== 'none') return '复核已发起，避免重复提交。'
  if (item.value.sentenceDays >= 365 || item.value.painLevel >= 8) return '刑期或痛感偏高，建议发起复核。'
  return '刑期处于常规范围，可继续执行。'
})

onLoad(query => {
  id.value = query.id
})

async function load() {
  if (!id.value) return
  try {
    item.value = await getHellSentence(id.value)
  } catch (error) {
    showError(error)
  }
}

async function review() {
  try {
    await reviewHellSentence(id.value, { note: '移动端发起刑期复核' })
    showSuccess('已发起')
    await load()
  } catch (error) {
    showError(error)
  }
}

onShow(load)
</script>

<style scoped>
.detail-page {
  padding-bottom: 140rpx;
}

.hell-detail {
  background:
    radial-gradient(circle at 80% 0%, rgba(239, 68, 68, 0.2), transparent 26%),
    linear-gradient(180deg, #0f0a08 0%, #2b0f0b 38%, #f7f1e3 38%, #f7f1e3 100%);
}

.hell-console {
  position: relative;
  overflow: hidden;
  padding: 28rpx;
  color: #fff7ed;
  border: 1rpx solid rgba(248, 113, 113, 0.28);
  border-radius: 28rpx;
  background:
    linear-gradient(90deg, rgba(248, 113, 113, 0.08) 1rpx, transparent 1rpx),
    linear-gradient(0deg, rgba(248, 113, 113, 0.08) 1rpx, transparent 1rpx),
    radial-gradient(circle at 80% 20%, rgba(249, 115, 22, 0.3), transparent 28%),
    linear-gradient(135deg, #120806, #450a0a 56%, #111827);
  background-size: 36rpx 36rpx, 36rpx 36rpx, auto, auto;
  box-shadow: 0 22rpx 46rpx rgba(69, 10, 10, 0.3);
}

.hell-console__glow {
  position: absolute;
  right: -60rpx;
  bottom: -80rpx;
  width: 260rpx;
  height: 260rpx;
  border-radius: 50%;
  background: rgba(249, 115, 22, 0.24);
  filter: blur(8rpx);
}

.hell-console__top {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16rpx;
}

.hell-console__eyebrow {
  display: block;
  color: #fb923c;
  font-size: 20rpx;
  font-weight: 900;
  letter-spacing: 2rpx;
}

.hell-console__title {
  display: block;
  margin-top: 8rpx;
  font-size: 48rpx;
  font-weight: 900;
  line-height: 1.15;
}

.hell-console__desc {
  position: relative;
  z-index: 1;
  display: block;
  margin-top: 16rpx;
  color: rgba(255, 247, 237, 0.76);
  font-size: 25rpx;
  line-height: 1.45;
}

.hell-gauges {
  position: relative;
  z-index: 1;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10rpx;
  margin-top: 24rpx;
}

.hell-gauge {
  min-width: 0;
  padding: 18rpx 12rpx;
  border: 1rpx solid rgba(251, 146, 60, 0.3);
  border-radius: 16rpx;
  background: rgba(15, 23, 42, 0.56);
}

.hell-gauge--hot {
  background: linear-gradient(135deg, rgba(127, 29, 29, 0.84), rgba(249, 115, 22, 0.32));
}

.hell-gauge__value {
  display: block;
  overflow: hidden;
  color: #ffedd5;
  font-size: 34rpx;
  font-weight: 900;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.hell-gauge__label {
  display: block;
  margin-top: 8rpx;
  color: rgba(255, 237, 213, 0.62);
  font-size: 20rpx;
}

.hell-warning-panel {
  display: grid;
  grid-template-columns: 66rpx minmax(0, 1fr);
  gap: 16rpx;
  align-items: center;
  padding: 20rpx;
  color: #fff7ed;
  border-radius: 20rpx;
  background: linear-gradient(135deg, #7f1d1d, #111827);
  box-shadow: 0 14rpx 30rpx rgba(127, 29, 29, 0.18);
}

.hell-warning-panel__bars {
  display: flex;
  align-items: flex-end;
  gap: 6rpx;
  height: 58rpx;
}

.hell-warning-panel__bars text {
  flex: 1;
  border-radius: 999rpx;
  background: #fb923c;
}

.hell-warning-panel__bars text:nth-child(1) {
  height: 28rpx;
}

.hell-warning-panel__bars text:nth-child(2) {
  height: 46rpx;
}

.hell-warning-panel__bars text:nth-child(3) {
  height: 58rpx;
}

.hell-warning-panel__title {
  display: block;
  font-size: 28rpx;
  font-weight: 900;
}

.hell-warning-panel__desc {
  display: block;
  margin-top: 6rpx;
  color: rgba(255, 247, 237, 0.72);
  font-size: 23rpx;
  line-height: 1.42;
}

.hell-card {
  border-color: rgba(127, 29, 29, 0.18);
  background: linear-gradient(180deg, #fffaf0, #fff7ed);
}
</style>
