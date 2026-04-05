<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Line } from 'vue-chartjs'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Filler,
  Tooltip,
  type ChartData,
  type ChartOptions,
} from 'chart.js'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Filler, Tooltip)

const route = useRoute()
const router = useRouter()
const siteId = route.params.id as string

// --- State ---
interface Stats { unique_visitors: number; pageviews: number }
interface TrafficPoint { hour: string; visitors: number }
interface TopPage { path: string; views: number; unique_visitors: number }

const stats = ref<Stats | null>(null)
const traffic = ref<TrafficPoint[]>([])
const pages = ref<TopPage[]>([])
const domain = ref('')
const isLoading = ref(true)

// --- Data fetch ---
onMounted(async () => {
  try {
    // Fetch domain name from sites list
    const sitesRes = await fetch('/api/sites')
    if (sitesRes.ok) {
      const sites = await sitesRes.json()
      const site = sites.find((s: { id: number }) => String(s.id) === siteId)
      if (site) domain.value = site.domain
    }

    const [statsRes, trafficRes, pagesRes] = await Promise.all([
      fetch(`/api/sites/${siteId}/stats`),
      fetch(`/api/sites/${siteId}/traffic`),
      fetch(`/api/sites/${siteId}/pages`),
    ])

    if (statsRes.ok) stats.value = await statsRes.json()
    if (trafficRes.ok) traffic.value = await trafficRes.json()
    if (pagesRes.ok) pages.value = await pagesRes.json()
  } catch (e) {
    console.error('Failed to load dashboard data', e)
  } finally {
    isLoading.value = false
  }
})

// --- Chart.js config ---
const chartData = computed<ChartData<'line'>>(() => ({
  labels: traffic.value.map((p) => {
    const d = new Date(p.hour)
    return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  }),
  datasets: [
    {
      label: 'Visitors',
      data: traffic.value.map((p) => p.visitors),
      borderColor: '#10b981',
      backgroundColor: 'rgba(16, 185, 129, 0.08)',
      borderWidth: 2,
      pointRadius: 0,
      tension: 0.4,
      fill: true,
    },
  ],
}))

const chartOptions: ChartOptions<'line'> = {
  responsive: true,
  maintainAspectRatio: false,
  animation: false,
  plugins: {
    legend: { display: false },
    tooltip: {
      backgroundColor: '#111827',
      titleColor: '#9ca3af',
      bodyColor: '#f9fafb',
      borderColor: '#374151',
      borderWidth: 1,
    },
  },
  scales: {
    x: {
      grid: { color: 'rgba(255,255,255,0.04)' },
      ticks: { color: '#6b7280', font: { size: 11 } },
    },
    y: {
      grid: { color: 'rgba(255,255,255,0.04)' },
      ticks: { color: '#6b7280', font: { size: 11 }, precision: 0 },
      beginAtZero: true,
    },
  },
}

const formatNum = (n: number) =>
  n >= 1000 ? `${(n / 1000).toFixed(1)}k` : String(n)
</script>

<template>
  <div class="min-h-screen text-gray-200">
    <main class="container mx-auto px-6 py-12 max-w-6xl">

      <!-- Header -->
      <div class="flex items-center justify-between mb-8">
        <div>
          <button
            class="text-sm text-gray-500 hover:text-gray-300 transition-colors mb-2 flex items-center gap-1"
            @click="router.push('/sites')"
          >
            ← Back to Sites
          </button>
          <h1 class="text-2xl font-bold text-white">{{ domain || `Site #${siteId}` }}</h1>
          <p class="text-xs text-gray-500 mt-1 uppercase tracking-widest">Last 24 hours</p>
        </div>
      </div>

      <!-- Loading skeleton -->
      <div v-if="isLoading" class="space-y-6">
        <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
          <div v-for="i in 4" :key="i" class="h-24 bg-gray-800/60 rounded-xl animate-pulse" />
        </div>
        <div class="h-64 bg-gray-800/60 rounded-xl animate-pulse" />
      </div>

      <template v-else>
        <!-- Stat Cards -->
        <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
          <div class="p-5 border border-gray-800 rounded-xl bg-gray-900/60">
            <p class="text-xs text-gray-500 uppercase tracking-widest mb-2">Unique Visitors</p>
            <p class="text-3xl font-bold text-white">{{ formatNum(stats?.unique_visitors ?? 0) }}</p>
          </div>
          <div class="p-5 border border-gray-800 rounded-xl bg-gray-900/60">
            <p class="text-xs text-gray-500 uppercase tracking-widest mb-2">Pageviews</p>
            <p class="text-3xl font-bold text-white">{{ formatNum(stats?.pageviews ?? 0) }}</p>
          </div>
          <div class="p-5 border border-gray-800 rounded-xl bg-gray-900/60 col-span-2 lg:col-span-2">
            <p class="text-xs text-gray-500 uppercase tracking-widest mb-2">Views / Visitor</p>
            <p class="text-3xl font-bold text-white">
              {{ stats && stats.unique_visitors > 0
                  ? (stats.pageviews / stats.unique_visitors).toFixed(1)
                  : '—' }}
            </p>
          </div>
        </div>

        <!-- Traffic Chart -->
        <div class="p-6 border border-gray-800 rounded-xl bg-gray-900/60 mb-8">
          <div class="flex items-center justify-between mb-6">
            <div>
              <h2 class="text-sm font-semibold uppercase tracking-widest text-gray-400">Traffic Overview</h2>
              <p class="text-xs text-gray-600 mt-0.5">Hourly unique visitors</p>
            </div>
            <span class="flex items-center gap-2 text-xs text-emerald-400">
              <span class="w-2 h-2 rounded-full bg-emerald-500 inline-block" />
              Current
            </span>
          </div>
          <div class="h-56" v-if="traffic.length > 0">
            <Line :data="chartData" :options="chartOptions" />
          </div>
          <div v-else class="h-56 flex items-center justify-center text-gray-600 text-sm">
            No traffic data yet for this period.
          </div>
        </div>

        <!-- Top Pages -->
        <div class="p-6 border border-gray-800 rounded-xl bg-gray-900/60">
          <h2 class="text-sm font-semibold uppercase tracking-widest text-gray-400 mb-6">Top Pages</h2>
          <div v-if="pages.length > 0">
            <div class="grid grid-cols-3 text-xs text-gray-500 uppercase tracking-widest pb-3 border-b border-gray-800 mb-2">
              <span>Path</span>
              <span class="text-right">Views</span>
              <span class="text-right">Unique</span>
            </div>
            <div
              v-for="page in pages"
              :key="page.path"
              class="grid grid-cols-3 py-3 border-b border-gray-800/50 text-sm hover:bg-gray-800/30 px-1 rounded transition-colors"
            >
              <span class="text-gray-300 font-mono truncate">{{ page.path }}</span>
              <span class="text-right text-white font-medium">{{ formatNum(page.views) }}</span>
              <span class="text-right text-gray-400">{{ formatNum(page.unique_visitors) }}</span>
            </div>
          </div>
          <div v-else class="text-center text-gray-600 text-sm py-8">
            No page data yet for this period.
          </div>
        </div>
      </template>

    </main>
  </div>
</template>
