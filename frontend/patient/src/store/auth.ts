import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as loginApi, getMyPatient } from '@/api'
import router from '@/router'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const userInfo = ref<any>(null)
  const patientId = ref<string>(localStorage.getItem('patientId') || '')

  const isLoggedIn = computed(() => !!token.value)
  const username = computed(() => userInfo.value?.username || userInfo.value?.realName || '')

  async function fetchPatientId() {
    try {
      const patient = await getMyPatient()
      patientId.value = patient.id
      localStorage.setItem('patientId', patient.id)
    } catch {
      // 非患者角色（doctor/admin 等）无患者档案，静默失败
      patientId.value = ''
    }
  }

  async function login(loginData: { username: string; password: string }) {
    const res = await loginApi(loginData)
    token.value = res.token
    userInfo.value = res.userInfo || res.user
    localStorage.setItem('token', res.token)
    localStorage.setItem('userInfo', JSON.stringify(res.userInfo || res.user))
    // 登录后自动获取患者 ID
    await fetchPatientId()
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    patientId.value = ''
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
    localStorage.removeItem('patientId')
    router.push('/login')
  }

  function restoreUserInfo() {
    const saved = localStorage.getItem('userInfo')
    if (saved) {
      try { userInfo.value = JSON.parse(saved) } catch { userInfo.value = null }
    }
    // 恢复缓存的 patientId
    const savedPatientId = localStorage.getItem('patientId')
    if (savedPatientId) {
      patientId.value = savedPatientId
    }
  }

  return { token, userInfo, patientId, isLoggedIn, username, login, logout, restoreUserInfo }
})
