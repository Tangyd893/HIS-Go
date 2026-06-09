<template>
  <div>
    <a-card title="预约挂号" size="small">
      <a-form layout="vertical" size="small">
        <a-form-item label="选择科室" required>
          <a-select
            v-model:value="selectedDeptId"
            placeholder="请选择科室"
            :loading="deptLoading"
            :options="deptOptions"
            show-search
            option-filter-prop="label"
            style="width: 100%"
            @change="fetchSchedules"
          />
        </a-form-item>
        <a-form-item label="就诊日期" required>
          <a-date-picker
            v-model:value="searchDate"
            :disabled-date="disabledDate"
            style="width: 100%"
            @change="fetchSchedules"
          />
        </a-form-item>
      </a-form>

      <a-spin :spinning="scheduleLoading">
        <div class="schedule-list">
          <div
            v-for="item in schedules"
            :key="item.id"
            class="schedule-item"
            :class="{ disabled: !item.remainingSlots }"
            @click="showBookModal(item)"
          >
            <div class="schedule-info">
              <div class="doctor-name">{{ item.doctorName }}</div>
              <div class="dept-name">{{ item.deptName }}</div>
            </div>
            <div class="schedule-meta">
              <div>{{ item.timeSlotLabel }}</div>
              <div class="fee">挂号费 ¥{{ item.fee ?? 0 }}</div>
              <a-tag :color="item.remainingSlots > 0 ? 'green' : 'red'">
                剩余 {{ item.remainingSlots }} / {{ item.totalCount ?? '-' }} 号
              </a-tag>
            </div>
          </div>
          <a-empty v-if="!scheduleLoading && !schedules.length" :description="emptyHint" />
        </div>
      </a-spin>
    </a-card>

    <a-card v-if="myRegistrations.length" title="我的预约" size="small" style="margin-top: 12px">
      <div v-for="reg in myRegistrations" :key="reg.id" class="reg-item">
        <div class="reg-left">
          <div>{{ reg.registrationDate }} · 排队号 {{ reg.queueNumber }}</div>
          <a-tag>{{ regStatusText[reg.status] ?? '已预约' }}</a-tag>
        </div>
        <a-popconfirm
          v-if="reg.status === 0"
          title="确定取消该预约吗？"
          ok-text="确定"
          cancel-text="再想想"
          @confirm="handleCancel(reg.id)"
        >
          <a-button size="small" danger type="link">取消</a-button>
        </a-popconfirm>
      </div>
    </a-card>

    <a-modal v-model:open="modalOpen" title="确认挂号" ok-text="确认挂号" @ok="handleBook">
      <a-descriptions :column="1" size="small">
        <a-descriptions-item label="科室">{{ selectedSchedule?.deptName }}</a-descriptions-item>
        <a-descriptions-item label="医生">{{ selectedSchedule?.doctorName }}</a-descriptions-item>
        <a-descriptions-item label="日期">{{ selectedSchedule?.date }}</a-descriptions-item>
        <a-descriptions-item label="时段">{{ selectedSchedule?.timeSlotLabel }}</a-descriptions-item>
        <a-descriptions-item label="挂号费">¥{{ selectedSchedule?.fee ?? 0 }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>

    <a-modal v-model:open="successOpen" title="挂号成功" :footer="null">
      <a-result status="success" title="预约成功">
        <template #subTitle>
          <div>排队号：<strong>{{ bookResult?.queueNumber }}</strong></div>
          <div>{{ bookResult?.registrationDate }} · 请按时到院就诊</div>
        </template>
        <template #extra>
          <a-button type="primary" @click="successOpen = false">知道了</a-button>
        </template>
      </a-result>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { message } from 'ant-design-vue'
import { getDepartments, getSchedules, getRegistrations, register, cancelRegistration } from '@/api'
import { useAuthStore } from '@/store/auth'
import { flattenDepartments, normalizeSchedule, resolvePatientId } from '@/utils/patient'
import dayjs, { type Dayjs } from 'dayjs'

const route = useRoute()
const authStore = useAuthStore()

const deptLoading = ref(false)
const scheduleLoading = ref(false)
const departments = ref<{ id: string; name: string }[]>([])
const schedules = ref<any[]>([])
const myRegistrations = ref<any[]>([])
const selectedDeptId = ref<string>()
const searchDate = ref<Dayjs>(dayjs())
const modalOpen = ref(false)
const successOpen = ref(false)
const selectedSchedule = ref<any>(null)
const bookResult = ref<any>(null)

const regStatusText: Record<number, string> = {
  0: '已预约', 1: '已签到', 2: '已就诊', 3: '已取消',
}

const deptOptions = computed(() =>
  departments.value.map(d => ({ value: d.id, label: d.name })),
)

const emptyHint = computed(() =>
  selectedDeptId.value ? '该科室当日暂无号源，请换一天试试' : '请先选择科室和日期',
)

function disabledDate(current: Dayjs) {
  return current && current < dayjs().startOf('day')
}

async function fetchDepartments() {
  deptLoading.value = true
  try {
    const list = await getDepartments()
    departments.value = flattenDepartments(Array.isArray(list) ? list : [])
    if (!selectedDeptId.value && departments.value.length) {
      selectedDeptId.value = departments.value[0].id
    }
  } catch {
    departments.value = []
    message.error('加载科室失败')
  } finally {
    deptLoading.value = false
  }
}

async function fetchSchedules() {
  if (!selectedDeptId.value) {
    schedules.value = []
    return
  }
  scheduleLoading.value = true
  try {
    const list = await getSchedules({
      deptId: selectedDeptId.value,
      date: searchDate.value?.format('YYYY-MM-DD'),
    })
    schedules.value = (list || []).map(normalizeSchedule)
  } catch {
    schedules.value = []
    message.error('加载号源失败')
  } finally {
    scheduleLoading.value = false
  }
}

function showBookModal(item: any) {
  if (!item.remainingSlots) {
    message.warning('该时段号源已满')
    return
  }
  selectedSchedule.value = item
  modalOpen.value = true
}

async function handleBook() {
  const patientId = resolvePatientId(authStore.userInfo)
  const patientName = authStore.userInfo?.realName || authStore.username || '当前患者'
  try {
    const reg = await register({
      patientId,
      patientName,
      scheduleId: selectedSchedule.value.id,
    })
    bookResult.value = reg
    modalOpen.value = false
    successOpen.value = true
    myRegistrations.value.unshift(reg)
    fetchSchedules()
  } catch {
    message.error('挂号失败，请稍后重试')
  }
}

async function fetchMyRegistrations() {
  const patientId = resolvePatientId(authStore.userInfo)
  try {
    const res: any = await getRegistrations({ patientId, page: 1, pageSize: 20 })
    const list = res?.list || []
    // 只展示待就诊的预约（状态0=已预约, 1=已签到）
    myRegistrations.value = list.filter((r: any) => r.status === 0 || r.status === 1)
  } catch {
    // 静默失败，不影响主流程
  }
}

async function handleCancel(regId: string) {
  try {
    await cancelRegistration(regId)
    message.success('已取消预约')
    myRegistrations.value = myRegistrations.value.filter(r => r.id !== regId)
    fetchSchedules()
  } catch {
    message.error('取消失败，请稍后重试')
  }
}

onMounted(async () => {
  authStore.restoreUserInfo()
  if (route.query.deptId) {
    selectedDeptId.value = String(route.query.deptId)
  }
  await fetchDepartments()
  await fetchSchedules()
  await fetchMyRegistrations()
})
</script>

<style scoped>
.schedule-list { margin-top: 4px; }
.schedule-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  border: 1px solid #f0f0f0;
  border-radius: 8px;
  margin-bottom: 8px;
  cursor: pointer;
  background: #fff;
}
.schedule-item.disabled { opacity: 0.55; cursor: not-allowed; }
.schedule-item:not(.disabled):active { background: #f5f5f5; }
.doctor-name { font-size: 15px; font-weight: 500; color: #333; }
.dept-name { font-size: 13px; color: #999; margin-top: 2px; }
.schedule-meta { text-align: right; }
.schedule-meta div { font-size: 13px; color: #666; margin-bottom: 4px; }
.fee { color: #fa8c16; }
.reg-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px dashed #f0f0f0;
  font-size: 14px;
}
.reg-left {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
</style>
