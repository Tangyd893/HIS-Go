<template>
  <div>
    <ServicePlaceholder
      v-if="serviceUnavailable"
      title="住院管理"
      description="住院管理模块在演示环境中暂未启用，如需演示请联系管理员。"
    />
    <a-card v-else title="住院管理">
      <template #extra>
        <a-button type="primary" @click="showAdmitModal"><PlusOutlined /> 入院登记</a-button>
      </template>
      <a-table :columns="columns" :data-source="dataSource" :loading="loading" :pagination="pagination" row-key="id" @change="onTableChange">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="statusColor[record.status]">{{ statusText[record.status] || record.status }}</a-tag>
          </template>
          <template v-if="column.key === 'action'">
            <a-space>
              <a @click="viewDetail(record.id)">详情</a>
              <a-button size="small" danger @click="handleDischarge(record.id)" v-if="record.status === 1">出院</a-button>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:open="modalOpen" title="入院登记" @ok="handleAdmit" width="600px">
      <a-form layout="vertical">
        <a-form-item label="患者ID"><a-input v-model:value="form.patientId" /></a-form-item>
        <a-form-item label="患者姓名"><a-input v-model:value="form.patientName" /></a-form-item>
        <a-form-item label="科室ID"><a-input v-model:value="form.deptId" /></a-form-item>
        <a-form-item label="床位号"><a-input v-model:value="form.bedNo" /></a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { inpatientApi } from '@/api/inpatient'
import ServicePlaceholder from '@/components/ServicePlaceholder.vue'

const loading = ref(false)
const serviceUnavailable = ref(false)
const dataSource = ref<any[]>([])
const modalOpen = ref(false)
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const form = reactive({ patientId: '', patientName: '', deptId: '', bedNo: '' })

const statusText: Record<number, string> = { 0: '待入院', 1: '在院', 2: '已出院' }
const statusColor: Record<number, string> = { 0: 'blue', 1: 'green', 2: 'default' }

const columns = [
  { title: '患者', dataIndex: 'patientName' },
  { title: '科室', dataIndex: 'deptName' },
  { title: '床位号', dataIndex: 'bedNo' },
  { title: '入院日期', dataIndex: 'admitDate' },
  { title: '出院日期', dataIndex: 'dischargeDate' },
  { title: '状态', key: 'status' },
  { title: '操作', key: 'action', width: 160 },
]

async function fetchData() {
  loading.value = true
  serviceUnavailable.value = false
  try {
    const res: any = await inpatientApi.getList({ page: pagination.current, pageSize: pagination.pageSize })
    dataSource.value = res?.list || []
    pagination.total = res?.total || 0
  } catch {
    serviceUnavailable.value = true
    dataSource.value = []
  } finally { loading.value = false }
}

function onTableChange(pag: any) { pagination.current = pag.current; fetchData() }
function viewDetail(id: string) { message.info(`查看住院记录: ${id}`) }
function showAdmitModal() { modalOpen.value = true }

async function handleAdmit() {
  try { await inpatientApi.admit(form); message.success('入院登记成功'); modalOpen.value = false; fetchData() } catch { }
}

async function handleDischarge(id: string) {
  try { await inpatientApi.discharge(id); message.success('出院办理成功'); fetchData() } catch { }
}

onMounted(fetchData)
</script>
