import { ref, type Ref } from 'vue'
import { message } from 'ant-design-vue'
import type { PageQuery, PageData } from '@/api/types'

/** 通用分页查询组合式函数 */
export function usePagination<T>(
  fetchFn: (params: PageQuery & Record<string, any>) => Promise<PageData<T>>,
  extraParams: Ref<Record<string, any>> = ref({}),
) {
  const loading = ref(false)
  const dataSource = ref<T[]>([])
  const total = ref(0)
  const current = ref(1)
  const pageSize = ref(10)

  async function fetch() {
    loading.value = true
    try {
      const res = await fetchFn({
        page: current.value,
        pageSize: pageSize.value,
        ...extraParams.value,
      })
      dataSource.value = res?.list || []
      total.value = res?.total || 0
    } catch {
      dataSource.value = []
      total.value = 0
    } finally {
      loading.value = false
    }
  }

  function onPageChange(page: number, size: number) {
    current.value = page
    pageSize.value = size
    fetch()
  }

  return { loading, dataSource, total, current, pageSize, fetch, onPageChange }
}

/** 通用表单提交组合式函数 */
export function useSubmit<T>(
  submitFn: (data: T) => Promise<void>,
  successMsg = '操作成功',
) {
  const submitting = ref(false)

  async function submit(data: T): Promise<boolean> {
    submitting.value = true
    try {
      await submitFn(data)
      message.success(successMsg)
      return true
    } catch {
      return false
    } finally {
      submitting.value = false
    }
  }

  return { submitting, submit }
}

/** 通用删除确认组合式函数 */
export function useDelete(
  deleteFn: (id: string) => Promise<void>,
  successMsg = '删除成功',
) {
  const deleting = ref(false)

  async function remove(id: string): Promise<boolean> {
    deleting.value = true
    try {
      await deleteFn(id)
      message.success(successMsg)
      return true
    } catch {
      return false
    } finally {
      deleting.value = false
    }
  }

  return { deleting, remove }
}

/** 窗口尺寸监听 */
export function useWindowSize() {
  const width = ref(window.innerWidth)
  const height = ref(window.innerHeight)

  function onResize() {
    width.value = window.innerWidth
    height.value = window.innerHeight
  }

  window.addEventListener('resize', onResize)

  const isMobile = computed(() => width.value < 768)

  return { width, height, isMobile }
}

import { computed } from 'vue'
