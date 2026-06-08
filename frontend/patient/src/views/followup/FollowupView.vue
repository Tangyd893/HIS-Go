<template>
  <a-card title="我的随访" size="small">
    <div v-for="item in plans" :key="item.id" class="followup-item">
      <div class="fu-header">
        <span class="fu-name">{{ item.planName || '随访计划' }}</span>
        <a-tag :color="item.status === 1 ? 'green' : 'blue'">
          {{ item.status === 1 ? '进行中' : '待开始' }}
        </a-tag>
      </div>
      <div class="fu-date">{{ item.createdAt }}</div>
      <a-button size="small" type="link" @click="viewDetail(item)">查看详情</a-button>
    </div>
    <a-empty v-if="!plans.length" description="暂无随访计划" />
  </a-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { getFollowupPlans } from '@/api'
import { useAuthStore } from '@/store/auth'
import { resolvePatientId } from '@/utils/patient'

const plans = ref<any[]>([])
const authStore = useAuthStore()

async function fetchData() {
  authStore.restoreUserInfo()
  const patientId = resolvePatientId(authStore.userInfo)
  try {
    const res: any = await getFollowupPlans({ patientId })
    plans.value = res?.list || []
  } catch { plans.value = [] }
}

function viewDetail(item: any) { message.info(`查看随访计划: ${item.id}`) }

onMounted(fetchData)
</script>

<style scoped>
.followup-item {
  padding: 12px;
  border: 1px solid #f0f0f0;
  border-radius: 8px;
  margin-bottom: 8px;
}
.fu-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.fu-name { font-size: 15px; font-weight: 500; }
.fu-date { font-size: 12px; color: #999; margin-top: 4px; margin-bottom: 8px; }
</style>
