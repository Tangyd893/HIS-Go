<template>
  <div>
    <a-card title="慢病管理" size="small">
      <template #extra>
        <a-button type="primary" size="small" @click="showReportModal">上报数据</a-button>
      </template>

      <a-descriptions v-if="contract" :column="2" bordered size="small" title="签约信息">
        <a-descriptions-item label="签约状态">{{ contract.status || '已签约' }}</a-descriptions-item>
        <a-descriptions-item label="签约日期">{{ contract.signDate }}</a-descriptions-item>
        <a-descriptions-item label="管理医生">{{ contract.doctorId }}</a-descriptions-item>
      </a-descriptions>
      <a-empty v-else description="暂无慢病签约信息" />
    </a-card>

    <a-card title="健康数据" size="small" style="margin-top: 12px">
      <div v-for="item in healthData" :key="item.id" class="data-item">
        <div class="data-type">{{ item.dataType }}</div>
        <div class="data-value">{{ item.dataValue }}</div>
        <div class="data-time">{{ item.createdAt }}</div>
      </div>
      <a-empty v-if="!healthData.length" description="暂无健康数据" />
    </a-card>

    <a-modal v-model:open="modalOpen" title="上报健康数据" @ok="handleReport">
      <a-form layout="vertical">
        <a-form-item label="数据类型">
          <a-select v-model:value="form.dataType">
            <a-select-option value="blood_pressure">血压</a-select-option>
            <a-select-option value="blood_sugar">血糖</a-select-option>
            <a-select-option value="weight">体重</a-select-option>
            <a-select-option value="heart_rate">心率</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="数值">
          <a-input v-model:value="form.dataValue" placeholder="请输入数值" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { reportHealthData, getHealthData } from '@/api'
import { useAuthStore } from '@/store/auth'

const authStore = useAuthStore()
const contract = ref<any>(null)
const healthData = ref<any[]>([])
const modalOpen = ref(false)
const form = reactive({ dataType: 'blood_pressure', dataValue: '' })

async function fetchData() {
  try {
    healthData.value = await getHealthData(authStore.userInfo?.id || 'current-patient')
  } catch { healthData.value = [] }
}

function showReportModal() { modalOpen.value = true }

async function handleReport() {
  try {
    await reportHealthData(form)
    message.success('上报成功')
    modalOpen.value = false
    fetchData()
  } catch { }
}

onMounted(fetchData)
</script>

<style scoped>
.data-item {
  display: flex;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #f5f5f5;
}
.data-type { width: 80px; font-weight: 500; }
.data-value { flex: 1; font-size: 18px; color: #1890ff; }
.data-time { font-size: 12px; color: #999; }
</style>
