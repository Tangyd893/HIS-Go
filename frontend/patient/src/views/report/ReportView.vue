<template>
  <a-card title="检查报告" size="small">
    <div v-for="item in reports" :key="item.id" class="report-item" @click="viewReport(item)">
      <div class="report-header">
        <span class="report-type">{{ item.examType }}-{{ item.examItem }}</span>
        <a-tag :color="statusColor[item.status]">{{ statusText[item.status] || item.status }}</a-tag>
      </div>
      <div class="report-date">{{ item.createdAt }}</div>
      <div v-if="item.conclusion || item.impression" class="report-result">{{ item.conclusion || item.impression }}</div>
    </div>
    <a-empty v-if="!reports.length" description="暂无检查报告" />

    <a-modal v-model:open="modalOpen" title="报告详情" :footer="null" width="90%">
      <template v-if="selectedReport">
        <a-descriptions :column="1" size="small" bordered>
          <a-descriptions-item label="检查类型">{{ selectedReport.examType }}</a-descriptions-item>
          <a-descriptions-item label="检查项目">{{ selectedReport.examItem }}</a-descriptions-item>
          <a-descriptions-item v-if="selectedReport.bodyPart" label="检查部位">{{ selectedReport.bodyPart }}</a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="statusColor[selectedReport.status]">{{ statusText[selectedReport.status] }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="日期">{{ selectedReport.createdAt }}</a-descriptions-item>
        </a-descriptions>
        <a-divider />
        <div v-if="selectedReport.findings" class="detail-section">
          <div class="detail-label">检查所见</div>
          <div class="detail-text">{{ selectedReport.findings }}</div>
        </div>
        <div v-if="selectedReport.impression" class="detail-section">
          <div class="detail-label">影像印象</div>
          <div class="detail-text">{{ selectedReport.impression }}</div>
        </div>
        <div v-if="selectedReport.conclusion" class="detail-section">
          <div class="detail-label">诊断结论</div>
          <div class="detail-text">{{ selectedReport.conclusion }}</div>
        </div>
      </template>
    </a-modal>
  </a-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getReports } from '@/api'
import { useAuthStore } from '@/store/auth'
import { resolvePatientId } from '@/utils/patient'

const reports = ref<any[]>([])
const authStore = useAuthStore()
const modalOpen = ref(false)
const selectedReport = ref<any>(null)
const statusText: Record<number, string> = { 0: '待审核', 1: '已审核', 2: '已出报告', 3: '已归档' }
const statusColor: Record<number, string> = { 0: 'orange', 1: 'blue', 2: 'green', 3: 'default' }

async function fetchData() {
  authStore.restoreUserInfo()
  const patientId = resolvePatientId(authStore.userInfo)
  try {
    const res: any = await getReports({ patientId })
    reports.value = res?.list || []
  } catch { reports.value = [] }
}

function viewReport(item: any) {
  selectedReport.value = item
  modalOpen.value = true
}

onMounted(fetchData)
</script>

<style scoped>
.report-item {
  padding: 12px;
  border: 1px solid #f0f0f0;
  border-radius: 8px;
  margin-bottom: 8px;
  cursor: pointer;
}
.report-item:active { background: #f5f5f5; }
.report-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.report-type { font-size: 15px; font-weight: 500; }
.report-date { font-size: 12px; color: #999; margin-top: 4px; }
.report-result { font-size: 13px; color: #666; margin-top: 8px; padding: 8px; background: #fafafa; border-radius: 4px; }
.detail-section { margin-bottom: 12px; }
.detail-label { font-size: 13px; font-weight: 500; color: #333; margin-bottom: 4px; }
.detail-text { font-size: 14px; color: #666; line-height: 1.6; padding: 8px; background: #fafafa; border-radius: 4px; }
</style>
