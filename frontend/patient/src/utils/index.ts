import dayjs from 'dayjs'

export function formatDate(date: string | Date, format = 'YYYY-MM-DD'): string {
  return dayjs(date).format(format)
}

export const genderMap: Record<string, string> = { M: '男', F: '女' }
