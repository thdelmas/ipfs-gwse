<template>
  <pre class="mt-2 bg-gray-100 p-4 rounded overflow-auto max-h-96">{{ content }}</pre>
</template>

<script setup lang="ts">
import { defineProps } from 'vue'
defineProps<{
  content: string
}>()
</script>
