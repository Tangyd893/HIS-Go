<template>
  <div>
    <a-card title="药品管理">
      <template #extra>
        <a-space>
          <a-input-search v-model:value="searchName" placeholder="搜索药品" style="width: 200px" @search="fetchData" />
          <a-button type="primary" @click="showStockModal()"><PlusOutlined /> 入库</a-button>
          <a-button @click="scanExpired">扫描过期药品</a-button>
        </a-space>
      </template>
      <a-table :columns="columns" :data-source="dataSource" :loading="loading" :pagination="pagination" row-key="id" @change="onTableChange">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'expiryDate'">
            <span :style="{ color: isExpired(record.expiryDate) ? 'red' : '' }">{{ record.expiryDate }}</span>
          </template>
          <template v-if="column.key === 'action'">
            <a-space>
              <a-button size="small" @click="showStockModal(record)">入库</a-button>
              <a-button size="small" type="primary" @click="showDispenseModal(record)">发药</a-button>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:open="stockModalOpen" title="药品入库" @ok="handleAddStock">
      <a-form layout="vertical">
        <a-form-item label="药品">{{ selectedDrug?.name }}</a-form-item>
        <a-form-item label="入库数量"><a-input-number v-model:value="stockQuantity" :min="1" style="width: 100%" /></a-form-item>
      </a-form>
    </a-modal>

    <a-modal v-model:open="dispenseModalOpen" title="发药" @ok="handleDispense">
      <a-form layout="vertical">
        <a-form-item label="药品">{{ selectedDrug?.name }}</a-form-item>
        <a-form-item label="关联处方" required>
          <a-select
            v-model:value="dispenseForm.prescriptionId"
            show-search
            placeholder="搜索处方（患者名/ID）"
            :filter-option="false"
            :options="prescriptionOptions"
            @search="searchPrescriptions"
            @focus="searchPrescriptions('')"
            style="width: 100%"
          />
        </a-form-item>
        <a-form-item label="数量"><a-input-number v-model:value="dispenseForm.quantity" :min="1" style="width: 100%" /></a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { pharmacyApi } from '@/api/pharmacy'
import { prescriptionApi } from '@/api/prescription'

const loading = ref(false)
const searchName = ref('')
const dataSource = ref<any[]>([])
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const stockModalOpen = ref(false)
const dispenseModalOpen = ref(false)
const selectedDrug = ref<any>(null)
const stockQuantity = ref(1)
const dispenseForm = reactive({ prescriptionId: '', quantity: 1 })
const prescriptions = ref<any[]>([])
const prescriptionOptions = ref<{ label: string; value: string }[]>([])

const columns = [
  { title: '药品名称', dataIndex: 'name' },
  { title: '规格', dataIndex: 'specification' },
  { title: '厂家', dataIndex: 'manufacturer' },
  { title: '库存', dataIndex: 'stock' },
  { title: '单价', dataIndex: 'price' },
  { title: '有效期', key: 'expiryDate' },
  { title: '操作', key: 'action', width: 160 },
]

function isExpired(date: string): boolean { return new Date(date) < new Date() }

async function fetchData() {
  loading.value = true
  try {
    const res: any = await pharmacyApi.getDrugs({ page: pagination.current, pageSize: pagination.pageSize, name: searchName.value || undefined })
    dataSource.value = res?.list || []
    pagination.total = res?.total || 0
  } catch { message.error('加载药品失败'); dataSource.value = [] } finally { loading.value = false }
}

function onTableChange(pag: any) { pagination.current = pag.current; fetchData() }

function showStockModal(record?: any) { selectedDrug.value = record || null; stockQuantity.value = 1; stockModalOpen.value = true }
function showDispenseModal(record: any) {
  selectedDrug.value = record
  dispenseForm.prescriptionId = ''
  dispenseForm.quantity = 1
  prescriptionOptions.value = []
  prescriptions.value = []
  dispenseModalOpen.value = true
}

async function searchPrescriptions(keyword: string) {
  try {
    const res: any = await prescriptionApi.getList({ page: 1, pageSize: 30 })
    prescriptions.value = (res?.list || []).filter((p: any) =>
      p.status === 1 && (!keyword || p.patientName?.includes(keyword) || p.id?.includes(keyword))
    )
    prescriptionOptions.value = prescriptions.value.map((p: any) => ({
      label: `#${p.id?.slice(0, 8)} - ${p.patientName || '未知患者'} (${p.createdAt?.slice(0, 10) || ''})`,
      value: p.id,
    }))
  } catch { prescriptionOptions.value = [] }
}

async function handleAddStock() {
  try { await pharmacyApi.addStock(selectedDrug.value.id, stockQuantity.value); message.success('入库成功'); stockModalOpen.value = false; fetchData() } catch { message.error('入库失败') }
}

async function handleDispense() {
  try {
    await pharmacyApi.dispense({ prescription_id: dispenseForm.prescriptionId, drug_id: selectedDrug.value.id, quantity: dispenseForm.quantity, dispenser_id: 'current' })
    message.success('发药成功'); dispenseModalOpen.value = false; fetchData()
  } catch { message.error('发药失败') }
}

async function scanExpired() {
  try { const expired = await pharmacyApi.getExpired(); message.warning(`发现 ${(expired as any[])?.length || 0} 个过期药品`) } catch { }
}

onMounted(fetchData)
</script>
