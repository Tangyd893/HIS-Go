<template>
  <a-card title="门诊诊疗">
    <template #extra>
      <a-button type="primary" @click="openCreateModal"><PlusOutlined /> 新建诊疗记录</a-button>
    </template>
    <a-table :columns="columns" :data-source="dataSource" :loading="loading" :pagination="pagination" row-key="id" @change="onTableChange">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <a @click="viewDetail(record)">查看详情</a>
        </template>
      </template>
    </a-table>

    <a-modal v-model:open="modalOpen" title="新建诊疗记录" @ok="handleCreate" width="700px">
      <a-form layout="vertical">
        <a-form-item label="选择患者" required>
          <a-select
            v-model:value="form.patientId"
            show-search
            placeholder="搜索患者姓名或手机号"
            :filter-option="false"
            :options="patientOptions"
            @search="searchPatients"
            @select="onPatientSelect"
            style="width: 100%"
          />
        </a-form-item>
        <a-form-item label="主诉"><a-textarea v-model:value="form.chiefComplaint" :rows="3" /></a-form-item>
        <a-form-item label="诊断"><a-textarea v-model:value="form.diagnosis" :rows="3" /></a-form-item>
      </a-form>
    </a-modal>

    <a-modal v-model:open="detailOpen" title="诊疗详情" :footer="null" width="600px">
      <a-descriptions v-if="detailRecord" :column="2" bordered size="small">
        <a-descriptions-item v-for="(v, k) in detailRecord" :key="k" :label="k">{{ v }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </a-card>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { clinicApi } from '@/api/clinic'
import { userApi } from '@/api/user'

const route = useRoute()
const loading = ref(false)
const dataSource = ref<any[]>([])
const modalOpen = ref(false)
const detailOpen = ref(false)
const detailRecord = ref<any>(null)
const form = reactive({ patientId: '', patientName: '', chiefComplaint: '', diagnosis: '' })
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const patients = ref<any[]>([])
const patientOptions = ref<{ label: string; value: string }[]>([])

async function searchPatients(keyword: string) {
  if (!keyword || keyword.length < 1) { patientOptions.value = []; return }
  try {
    const res = await userApi.getPatients({ name: keyword, page: 1, pageSize: 20 })
    patients.value = res?.list || []
    patientOptions.value = patients.value.map((p: any) => ({
      label: `${p.name} (${p.phone || p.idCard || ''})`,
      value: p.id,
    }))
  } catch { patientOptions.value = [] }
}

function onPatientSelect(value: string) {
  const p = patients.value.find((x: any) => x.id === value)
  form.patientName = p?.name || ''
}

function openCreateModal() {
  form.patientId = ''
  form.patientName = ''
  form.chiefComplaint = ''
  form.diagnosis = ''
  patientOptions.value = []
  patients.value = []
  modalOpen.value = true
}

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
function viewDetail(record: any) { detailRecord.value = record; detailOpen.value = true }

async function handleCreate() {
  try {
    await clinicApi.create({
      patientId: form.patientId,
      patientName: form.patientName,
      chiefComplaint: form.chiefComplaint,
      diagnosis: form.diagnosis,
    })
    message.success('创建成功'); modalOpen.value = false; fetchData()
  } catch { }
}

onMounted(() => {
  fetchData()
  // 从挂号记录跳转时自动填入患者信息
  const qPatientId = route.query.patientId as string
  const qRegId = route.query.registrationId as string
  if (qPatientId) {
    form.patientId = qPatientId
    modalOpen.value = true
  }
})
</script>
