import http from './client'
import type { LoginRequest, LoginResponse, UserInfo } from './types'

export const authApi = {
  login(data: LoginRequest): Promise<LoginResponse> {
    return http.post('/auth/login', data)
  },
  refresh(token: string): Promise<LoginResponse> {
    return http.post('/auth/refresh', { token })
  },
  logout(): Promise<void> {
    return http.post('/auth/logout')
  },
  getCurrentUser(): Promise<UserInfo> {
    return http.get('/auth/current')
  },
}
