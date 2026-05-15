import { create } from 'zustand'
import type { UserInfo, LoginRequest } from '../api/types'
import { authApi } from '../api/auth'

interface AuthState {
  token: string
  userInfo: UserInfo | null
  isLoggedIn: boolean
  username: string
  role: string
  login: (data: LoginRequest) => Promise<void>
  logout: () => void
  restoreUserInfo: () => void
}

export const useAuthStore = create<AuthState>((set, get) => ({
  token: localStorage.getItem('token') || '',
  userInfo: (() => {
    try {
      const saved = localStorage.getItem('userInfo')
      return saved ? JSON.parse(saved) : null
    } catch {
      return null
    }
  })(),
  get isLoggedIn() {
    return !!get().token
  },
  get username() {
    return get().userInfo?.username || ''
  },
  get role() {
    return get().userInfo?.role || ''
  },
  async login(data: LoginRequest) {
    const res = await authApi.login(data)
    localStorage.setItem('token', res.token)
    localStorage.setItem('userInfo', JSON.stringify(res.user))
    set({ token: res.token, userInfo: res.user })
  },
  logout() {
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
    set({ token: '', userInfo: null })
  },
  restoreUserInfo() {
    const saved = localStorage.getItem('userInfo')
    if (saved) {
      try {
        set({ userInfo: JSON.parse(saved) })
      } catch {
        set({ userInfo: null })
      }
    }
  },
}))
