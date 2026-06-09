import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as loginApi } from '@/api'
import router from '@/router'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const userInfo = ref<any>(null)

  const isLoggedIn = computed(() => !!token.value)
  const username = computed(() => userInfo.value?.username || userInfo.value?.realName || '')

  async function login(loginData: { username: string; password: string }) {
    const res = await loginApi(loginData)
    token.value = res.token
    userInfo.value = res.userInfo || res.user
    localStorage.setItem('token', res.token)
    localStorage.setItem('userInfo', JSON.stringify(res.userInfo || res.user))
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
    router.push('/login')
  }

  function restoreUserInfo() {
    const saved = localStorage.getItem('userInfo')
    if (saved) {
      try { userInfo.value = JSON.parse(saved) } catch { userInfo.value = null }
    }
  }

  return { token, userInfo, isLoggedIn, username, login, logout, restoreUserInfo }
})
