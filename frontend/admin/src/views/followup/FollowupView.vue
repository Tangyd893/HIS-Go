<template>
  <ServicePlaceholder v-if="serviceUnavailable" title="随访管理" />
  <a-card v-else title="随访管理">
    <a-table :columns="columns" :data-source="dataSource" :loading="loading" :pagination="pagination" row-key="id" @change="onTableChange" />
  </a-card>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { followupApi } from '@/api/others'
import ServicePlaceholder from '@/components/ServicePlaceholder.vue'

const serviceUnavailable = ref(false)
const loading = ref(false)
const dataSource = ref<any[]>([])
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })

const columns = [
  { title: '患者', dataIndex: 'patientId' },
  { title: '计划名称', dataIndex: 'planName' },
  { title: '状态', dataIndex: 'status' },
  { title: '创建时间', dataIndex: 'createdAt' },
]

async function fetchData() {
  loading.value = true
  try {
    const res: any = await followupApi.getPlans({ page: pagination.current, pageSize: pagination.pageSize })
    dataSource.value = res?.list || []
    pagination.total = res?.total || 0
  } catch { serviceUnavailable.value = true; dataSource.value = [] } finally { loading.value = false }
}

function onTableChange(pag: any) { pagination.current = pag.current; fetchData() }

onMounted(fetchData)
</script>
