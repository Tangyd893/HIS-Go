<template>
  <span :style="{ color }">{{ text }}</span>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  status: number
  map: Record<number, { label: string; color: string }>
}>()

const text = computed(() => props.map[props.status]?.label ?? `状态${props.status}`)
const color = computed(() => props.map[props.status]?.color ?? '#999')
</script>
