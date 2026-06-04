<template>
  <div>
    <a-card title="健康档案" size="small">
      <a-descriptions v-if="summary" :column="1" bordered size="small">
        <a-descriptions-item v-for="(v, k) in summary" :key="k" :label="k">{{ v }}</a-descriptions-item>
      </a-descriptions>
      <a-empty v-else description="暂无档案数据" />
    </a-card>

    <a-card title="健康时间轴" size="small" style="margin-top: 12px">
      <a-timeline>
        <a-timeline-item v-for="(item, idx) in timeline" :key="idx" :color="item.color || 'blue'">
          <div class="tl-event">{{ item.event }}</div>
          <div class="tl-date">{{ item.date }}</div>
        </a-timeline-item>
      </a-timeline>
      <a-empty v-if="!timeline.length" description="暂无时间轴数据" />
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getHealthSummary, getHealthTimeline } from '@/api'
import { useAuthStore } from '@/store/auth'

const summary = ref<any>(null)
const timeline = ref<any[]>([])
const authStore = useAuthStore()

async function fetchData() {
  const patientId = authStore.userInfo?.id || 'current-patient'
  try {
    summary.value = await getHealthSummary(patientId)
    timeline.value = await getHealthTimeline(patientId)
  } catch { summary.value = null; timeline.value = [] }
}

onMounted(fetchData)
</script>

<style scoped>
.tl-event { font-size: 14px; }
.tl-date { font-size: 12px; color: #999; }
</style>
