<template>
  <a-modal
    v-model:open="visible"
    :title="title"
    :confirm-loading="loading"
    @ok="handleOk"
    @cancel="handleCancel"
    :width="width"
    :destroyOnClose="true"
  >
    <slot />
  </a-modal>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

const props = withDefaults(defineProps<{
  title?: string
  width?: string | number
  modelValue?: boolean
}>(), {
  title: '表单',
  width: 600,
  modelValue: false,
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'submit': []
  'cancel': []
}>()

const visible = ref(props.modelValue)
const loading = ref(false)

watch(() => props.modelValue, (val) => { visible.value = val })
watch(visible, (val) => { emit('update:modelValue', val) })

async function handleOk() {
  loading.value = true
  try {
    emit('submit')
  } finally {
    loading.value = false
  }
}

function handleCancel() {
  emit('cancel')
  visible.value = false
}

defineExpose({ setLoading: (val: boolean) => { loading.value = val }, close: () => { visible.value = false } })
</script>
