<template>
  <a-card title="我的随访" size="small">
    <div v-for="item in plans" :key="item.id" class="followup-item" @click="viewDetail(item)">
      <div class="fu-header">
        <span class="fu-name">{{ item.planName || '随访计划' }}</span>
        <a-tag :color="item.status === 1 ? 'green' : 'blue'">
          {{ item.status === 1 ? '进行中' : '待开始' }}
        </a-tag>
      </div>
      <div class="fu-date">{{ formatDate(item.startDate) }} ~ {{ formatDate(item.endDate) }}</div>
      <a-button size="small" type="link" @click.stop="viewDetail(item)">查看详情</a-button>
    </div>
    <a-empty v-if="!plans.length" description="暂无随访计划" />

    <a-modal v-model:open="modalOpen" title="随访详情" :footer="null" width="90%">
      <template v-if="selectedPlan">
        <a-descriptions :column="1" size="small" bordered>
          <a-descriptions-item label="计划名称">{{ selectedPlan.planName || '随访计划' }}</a-descriptions-item>
          <a-descriptions-item label="开始日期">{{ formatDate(selectedPlan.startDate) }}</a-descriptions-item>
          <a-descriptions-item label="结束日期">{{ formatDate(selectedPlan.endDate) }}</a-descriptions-item>
          <a-descriptions-item label="频次">{{ freqLabel[selectedPlan.frequency] || '每' + selectedPlan.frequency + '周' }}</a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="selectedPlan.status === 1 ? 'green' : 'blue'">{{ selectedPlan.status === 1 ? '进行中' : '待开始' }}</a-tag>
          </a-descriptions-item>
        </a-descriptions>
        <a-divider v-if="selectedPlan.tasks?.length">随访任务</a-divider>
        <a-timeline v-if="selectedPlan.tasks?.length">
          <a-timeline-item
            v-for="task in selectedPlan.tasks"
            :key="task.id"
            :color="task.status === 1 ? 'green' : 'gray'"
          >
            <div>{{ task.content }}</div>
            <div class="task-meta">
              <span>{{ formatDate(task.executeDate) }}</span>
              <a-tag :color="task.status === 1 ? 'green' : 'default'" size="small">
                {{ task.status === 1 ? '已完成' : '待执行' }}
              </a-tag>
            </div>
          </a-timeline-item>
        </a-timeline>
      </template>
    </a-modal>
  </a-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getFollowupPlans } from '@/api'
import { useAuthStore } from '@/store/auth'
import { resolvePatientId } from '@/utils/patient'

const plans = ref<any[]>([])
const authStore = useAuthStore()
const modalOpen = ref(false)
const selectedPlan = ref<any>(null)

const freqLabel: Record<number, string> = { 1: '每周', 2: '每两周', 3: '每月', 4: '每季度' }

function formatDate(value?: string) {
  if (!value) return ''
  return value.replace('T', ' ').slice(0, 10)
}

async function fetchData() {
  authStore.restoreUserInfo()
  const patientId = resolvePatientId(authStore.userInfo)
  try {
    const res: any = await getFollowupPlans({ patientId })
    plans.value = res?.list || []
  } catch { plans.value = [] }
}

function viewDetail(item: any) {
  selectedPlan.value = item
  modalOpen.value = true
}

onMounted(fetchData)
</script>

<style scoped>
.followup-item {
  padding: 12px;
  border: 1px solid #f0f0f0;
  border-radius: 8px;
  margin-bottom: 8px;
  cursor: pointer;
}
.followup-item:active { background: #f5f5f5; }
.fu-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.fu-name { font-size: 15px; font-weight: 500; }
.fu-date { font-size: 12px; color: #999; margin-top: 4px; margin-bottom: 8px; }
.task-meta { display: flex; gap: 8px; align-items: center; font-size: 12px; color: #999; margin-top: 4px; }
</style>
