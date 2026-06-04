<template>
  <a-card title="科室管理">
    <template #extra>
      <a-button type="primary"><PlusOutlined /> 新增科室</a-button>
    </template>
    <a-table :columns="columns" :data-source="dataSource" :loading="loading" row-key="id" />
  </a-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { userApi } from '@/api/user'

const loading = ref(false)
const dataSource = ref<any[]>([])

const columns = [
  { title: '科室名称', dataIndex: 'name' },
  { title: '上级科室', dataIndex: 'parentId' },
  { title: '描述', dataIndex: 'description' },
]

async function fetchData() {
  loading.value = true
  try {
    dataSource.value = await userApi.getDepartments()
  } catch { dataSource.value = [] } finally { loading.value = false }
}

onMounted(fetchData)
</script>
