<template>
  <div>
    <a-card title="患者管理">
      <template #extra>
        <a-space>
          <a-input-search
            v-model:value="searchName"
            placeholder="搜索患者姓名"
            style="width: 200px"
            @search="fetchData"
          />
          <a-button type="primary" @click="showCreateModal">
            <PlusOutlined /> 新增患者
          </a-button>
        </a-space>
      </template>
      <a-table
        :columns="columns"
        :data-source="dataSource"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        @change="onTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'gender'">
            {{ genderMap[record.gender] || record.gender }}
          </template>
          <template v-if="column.key === 'action'">
            <a-space>
              <a @click="viewDetail(record)">详情</a>
              <a @click="editRecord(record)">编辑</a>
              <a-popconfirm title="确认删除？" @confirm="deleteRecord(record.id)">
                <a style="color: red">删除</a>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:open="modalOpen" :title="modalTitle" @ok="handleSubmit" :confirm-loading="submitting">
      <a-form :model="formState" layout="vertical">
        <a-form-item label="姓名" required>
          <a-input v-model:value="formState.name" />
        </a-form-item>
        <a-form-item label="性别">
          <a-select v-model:value="formState.gender">
            <a-select-option value="M">男</a-select-option>
            <a-select-option value="F">女</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="手机号">
          <a-input v-model:value="formState.phone" />
        </a-form-item>
        <a-form-item label="身份证号">
          <a-input v-model:value="formState.idCard" />
        </a-form-item>
        <a-form-item label="地址">
          <a-input v-model:value="formState.address" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { userApi } from '@/api/user'
import { genderMap } from '@/utils'

const loading = ref(false)
const searchName = ref('')
const dataSource = ref<any[]>([])
const modalOpen = ref(false)
const modalTitle = ref('新增患者')
const submitting = ref(false)
const isEdit = ref(false)

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
})

const formState = reactive({
  id: '',
  name: '',
  gender: 'M',
  phone: '',
  idCard: '',
  address: '',
})

const columns = [
  { title: '姓名', dataIndex: 'name', key: 'name' },
  { title: '性别', key: 'gender' },
  { title: '手机号', dataIndex: 'phone', key: 'phone' },
  { title: '身份证号', dataIndex: 'idCard', key: 'idCard' },
  { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt' },
  { title: '操作', key: 'action', width: 200 },
]

async function fetchData() {
  loading.value = true
  try {
    const res: any = await userApi.getPatients({
      page: pagination.current,
      pageSize: pagination.pageSize,
      name: searchName.value || undefined,
    })
    dataSource.value = res?.list || res || []
    pagination.total = res?.total || 0
  } catch {
    dataSource.value = []
  } finally {
    loading.value = false
  }
}

function onTableChange(pag: any) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchData()
}

function showCreateModal() {
  isEdit.value = false
  modalTitle.value = '新增患者'
  Object.assign(formState, { id: '', name: '', gender: 'M', phone: '', idCard: '', address: '' })
  modalOpen.value = true
}

function editRecord(record: any) {
  isEdit.value = true
  modalTitle.value = '编辑患者'
  Object.assign(formState, record)
  modalOpen.value = true
}

function viewDetail(record: any) {
  message.info(`查看患者 ${record.name} 的详细信息`)
}

async function deleteRecord(id: string) {
  try {
    await userApi.deletePatient(id)
    message.success('删除成功')
    fetchData()
  } catch { message.error('删除失败') }
}

async function handleSubmit() {
  submitting.value = true
  try {
    if (isEdit.value) {
      await userApi.updatePatient(formState.id, {
        name: formState.name,
        gender: formState.gender,
        phone: formState.phone,
        idCard: formState.idCard,
        address: formState.address,
      })
    } else {
      await userApi.createPatient({
        name: formState.name,
        gender: formState.gender,
        phone: formState.phone,
        idCard: formState.idCard,
        address: formState.address,
      })
    }
    modalOpen.value = false
    message.success(isEdit.value ? '保存成功' : '创建成功')
    fetchData()
  } catch {
    message.error('操作失败')
  } finally {
    submitting.value = false
  }
}

onMounted(fetchData)
</script>
