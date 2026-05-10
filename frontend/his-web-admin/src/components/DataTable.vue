<template>
  <a-card :title="title" size="small">
    <a-table
      :columns="columns"
      :data-source="dataSource"
      :loading="loading"
      :pagination="{
        current,
        pageSize,
        total,
        showSizeChanger: true,
        showTotal: (t: number) => `共 ${t} 条`,
      }"
      row-key="id"
      @change="(pag: any) => $emit('pageChange', pag.current, pag.pageSize)"
    >
      <template #bodyCell="{ column, record }">
        <slot name="bodyCell" :column="column" :record="record" />
      </template>
    </a-table>
  </a-card>
</template>

<script setup lang="ts">
defineProps<{
  title: string
  columns: any[]
  dataSource: any[]
  loading: boolean
  current: number
  pageSize: number
  total: number
}>()

defineEmits<{
  pageChange: [page: number, pageSize: number]
}>()
</script>
