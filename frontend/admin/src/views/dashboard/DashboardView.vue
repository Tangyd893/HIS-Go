<template>
  <div class="dashboard">
    <a-row :gutter="16">
      <a-col :span="6" v-for="card in statCards" :key="card.title">
        <a-card hoverable>
          <a-statistic
            :title="card.title"
            :value="card.value"
            :prefix="card.icon"
            :value-style="{ color: card.color }"
          />
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16" style="margin-top: 24px">
      <a-col :span="12">
        <a-card title="快捷操作">
          <a-space wrap>
            <a-button type="primary" @click="router.push('/registration')">
              <FormOutlined /> 预约挂号
            </a-button>
            <a-button @click="router.push('/clinic')">
              <MedicineBoxOutlined /> 门诊接诊
            </a-button>
            <a-button @click="router.push('/prescription')">
              <FileTextOutlined /> 开具处方
            </a-button>
            <a-button @click="router.push('/billing')">
              <DollarOutlined /> 收费结算
            </a-button>
            <a-button @click="router.push('/pharmacy')">
              <ExperimentOutlined /> 药房发药
            </a-button>
            <a-button @click="router.push('/inpatient')">
              <HomeOutlined /> 住院登记
            </a-button>
          </a-space>
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card title="系统信息">
          <a-descriptions :column="1" size="small">
            <a-descriptions-item label="系统版本">HIS-Go v1.0.0</a-descriptions-item>
            <a-descriptions-item label="当前用户">{{ authStore.username }}</a-descriptions-item>
            <a-descriptions-item label="角色">{{ authStore.role }}</a-descriptions-item>
            <a-descriptions-item label="后端服务">Go + Gin + gRPC</a-descriptions-item>
            <a-descriptions-item label="前端框架">Vue 3 + Ant Design Vue 4</a-descriptions-item>
            <a-descriptions-item label="数据库">PostgreSQL 17</a-descriptions-item>
          </a-descriptions>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16" style="margin-top: 24px" v-if="authStore.role === 'admin'">
      <a-col :span="24">
        <a-card title="服务模块状态（仅管理员可见）">
          <a-row :gutter="[12, 12]">
            <a-col :span="3" v-for="svc in services" :key="svc.name">
              <a-tag :color="svc.status === 'running' ? 'green' : 'orange'" class="service-tag">
                {{ svc.name }}
              </a-tag>
            </a-col>
          </a-row>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'
import { h } from 'vue'
import {
  FormOutlined, MedicineBoxOutlined, FileTextOutlined,
  DollarOutlined, ExperimentOutlined, HomeOutlined,
} from '@ant-design/icons-vue'

const router = useRouter()
const authStore = useAuthStore()

const statCards = ref([
  { title: '今日挂号', value: 128, icon: h(FormOutlined), color: '#1890ff' },
  { title: '今日门诊', value: 96, icon: h(MedicineBoxOutlined), color: '#52c41a' },
  { title: '住院患者', value: 45, icon: h(HomeOutlined), color: '#faad14' },
  { title: '今日收入', value: '¥12,580', icon: h(DollarOutlined), color: '#f5222d' },
])

const services = ref([
  { name: 'Gateway', status: 'running' },
  { name: 'Auth', status: 'running' },
  { name: 'User', status: 'running' },
  { name: 'Registration', status: 'running' },
  { name: 'Clinic', status: 'running' },
  { name: 'Prescription', status: 'running' },
  { name: 'Billing', status: 'running' },
  { name: 'Pharmacy', status: 'running' },
  { name: 'Examination', status: 'running' },
  { name: 'Inpatient', status: 'running' },
  { name: 'Schedule', status: 'running' },
  { name: 'Outpatient', status: 'running' },
  { name: 'Followup', status: 'running' },
  { name: 'HealthRecord', status: 'running' },
  { name: 'Notification', status: 'running' },
  { name: 'Statistics', status: 'running' },
  { name: 'System', status: 'running' },
  { name: 'EMR', status: 'running' },
])
</script>

<style scoped>
.dashboard .service-tag {
  font-size: 13px;
  padding: 4px 12px;
  margin: 2px;
  display: inline-block;
  width: 100%;
  text-align: center;
}
</style>
