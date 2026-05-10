import http from './client'
import type { LoginRequest, LoginResponse, UserInfo } from './types'

export const authApi = {
  /** 登录 */
  login(data: LoginRequest): Promise<LoginResponse> {
    return http.post('/auth/login', data)
  },

  /** 刷新令牌 */
  refresh(token: string): Promise<LoginResponse> {
    return http.post('/auth/refresh', { token })
  },

  /** 登出 */
  logout(): Promise<void> {
    return http.post('/auth/logout')
  },

  /** 获取当前用户信息 */
  getCurrentUser(): Promise<UserInfo> {
    return http.get('/auth/current')
  },
}
