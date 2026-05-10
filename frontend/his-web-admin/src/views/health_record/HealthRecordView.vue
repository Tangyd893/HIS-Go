<template>
  <a-card title="健康档案">
    <a-form layout="inline" style="margin-bottom: 16px">
      <a-form-item label="患者ID">
        <a-input-search v-model:value="patientId" placeholder="输入患者ID查询" @search="loadSummary" />
      </a-form-item>
    </a-form>

    <a-descriptions v-if="summary" :column="2" bordered size="small" title="档案摘要">
      <a-descriptions-item v-for="(v, k) in summary" :key="k" :label="k">{{ v }}</a-descriptions-item>
    </a-descriptions>
    <a-empty v-else description="请输入患者ID查看健康档案" />

    <a-card v-if="timeline.length" title="健康时间轴" style="margin-top: 16px" size="small">
      <a-timeline>
        <a-timeline-item v-for="(item, index) in timeline" :key="index" :color="item.color || 'blue'">
          {{ item.event }} — {{ item.date }}
        </a-timeline-item>
      </a-timeline>
    </a-card>
  </a-card>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { healthRecordApi } from '@/api/others'

const patientId = ref('')
const summary = ref<any>(null)
const timeline = ref<any[]>([])

async function loadSummary() {
  if (!patientId.value) return
  try {
    summary.value = await healthRecordApi.getSummary(patientId.value)
    timeline.value = await healthRecordApi.getTimeline(patientId.value)
  } catch { summary.value = null; timeline.value = [] }
}
</script>
