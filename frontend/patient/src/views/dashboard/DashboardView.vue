<template>
  <div class="dashboard">
    <div class="user-card" v-if="authStore.userInfo">
      <a-avatar :size="64" :style="{ backgroundColor: '#1890ff' }">{{ authStore.username.charAt(0).toUpperCase() }}</a-avatar>
      <div class="user-info">
        <h2>{{ authStore.userInfo.realName || authStore.username }}</h2>
        <p>欢迎使用患者服务中心</p>
      </div>
    </div>

    <div class="triage-banner" @click="router.push('/triage')">
      <div class="triage-banner-icon">🩺</div>
      <div class="triage-banner-text">
        <strong>就诊助手</strong>
        <span>不确定挂哪个科？描述症状，AI 为您推荐</span>
      </div>
      <span class="triage-banner-arrow">→</span>
    </div>

    <a-row :gutter="[12, 12]" class="service-grid">
      <a-col :span="8" v-for="svc in services" :key="svc.path">
        <div class="service-card" @click="router.push(svc.path)">
          <component :is="svc.icon" class="service-icon" :style="{ color: svc.color }" />
          <span class="service-label">
            {{ svc.label }}
            <span v-if="svc.count != null && svc.count > 0" class="count-badge">{{ svc.count }}</span>
          </span>
        </div>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'
import { getPrescriptions, getReports, getFollowupPlans } from '@/api'
import { resolvePatientId } from '@/utils/patient'
import {
  ScheduleOutlined, FileTextOutlined, FileSearchOutlined,
  HeartOutlined, SafetyCertificateOutlined, PhoneOutlined,
} from '@ant-design/icons-vue'

const router = useRouter()
const authStore = useAuthStore()

const prescriptionCount = ref<number | null>(null)
const reportCount = ref<number | null>(null)
const followupCount = ref<number | null>(null)

onMounted(async () => {
  authStore.restoreUserInfo()
  const patientId = resolvePatientId(authStore.userInfo)
  try { const r: any = await getPrescriptions({ patientId, page: 1, pageSize: 1 }); prescriptionCount.value = r?.total ?? 0 } catch { prescriptionCount.value = 0 }
  try { const r: any = await getReports({ patientId, page: 1, pageSize: 1 }); reportCount.value = r?.total ?? 0 } catch { reportCount.value = 0 }
  try { const r: any = await getFollowupPlans({ patientId, page: 1, pageSize: 1 }); followupCount.value = r?.total ?? 0 } catch { followupCount.value = 0 }
})

const services = computed(() => [
  { path: '/appointment', label: '预约挂号', icon: ScheduleOutlined, color: '#1890ff', count: null as number | null },
  { path: '/prescription', label: '我的处方', icon: FileTextOutlined, color: '#faad14', count: prescriptionCount.value },
  { path: '/report', label: '检查报告', icon: FileSearchOutlined, color: '#722ed1', count: reportCount.value },
  { path: '/health-record', label: '健康档案', icon: HeartOutlined, color: '#eb2f96', count: null as number | null },
  { path: '/chronic', label: '慢病管理', icon: SafetyCertificateOutlined, color: '#13c2c2', count: null as number | null },
  { path: '/followup', label: '我的随访', icon: PhoneOutlined, color: '#52c41a', count: followupCount.value },
])
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

.service-grid { margin: 0; }

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

.service-icon { font-size: 32px; }

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

.triage-banner {
  background: linear-gradient(135deg, #fff7e6, #fff1f0);
  border: 1px solid #ffd591;
  border-radius: 12px;
  padding: 14px 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
  cursor: pointer;
}
.triage-banner:hover { box-shadow: 0 4px 12px rgba(245, 34, 45, 0.12); }
.triage-banner-icon { font-size: 32px; flex-shrink: 0; }
.triage-banner-text { flex: 1; display: flex; flex-direction: column; gap: 2px; }
.triage-banner-text strong { font-size: 15px; color: #cf1322; }
.triage-banner-text span { font-size: 12px; color: #999; }
.triage-banner-arrow { font-size: 20px; color: #cf1322; flex-shrink: 0; }
</style>
