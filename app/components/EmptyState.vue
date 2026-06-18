<template>
  <view class="empty-state">
    <up-empty :text="text" mode="data" />
  </view>
</template>

<script setup>
defineProps({
  text: { type: String, default: '暂无数据' }
})
</script>

<style scoped>
.empty-state {
  padding: 80rpx 0;
}
</style>

