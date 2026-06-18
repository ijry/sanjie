<template>
  <view class="fixed-action-bar">
    <slot />
  </view>
</template>

<style scoped>
.fixed-action-bar {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  gap: 16rpx;
  padding: 18rpx 24rpx calc(18rpx + env(safe-area-inset-bottom));
  background: #fffaf0;
  border-top: 1rpx solid #e5d8bd;
  box-sizing: border-box;
  z-index: 10;
}
</style>

