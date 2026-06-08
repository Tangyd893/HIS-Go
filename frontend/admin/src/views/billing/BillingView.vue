<template>
  <div>
    <a-card title="收费结算">
      <a-table :columns="columns" :data-source="dataSource" :loading="loading" :pagination="pagination" row-key="id" @change="onTableChange">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="statusColor[record.status]">{{ statusText[record.status] || record.status }}</a-tag>
          </template>
          <template v-if="column.key === 'action'">
            <a-space>
              <a @click="viewDetail(record)">详情</a>
              <a-button size="small" type="primary" @click="handlePay(record)" v-if="record.status === 0">收费</a-button>
              <a-button size="small" danger @click="handleRefund(record.id)" v-if="record.status === 1">退费</a-button>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:open="detailOpen" title="账单详情" :footer="null" width="600px">
      <a-descriptions v-if="detailRecord" :column="2" bordered size="small">
        <a-descriptions-item v-for="(v, k) in detailRecord" :key="k" :label="k">{{ v }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { billingApi } from '@/api/billing'

const loading = ref(false)
const dataSource = ref<any[]>([])
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const detailOpen = ref(false)
const detailRecord = ref<any>(null)

const statusText: Record<number, string> = { 0: '待支付', 1: '已支付', 2: '已退费' }
const statusColor: Record<number, string> = { 0: 'orange', 1: 'green', 2: 'red' }

const columns = [
  { title: '账单号', dataIndex: 'billNo' },
  { title: '患者', dataIndex: 'patientName' },
  { title: '总金额', dataIndex: 'totalAmount' },
  { title: '已付金额', dataIndex: 'paidAmount' },
  { title: '状态', key: 'status' },
  { title: '创建时间', dataIndex: 'createdAt' },
  { title: '操作', key: 'action', width: 200 },
]

async function fetchData() {
  loading.value = true
  try {
    const res: any = await billingApi.getList({ page: pagination.current, pageSize: pagination.pageSize })
    dataSource.value = res?.list || []
    pagination.total = res?.total || 0
  } catch { message.error('加载账单失败'); dataSource.value = [] } finally { loading.value = false }
}

function onTableChange(pag: any) { pagination.current = pag.current; fetchData() }
function viewDetail(record: any) { detailRecord.value = record; detailOpen.value = true }

async function handlePay(record: any) {
  try { await billingApi.pay(record.id, 0); message.success('支付成功'); fetchData() } catch { message.error('支付失败') }
}

async function handleRefund(id: string) {
  try { await billingApi.refund(id); message.success('退费成功'); fetchData() } catch { message.error('退费失败') }
}

onMounted(fetchData)
</script>
