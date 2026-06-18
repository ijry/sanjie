<template>
  <view class="page stack">
    <view class="panel">
      <SectionTitle title="选择演示身份" />
      <up-cell-group :border="false">
        <up-cell v-for="item in users" :key="item.id" :title="item.name" :label="item.role" clickable @click="selectUser(item.id)" />
      </up-cell-group>
    </view>
  </view>
</template>

<script setup>
import SectionTitle from '../../components/SectionTitle.vue'
import { useUserStore } from '../../stores/user'
import { showError, showSuccess } from '../../utils/format'

const store = useUserStore()
const users = [
  { id: 1, name: '玉帝', role: '三界超级管理员' },
  { id: 2, name: '阎王', role: '地府域管理员' },
  { id: 3, name: '判官', role: '审批与复核' },
  { id: 4, name: '黑无常', role: '勾魂外勤' },
  { id: 5, name: '孟婆', role: '投胎前置流程' }
]

async function selectUser(id) {
  try {
    await store.switchUser(id)
    showSuccess('已切换')
    uni.navigateBack()
  } catch (error) {
    showError(error)
  }
}
</script>

