<template>
  <a-card title="员工管理">
    <template #extra>
      <a-space>
        <a-input-search v-model:value="searchName" placeholder="搜索员工姓名" style="width: 200px" @search="fetchData" />
        <a-button type="primary"><PlusOutlined /> 新增员工</a-button>
      </a-space>
    </template>
    <a-table :columns="columns" :data-source="dataSource" :loading="loading" :pagination="pagination" :scroll="{ x: 800 }" row-key="id" @change="onTableChange">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          {{ record.status === 1 ? '在职' : '停用' }}
        </template>
      </template>
    </a-table>
  </a-card>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { userApi } from '@/api/user'

const loading = ref(false)
const searchName = ref('')
const dataSource = ref<any[]>([])

const pagination = reactive({ current: 1, pageSize: 10, total: 0 })

const columns = [
  { title: '姓名', dataIndex: 'name', key: 'name', width: 100 },
  { title: '手机号', dataIndex: 'phone', key: 'phone', width: 130 },
  { title: '科室', dataIndex: 'deptName', key: 'deptName', width: 120 },
  { title: '职称', dataIndex: 'title', key: 'title', width: 120 },
  { title: '状态', key: 'status', width: 80 },
]

async function fetchData() {
  loading.value = true
  try {
    const res: any = await userApi.getEmployees({
      page: pagination.current,
      pageSize: pagination.pageSize,
      name: searchName.value || undefined,
    })
    dataSource.value = res?.list || res || []
    pagination.total = res?.total || 0
  } catch { message.error('加载员工失败'); dataSource.value = [] } finally { loading.value = false }
}

function onTableChange(pag: any) {
  pagination.current = pag.current
  fetchData()
}

onMounted(fetchData)
</script>
