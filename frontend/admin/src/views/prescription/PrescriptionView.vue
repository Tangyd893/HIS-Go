<template>
  <div>
    <a-card title="处方管理">
      <template #extra>
        <a-button type="primary" @click="showCreateModal"><PlusOutlined /> 开具处方</a-button>
      </template>
      <a-table :columns="columns" :data-source="dataSource" :loading="loading" :pagination="pagination" row-key="id" @change="onTableChange">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="statusColor[record.status]">{{ statusText[record.status] || record.status }}</a-tag>
          </template>
          <template v-if="column.key === 'action'">
            <a-space>
              <a @click="viewDetail(record.id)">详情</a>
              <a-button size="small" type="primary" @click="reviewPrescription(record)" v-if="record.status === 0">审核</a-button>
              <a-button size="small" danger @click="cancelPrescription(record.id)" v-if="record.status === 0">取消</a-button>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:open="modalOpen" title="开具处方" @ok="handleCreate" width="700px">
      <a-form layout="vertical">
        <a-form-item label="患者ID"><a-input v-model:value="form.patientId" /></a-form-item>
        <a-form-item label="诊断"><a-textarea v-model:value="form.diagnosis" :rows="2" /></a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { prescriptionApi } from '@/api/prescription'

const loading = ref(false)
const dataSource = ref<any[]>([])
const modalOpen = ref(false)
const form = reactive({ patientId: '', diagnosis: '' })
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })

const statusText: Record<number, string> = { 0: '待审核', 1: '已审核', 2: '已发药', 3: '已取消' }
const statusColor: Record<number, string> = { 0: 'orange', 1: 'blue', 2: 'green', 3: 'red' }

const columns = [
  { title: '患者', dataIndex: 'patientName' },
  { title: '医生', dataIndex: 'doctorName' },
  { title: '创建时间', dataIndex: 'createdAt' },
  { title: '状态', key: 'status' },
  { title: '操作', key: 'action', width: 200 },
]

async function fetchData() {
  loading.value = true
  try {
    const res: any = await prescriptionApi.getList({ page: pagination.current, pageSize: pagination.pageSize })
    dataSource.value = res?.list || []
    pagination.total = res?.total || 0
  } catch { dataSource.value = [] } finally { loading.value = false }
}

function onTableChange(pag: any) { pagination.current = pag.current; fetchData() }
function showCreateModal() { modalOpen.value = true }
function viewDetail(id: string) { message.info(`查看处方: ${id}`) }

async function handleCreate() {
  try { await prescriptionApi.create({ ...form, details: [] }); message.success('创建成功'); modalOpen.value = false; fetchData() } catch { }
}

async function reviewPrescription(record: any) {
  try { await prescriptionApi.review({ id: record.id, approved: true }); message.success('审核通过'); fetchData() } catch { }
}

async function cancelPrescription(id: string) {
  try { await prescriptionApi.cancel(id); message.success('已取消'); fetchData() } catch { }
}

onMounted(fetchData)
</script>
