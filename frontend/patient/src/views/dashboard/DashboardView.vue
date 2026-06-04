<template>
  <div class="dashboard">
    <div class="user-card" v-if="authStore.userInfo">
      <a-avatar :size="64" :style="{ backgroundColor: '#1890ff' }">{{ authStore.username.charAt(0).toUpperCase() }}</a-avatar>
      <div class="user-info">
        <h2>{{ authStore.userInfo.realName || authStore.username }}</h2>
        <p>{{ authStore.userInfo.deptName || '' }}</p>
      </div>
    </div>

    <a-row :gutter="[12, 12]" class="service-grid">
      <a-col :span="8" v-for="svc in services" :key="svc.path">
        <div class="service-card" @click="router.push(svc.path)">
          <component :is="svc.icon" class="service-icon" :style="{ color: svc.color }" />
          <span class="service-label">
            {{ svc.label }}
            <span v-if="svc.count != null" class="count-badge">{{ svc.count }}</span>
          </span>
        </div>
      </a-col>
    </a-row>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'
import { getPrescriptions, getReports, getFollowupPlans } from '@/api'
import {
  ScheduleOutlined, MessageOutlined, FileTextOutlined,
  FileSearchOutlined, HeartOutlined, SafetyCertificateOutlined,
} from '@ant-design/icons-vue'

const router = useRouter()
const authStore = useAuthStore()

const prescriptionCount = ref<number | null>(null)
const reportCount = ref<number | null>(null)
const followupCount = ref<number | null>(null)

onMounted(async () => {
  try { const r = await getPrescriptions({ page: 1, pageSize: 1 }); prescriptionCount.value = r?.total ?? 0 } catch { prescriptionCount.value = 0 }
  try { const r = await getReports({ page: 1, pageSize: 1 }); reportCount.value = r?.total ?? 0 } catch { reportCount.value = 0 }
  try { const r = await getFollowupPlans({ page: 1, pageSize: 1 }); followupCount.value = r?.total ?? 0 } catch { followupCount.value = 0 }
})

const services = [
  { path: '/appointment', label: '预约挂号', icon: ScheduleOutlined, color: '#1890ff', count: null },
  { path: '/consultation', label: '在线问诊', icon: MessageOutlined, color: '#52c41a', count: null },
  { path: '/prescription', label: '我的处方', icon: FileTextOutlined, color: '#faad14', count: prescriptionCount },
  { path: '/report', label: '检查报告', icon: FileSearchOutlined, color: '#722ed1', count: reportCount },
  { path: '/health-record', label: '健康档案', icon: HeartOutlined, color: '#eb2f96', count: null },
  { path: '/chronic', label: '慢病管理', icon: SafetyCertificateOutlined, color: '#13c2c2', count: null },
]
</script>

<style scoped>
.user-card {
  display: flex;
  align-items: center;
  gap: 16px;
  background: linear-gradient(135deg, #1890ff, #36cfc9);
  padding: 24px;
  border-radius: 12px;
  color: #fff;
  margin-bottom: 16px;
}

.user-info h2 {
  margin: 0;
  font-size: 20px;
}

.user-info p {
  margin: 4px 0 0;
  opacity: 0.8;
  font-size: 13px;
}

.service-grid {
  margin: 0;
}

.service-card {
  background: #fff;
  border-radius: 12px;
  padding: 16px 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
}

.service-icon {
  font-size: 32px;
}

.service-label {
  font-size: 13px;
  color: #333;
  display: flex;
  align-items: center;
  gap: 4px;
}

.count-badge {
  background: #1890ff;
  color: #fff;
  font-size: 11px;
  border-radius: 10px;
  padding: 0 6px;
  min-width: 18px;
  text-align: center;
  line-height: 18px;
}
</style>
