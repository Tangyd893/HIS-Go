<template>
  <div>
    <a-card title="号源管理">
      <a-form layout="inline" style="margin-bottom: 16px">
        <a-form-item label="日期">
          <a-date-picker v-model:value="searchDate" @change="fetchSchedules" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="fetchSchedules">查询</a-button>
        </a-form-item>
      </a-form>
      <a-table :columns="scheduleColumns" :data-source="schedules" :loading="loading" row-key="id">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'action'">
            <a-button type="primary" size="small" @click="showRegisterModal(record)">挂号</a-button>
          </template>
          <template v-if="column.key === 'timeSlot'">
            {{ ({ 1: '上午', 2: '下午' } as Record<number, string>)[record.timeSlot] || record.timeSlot }}
          </template>
        </template>
      </a-table>
    </a-card>

    <a-card title="挂号记录" style="margin-top: 16px">
      <a-table :columns="regColumns" :data-source="registrations" :loading="regLoading" row-key="id">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="record.status === 0 ? 'blue' : record.status === 1 ? 'green' : 'red'">
              {{ ({ 0: '待签到', 1: '已签到', 3: '已取消' } as Record<number, string>)[record.status] || record.status }}
            </a-tag>
          </template>
          <template v-if="column.key === 'action'">
            <a-space>
              <a-button size="small" @click="signin(record.id)" v-if="record.status === 0">签到</a-button>
              <a-button size="small" type="primary" @click="startClinic(record)" v-if="record.status === 1">开始接诊</a-button>
              <a-button size="small" danger @click="cancelReg(record.id)" v-if="record.status === 0">取消</a-button>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:open="registerModalOpen" title="患者挂号" @ok="doRegister">
      <a-form layout="vertical">
        <a-form-item label="选择患者" required>
          <a-select
            v-model:value="registerForm.patientId"
            show-search
            placeholder="搜索患者姓名或手机号"
            :filter-option="false"
            :options="patientOptions"
            @search="fetchPatients"
            @select="onPatientSelect"
            style="width: 100%"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { registrationApi } from '@/api/registration'
import { userApi } from '@/api/user'
import dayjs from 'dayjs'

const router = useRouter()
const loading = ref(false)
const regLoading = ref(false)
const schedules = ref<any[]>([])
const registrations = ref<any[]>([])
const searchDate = ref<dayjs.Dayjs | null>(dayjs())
const registerModalOpen = ref(false)
const selectedSchedule = ref<any>(null)

const patients = ref<any[]>([])
const patientOptions = ref<{ label: string; value: string }[]>([])

const registerForm = reactive({ patientId: '', patientName: '' })

async function fetchPatients(keyword: string) {
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

function onPatientSelect(value: string, option: any) {
  const p = patients.value.find((x: any) => x.id === value)
  registerForm.patientName = p?.name || ''
}

const scheduleColumns = [
  { title: '医生', dataIndex: 'doctorName' },
  { title: '科室', dataIndex: 'deptName' },
  { title: '日期', dataIndex: 'date' },
  { title: '时段', dataIndex: 'timeSlot', key: 'timeSlot' },
  { title: '剩余号源', dataIndex: 'remainCount' },
  { title: '操作', key: 'action', width: 100 },
]

const regColumns = [
  { title: '患者', dataIndex: 'patientName' },
  { title: '科室', dataIndex: 'deptName' },
  { title: '医生', dataIndex: 'doctorName' },
  { title: '挂号日期', dataIndex: 'registrationDate' },
  { title: '状态', key: 'status' },
  { title: '操作', key: 'action', width: 160 },
]

async function fetchSchedules() {
  loading.value = true
  try {
    schedules.value = await registrationApi.getSchedules({
      date: searchDate.value?.format('YYYY-MM-DD') || undefined,
    })
  } catch { schedules.value = [] } finally { loading.value = false }
}

async function fetchRegistrations() {
  regLoading.value = true
  try {
    const res = await registrationApi.getRegistrations({
      page: 1,
      pageSize: 50,
      date: searchDate.value?.format('YYYY-MM-DD') || undefined,
    })
    registrations.value = res?.list || []
  } catch { registrations.value = [] } finally { regLoading.value = false }
}

function showRegisterModal(record: any) {
  selectedSchedule.value = record
  registerForm.patientId = ''
  registerForm.patientName = ''
  patientOptions.value = []
  patients.value = []
  registerModalOpen.value = true
}

async function doRegister() {
  try {
    await registrationApi.register({
      patientId: registerForm.patientId,
      patientName: registerForm.patientName,
      scheduleId: selectedSchedule.value.id,
    })
    message.success('挂号成功')
    registerModalOpen.value = false
    fetchSchedules()
    fetchRegistrations()
  } catch { }
}

async function signin(id: string) {
  try {
    await registrationApi.signin(id)
    message.success('签到成功')
    fetchRegistrations()
    fetchSchedules()
  } catch { }
}

async function cancelReg(id: string) {
  try {
    await registrationApi.cancel(id)
    message.success('已取消')
    fetchRegistrations()
    fetchSchedules()
  } catch { }
}

function startClinic(record: any) {
  router.push({ path: '/clinic', query: { patientId: record.patientId, registrationId: record.id } })
}

onMounted(() => { fetchSchedules(); fetchRegistrations() })
</script>
