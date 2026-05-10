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
          <span class="service-label">{{ svc.label }}</span>
        </div>
      </a-col>
    </a-row>

    <a-card title="最新动态" size="small" style="margin-top: 16px">
      <a-list item-layout="horizontal" :data-source="news" size="small">
        <template #renderItem="{ item }">
          <a-list-item>
            <a-list-item-meta :title="item.title" :description="item.time" />
          </a-list-item>
        </template>
      </a-list>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, h } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'
import {
  ScheduleOutlined, MessageOutlined, FileTextOutlined,
  FileSearchOutlined, HeartOutlined, SafetyCertificateOutlined,
  PhoneOutlined,
} from '@ant-design/icons-vue'

const router = useRouter()
const authStore = useAuthStore()

const services = [
  { path: '/appointment', label: '预约挂号', icon: ScheduleOutlined, color: '#1890ff' },
  { path: '/consultation', label: '在线问诊', icon: MessageOutlined, color: '#52c41a' },
  { path: '/prescription', label: '我的处方', icon: FileTextOutlined, color: '#faad14' },
  { path: '/report', label: '检查报告', icon: FileSearchOutlined, color: '#722ed1' },
  { path: '/health-record', label: '健康档案', icon: HeartOutlined, color: '#eb2f96' },
  { path: '/chronic', label: '慢病管理', icon: SafetyCertificateOutlined, color: '#13c2c2' },
]

const news = ref([
  { title: '您有一条新的随访提醒', time: '2026-05-10' },
  { title: '体检报告已生成，请查看', time: '2026-05-08' },
  { title: '慢病用药提醒：请按时服药', time: '2026-05-07' },
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
}
</style>
