<template>
  <view class="metric-card">
    <text class="metric-card__label">{{ label }}</text>
    <text class="metric-card__value">{{ value }}</text>
  </view>
</template>

<script setup>
defineProps({
  label: { type: String, required: true },
  value: { type: [String, Number], required: true }
})
</script>

<style scoped lang="scss">
.metric-card {
  background: #fffaf0;
  border: 1rpx solid #e5d8bd;
  border-radius: 8rpx;
  padding: 22rpx;
}

.metric-card__label {
  display: block;
  color: #756f63;
  font-size: 24rpx;
}

.metric-card__value {
  display: block;
  margin-top: 10rpx;
  color: #25211b;
  font-size: 44rpx;
  font-weight: 700;
}
</style>

