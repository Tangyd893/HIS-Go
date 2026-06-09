<template>
  <ServicePlaceholder v-if="serviceUnavailable" title="检查检验" />
  <a-card v-else title="检查检验">
    <a-table :columns="columns" :data-source="dataSource" :loading="loading" :pagination="pagination" row-key="id" @change="onTableChange">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <a-tag :color="statusColor[record.status]">{{ statusText[record.status] || record.status }}</a-tag>
        </template>
        <template v-if="column.key === 'action'">
          <a-space>
            <a @click="viewDetail(record)">详情</a>
            <a-button size="small" type="primary" @click="reviewReport(record)" v-if="record.status === 1">审核</a-button>
          </a-space>
        </template>
      </template>
    </a-table>
  </a-card>

  <a-modal v-model:open="detailOpen" title="报告详情" :footer="null" width="600px">
    <a-descriptions v-if="detailRecord" :column="1" bordered size="small">
      <a-descriptions-item v-for="(v, k) in detailRecord" :key="k" :label="k">{{ v }}</a-descriptions-item>
    </a-descriptions>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { examinationApi } from '@/api/examination'
import ServicePlaceholder from '@/components/ServicePlaceholder.vue'

const serviceUnavailable = ref(false)
const loading = ref(false)
const dataSource = ref<any[]>([])
const detailOpen = ref(false)
const detailRecord = ref<any>(null)
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })

const statusText: Record<number, string> = { 0: '待检查', 1: '已检查', 2: '已审核', 3: '已发布' }
const statusColor: Record<number, string> = { 0: 'orange', 1: 'blue', 2: 'green', 3: 'default' }

const columns = [
  { title: '患者', dataIndex: 'patientName' },
  { title: '检查类型', dataIndex: 'examType' },
  { title: '检查项目', dataIndex: 'examItem' },
  { title: '状态', key: 'status' },
  { title: '创建时间', dataIndex: 'createdAt' },
  { title: '操作', key: 'action', width: 160 },
]

async function fetchData() {
  loading.value = true
  try {
    const res: any = await examinationApi.getList({ page: pagination.current, pageSize: pagination.pageSize })
    dataSource.value = res?.list || []
    pagination.total = res?.total || 0
  } catch { serviceUnavailable.value = true; dataSource.value = [] } finally { loading.value = false }
}

function onTableChange(pag: any) { pagination.current = pag.current; fetchData() }
function viewDetail(record: any) { detailRecord.value = record; detailOpen.value = true }

async function reviewReport(record: any) {
  try { await examinationApi.review({ report_id: record.id, reviewer_id: 'current', approved: true }); message.success('审核通过'); fetchData() } catch { }
}

onMounted(fetchData)
</script>
