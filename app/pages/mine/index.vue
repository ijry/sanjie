<template>
  <view class="page stack">
    <view class="profile panel">
      <up-avatar :text="avatarText" size="56" bg-color="#a93226" />
      <view>
        <text class="profile__name">{{ user.current?.nickname || '未登录' }}</text>
        <text class="profile__role">{{ roleName(user.current?.role) }}</text>
      </view>
    </view>

    <up-cell-group>
      <up-cell title="切换身份" is-link url="/pages/login/index" />
      <up-cell title="审计日志" is-link url="/pages/audit/index" />
      <up-cell title="接口地址" :value="API_BASE_URL" />
      <up-cell title="项目文档" value="docs-site" />
    </up-cell-group>
  </view>
</template>

<script setup>
import { computed } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { API_BASE_URL } from '../../api/http'
import { useUserStore } from '../../stores/user'
import { roleName } from '../../utils/auth'

const user = useUserStore()
const avatarText = computed(() => (user.current?.nickname || '三界').slice(0, 1))

onShow(() => {
  user.fetchMe()
})
</script>

<style scoped>
.profile {
  display: flex;
  align-items: center;
  gap: 24rpx;
}

.profile__name {
  display: block;
  font-size: 36rpx;
  font-weight: 800;
}

.profile__role {
  display: block;
  margin-top: 8rpx;
  color: #756f63;
}
</style>

