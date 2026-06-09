<template>
  <div>
    <ServicePlaceholder
      v-if="serviceUnavailable"
      title="电子病历"
      description="电子病历模块在演示环境中暂未启用，如需演示请联系管理员。"
    />
    <a-card v-else title="电子病历">
      <template #extra>
        <a-button type="primary" @click="showCreateModal"><PlusOutlined /> 新建病历</a-button>
      </template>
      <a-table :columns="columns" :data-source="dataSource" :loading="loading" :pagination="pagination" row-key="id" @change="onTableChange">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="statusColor[record.status]">{{ statusText[record.status] || record.status }}</a-tag>
          </template>
          <template v-if="column.key === 'action'">
            <a-space>
              <a @click="viewDetail(record)">详情</a>
              <a-button size="small" type="primary" @click="doQualityControl(record)" v-if="record.status === 0">质控</a-button>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:open="detailOpen" title="病历详情" :footer="null" width="600px">
      <a-descriptions v-if="detailRecord" :column="1" bordered size="small">
        <a-descriptions-item v-for="(v, k) in detailRecord" :key="k" :label="k">{{ v }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>

    <a-modal v-model:open="modalOpen" title="新建病历" @ok="handleCreate" width="800px">
      <a-form layout="vertical">
        <a-form-item label="患者ID"><a-input v-model:value="form.patientId" /></a-form-item>
        <a-form-item label="主观资料 (S)"><a-textarea v-model:value="form.subjective" :rows="3" placeholder="主诉、病史等" /></a-form-item>
        <a-form-item label="客观资料 (O)"><a-textarea v-model:value="form.objective" :rows="3" placeholder="体格检查、辅助检查等" /></a-form-item>
        <a-form-item label="评估 (A)"><a-textarea v-model:value="form.assessment" :rows="3" placeholder="诊断、评估" /></a-form-item>
        <a-form-item label="计划 (P)"><a-textarea v-model:value="form.plan" :rows="3" placeholder="治疗方案、随访计划" /></a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { emrApi } from '@/api/emr'
import ServicePlaceholder from '@/components/ServicePlaceholder.vue'

const loading = ref(false)
const serviceUnavailable = ref(false)
const dataSource = ref<any[]>([])
const modalOpen = ref(false)
const detailOpen = ref(false)
const detailRecord = ref<any>(null)
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const form = reactive({ patientId: '', subjective: '', objective: '', assessment: '', plan: '' })

const statusText: Record<number, string> = { 0: '草稿', 1: '已审核', 2: '已质控' }
const statusColor: Record<number, string> = { 0: 'orange', 1: 'blue', 2: 'green' }

const columns = [
  { title: '患者', dataIndex: 'patientName' },
  { title: '医生', dataIndex: 'doctorId' },
  { title: '模板', dataIndex: 'templateId' },
  { title: '状态', key: 'status' },
  { title: '创建时间', dataIndex: 'createdAt' },
  { title: '操作', key: 'action', width: 160 },
]

async function fetchData() {
  loading.value = true
  serviceUnavailable.value = false
  try {
    const res: any = await emrApi.getList({ page: pagination.current, pageSize: pagination.pageSize })
    dataSource.value = res?.list || []
    pagination.total = res?.total || 0
  } catch {
    serviceUnavailable.value = true
    dataSource.value = []
  } finally { loading.value = false }
}

function onTableChange(pag: any) { pagination.current = pag.current; fetchData() }
function viewDetail(record: any) { detailRecord.value = record; detailOpen.value = true }
function showCreateModal() { modalOpen.value = true }

async function handleCreate() {
  try {
    await emrApi.create({
      patientId: form.patientId,
      chiefComplaint: form.subjective,
      presentIllness: form.subjective,
      physicalExam: form.objective,
      diagnosis: form.assessment,
      treatmentPlan: form.plan,
    })
    message.success('创建成功')
    modalOpen.value = false
    fetchData()
  } catch { message.error('创建失败') }
}

async function doQualityControl(record: any) {
  try { await emrApi.qualityControl(record.id, { reviewer_id: 'current', level: 1 }); message.success('质控完成'); fetchData() } catch { message.error('质控失败') }
}

onMounted(fetchData)
</script>
