import axios, { type AxiosInstance, type AxiosResponse, type InternalAxiosRequestConfig } from 'axios'
import { message } from 'ant-design-vue'
import { useAuthStore } from '@/store/auth'
import router from '@/router'

const http: AxiosInstance = axios.create({
  baseURL: '/api',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
})

http.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const authStore = useAuthStore()
    if (authStore.token) {
      config.headers.Authorization = `Bearer ${authStore.token}`
    }
    return config
  },
  (error) => Promise.reject(error),
)

http.interceptors.response.use(
  (response: AxiosResponse) => {
    const { code, message: msg, data } = response.data
    if (code === 0) {
      return data
    }
    if (code === 10001 || code === 10002) {
      const authStore = useAuthStore()
      authStore.logout()
      router.push('/login')
      return Promise.reject(new Error(msg || '登录已过期'))
    }
    message.error(msg || '请求失败')
    return Promise.reject(new Error(msg || '请求失败'))
  },
  (error) => {
    if (error.response?.status === 401) {
      const authStore = useAuthStore()
      authStore.logout()
      router.push('/login')
    }
    message.error(error.message || '网络异常')
    return Promise.reject(error)
  },
)

export default http
