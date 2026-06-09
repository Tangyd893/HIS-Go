<template>
  <div>
    <a-card title="我的处方" size="small">
      <div v-for="item in prescriptions" :key="item.id" class="prescription-item" @click="showDetail(item)">
        <div class="pres-header">
          <span class="pres-doctor">{{ item.doctorName || '张医生' }}</span>
          <a-tag :color="statusColor[item.status]">{{ statusText[item.status] || item.status }}</a-tag>
        </div>
        <div class="pres-meta">
          <span>{{ formatDate(item.createdAt) }}</span>
          <span v-if="item.note">{{ item.note }}</span>
        </div>
        <div v-if="item.details?.length" class="pres-details">
          <div v-for="d in item.details.slice(0, 2)" :key="d.id" class="drug-item">
            <span class="drug-name-preview"><strong>{{ d.drugName }}</strong> {{ d.specification }}</span>
            <span class="drug-usage-preview">{{ d.dosage }}{{ d.unit || '' }} · {{ d.frequency }} · {{ d.usage }} · {{ d.days }}天 ×{{ d.quantity }}</span>
          </div>
          <div v-if="item.details.length > 2" class="drug-more">…共 {{ item.details.length }} 种药品，点击查看详情</div>
        </div>
      </div>
      <a-empty v-if="!prescriptions.length" description="暂无处方记录" />
    </a-card>

    <a-drawer
      v-if="selectedPres"
      :open="drawerOpen"
      title="处方详情"
      placement="bottom"
      height="70%"
      @close="drawerOpen = false"
    >
      <a-descriptions :column="2" size="small" bordered>
        <a-descriptions-item label="医生">{{ selectedPres.doctorName || '张医生' }}</a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag :color="statusColor[selectedPres.status]">{{ statusText[selectedPres.status] || selectedPres.status }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="日期">{{ formatDate(selectedPres.createdAt) }}</a-descriptions-item>
        <a-descriptions-item label="诊断">{{ selectedPres.note || '—' }}</a-descriptions-item>
      </a-descriptions>
      <a-divider>药品明细</a-divider>
      <div v-for="d in selectedPres.details" :key="d.id" class="detail-drug">
        <div class="detail-drug-header">
          <strong>{{ d.drugName }}</strong>
          <span class="detail-drug-spec">{{ d.specification }}</span>
        </div>
        <div class="detail-drug-info">
          用量: {{ d.dosage }}{{ d.unit || '' }} | 用法: {{ d.usage }} | 频次: {{ d.frequency }} | 天数: {{ d.days }}天
        </div>
        <div class="detail-drug-footer">
          <span>数量: {{ d.quantity }}</span>
          <span class="detail-drug-price">单价: ¥{{ d.unitPrice ?? '—' }}</span>
        </div>
      </div>
    </a-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getPrescriptions } from '@/api'
import { useAuthStore } from '@/store/auth'
import { resolvePatientId } from '@/utils/patient'

const authStore = useAuthStore()
const prescriptions = ref<any[]>([])
const drawerOpen = ref(false)
const selectedPres = ref<any>(null)
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

function showDetail(item: any) {
  selectedPres.value = item
  drawerOpen.value = true
}

onMounted(fetchData)
</script>

<style scoped>
.prescription-item {
  padding: 12px;
  border: 1px solid #f0f0f0;
  border-radius: 8px;
  margin-bottom: 8px;
  cursor: pointer;
}
.prescription-item:active { background: #f5f5f5; }
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
.drug-name-preview { font-size: 13px; }
.drug-usage-preview { font-size: 12px; color: #999; display: block; margin-top: 2px; }
.drug-more { font-size: 12px; color: #1890ff; margin-top: 6px; text-align: center; }
.detail-drug {
  padding: 10px;
  border: 1px solid #f0f0f0;
  border-radius: 6px;
  margin-bottom: 8px;
}
.detail-drug-header { display: flex; gap: 8px; align-items: baseline; margin-bottom: 4px; }
.detail-drug-spec { font-size: 12px; color: #999; }
.detail-drug-info { font-size: 13px; color: #666; margin-bottom: 4px; }
.detail-drug-footer { display: flex; justify-content: space-between; font-size: 13px; }
.detail-drug-price { color: #fa8c16; }
</style>
