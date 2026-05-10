<template>
  <div>
    <a-card title="排班管理">
      <template #extra>
        <a-space>
          <a-date-picker v-model:value="filterDate" placeholder="选择日期" />
          <a-button type="primary" @click="showGenerateModal"><PlusOutlined /> 生成排班</a-button>
        </a-space>
      </template>
      <a-table :columns="columns" :data-source="dataSource" :loading="loading" row-key="id">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'action'">
            <a-button size="small" danger @click="cancelSchedule(record.id)">取消</a-button>
          </template>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:open="modalOpen" title="生成排班" @ok="handleGenerate">
      <a-form layout="vertical">
        <a-form-item label="开始日期" required>
          <a-date-picker v-model:value="genForm.startDate" style="width: 100%" />
        </a-form-item>
        <a-form-item label="结束日期" required>
          <a-date-picker v-model:value="genForm.endDate" style="width: 100%" />
        </a-form-item>
        <a-form-item label="科室">
          <a-input v-model:value="genForm.deptId" placeholder="科室ID" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { scheduleApi } from '@/api/schedule'
import dayjs from 'dayjs'

const loading = ref(false)
const dataSource = ref<any[]>([])
const filterDate = ref<dayjs.Dayjs | null>(dayjs())
const modalOpen = ref(false)
const genForm = reactive({ startDate: null as dayjs.Dayjs | null, endDate: null as dayjs.Dayjs | null, deptId: '' })

const columns = [
  { title: '医生', dataIndex: 'doctorName' },
  { title: '科室', dataIndex: 'deptName' },
  { title: '日期', dataIndex: 'date' },
  { title: '时段', dataIndex: 'timeSlot' },
  { title: '总号源', dataIndex: 'totalSlots' },
  { title: '剩余', dataIndex: 'remainingSlots' },
  { title: '操作', key: 'action', width: 100 },
]

async function fetchData() {
  loading.value = true
  try {
    dataSource.value = await scheduleApi.getList({ date: filterDate.value?.format('YYYY-MM-DD') })
  } catch { dataSource.value = [] } finally { loading.value = false }
}

function showGenerateModal() { modalOpen.value = true }

async function handleGenerate() {
  try {
    await scheduleApi.generate({
      startDate: genForm.startDate?.format('YYYY-MM-DD') || '',
      endDate: genForm.endDate?.format('YYYY-MM-DD') || '',
      deptId: genForm.deptId,
    })
    message.success('排班生成成功')
    modalOpen.value = false
    fetchData()
  } catch { }
}

async function cancelSchedule(id: string) {
  try { await scheduleApi.cancel(id); message.success('已取消'); fetchData() } catch { }
}

onMounted(fetchData)
</script>
