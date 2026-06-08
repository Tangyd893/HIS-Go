<template>
  <div class="triage-page">
    <!-- 页头说明 -->
    <a-alert
      type="info"
      show-icon
      style="margin-bottom: 16px"
    >
      <template #message>
        <strong>就诊助手</strong> — 描述您的症状，AI 将推荐合适的科室
      </template>
    </a-alert>

    <!-- 聊天区域 -->
    <div class="chat-area" ref="chatArea">
      <!-- 空状态引导 -->
      <a-empty
        v-if="messages.length === 0 && !loading"
        description="请描述您的症状，例如：咳嗽三天、头晕恶心"
        style="margin-top: 60px"
      >
        <template #image>
          <span style="font-size: 48px">🩺</span>
        </template>
      </a-empty>

      <!-- 消息列表 -->
      <div
        v-for="(msg, idx) in messages"
        :key="idx"
        :class="['message', msg.role]"
      >
        <!-- 用户消息 -->
        <div v-if="msg.role === 'user'" class="msg-bubble user">
          {{ msg.content }}
        </div>

        <!-- 助手消息 -->
        <div v-else class="msg-bubble assistant">
          <!-- Markdown 风格渲染 -->
          <div class="advice-text" v-html="renderAdvice(msg.content)"></div>

          <!-- 推荐科室操作按钮 -->
          <div v-if="msg.depts && msg.depts.length" class="dept-actions">
            <a-space wrap style="margin-top: 12px">
              <a-button
                v-for="dept in msg.depts"
                :key="dept.id"
                type="primary"
                size="small"
                @click="goAppointment(dept)"
              >
                预约{{ dept.name }}
              </a-button>
            </a-space>
          </div>

          <!-- 模式标签 -->
          <div class="msg-meta">
            <a-tag v-if="msg.mode === 'llm'" color="purple">AI 建议</a-tag>
            <a-tag v-else color="blue">关键词匹配</a-tag>
            <a-tag v-if="msg.urgency === 'high'" color="red">紧急</a-tag>
            <a-tag v-else-if="msg.urgency === 'medium'" color="orange">建议就医</a-tag>
          </div>
        </div>
      </div>

      <!-- 加载态 -->
      <div v-if="loading" class="message assistant">
        <div class="msg-bubble assistant loading">
          <span class="dot">●</span>
          <span class="dot">●</span>
          <span class="dot">●</span>
          <span style="margin-left: 8px; color: #999">正在分析症状...</span>
        </div>
      </div>
    </div>

    <!-- 底部输入区 -->
    <div class="input-area">
      <a-textarea
        v-model:value="inputText"
        :auto-size="{ minRows: 1, maxRows: 3 }"
        placeholder="请描述您的症状，例如：咳嗽三天、头晕恶心..."
        @pressEnter="handleSend"
        :disabled="loading"
      />
      <a-button
        type="primary"
        :loading="loading"
        @click="handleSend"
        :disabled="!inputText.trim()"
        style="margin-left: 8px"
      >
        发送
      </a-button>
    </div>

    <!-- 免责声明 -->
    <div class="disclaimer-footer">
      ⚠️ 本建议仅供参考，不能替代专业医疗诊断。如症状严重请及时就医。
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { triageChat, type TriageResponse } from '@/api'

interface ChatMessage {
  role: 'user' | 'assistant'
  content: string
  depts?: { id: string; name: string }[]
  urgency?: string
  mode?: string
}

const router = useRouter()
const chatArea = ref<HTMLElement>()
const inputText = ref('')
const messages = ref<ChatMessage[]>([])
const loading = ref(false)

// 渲染建议文本（简单 Markdown → HTML）
function renderAdvice(text: string): string {
  if (!text) return ''
  let html = text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
    .replace(/\n\n/g, '</p><p>')
    .replace(/\n/g, '<br/>')
  return '<p>' + html + '</p>'
}

// 滚动到底部
async function scrollToBottom() {
  await nextTick()
  if (chatArea.value) {
    chatArea.value.scrollTop = chatArea.value.scrollHeight
  }
}

// 发送消息
async function handleSend() {
  const text = inputText.value.trim()
  if (!text || loading.value) return

  // 添加用户消息
  messages.value.push({ role: 'user', content: text })
  inputText.value = ''
  loading.value = true
  await scrollToBottom()

  try {
    const res: TriageResponse = await triageChat(text)
    messages.value.push({
      role: 'assistant',
      content: res.advice || '抱歉，暂时无法分析该症状，请尝试更详细地描述。',
      depts: res.depts || [],
      urgency: res.urgency,
      mode: res.mode,
    })
  } catch {
    message.error('服务暂时不可用，请稍后重试')
    messages.value.push({
      role: 'assistant',
      content: '抱歉，就诊助手服务暂时不可用。请检查服务配置或稍后重试。',
      depts: [],
      urgency: 'medium',
      mode: 'keyword',
    })
  } finally {
    loading.value = false
    await scrollToBottom()
  }
}

// 跳转预约页（预选科室）
function goAppointment(dept: { id: string; name: string }) {
  if (dept.id) {
    router.push({ path: '/appointment', query: { deptId: dept.id, deptName: dept.name } })
  } else {
    router.push('/appointment')
  }
}
</script>

<style scoped>
.triage-page {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 140px);
  max-width: 720px;
  margin: 0 auto;
}

.chat-area {
  flex: 1;
  overflow-y: auto;
  padding: 16px 0;
}

.message {
  display: flex;
  margin-bottom: 16px;
}

.message.user {
  justify-content: flex-end;
}

.message.assistant {
  justify-content: flex-start;
}

.msg-bubble {
  max-width: 85%;
  padding: 12px 16px;
  border-radius: 12px;
  line-height: 1.6;
  word-break: break-word;
}

.msg-bubble.user {
  background: #1890ff;
  color: #fff;
  border-bottom-right-radius: 4px;
}

.msg-bubble.assistant {
  background: #f5f5f5;
  border-bottom-left-radius: 4px;
}

.msg-bubble.loading {
  background: #f0f0f0;
  display: flex;
  align-items: center;
}

.msg-bubble.loading .dot {
  animation: blink 1.4s infinite both;
  font-size: 18px;
  color: #bbb;
  margin-right: 4px;
}

.msg-bubble.loading .dot:nth-child(2) { animation-delay: 0.2s; }
.msg-bubble.loading .dot:nth-child(3) { animation-delay: 0.4s; }

@keyframes blink {
  0%, 80%, 100% { opacity: 0.2; }
  40% { opacity: 1; }
}

.advice-text :deep(p) {
  margin: 0 0 8px;
}

.advice-text :deep(p:last-child) {
  margin-bottom: 0;
}

.msg-meta {
  margin-top: 8px;
  display: flex;
  gap: 6px;
}

.input-area {
  display: flex;
  align-items: flex-end;
  padding: 12px 0;
  border-top: 1px solid #f0f0f0;
  background: #fff;
}

.input-area :deep(.ant-input) {
  border-radius: 8px;
}

.disclaimer-footer {
  text-align: center;
  color: #999;
  font-size: 12px;
  padding: 8px 0;
  border-top: 1px solid #f5f5f5;
}
</style>
