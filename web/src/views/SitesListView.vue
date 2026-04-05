<script setup lang="ts">
import { useRouter } from 'vue-router'
import Globe from '../components/icons/Globe.vue'
import Plus from '../components/icons/Plus.vue'
import ExternalLink from '../components/icons/ExternalLink.vue'
import { onMounted, ref } from 'vue'

const router = useRouter()

interface Site {
  id: number;
  domain: string;
  visitors: number;
}

const sites = ref<Site[]>([])
const isLoading = ref(true)

onMounted(async () => {
  try {
    const response = await fetch('/api/sites')
    if (response.ok) {
      sites.value = await response.json()
    }
  } catch (error) {
    console.error("Failed to load sites:", error)
  } finally {
    isLoading.value = false
  }
})
</script>

<template>
  <!-- Main Content -->
  <main class="container mx-auto px-6 py-12 max-w-5xl">
    <!-- Page Header -->
    <div class="flex items-center justify-between mb-8">
      <div>
        <h1 class="text-3xl font-bold text-white mb-2">Your Sites</h1>
        <p class="text-gray-400">Manage your tracked properties.</p>
      </div>
      <button
        class="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-emerald-500 hover:bg-emerald-400 text-emerald-950 font-bold transition-all shadow-lg shadow-emerald-500/20 text-sm"
        @click="router.push('/onboarding')">
        <Plus class="w-4 h-4" />
        Add site
      </button>
    </div>

    <!-- Sites Grid -->
    <div class="grid sm:grid-cols-2 lg:grid-cols-3 gap-6">

      <!-- Site Card -->
      <div v-for="site in sites" :key="site.id"
        class="p-6 border border-gray-800 bg-gray-900 hover:bg-gray-800/80 hover:border-emerald-500/50 transition-all cursor-pointer group relative overflow-hidden rounded-xl hover:shadow-xl hover:shadow-emerald-500/10"
        @click="router.push(`/sites/${site.id}`)">
        <!-- External link icon on hover -->
        <div class="absolute top-0 right-0 p-4 opacity-0 group-hover:opacity-100 transition-opacity">
          <ExternalLink class="w-5 h-5 text-emerald-400" />
        </div>

        <!-- Globe Icon -->
        <div
          class="w-12 h-12 rounded-full bg-gray-800 border border-gray-700 flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300">
          <Globe class="w-6 h-6 text-emerald-400" />
        </div>

        <h3 class="text-lg font-semibold text-white group-hover:text-emerald-400 transition-colors mb-2">
          {{ site.domain }}
        </h3>

        <!-- Live Visitors Badge -->
        <div class="flex items-center gap-2 mt-4 text-sm font-medium">
          <span class="relative flex h-2.5 w-2.5">
            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
            <span class="relative inline-flex rounded-full h-2.5 w-2.5 bg-emerald-500"></span>
          </span>
          <span class="text-gray-300">{{ site.visitors || 0 }} current visitors</span>
        </div>
      </div>

      <!-- Add New Site Card (Dashed) -->
      <button
        class="rounded-xl border-2 border-dashed border-gray-800 hover:border-emerald-500/50 bg-transparent flex flex-col items-center justify-center gap-4 text-gray-400 hover:text-emerald-400 hover:bg-gray-900/50 transition-all min-h-[220px] focus:outline-none focus:ring-2 focus:ring-emerald-500"
        @click="router.push('/onboarding')">
        <div class="w-12 h-12 rounded-full bg-gray-800 flex items-center justify-center">
          <Plus class="w-6 h-6" />
        </div>
        <span class="font-medium text-lg">Add new site</span>
      </button>

    </div>
  </main>
</template>
