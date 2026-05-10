<template>
  <div>
    <a-card title="在线问诊" size="small">
      <template #extra>
        <a-button type="primary" size="small" @click="showCreateModal">新建问诊</a-button>
      </template>

      <div v-for="item in consultations" :key="item.id" class="consult-item" @click="openChat(item.id)">
        <a-avatar size="40" style="background: #1890ff">医</a-avatar>
        <div class="consult-info">
          <div class="consult-doctor">{{ item.doctorId || '在线医生' }}</div>
          <div class="consult-time">{{ item.createdAt }}</div>
        </div>
        <a-tag :color="item.status === 1 ? 'green' : 'blue'">
          {{ item.status === 1 ? '已回复' : '待回复' }}
        </a-tag>
      </div>
      <a-empty v-if="!consultations.length" description="暂无问诊记录" />
    </a-card>

    <a-modal v-model:open="modalOpen" title="新建问诊" @ok="handleCreate" width="90%">
      <a-form layout="vertical">
        <a-form-item label="症状描述">
          <a-textarea v-model:value="form.symptom" :rows="4" placeholder="请描述您的症状..." />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { createConsultation, getConsultations } from '@/api'

const consultations = ref<any[]>([])
const modalOpen = ref(false)
const form = reactive({ symptom: '' })

async function fetchData() {
  try {
    const res: any = await getConsultations({})
    consultations.value = res?.list || []
  } catch { consultations.value = [] }
}

function showCreateModal() { modalOpen.value = true }
function openChat(id: string) { message.info(`进入问诊对话: ${id}`) }

async function handleCreate() {
  try {
    await createConsultation({ symptom: form.symptom })
    message.success('问诊创建成功，医生将尽快回复')
    modalOpen.value = false
    fetchData()
  } catch { }
}

onMounted(fetchData)
</script>

<style scoped>
.consult-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-bottom: 1px solid #f0f0f0;
  cursor: pointer;
}
.consult-item:active { background: #f5f5f5; }
.consult-info { flex: 1; }
.consult-doctor { font-size: 15px; font-weight: 500; }
.consult-time { font-size: 12px; color: #999; margin-top: 2px; }
</style>
