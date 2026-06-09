<template>
  <div class="patient-app">
    <a-layout-header class="patient-header">
      <div class="header-title">医院患者端</div>
      <a-dropdown>
        <a-space>
          <a-avatar size="small" :style="{ backgroundColor: '#1890ff' }">{{ authStore.username.charAt(0).toUpperCase() }}</a-avatar>
          <span>{{ authStore.userInfo?.realName || authStore.username }}</span>
        </a-space>
        <template #overlay>
          <a-menu>
            <a-menu-item @click="authStore.logout">退出登录</a-menu-item>
          </a-menu>
        </template>
      </a-dropdown>
    </a-layout-header>

    <div class="patient-content">
      <router-view />
    </div>

    <div class="patient-tabbar">
      <div
        v-for="tab in tabs"
        :key="tab.path"
        class="tab-item"
        :class="{ active: currentPath === tab.path }"
        @click="navigate(tab.path)"
      >
        <component :is="tab.icon" />
        <span>{{ tab.label }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/store/auth'
import {
  HomeOutlined, ScheduleOutlined,
  FileTextOutlined, HeartOutlined, RobotOutlined,
  FileSearchOutlined,
} from '@ant-design/icons-vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const currentPath = computed(() => {
  const path = route.path.replace('/', '')
  return path || 'dashboard'
})

const tabs = [
  { path: 'dashboard', label: '首页', icon: HomeOutlined },
  { path: 'triage', label: '助手', icon: RobotOutlined },
  { path: 'appointment', label: '挂号', icon: ScheduleOutlined },
  { path: 'prescription', label: '处方', icon: FileTextOutlined },
  { path: 'report', label: '报告', icon: FileSearchOutlined },
  { path: 'health-record', label: '档案', icon: HeartOutlined },
]

function navigate(path: string) {
  router.push(`/${path}`)
}
</script>

<style scoped>
.patient-app {
  min-height: 100vh;
  background: #f5f5f5;
  display: flex;
  flex-direction: column;
}

.patient-header {
  background: #fff;
  padding: 12px 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #f0f0f0;
  position: sticky;
  top: 0;
  z-index: 10;
}

.header-title {
  font-size: 16px;
  font-weight: bold;
  color: #1890ff;
}

.patient-content {
  flex: 1;
  padding: 12px;
  padding-bottom: 72px;
  overflow-y: auto;
}

.patient-tabbar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: #fff;
  display: flex;
  border-top: 1px solid #f0f0f0;
  padding: 6px 0;
  padding-bottom: env(safe-area-inset-bottom, 6px);
  z-index: 10;
}

.tab-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2px;
  font-size: 12px;
  color: #999;
  cursor: pointer;
  padding: 4px 0;
}

.tab-item.active {
  color: #1890ff;
}

.tab-item .anticon {
  font-size: 20px;
}
</style>
