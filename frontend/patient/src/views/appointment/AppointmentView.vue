<template>
  <div>
    <a-card title="预约挂号" size="small">
      <a-form layout="horizontal" size="small">
        <a-form-item label="日期">
          <a-date-picker v-model:value="searchDate" @change="fetchSchedules" style="width: 100%" />
        </a-form-item>
      </a-form>

      <div class="schedule-list">
        <div v-for="item in schedules" :key="item.id" class="schedule-item" @click="showBookModal(item)">
          <div class="schedule-info">
            <div class="doctor-name">{{ item.doctorName }}</div>
            <div class="dept-name">{{ item.deptName }}</div>
          </div>
          <div class="schedule-meta">
            <div>{{ item.timeSlot }}</div>
            <a-tag :color="item.remainingSlots > 0 ? 'green' : 'red'">
              剩余{{ item.remainingSlots }}号
            </a-tag>
          </div>
        </div>
        <a-empty v-if="!schedules.length" description="暂无号源，请选择日期" />
      </div>
    </a-card>

    <a-modal v-model:open="modalOpen" title="确认挂号" @ok="handleBook">
      <a-descriptions :column="1" size="small">
        <a-descriptions-item label="医生">{{ selectedSchedule?.doctorName }}</a-descriptions-item>
        <a-descriptions-item label="科室">{{ selectedSchedule?.deptName }}</a-descriptions-item>
        <a-descriptions-item label="日期">{{ selectedSchedule?.date }}</a-descriptions-item>
        <a-descriptions-item label="时段">{{ selectedSchedule?.timeSlot }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { getSchedules, register } from '@/api'
import dayjs from 'dayjs'

const schedules = ref<any[]>([])
const searchDate = ref(dayjs())
const modalOpen = ref(false)
const selectedSchedule = ref<any>(null)

async function fetchSchedules() {
  try {
    schedules.value = await getSchedules({
      date: searchDate.value?.format('YYYY-MM-DD') || undefined,
    })
  } catch { schedules.value = [] }
}

function showBookModal(item: any) {
  selectedSchedule.value = item
  modalOpen.value = true
}

async function handleBook() {
  try {
    await register({
      patientId: 'current-patient',
      patientName: '当前患者',
      scheduleId: selectedSchedule.value.id,
    })
    message.success('挂号成功')
    modalOpen.value = false
    fetchSchedules()
  } catch { }
}

onMounted(fetchSchedules)
</script>

<style scoped>
.schedule-list { margin-top: 12px; }
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
.schedule-item:active { background: #f5f5f5; }
.doctor-name { font-size: 15px; font-weight: 500; }
.dept-name { font-size: 13px; color: #999; margin-top: 2px; }
.schedule-meta { text-align: right; }
.schedule-meta div { font-size: 13px; color: #666; margin-bottom: 4px; }
</style>
