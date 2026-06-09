<template>
  <div>
    <a-card title="系统设置">
      <a-tabs v-model:activeKey="activeTab">
        <a-tab-pane key="dict" tab="字典管理">
          <a-select v-model:value="currentDictType" placeholder="选择字典类型" style="width: 200px; margin-bottom: 16px" @change="loadDictItems">
            <a-select-option v-for="d in dictTypes" :key="d.id" :value="d.dictType">{{ d.dictLabel }}</a-select-option>
          </a-select>
          <a-table :columns="dictColumns" :data-source="dictItems" row-key="id" size="small" />
        </a-tab-pane>
        <a-tab-pane key="params" tab="参数设置">
          <a-table :columns="paramColumns" :data-source="params" row-key="id" size="small">
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'action'">
                <a @click="editParam(record)">编辑</a>
                <a-divider type="vertical" />
                <a @click="viewDetail(record)">详情</a>
              </template>
            </template>
          </a-table>
        </a-tab-pane>
        <a-tab-pane key="logs" tab="操作日志">
          <a-table :columns="logColumns" :data-source="logs" :loading="logLoading" :pagination="logPagination" row-key="id" size="small" />
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <a-modal v-model:open="detailOpen" title="参数详情" :footer="null" width="500px">
      <a-descriptions v-if="detailRecord" :column="1" bordered size="small">
        <a-descriptions-item v-for="(v, k) in detailRecord" :key="k" :label="k">{{ v }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { systemApi } from '@/api/others'

const activeTab = ref('dict')
const dictTypes = ref<any[]>([])
const currentDictType = ref('')
const dictItems = ref<any[]>([])
const params = ref<any[]>([])
const logs = ref<any[]>([])
const logLoading = ref(false)
const detailOpen = ref(false)
const detailRecord = ref<any>(null)
const logPagination = reactive({ current: 1, pageSize: 10, total: 0 })

const dictColumns = [
  { title: '标签', dataIndex: 'dictLabel' },
  { title: '值', dataIndex: 'dictValue' },
  { title: '排序', dataIndex: 'sort' },
]

const paramColumns = [
  { title: '参数名', dataIndex: 'paramKey' },
  { title: '参数值', dataIndex: 'paramValue' },
  { title: '描述', dataIndex: 'description' },
  { title: '操作', key: 'action', width: 80 },
]

const logColumns = [
  { title: '用户ID', dataIndex: 'userId' },
  { title: '模块', dataIndex: 'module' },
  { title: '操作', dataIndex: 'operation' },
  { title: '时间', dataIndex: 'createdAt' },
]

async function loadDictTypes() {
  try { dictTypes.value = await systemApi.getDictTypes() } catch { }
}

async function loadDictItems() {
  if (!currentDictType.value) return
  try { dictItems.value = await systemApi.getDictItems(currentDictType.value) } catch { }
}

async function loadParams() {
  try { params.value = await systemApi.getParams() } catch { }
}

async function loadLogs() {
  logLoading.value = true
  try {
    const res: any = await systemApi.getOperationLogs({ page: logPagination.current, pageSize: logPagination.pageSize })
    logs.value = res?.list || []
    logPagination.total = res?.total || 0
  } catch { } finally { logLoading.value = false }
}

function editParam(record: any) { message.info(`编辑参数: ${record.paramKey}`) }
function viewDetail(record: any) { detailRecord.value = record; detailOpen.value = true }

onMounted(() => { loadDictTypes(); loadParams(); loadLogs() })
</script>
