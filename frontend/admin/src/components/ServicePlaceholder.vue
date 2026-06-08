<template>
  <div class="service-placeholder">
    <div class="placeholder-content">
      <div class="placeholder-icon">
        <ToolOutlined />
      </div>
      <h3>{{ title }}</h3>
      <p>{{ description }}</p>
      <a-tag color="orange">演示环境暂未启用</a-tag>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ToolOutlined } from '@ant-design/icons-vue'

withDefaults(defineProps<{
  title?: string
  description?: string
}>(), {
  title: '功能开发中',
  description: '该功能模块在演示环境中暂未启用，请联系管理员或等待后续版本更新。',
})
</script>

<style scoped>
.service-placeholder {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
  background: #fafafa;
  border-radius: 8px;
  border: 1px dashed #d9d9d9;
}

.placeholder-content {
  text-align: center;
  color: #8c8c8c;
}

.placeholder-icon {
  font-size: 48px;
  margin-bottom: 16px;
  color: #bfbfbf;
}

.placeholder-content h3 {
  color: #595959;
  margin-bottom: 8px;
}

.placeholder-content p {
  max-width: 320px;
  margin-bottom: 12px;
}
</style>
