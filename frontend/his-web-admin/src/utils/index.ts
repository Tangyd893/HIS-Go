import type { App } from 'vue'
import dayjs from 'dayjs'

export function setupUtils(app: App) {
  app.config.globalProperties.$dayjs = dayjs
}

export function formatDate(date: string | Date, format = 'YYYY-MM-DD HH:mm:ss'): string {
  return dayjs(date).format(format)
}

export function formatMoney(amount: number): string {
  return `¥${(amount / 100).toFixed(2)}`
}

export const statusMap: Record<number, { label: string; color: string }> = {
  0: { label: '待处理', color: 'blue' },
  1: { label: '处理中', color: 'orange' },
  2: { label: '已完成', color: 'green' },
  3: { label: '已取消', color: 'red' },
  4: { label: '已退费', color: 'magenta' },
}

export const genderMap: Record<string, string> = {
  M: '男',
  F: '女',
}

export const payMethodMap: Record<number, string> = {
  0: '现金',
  1: '微信支付',
  2: '支付宝',
  3: '银行卡',
  4: '医保卡',
}
