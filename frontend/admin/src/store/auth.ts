import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { UserInfo, LoginRequest } from '@/api/types'
import { authApi } from '@/api/auth'
import router from '@/router'

function normalizeUserInfo(raw: unknown): UserInfo | null {
  if (!raw || typeof raw !== 'object') return null
  const obj = raw as Record<string, unknown>
  const username = (obj.username as string) || ''
  if (!username) return null
  return {
    userId: (obj.userId as string) || (obj.id as string) || '',
    username,
    realName: (obj.realName as string) || '',
    avatar: (obj.avatar as string) || '',
    role: (obj.role as string) || '',
    deptId: (obj.deptId as string) || '',
    deptName: (obj.deptName as string) || '',
    perms: (obj.perms as string[]) || [],
  }
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const userInfo = ref<UserInfo | null>(null)

  const isLoggedIn = computed(() => !!token.value)
  const username = computed(() => userInfo.value?.username || userInfo.value?.realName || '')
  const role = computed(() => userInfo.value?.role || '')

  async function login(loginData: LoginRequest) {
    const res = await authApi.login(loginData)
    token.value = res.token
    userInfo.value = normalizeUserInfo(res.userInfo) || null
    localStorage.setItem('token', res.token)
    if (userInfo.value) {
      localStorage.setItem('userInfo', JSON.stringify(userInfo.value))
    }
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
    if (!saved) return
    try {
      const parsed = JSON.parse(saved)
      userInfo.value = normalizeUserInfo(parsed.userInfo ?? parsed.user ?? parsed)
    } catch {
      userInfo.value = null
    }
  }

  async function ensureUserInfo() {
    if (userInfo.value?.username) return
    restoreUserInfo()
    if (userInfo.value?.username || !token.value) return
    try {
      const info = await authApi.getCurrentUser()
      userInfo.value = normalizeUserInfo(info)
      if (userInfo.value) {
        localStorage.setItem('userInfo', JSON.stringify(userInfo.value))
      }
    } catch {
      // token 失效时由 API 拦截器跳转登录
    }
  }

  return { token, userInfo, isLoggedIn, username, role, login, logout, restoreUserInfo, ensureUserInfo }
})
