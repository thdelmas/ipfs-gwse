<template>
  <p class="text-gray-600 mt-2 italic">
    No preview available for this file type ({{ type }})
  </p>
</template>

<script setup lang="ts">
import { defineProps } from 'vue'
defineProps<{
  type: string
}>()
</script>
