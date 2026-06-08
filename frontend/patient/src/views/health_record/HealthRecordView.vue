<template>
  <div>
    <a-card title="健康档案" size="small">
      <a-spin :spinning="loading">
        <div v-if="summary" class="summary-grid">
          <div class="stat-card">
            <div class="stat-num">{{ summary.totalVisits ?? 0 }}</div>
            <div class="stat-label">就诊次数</div>
          </div>
          <div class="stat-card">
            <div class="stat-num">{{ summary.totalPrescriptions ?? 0 }}</div>
            <div class="stat-label">处方记录</div>
          </div>
          <div class="stat-card">
            <div class="stat-num">{{ summary.totalExaminations ?? 0 }}</div>
            <div class="stat-label">检查报告</div>
          </div>
        </div>
        <a-empty v-else-if="!loading" description="暂无档案摘要" />
      </a-spin>
    </a-card>

    <a-card title="就诊时间轴" size="small" style="margin-top: 12px">
      <a-timeline v-if="timeline.length">
        <a-timeline-item
          v-for="item in timeline"
          :key="item.id"
          :color="eventColor[item.eventType] || 'blue'"
        >
          <div class="tl-head">
            <a-tag>{{ formatEventType(item.eventType) }}</a-tag>
            <span class="tl-date">{{ item.date }}</span>
          </div>
          <div class="tl-desc">{{ item.description }}</div>
        </a-timeline-item>
      </a-timeline>
      <a-empty v-else description="暂无就诊记录，完成挂号后将自动汇总" />
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getHealthSummary, getHealthTimeline } from '@/api'
import { useAuthStore } from '@/store/auth'
import { formatEventType, resolvePatientId } from '@/utils/patient'

const loading = ref(false)
const summary = ref<any>(null)
const timeline = ref<any[]>([])
const authStore = useAuthStore()

const eventColor: Record<string, string> = {
  visit: 'green',
  prescription: 'orange',
  examination: 'purple',
  followup: 'blue',
}

async function fetchData() {
  authStore.restoreUserInfo()
  const patientId = resolvePatientId(authStore.userInfo)
  loading.value = true
  try {
    summary.value = await getHealthSummary(patientId)
    timeline.value = await getHealthTimeline(patientId) || []
  } catch {
    summary.value = null
    timeline.value = []
  } finally {
    loading.value = false
  }
}

onMounted(fetchData)
</script>

<style scoped>
.summary-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
}
.stat-card {
  background: #f6ffed;
  border: 1px solid #b7eb8f;
  border-radius: 8px;
  padding: 12px 8px;
  text-align: center;
}
.stat-card:nth-child(2) { background: #fff7e6; border-color: #ffd591; }
.stat-card:nth-child(3) { background: #f9f0ff; border-color: #d3adf7; }
.stat-num { font-size: 22px; font-weight: 600; color: #333; }
.stat-label { font-size: 12px; color: #666; margin-top: 4px; }
.tl-head { display: flex; align-items: center; gap: 8px; margin-bottom: 4px; }
.tl-date { font-size: 12px; color: #999; }
.tl-desc { font-size: 14px; color: #333; }
</style>
