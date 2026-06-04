<template>
  <a-card title="检查报告" size="small">
    <div v-for="item in reports" :key="item.id" class="report-item" @click="viewReport(item)">
      <div class="report-header">
        <span class="report-type">{{ item.examType }}-{{ item.examItem }}</span>
        <a-tag :color="statusColor[item.status]">{{ statusText[item.status] || item.status }}</a-tag>
      </div>
      <div class="report-date">{{ item.createdAt }}</div>
      <div v-if="item.result" class="report-result">{{ item.result }}</div>
    </div>
    <a-empty v-if="!reports.length" description="暂无检查报告" />
  </a-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { getReports } from '@/api'

const reports = ref<any[]>([])
const statusText: Record<number, string> = { 0: '待审核', 1: '已审核', 2: '已取消' }
const statusColor: Record<number, string> = { 0: 'orange', 1: 'green', 2: 'red' }

async function fetchData() {
  try {
    const res: any = await getReports({})
    reports.value = res?.list || []
  } catch { reports.value = [] }
}

function viewReport(item: any) {
  message.info(`查看报告详情: ${item.examItem}`)
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
</style>
