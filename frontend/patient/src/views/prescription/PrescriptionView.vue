<template>
  <div>
    <a-card title="我的处方" size="small">
      <div v-for="item in prescriptions" :key="item.id" class="prescription-item">
        <div class="pres-header">
          <span class="pres-doctor">{{ item.doctorName || '医生' }}</span>
          <a-tag :color="statusColor[item.status]">{{ statusText[item.status] || item.status }}</a-tag>
        </div>
        <div class="pres-date">{{ item.createdAt }}</div>
        <div v-if="item.details" class="pres-details">
          <div v-for="d in item.details" :key="d.id" class="drug-item">
            <span>{{ d.drugName }}</span>
            <span>{{ d.quantity }} × {{ d.usage }}</span>
          </div>
        </div>
      </div>
      <a-empty v-if="!prescriptions.length" description="暂无处方记录" />
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getPrescriptions } from '@/api'

const prescriptions = ref<any[]>([])
const statusText: Record<number, string> = { 0: '待审核', 1: '已审核', 2: '已发药', 3: '已取消' }
const statusColor: Record<number, string> = { 0: 'orange', 1: 'blue', 2: 'green', 3: 'red' }

async function fetchData() {
  try {
    const res: any = await getPrescriptions({})
    prescriptions.value = res?.list || []
  } catch { prescriptions.value = [] }
}

onMounted(fetchData)
</script>

<style scoped>
.prescription-item {
  padding: 12px;
  border: 1px solid #f0f0f0;
  border-radius: 8px;
  margin-bottom: 8px;
}
.pres-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.pres-doctor { font-size: 15px; font-weight: 500; }
.pres-date { font-size: 12px; color: #999; margin-top: 4px; }
.pres-details { margin-top: 8px; }
.drug-item {
  display: flex;
  justify-content: space-between;
  padding: 4px 0;
  font-size: 14px;
  border-bottom: 1px dashed #f0f0f0;
}
</style>
