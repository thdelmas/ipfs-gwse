<script setup lang="ts">
import { ref } from 'vue'
import ImagePreview from './previews/ImagePreview.vue'
import PdfPreview from './previews/PdfPreview.vue'
import TextPreview from './previews/TextPreview.vue'
import UnsupportedPreview from './previews/UnsupportedPreview.vue'

const cid = ref('')
const isLoading = ref(false)
const error = ref('')
const fileUrl = ref('')
const fileName = ref('')
const contentType = ref('')
const textPreview = ref('')

const fetchFile = async () => {
  if (!cid.value) return

  isLoading.value = true
  error.value = ''
  fileUrl.value = ''
  fileName.value = ''
  contentType.value = ''
  textPreview.value = ''

  try {
    const res = await fetch(`http://localhost:8000/${cid.value}`)
    if (!res.ok) throw new Error(`Error: ${res.status}`)

    const blob = await res.blob()
    contentType.value = res.headers.get('Content-Type') || ''
    fileName.value = res.headers.get('X-File-Name') || cid.value
    fileUrl.value = URL.createObjectURL(blob)

    if (contentType.value.startsWith('text/')) {
      const text = await blob.text()
      textPreview.value = text
    }
  } catch (err: any) {
    error.value = err.message
  } finally {
    isLoading.value = false
  }
}
</script>


<template>
<div class="w-full  mt-10 p-6 sm:p-8 bg-white shadow rounded-2xl">
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

      <div class="mt-6">
        <p class="font-semibold">Preview:</p>


        <img
          v-if="contentType.startsWith('image/')"
          :src="fileUrl"
          :alt="fileName"
          class="mt-2 rounded border"
        />

        <iframe
          v-else-if="contentType === 'application/pdf'"
          :src="fileUrl"
          class="mt-2 w-full h-96 border rounded"
        ></iframe>

        <pre
          v-else-if="contentType.startsWith('text/')"
          class="mt-2 bg-gray-100 p-4 rounded overflow-auto max-h-96"
        >{{ textPreview }}</pre>

        <p v-else class="text-gray-600 mt-2 italic">
          No preview available for this file type ({{ contentType }})
        </p>
      </div>
    </div>
  </div>
</template>
