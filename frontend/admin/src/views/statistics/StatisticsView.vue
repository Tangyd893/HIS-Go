<template>
  <a-card title="数据统计">
    <a-row :gutter="16">
      <a-col :span="8">
        <a-card hoverable>
          <a-statistic title="总挂号量" :value="stats.totalRegistrations" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card hoverable>
          <a-statistic title="总门诊量" :value="stats.totalClinic" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card hoverable>
          <a-statistic title="总处方数" :value="stats.totalPrescriptions" />
        </a-card>
      </a-col>
    </a-row>
    <a-row :gutter="16" style="margin-top: 16px">
      <a-col :span="8">
        <a-card hoverable>
          <a-statistic title="总结算金额" :value="stats.totalBilling" prefix="¥" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card hoverable>
          <a-statistic title="住院人数" :value="stats.totalInpatient" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card hoverable>
          <a-statistic title="运营天数" :value="30" />
        </a-card>
      </a-col>
    </a-row>

    <a-card title="运营报表查询" style="margin-top: 24px">
      <a-space style="margin-bottom: 16px">
        <a-date-picker v-model:value="startDate" placeholder="开始日期" />
        <a-date-picker v-model:value="endDate" placeholder="结束日期" />
        <a-button type="primary" @click="loadOperation">查询运营数据</a-button>
        <a-button @click="loadRevenue">收入趋势</a-button>
        <a-button @click="loadWorkload">科室工作量</a-button>
      </a-space>
      <a-empty v-if="!chartData.length" description="请选择日期范围查询" />
      <a-table v-else :columns="chartColumns" :data-source="chartData" row-key="id" size="small" />
    </a-card>
  </a-card>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { statisticsApi } from '@/api/others'
import dayjs from 'dayjs'

const stats = reactive({ totalRegistrations: 1280, totalClinic: 960, totalPrescriptions: 740, totalBilling: 158000, totalInpatient: 45 })
const startDate = ref(dayjs().subtract(30, 'day'))
const endDate = ref(dayjs())
const chartData = ref<any[]>([])
const chartColumns = ref<any[]>([])

async function loadOperation() {
  try {
    const data = await statisticsApi.getOperation({
      start_date: startDate.value?.format('YYYY-MM-DD') || '',
      end_date: endDate.value?.format('YYYY-MM-DD') || '',
    })
    if (data) {
      Object.assign(stats, data)
      chartData.value = [{ ...data, id: '1' }]
      chartColumns.value = Object.keys(data).map(k => ({ title: k, dataIndex: k }))
    }
  } catch { }
}

async function loadRevenue() {
  try {
    const data = await statisticsApi.getRevenueTrend({
      start_date: startDate.value?.format('YYYY-MM-DD') || '',
      end_date: endDate.value?.format('YYYY-MM-DD') || '',
    })
    if (Array.isArray(data)) {
      chartData.value = data
      chartColumns.value = data.length ? Object.keys(data[0]).map(k => ({ title: k, dataIndex: k })) : []
    }
  } catch { }
}

async function loadWorkload() {
  try {
    const data = await statisticsApi.getDeptWorkload({
      start_date: startDate.value?.format('YYYY-MM-DD') || '',
      end_date: endDate.value?.format('YYYY-MM-DD') || '',
    })
    if (Array.isArray(data)) {
      chartData.value = data
      chartColumns.value = data.length ? Object.keys(data[0]).map(k => ({ title: k, dataIndex: k })) : []
    }
  } catch { }
}

onMounted(loadOperation)
</script>
