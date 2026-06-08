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
              <a @click="viewDetail(record)">详情</a>
              <a-button size="small" type="primary" @click="reviewPrescription(record)" v-if="record.status === 0">审核</a-button>
              <a-button size="small" danger @click="cancelPrescription(record.id)" v-if="record.status === 0">取消</a-button>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:open="modalOpen" title="开具处方" @ok="handleCreate" width="750px">
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
        <a-form-item label="诊断"><a-textarea v-model:value="form.diagnosis" :rows="2" /></a-form-item>
        <a-form-item label="药品明细">
          <div v-for="(item, idx) in form.details" :key="idx" style="display: flex; gap: 8px; margin-bottom: 8px; align-items: center">
            <a-select
              v-model:value="item.drugId"
              show-search
              placeholder="搜索药品"
              :filter-option="false"
              :options="drugOptions"
              @search="searchDrugs"
              @select="(val: string) => onDrugSelect(val, idx)"
              style="flex: 2"
            />
            <a-input-number v-model:value="item.quantity" :min="1" placeholder="数量" style="width: 70px" />
            <a-input v-model:value="item.dosage" placeholder="用量" style="width: 80px" />
            <a-input v-model:value="item.usage" placeholder="用法" style="width: 100px" />
            <a-button size="small" danger @click="removeDetail(idx)">×</a-button>
          </div>
          <a-button type="dashed" size="small" @click="addDetail" style="margin-top: 4px">+ 添加药品</a-button>
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal v-model:open="detailOpen" title="处方详情" :footer="null" width="650px">
      <a-descriptions v-if="detailRecord" :column="2" bordered size="small">
        <template v-for="(v, k) in detailRecord" :key="k">
          <a-descriptions-item v-if="String(k) !== 'details'" :label="k">{{ v }}</a-descriptions-item>
        </template>
      </a-descriptions>
      <a-card v-if="detailRecord?.details?.length" title="药品明细" size="small" style="margin-top: 12px">
        <a-table :columns="[{title:'药品',dataIndex:'drugName'},{title:'数量',dataIndex:'quantity'},{title:'用量',dataIndex:'dosage'},{title:'用法',dataIndex:'usage'}]" :data-source="detailRecord.details" row-key="id" size="small" :pagination="false" />
      </a-card>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { prescriptionApi } from '@/api/prescription'
import { userApi } from '@/api/user'
import { pharmacyApi } from '@/api/pharmacy'

const loading = ref(false)
const dataSource = ref<any[]>([])
const modalOpen = ref(false)
const detailOpen = ref(false)
const detailRecord = ref<any>(null)
const form = reactive({ patientId: '', patientName: '', diagnosis: '', details: [] as any[] })
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const patients = ref<any[]>([])
const patientOptions = ref<{ label: string; value: string }[]>([])
const drugs = ref<any[]>([])
const drugOptions = ref<{ label: string; value: string }[]>([])

const statusText: Record<number, string> = { 0: '待审核', 1: '已审核', 2: '已发药', 3: '已取消' }
const statusColor: Record<number, string> = { 0: 'orange', 1: 'blue', 2: 'green', 3: 'red' }

const columns = [
  { title: '患者', dataIndex: 'patientName' },
  { title: '医生', dataIndex: 'doctorName' },
  { title: '创建时间', dataIndex: 'createdAt' },
  { title: '状态', key: 'status' },
  { title: '操作', key: 'action', width: 200 },
]

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

async function searchDrugs(keyword: string) {
  if (!keyword || keyword.length < 1) { drugOptions.value = []; return }
  try {
    const res = await pharmacyApi.getDrugs({ name: keyword, page: 1, pageSize: 20 })
    drugs.value = (res as any)?.list || []
    drugOptions.value = drugs.value.map((d: any) => ({
      label: `${d.name} (${d.specification || ''}) - ¥${d.price}`,
      value: d.id,
    }))
  } catch { drugOptions.value = [] }
}

function onDrugSelect(value: string, idx: number) {
  const d = drugs.value.find((x: any) => x.id === value)
  if (d) {
    form.details[idx].drugName = d.name
    form.details[idx].price = d.price
  }
}

function addDetail() {
  form.details.push({ drugId: '', drugName: '', quantity: 1, dosage: '', usage: '' })
}

function removeDetail(idx: number) {
  form.details.splice(idx, 1)
}

async function fetchData() {
  loading.value = true
  try {
    const res: any = await prescriptionApi.getList({ page: pagination.current, pageSize: pagination.pageSize })
    dataSource.value = res?.list || []
    pagination.total = res?.total || 0
  } catch { dataSource.value = [] } finally { loading.value = false }
}

function onTableChange(pag: any) { pagination.current = pag.current; fetchData() }

function showCreateModal() {
  form.patientId = ''
  form.patientName = ''
  form.diagnosis = ''
  form.details = []
  patientOptions.value = []
  patients.value = []
  drugOptions.value = []
  drugs.value = []
  modalOpen.value = true
}

function viewDetail(record: any) { detailRecord.value = record; detailOpen.value = true }

async function handleCreate() {
  try {
    await prescriptionApi.create({
      prescription: {
        patientId: form.patientId,
        patientName: form.patientName,
      },
      details: form.details,
    })
    message.success('创建成功'); modalOpen.value = false; fetchData()
  } catch { }
}

async function reviewPrescription(record: any) {
  try { await prescriptionApi.review({ id: record.id, approved: true }); message.success('审核通过'); fetchData() } catch { }
}

async function cancelPrescription(id: string) {
  try { await prescriptionApi.cancel(id); message.success('已取消'); fetchData() } catch { }
}

onMounted(fetchData)
</script>
