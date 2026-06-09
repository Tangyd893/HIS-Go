<template>
  <div>
    <ServicePlaceholder v-if="serviceUnavailable" title="消息通知" />
    <template v-else>
      <a-card title="消息通知">
        <template #extra>
          <a-button type="primary" @click="showSendModal"><PlusOutlined /> 发送通知</a-button>
        </template>

        <a-tabs v-model:activeKey="activeTab">
          <a-tab-pane key="templates" tab="模板管理">
            <a-table :columns="templateColumns" :data-source="templates" row-key="id" size="small" />
          </a-tab-pane>
          <a-tab-pane key="send" tab="发送通知" />
        </a-tabs>
      </a-card>

      <a-modal v-model:open="modalOpen" title="发送通知" @ok="handleSend" width="600px">
        <a-form layout="vertical">
          <a-form-item label="接收人ID"><a-input v-model:value="sendForm.receiverId" /></a-form-item>
          <a-form-item label="标题"><a-input v-model:value="sendForm.title" /></a-form-item>
          <a-form-item label="内容"><a-textarea v-model:value="sendForm.content" :rows="4" /></a-form-item>
          <a-form-item label="渠道">
            <a-select v-model:value="sendForm.channel">
              <a-select-option :value="1">站内信</a-select-option>
              <a-select-option :value="2">短信</a-select-option>
              <a-select-option :value="3">邮件</a-select-option>
            </a-select>
          </a-form-item>
        </a-form>
      </a-modal>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { notificationApi } from '@/api/others'
import ServicePlaceholder from '@/components/ServicePlaceholder.vue'

const serviceUnavailable = ref(false)
const activeTab = ref('templates')
const templates = ref<any[]>([])
const modalOpen = ref(false)
const sendForm = reactive({ receiverId: '', title: '', content: '', channel: 1 })

const templateColumns = [
  { title: '模板名称', dataIndex: 'name' },
  { title: '渠道', dataIndex: 'channel' },
  { title: '创建时间', dataIndex: 'createdAt' },
]

async function loadTemplates() {
  try { templates.value = await notificationApi.getTemplates() } catch { serviceUnavailable.value = true }
}

function showSendModal() { modalOpen.value = true }

async function handleSend() {
  try { await notificationApi.send(sendForm); message.success('发送成功'); modalOpen.value = false } catch { }
}

onMounted(loadTemplates)
</script>
