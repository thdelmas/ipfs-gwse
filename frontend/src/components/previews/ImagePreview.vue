<script setup lang="ts">
import { defineProps } from 'vue'
defineProps<{ src: string; alt: string }>()
</script>

<template>
  <img :src="src" :alt="alt" class="mt-2 rounded border max-w-full" />
</template>
