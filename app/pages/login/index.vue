<template>
  <view class="page stack dense-page login-page">
    <view class="identity-hero">
      <text class="identity-hero__eyebrow">ROLE SWITCHER</text>
      <text class="identity-hero__title">仙籍身份牌</text>
      <text class="identity-hero__desc">演示系统不做真实登录，直接切换三界角色，观察不同岗位的推荐入口和业务闭环。</text>
      <view class="identity-hero__stats">
        <view>
          <text>{{ users.length }}</text>
          <text>演示身份</text>
        </view>
        <view>
          <text>{{ activeUser?.name || '-' }}</text>
          <text>当前值守</text>
        </view>
      </view>
    </view>

    <view class="role-deck">
      <view
        v-for="item in users"
        :key="item.id"
        class="role-card"
        :class="[{ 'role-card--active': item.id === activeUserId }, `role-card--${item.tone}`]"
        @click="selectUser(item.id)"
      >
        <view class="role-card__top">
          <view class="role-card__avatar">
            <text>{{ item.name.slice(0, 1) }}</text>
          </view>
          <text class="role-card__status">{{ item.id === activeUserId ? '当前' : '切换' }}</text>
        </view>
        <text class="role-card__name">{{ item.name }}</text>
        <text class="role-card__role">{{ item.role }}</text>
        <view class="role-card__route">
          <text>{{ item.scope }}</text>
        </view>
      </view>
    </view>

    <view class="switch-guide">
      <text class="switch-guide__title">建议演示方式</text>
      <view class="switch-guide__row">
        <text>玉帝 / 判官</text>
        <text>处理审批、生死簿、审计闭环。</text>
      </view>
      <view class="switch-guide__row">
        <text>黑无常 / 城隍</text>
        <text>走勾魂外勤任务，体验异常上报。</text>
      </view>
      <view class="switch-guide__row">
        <text>孟婆</text>
        <text>进入汤坊库存台和投胎发汤流程。</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useUserStore } from '../../stores/user'
import { showError, showSuccess } from '../../utils/format'

const store = useUserStore()
const users = [
  { id: 1, name: '玉帝', role: '三界超级管理员', scope: '总控、审批、命簿', tone: 'gold' },
  { id: 2, name: '阎王', role: '地府域管理员', scope: '地狱、刑期复核', tone: 'red' },
  { id: 3, name: '判官', role: '审批与复核', scope: '案牍、改寿、转交', tone: 'amber' },
  { id: 4, name: '黑无常', role: '勾魂外勤', scope: '夜巡、异常上报', tone: 'blue' },
  { id: 5, name: '孟婆', role: '投胎前置流程', scope: '汤坊、发汤、残留', tone: 'teal' },
  { id: 6, name: '巨灵神', role: '天兵统领', scope: '巡防、告警、值守', tone: 'ink' },
  { id: 7, name: '城隍爷', role: '地方协同', scope: '协查、勾魂、投胎', tone: 'paper' },
  { id: 8, name: '巡天审计使', role: '审计神仙', scope: '天眼、追责、风险', tone: 'violet' }
]

const activeUserId = computed(() => Number(store.current?.id || 1))
const activeUser = computed(() => users.find(item => item.id === activeUserId.value))

async function selectUser(id) {
  try {
    await store.switchUser(id)
    showSuccess('已切换')
    uni.switchTab({ url: '/pages/mine/index' })
  } catch (error) {
    showError(error)
  }
}

onShow(() => {
  store.fetchMe()
})
</script>

<style scoped>
.login-page {
  background:
    radial-gradient(circle at 20% 0%, rgba(234, 179, 8, 0.18), transparent 27%),
    radial-gradient(circle at 86% 10%, rgba(124, 58, 237, 0.14), transparent 25%),
    linear-gradient(180deg, #17120c 0%, #3f2a12 32%, #f7f1e3 32%, #f7f1e3 100%);
}

.identity-hero {
  position: relative;
  overflow: hidden;
  padding: 28rpx;
  color: #fffaf0;
  border: 1rpx solid rgba(251, 191, 36, 0.26);
  border-radius: 32rpx;
  background:
    linear-gradient(90deg, rgba(251, 191, 36, 0.08) 1rpx, transparent 1rpx),
    linear-gradient(135deg, #111827, #78350f 58%, #7f1d1d);
  background-size: 36rpx 36rpx;
  box-shadow: 0 20rpx 46rpx rgba(49, 28, 13, 0.22);
}

.identity-hero::after {
  content: "";
  position: absolute;
  right: -48rpx;
  top: -52rpx;
  width: 180rpx;
  height: 180rpx;
  border: 8rpx double rgba(255, 250, 240, 0.14);
  border-radius: 50%;
}

.identity-hero__eyebrow,
.identity-hero__title,
.identity-hero__desc {
  position: relative;
  z-index: 1;
  display: block;
}

.identity-hero__eyebrow {
  color: #fbbf24;
  font-size: 20rpx;
  font-weight: 900;
  letter-spacing: 2rpx;
}

.identity-hero__title {
  margin-top: 6rpx;
  font-size: 46rpx;
  font-weight: 900;
}

.identity-hero__desc {
  width: 78%;
  margin-top: 10rpx;
  color: rgba(255, 250, 240, 0.72);
  font-size: 24rpx;
  line-height: 1.45;
}

.identity-hero__stats {
  position: relative;
  z-index: 1;
  display: grid;
  grid-template-columns: 160rpx minmax(0, 1fr);
  gap: 10rpx;
  margin-top: 22rpx;
}

.identity-hero__stats view {
  min-width: 0;
  padding: 14rpx 12rpx;
  border: 1rpx solid rgba(255, 250, 240, 0.14);
  border-radius: 16rpx;
  background: rgba(255, 250, 240, 0.08);
}

.identity-hero__stats text:first-child,
.identity-hero__stats text:last-child {
  display: block;
}

.identity-hero__stats text:first-child {
  overflow: hidden;
  color: #fef3c7;
  font-size: 30rpx;
  font-weight: 900;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.identity-hero__stats text:last-child {
  margin-top: 4rpx;
  color: rgba(255, 250, 240, 0.66);
  font-size: 20rpx;
}

.role-deck {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14rpx;
}

.role-card {
  position: relative;
  overflow: hidden;
  min-width: 0;
  padding: 20rpx;
  color: #fffaf0;
  border: 1rpx solid rgba(255, 250, 240, 0.14);
  border-radius: 24rpx;
  background: linear-gradient(135deg, #25211b, #0f0d0a);
  box-shadow: 0 12rpx 28rpx rgba(37, 33, 27, 0.12);
}

.role-card--active {
  outline: 5rpx solid rgba(251, 191, 36, 0.42);
}

.role-card--gold { background: linear-gradient(135deg, #b45309, #7c2d12); }
.role-card--red { background: linear-gradient(135deg, #a93226, #450a0a); }
.role-card--amber { background: linear-gradient(135deg, #92400e, #3f2a12); }
.role-card--blue { background: linear-gradient(135deg, #0284c7, #0f172a); }
.role-card--teal { background: linear-gradient(135deg, #0f766e, #134e4a); }
.role-card--ink { background: linear-gradient(135deg, #25211b, #111827); }
.role-card--paper { background: linear-gradient(135deg, #78350f, #7c2d12); }
.role-card--violet { background: linear-gradient(135deg, #7c3aed, #312e81); }

.role-card__top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12rpx;
}

.role-card__avatar {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 66rpx;
  height: 66rpx;
  border: 3rpx double rgba(255, 250, 240, 0.36);
  border-radius: 50%;
  font-size: 30rpx;
  font-weight: 900;
}

.role-card__status {
  padding: 6rpx 12rpx;
  border: 1rpx solid rgba(255, 250, 240, 0.2);
  border-radius: 999rpx;
  background: rgba(255, 250, 240, 0.1);
  font-size: 20rpx;
  font-weight: 900;
}

.role-card__name,
.role-card__role {
  display: block;
}

.role-card__name {
  margin-top: 18rpx;
  font-size: 34rpx;
  font-weight: 900;
}

.role-card__role {
  margin-top: 6rpx;
  color: rgba(255, 250, 240, 0.72);
  font-size: 22rpx;
}

.role-card__route {
  margin-top: 16rpx;
  padding: 12rpx;
  border: 1rpx dashed rgba(255, 250, 240, 0.22);
  border-radius: 14rpx;
  background: rgba(255, 250, 240, 0.08);
}

.role-card__route text {
  color: rgba(255, 250, 240, 0.82);
  font-size: 21rpx;
  font-weight: 800;
}

.switch-guide {
  padding: 22rpx;
  border: 1rpx solid rgba(120, 53, 15, 0.16);
  border-radius: 24rpx;
  background: rgba(255, 250, 240, 0.94);
}

.switch-guide__title {
  display: block;
  color: #25211b;
  font-size: 31rpx;
  font-weight: 900;
}

.switch-guide__row {
  display: grid;
  grid-template-columns: 180rpx minmax(0, 1fr);
  gap: 14rpx;
  margin-top: 12rpx;
  padding: 14rpx;
  border: 1rpx dashed rgba(169, 50, 38, 0.18);
  border-radius: 16rpx;
  background: #fff7ed;
}

.switch-guide__row text:first-child {
  color: #a93226;
  font-size: 22rpx;
  font-weight: 900;
}

.switch-guide__row text:last-child {
  color: #756f63;
  font-size: 22rpx;
  line-height: 1.4;
}
</style>
