<template>
  <a-card title="在线问诊">
    <a-table :columns="columns" :data-source="dataSource" :loading="loading" :pagination="pagination" row-key="id" @change="onTableChange">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <a @click="viewDetail(record.id)">查看详情</a>
        </template>
      </template>
    </a-table>
  </a-card>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { outpatientApi } from '@/api/others'

const loading = ref(false)
const dataSource = ref<any[]>([])
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })

const columns = [
  { title: '患者', dataIndex: 'patientId' },
  { title: '医生', dataIndex: 'doctorId' },
  { title: '状态', dataIndex: 'status' },
  { title: '创建时间', dataIndex: 'createdAt' },
  { title: '操作', key: 'action', width: 100 },
]

async function fetchData() {
  loading.value = true
  try {
    const res: any = await outpatientApi.getConsultations({ page: pagination.current, pageSize: pagination.pageSize })
    dataSource.value = res?.list || []
    pagination.total = res?.total || 0
  } catch { dataSource.value = [] } finally { loading.value = false }
}

function onTableChange(pag: any) { pagination.current = pag.current; fetchData() }
function viewDetail(id: string) { message.info(`查看问诊: ${id}`) }

onMounted(fetchData)
</script>
