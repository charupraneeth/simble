<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import AnalyticsDashboard from '../components/AnalyticsDashboard.vue'

const route = useRoute()
const router = useRouter()
const siteId = route.params.id as string
const domain = ref('')

onMounted(async () => {
  const res = await fetch('/api/sites')
  if (res.ok) {
    const sites = await res.json()
    const site = sites.find((s: { id: number }) => String(s.id) === siteId)
    if (site) domain.value = site.domain
  }
})
</script>

<template>
  <div class="min-h-screen text-gray-200">
    <main class="container mx-auto px-6 py-12 max-w-6xl">
      <AnalyticsDashboard :base-url="`/api/sites/${siteId}`" :title="domain || `Site #${siteId}`">
        <template #above-title>
          <button
            class="text-sm text-gray-500 hover:text-gray-300 transition-colors mb-2 flex items-center gap-1"
            @click="router.push('/sites')"
          >
            ← Back to Sites
          </button>
        </template>
      </AnalyticsDashboard>
    </main>
  </div>
</template>
