<template>
  <view class="page stack dense-page audit-page">
    <view class="audit-hero">
      <view class="audit-eye">
        <view class="audit-eye__iris">
          <text>{{ items.length }}</text>
        </view>
      </view>
      <view class="audit-hero__main">
        <text class="audit-hero__eyebrow">CELESTIAL TRACE</text>
        <text class="audit-hero__title">天眼留痕</text>
        <text class="audit-hero__desc">改寿、插队、分流、发汤都留下轨迹；谁动了命簿，一眼能查。</text>
      </view>
      <view class="audit-hero__stats">
        <view class="audit-stat">
          <text>{{ approvalCount }}</text>
          <text>审批轨迹</text>
        </view>
        <view class="audit-stat">
          <text>{{ lifeBookCount }}</text>
          <text>命簿操作</text>
        </view>
        <view class="audit-stat">
          <text>{{ captureCount }}</text>
          <text>勾魂外勤</text>
        </view>
      </view>
    </view>

    <view class="audit-filter">
      <view
        v-for="tab in targetTabs"
        :key="tab.value"
        class="audit-filter__chip"
        :class="{ 'audit-filter__chip--active': currentTarget === tab.value }"
        @click="currentTarget = tab.value"
      >
        <text>{{ tab.name }}</text>
        <text>{{ countByTarget(tab.value) }}</text>
      </view>
    </view>

    <EmptyState v-if="filteredItems.length === 0" text="暂无审计日志" />

    <view v-else class="trace-timeline">
      <view v-for="item in filteredItems" :key="item.id" class="trace-card" :class="`trace-card--${actionTone(item.action)}`">
        <view class="trace-card__rail">
          <text>{{ actionMark(item.action) }}</text>
        </view>
        <view class="trace-card__main">
          <view class="trace-card__head">
            <view>
              <text class="trace-card__action">{{ actionLabel(item.action) }}</text>
              <text class="trace-card__target">{{ targetLabel(item.targetType) }} #{{ item.targetId }}</text>
            </view>
            <text class="trace-card__actor">{{ actorName(item.actorId) }}</text>
          </view>
          <text class="trace-card__note">{{ item.note || '无备注' }}</text>
          <view class="trace-card__foot">
            <text>{{ displayTime(item.createdAt) }}</text>
            <text>LOG-{{ item.id }}</text>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import EmptyState from '../../components/EmptyState.vue'
import { getAuditLogs } from '../../api/audit'
import { displayTime, showError } from '../../utils/format'

const items = ref([])
const currentTarget = ref('all')
const targetTabs = [
  { name: '全部', value: 'all' },
  { name: '审批', value: 'approval' },
  { name: '命簿', value: 'life_book' },
  { name: '勾魂', value: 'capture_task' },
  { name: '地狱', value: 'hell_floor' },
  { name: '愿望', value: 'wish' },
  { name: '汤坊', value: 'soup_inventory' }
]

const actorLabels = {
  1: '玉帝',
  2: '阎王',
  3: '判官',
  4: '黑无常',
  5: '孟婆',
  6: '巨灵神',
  7: '城隍爷',
  8: '巡天审计使'
}

const approvalCount = computed(() => countByTarget('approval'))
const captureCount = computed(() => countByTarget('capture_task'))
const lifeBookCount = computed(() => countByTarget('life_book'))
const filteredItems = computed(() => {
  if (currentTarget.value === 'all') return items.value
  return items.value.filter(item => item.targetType === currentTarget.value)
})

async function load() {
  try {
    items.value = await getAuditLogs()
  } catch (error) {
    showError(error)
  }
}

function countByTarget(targetType) {
  if (targetType === 'all') return items.value.length
  return items.value.filter(item => item.targetType === targetType).length
}

function actorName(actorId) {
  return actorLabels[actorId] || `用户 ${actorId || '-'}`
}

function targetLabel(targetType) {
  const labels = {
    approval: '审批单',
    life_book: '生死簿',
    capture_task: '勾魂任务',
    hell_floor: '地狱楼层',
    hell_sentence: '刑期记录',
    reincarnation: '投胎单',
    wish: '愿望工单',
    soup_inventory: '汤坊库存'
  }
  return labels[targetType] || targetType || '-'
}

function actionLabel(action) {
  const labels = {
    'life_book.freeze': '命簿冻结',
    'approval.create': '创建审批',
    'approval.approved': '审批通过',
    'approval.rejected': '审批驳回',
    'approval.transferred': '审批转交',
    'capture.exception': '勾魂异常',
    'capture.start': '勾魂开始',
    'hell.dispatch': '地狱分流',
    'hell.review': '刑期复核',
    'soup.issue': '孟婆汤发放',
    'soup.restock': '汤坊补货',
    'soup.adjust': '汤坊扣减',
    'soup.warning': '汤坊预警',
    'wish.route': '愿望转派',
    'wish.resolve': '愿望办结',
    'wish.rejected': '愿望驳回',
    'reincarnation.approved': '投胎通过',
    'audit.warning': '审计预警'
  }
  return labels[action] || action || '-'
}

function actionTone(action) {
  if (action.includes('rejected') || action.includes('exception') || action.includes('warning')) return 'red'
  if (action.includes('approved') || action.includes('resolve') || action.includes('issue')) return 'green'
  if (action.includes('freeze') || action.includes('review') || action.includes('transferred')) return 'amber'
  return 'blue'
}

function actionMark(action) {
  if (action.includes('life_book')) return '命'
  if (action.includes('approval')) return '审'
  if (action.includes('capture')) return '魂'
  if (action.includes('hell')) return '狱'
  if (action.includes('soup')) return '汤'
  if (action.includes('wish')) return '愿'
  return '眼'
}

onShow(load)
</script>

<style scoped>
.audit-page {
  background:
    radial-gradient(circle at 50% 0%, rgba(59, 130, 246, 0.18), transparent 27%),
    linear-gradient(180deg, #0b1120 0%, #1f2937 32%, #f7f1e3 32%, #f7f1e3 100%);
}

.audit-hero {
  display: grid;
  grid-template-columns: 150rpx minmax(0, 1fr);
  gap: 18rpx;
  padding: 26rpx;
  color: #eff6ff;
  border: 1rpx solid rgba(147, 197, 253, 0.24);
  border-radius: 30rpx;
  background:
    linear-gradient(90deg, rgba(147, 197, 253, 0.08) 1rpx, transparent 1rpx),
    linear-gradient(0deg, rgba(147, 197, 253, 0.08) 1rpx, transparent 1rpx),
    linear-gradient(135deg, #020617, #1e3a8a 58%, #111827);
  background-size: 38rpx 38rpx;
  box-shadow: 0 20rpx 46rpx rgba(15, 23, 42, 0.24);
}

.audit-eye {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 138rpx;
  height: 138rpx;
  border: 2rpx solid rgba(147, 197, 253, 0.32);
  border-radius: 72rpx 18rpx 72rpx 18rpx;
  transform: rotate(45deg);
  background: rgba(30, 64, 175, 0.3);
}

.audit-eye__iris {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 78rpx;
  height: 78rpx;
  border-radius: 50%;
  background: radial-gradient(circle, #93c5fd, #1d4ed8);
  transform: rotate(-45deg);
}

.audit-eye__iris text {
  color: #eff6ff;
  font-size: 34rpx;
  font-weight: 900;
}

.audit-hero__main {
  min-width: 0;
}

.audit-hero__eyebrow,
.audit-hero__title,
.audit-hero__desc {
  display: block;
}

.audit-hero__eyebrow {
  color: #93c5fd;
  font-size: 20rpx;
  font-weight: 900;
  letter-spacing: 2rpx;
}

.audit-hero__title {
  margin-top: 6rpx;
  font-size: 44rpx;
  font-weight: 900;
}

.audit-hero__desc {
  margin-top: 10rpx;
  color: rgba(239, 246, 255, 0.72);
  font-size: 24rpx;
  line-height: 1.45;
}

.audit-hero__stats {
  grid-column: 1 / -1;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10rpx;
}

.audit-stat {
  padding: 14rpx 12rpx;
  border: 1rpx solid rgba(219, 234, 254, 0.15);
  border-radius: 16rpx;
  background: rgba(15, 23, 42, 0.42);
}

.audit-stat text:first-child {
  display: block;
  color: #dbeafe;
  font-size: 30rpx;
  font-weight: 900;
}

.audit-stat text:last-child {
  display: block;
  margin-top: 4rpx;
  color: rgba(219, 234, 254, 0.66);
  font-size: 20rpx;
}

.audit-filter {
  display: flex;
  gap: 10rpx;
  overflow-x: auto;
  padding: 6rpx 0;
}

.audit-filter__chip {
  display: flex;
  align-items: center;
  gap: 10rpx;
  flex-shrink: 0;
  padding: 12rpx 16rpx;
  border: 1rpx solid rgba(37, 99, 235, 0.14);
  border-radius: 999rpx;
  background: rgba(255, 250, 240, 0.86);
  color: #475569;
  font-size: 22rpx;
  font-weight: 900;
}

.audit-filter__chip--active {
  color: #eff6ff;
  background: #1d4ed8;
}

.trace-timeline {
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 14rpx;
}

.trace-card {
  display: grid;
  grid-template-columns: 70rpx minmax(0, 1fr);
  gap: 14rpx;
}

.trace-card__rail {
  position: relative;
  display: flex;
  justify-content: center;
}

.trace-card__rail::before {
  content: "";
  position: absolute;
  top: 58rpx;
  bottom: -24rpx;
  width: 2rpx;
  background: rgba(30, 64, 175, 0.18);
}

.trace-card__rail text {
  z-index: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 58rpx;
  height: 58rpx;
  color: #fffaf0;
  border-radius: 50%;
  background: #2563eb;
  font-size: 24rpx;
  font-weight: 900;
}

.trace-card--red .trace-card__rail text { background: #a93226; }
.trace-card--green .trace-card__rail text { background: #15803d; }
.trace-card--amber .trace-card__rail text { background: #b45309; }

.trace-card__main {
  min-width: 0;
  padding: 18rpx;
  border: 1rpx solid rgba(37, 99, 235, 0.14);
  border-radius: 20rpx;
  background: linear-gradient(135deg, #fffaf0, #eff6ff);
  box-shadow: 0 10rpx 26rpx rgba(15, 23, 42, 0.06);
}

.trace-card--red .trace-card__main {
  border-color: rgba(169, 50, 38, 0.2);
  background: linear-gradient(135deg, #fffaf0, #fef2f2);
}

.trace-card--green .trace-card__main {
  border-color: rgba(21, 128, 61, 0.18);
  background: linear-gradient(135deg, #fffaf0, #f0fdf4);
}

.trace-card--amber .trace-card__main {
  border-color: rgba(180, 83, 9, 0.2);
  background: linear-gradient(135deg, #fffaf0, #fffbeb);
}

.trace-card__head,
.trace-card__foot {
  display: flex;
  justify-content: space-between;
  gap: 12rpx;
}

.trace-card__action,
.trace-card__target,
.trace-card__actor,
.trace-card__note {
  display: block;
}

.trace-card__action {
  color: #0f172a;
  font-size: 30rpx;
  font-weight: 900;
}

.trace-card__target {
  margin-top: 5rpx;
  color: #64748b;
  font-size: 22rpx;
}

.trace-card__actor {
  flex-shrink: 0;
  color: #1d4ed8;
  font-size: 22rpx;
  font-weight: 900;
}

.trace-card__note {
  margin-top: 14rpx;
  padding: 12rpx;
  color: #334155;
  border-radius: 14rpx;
  background: rgba(255, 250, 240, 0.72);
  font-size: 23rpx;
  line-height: 1.45;
}

.trace-card__foot {
  margin-top: 12rpx;
  color: #64748b;
  font-size: 20rpx;
  font-weight: 800;
}
</style>
