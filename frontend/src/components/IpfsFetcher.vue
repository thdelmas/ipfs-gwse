<script setup lang="ts">
import { ref } from 'vue'

const cid = ref('')
const isLoading = ref(false)
const error = ref('')
const fileUrl = ref('')
const fileName = ref('')

const fetchFile = async () => {
  if (!cid.value) return

  isLoading.value = true
  error.value = ''
  fileUrl.value = ''
  fileName.value = ''

  try {
    const res = await fetch(`http://localhost:8000/${cid.value}`)
    if (!res.ok) throw new Error(`Error: ${res.status}`)

    const blob = await res.blob()
    fileName.value = res.headers.get('X-File-Name') || cid.value
    fileUrl.value = URL.createObjectURL(blob)
  } catch (err: any) {
    error.value = err.message
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="max-w-xl mx-auto mt-10 p-6 bg-white shadow rounded-xl">
    <h1 class="text-xl font-semibold mb-4">🔍 Search IPFS by CID</h1>

    <input
      v-model="cid"
      @keyup.enter="fetchFile"
      type="text"
      placeholder="Enter CID (e.g. Qm...TTNT)"
      class="w-full px-4 py-2 border rounded-md mb-4"
    />

    <button
      @click="fetchFile"
      :disabled="isLoading"
      class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
    >
      {{ isLoading ? 'Fetching...' : 'Fetch File' }}
    </button>

    <div v-if="error" class="mt-4 text-red-600">{{ error }}</div>

    <div v-if="fileUrl" class="mt-6">
      <p class="mb-2 font-medium">File: {{ fileName }}</p>
      <a
        :href="fileUrl"
        :download="fileName"
        class="text-blue-600 underline"
        target="_blank"
        rel="noopener"
      >
        ⬇️ Download File
      </a>

      <div v-if="fileUrl.includes('image')" class="mt-4">
        <img :src="fileUrl" :alt="fileName" class="rounded border" />
      </div>
    </div>
  </div>
</template>
