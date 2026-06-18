<template>
  <view class="section-title">
    <text>{{ title }}</text>
  </view>
</template>

<script setup>
defineProps({
  title: { type: String, required: true }
})
</script>

<style scoped>
.section-title {
  margin: 12rpx 0;
  color: #25211b;
  font-size: 30rpx;
  font-weight: 700;
}
</style>

