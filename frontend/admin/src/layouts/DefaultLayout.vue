<template>
  <a-layout style="min-height: 100vh">
    <a-layout-sider v-model:collapsed="appStore.collapsed" collapsible theme="dark" breakpoint="lg">
      <div class="logo">
        <span v-if="!appStore.collapsed">HIS-Go 管理系统</span>
        <span v-else>HIS</span>
      </div>
      <a-menu
        v-model:selectedKeys="selectedKeys"
        v-model:openKeys="openKeys"
        theme="dark"
        mode="inline"
        @click="onMenuClick"
      >
        <a-menu-item key="dashboard">
          <DashboardOutlined />
          <span>工作台</span>
        </a-menu-item>
        <a-sub-menu key="user">
          <template #title>
            <TeamOutlined />
            <span>用户管理</span>
          </template>
          <a-menu-item key="user/patients">患者管理</a-menu-item>
          <a-menu-item key="user/employees">员工管理</a-menu-item>
          <a-menu-item key="user/departments">科室管理</a-menu-item>
        </a-sub-menu>
        <a-menu-item key="registration">
          <FormOutlined />
          <span>挂号管理</span>
        </a-menu-item>
        <a-menu-item key="schedule">
          <ScheduleOutlined />
          <span>排班管理</span>
        </a-menu-item>
        <a-menu-item key="clinic">
          <MedicineBoxOutlined />
          <span>门诊诊疗</span>
        </a-menu-item>
        <a-menu-item key="prescription">
          <FileTextOutlined />
          <span>处方管理</span>
        </a-menu-item>
        <a-menu-item key="pharmacy">
          <ExperimentOutlined />
          <span>药房管理</span>
        </a-menu-item>
        <a-menu-item key="billing">
          <DollarOutlined />
          <span>收费结算</span>
        </a-menu-item>
        <a-sub-menu key="inpatient-group">
          <template #title>
            <HomeOutlined />
            <span>住院管理</span>
          </template>
          <a-menu-item key="inpatient">住院记录</a-menu-item>
          <a-menu-item key="emr">电子病历</a-menu-item>
        </a-sub-menu>
        <a-menu-item key="examination">
          <FundProjectionScreenOutlined />
          <span>检查检验</span>
        </a-menu-item>
        <a-sub-menu key="outpatient-group">
          <template #title>
            <GlobalOutlined />
            <span>院外服务</span>
          </template>
          <a-menu-item key="outpatient">院外患者</a-menu-item>
          <a-menu-item key="followup">随访管理</a-menu-item>
          <a-menu-item key="health-record">健康档案</a-menu-item>
        </a-sub-menu>
        <a-menu-item key="notification">
          <BellOutlined />
          <span>消息通知</span>
        </a-menu-item>
        <a-menu-item key="statistics">
          <BarChartOutlined />
          <span>数据统计</span>
        </a-menu-item>
        <a-menu-item key="system">
          <SettingOutlined />
          <span>系统设置</span>
        </a-menu-item>
      </a-menu>
    </a-layout-sider>
    <a-layout>
      <a-layout-header class="header">
        <div class="header-left">
          <MenuFoldOutlined v-if="!appStore.collapsed" @click="appStore.toggleCollapsed" class="trigger" />
          <MenuUnfoldOutlined v-else @click="appStore.toggleCollapsed" class="trigger" />
          <a-breadcrumb class="breadcrumb">
            <a-breadcrumb-item v-for="item in breadcrumbs" :key="item">{{ item }}</a-breadcrumb-item>
          </a-breadcrumb>
        </div>
        <div class="header-right">
          <a-dropdown>
            <a-space>
              <a-avatar size="small" :style="{ backgroundColor: '#1890ff' }">
                {{ (authStore.username || '?').charAt(0).toUpperCase() }}
              </a-avatar>
              <span>{{ authStore.username }}</span>
              <DownOutlined />
            </a-space>
            <template #overlay>
              <a-menu>
                <a-menu-item @click="authStore.logout">
                  <LogoutOutlined /> 退出登录
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </a-layout-header>
      <a-layout-content class="content">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/store/auth'
import { useAppStore } from '@/store/app'
import {
  DashboardOutlined, TeamOutlined, FormOutlined, ScheduleOutlined,
  MedicineBoxOutlined, FileTextOutlined, ExperimentOutlined, DollarOutlined,
  HomeOutlined, FundProjectionScreenOutlined, GlobalOutlined, BellOutlined,
  BarChartOutlined, SettingOutlined,
  MenuFoldOutlined, MenuUnfoldOutlined, DownOutlined, LogoutOutlined,
} from '@ant-design/icons-vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const appStore = useAppStore()

const selectedKeys = ref<string[]>([route.path.replace('/', '') || 'dashboard'])
const openKeys = ref<string[]>(['user'])

const breadcrumbs = computed(() => {
  const matched = route.matched.filter(r => r.meta.title)
  return matched.map(r => r.meta.title as string)
})

function onMenuClick({ key }: { key: string }) {
  router.push(`/${key}`)
}

watch(() => route.path, (path) => {
  selectedKeys.value = [path.replace('/', '') || 'dashboard']
})
</script>

<style scoped>
.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 18px;
  font-weight: bold;
  background: rgba(255, 255, 255, 0.1);
  overflow: hidden;
  white-space: nowrap;
}

.header {
  background: #fff;
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
  z-index: 1;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-right {
  cursor: pointer;
}

.trigger {
  font-size: 18px;
  cursor: pointer;
}

.content {
  margin: 16px;
  padding: 24px;
  background: #fff;
  border-radius: 4px;
  min-height: 280px;
  overflow-y: auto;
}
</style>
