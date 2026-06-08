<template>
  <div>
    <a-card title="我的处方" size="small">
      <div v-for="item in prescriptions" :key="item.id" class="prescription-item">
        <div class="pres-header">
          <span class="pres-doctor">{{ item.doctorName || '张医生' }}</span>
          <a-tag :color="statusColor[item.status]">{{ statusText[item.status] || item.status }}</a-tag>
        </div>
        <div class="pres-meta">
          <span>{{ formatDate(item.createdAt) }}</span>
          <span v-if="item.note">{{ item.note }}</span>
        </div>
        <div v-if="item.details?.length" class="pres-details">
          <div v-for="d in item.details" :key="d.id" class="drug-item">
            <div class="drug-name">
              <strong>{{ d.drugName }}</strong>
              <span class="spec">{{ d.specification }}</span>
            </div>
            <div class="drug-usage">
              {{ d.dosage }}{{ d.unit || '' }} · {{ d.frequency }} · {{ d.usage }} · {{ d.days }}天
            </div>
            <div class="drug-qty">× {{ d.quantity }}</div>
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
import { useAuthStore } from '@/store/auth'
import { resolvePatientId } from '@/utils/patient'

const authStore = useAuthStore()
const prescriptions = ref<any[]>([])
const statusText: Record<number, string> = {
  0: '草稿',
  1: '待审核',
  2: '已审核',
  3: '已收费',
  4: '已发药',
}
const statusColor: Record<number, string> = {
  0: 'default',
  1: 'orange',
  2: 'blue',
  3: 'cyan',
  4: 'green',
}

function formatDate(value?: string) {
  if (!value) return ''
  return value.replace('T', ' ').slice(0, 16)
}

async function fetchData() {
  authStore.restoreUserInfo()
  const patientId = resolvePatientId(authStore.userInfo)
  try {
    const res: any = await getPrescriptions({ patientId })
    prescriptions.value = res?.list || res || []
  } catch {
    prescriptions.value = []
  }
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
.pres-meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}
.pres-details { margin-top: 8px; }
.drug-item {
  padding: 8px 0;
  border-bottom: 1px dashed #f0f0f0;
}
.drug-item:last-child { border-bottom: none; }
.drug-name { display: flex; gap: 8px; align-items: baseline; }
.spec { font-size: 12px; color: #999; }
.drug-usage { font-size: 13px; color: #666; margin-top: 2px; }
.drug-qty { font-size: 13px; color: #333; margin-top: 2px; }
</style>
