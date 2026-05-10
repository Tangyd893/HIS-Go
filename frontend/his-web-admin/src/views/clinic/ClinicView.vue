<template>
  <a-card title="门诊诊疗">
    <template #extra>
      <a-button type="primary" @click="modalOpen = true"><PlusOutlined /> 新建诊疗记录</a-button>
    </template>
    <a-table :columns="columns" :data-source="dataSource" :loading="loading" :pagination="pagination" row-key="id" @change="onTableChange">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <a @click="viewDetail(record.id)">查看详情</a>
        </template>
      </template>
    </a-table>

    <a-modal v-model:open="modalOpen" title="新建诊疗记录" @ok="handleCreate" width="700px">
      <a-form layout="vertical">
        <a-form-item label="患者ID"><a-input v-model:value="form.patientId" /></a-form-item>
        <a-form-item label="主诉"><a-textarea v-model:value="form.chiefComplaint" :rows="3" /></a-form-item>
        <a-form-item label="诊断"><a-textarea v-model:value="form.diagnosis" :rows="3" /></a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { clinicApi } from '@/api/clinic'

const loading = ref(false)
const dataSource = ref<any[]>([])
const modalOpen = ref(false)
const form = reactive({ patientId: '', chiefComplaint: '', diagnosis: '' })
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })

const columns = [
  { title: '患者', dataIndex: 'patientName' },
  { title: '医生', dataIndex: 'doctorName' },
  { title: '主诉', dataIndex: 'chiefComplaint' },
  { title: '诊断', dataIndex: 'diagnosis' },
  { title: '就诊时间', dataIndex: 'createdAt' },
  { title: '操作', key: 'action', width: 100 },
]

async function fetchData() {
  loading.value = true
  try {
    const res: any = await clinicApi.getList({ page: pagination.current, pageSize: pagination.pageSize })
    dataSource.value = res?.list || []
    pagination.total = res?.total || 0
  } catch { dataSource.value = [] } finally { loading.value = false }
}

function onTableChange(pag: any) { pagination.current = pag.current; fetchData() }
function viewDetail(id: string) { message.info(`查看诊疗记录: ${id}`) }

async function handleCreate() {
  try { await clinicApi.create(form); message.success('创建成功'); modalOpen.value = false; fetchData() } catch { }
}

onMounted(fetchData)
</script>
