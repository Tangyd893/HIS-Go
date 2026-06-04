import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { UserInfo, LoginRequest } from '@/api/types'
import { authApi } from '@/api/auth'
import router from '@/router'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const userInfo = ref<UserInfo | null>(null)

  const isLoggedIn = computed(() => !!token.value)
  const username = computed(() => userInfo.value?.username || '')
  const role = computed(() => userInfo.value?.role || '')

  async function login(loginData: LoginRequest) {
    const res = await authApi.login(loginData)
    token.value = res.token
    userInfo.value = res.user
    localStorage.setItem('token', res.token)
    localStorage.setItem('userInfo', JSON.stringify(res.user))
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
      try {
        userInfo.value = JSON.parse(saved)
      } catch {
        userInfo.value = null
      }
    }
  }

  return { token, userInfo, isLoggedIn, username, role, login, logout, restoreUserInfo }
})
