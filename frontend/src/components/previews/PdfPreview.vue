<template>
  <iframe :src="url" class="mt-2 w-full h-96 border rounded"></iframe>
</template>

<script setup lang="ts">
import { defineProps } from 'vue'

defineProps<{
  url: string
}>()
</script>
