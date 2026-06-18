import { defineStore } from 'pinia'
import { getMe } from '../api/users'

export const useUserStore = defineStore('user', {
  state: () => ({
    current: null
  }),
  actions: {
    async fetchMe() {
      this.current = await getMe()
      return this.current
    },
    async switchUser(userId) {
      uni.setStorageSync('sanjie_user_id', userId)
      return this.fetchMe()
    }
  }
})

