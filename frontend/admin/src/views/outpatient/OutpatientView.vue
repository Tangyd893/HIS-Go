<template>
  <ServicePlaceholder v-if="serviceUnavailable" title="院外服务" />
  <a-card v-else title="院外患者管理">
    <a-table :columns="columns" :data-source="dataSource" :loading="loading" :pagination="pagination" row-key="id" @change="onTableChange">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <a @click="viewDetail(record)">查看详情</a>
        </template>
      </template>
    </a-table>
  </a-card>

  <a-modal v-model:open="detailOpen" title="问诊详情" :footer="null" width="600px">
    <a-descriptions v-if="detailRecord" :column="1" bordered size="small">
      <a-descriptions-item v-for="(v, k) in detailRecord" :key="k" :label="k">{{ v }}</a-descriptions-item>
    </a-descriptions>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { outpatientApi } from '@/api/others'
import ServicePlaceholder from '@/components/ServicePlaceholder.vue'

const serviceUnavailable = ref(false)
const loading = ref(false)
const dataSource = ref<any[]>([])
const detailOpen = ref(false)
const detailRecord = ref<any>(null)
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
  } catch { serviceUnavailable.value = true; dataSource.value = [] } finally { loading.value = false }
}

function onTableChange(pag: any) { pagination.current = pag.current; fetchData() }
function viewDetail(record: any) { detailRecord.value = record; detailOpen.value = true }

onMounted(fetchData)
</script>
