<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
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

const props = defineProps<{
  baseUrl: string    // e.g. '/api/sites/1' or '/api/demo'
  title: string      // displayed in the header h1
}>()

// --- State ---
interface Stats { unique_visitors: number; pageviews: number }
interface TrafficPoint { hour: string; visitors: number }
interface TopPage { path: string; views: number; unique_visitors: number }
interface TopCountry { country_code: string; views: number; unique_visitors: number }
interface TopReferrer { referrer: string; views: number; unique_visitors: number }
interface TopEvent { name: string; events: number; unique_visitors: number }

const stats = ref<Stats | null>(null)
const traffic = ref<TrafficPoint[]>([])
const pages = ref<TopPage[]>([])
const countries = ref<TopCountry[]>([])
const referrers = ref<TopReferrer[]>([])
const events = ref<TopEvent[]>([])
const isLoading = ref(true)
const range = ref<'24h' | '7d' | '30d'>('7d')

const RANGES: { label: string; value: '24h' | '7d' | '30d' }[] = [
  { label: '24h', value: '24h' },
  { label: '7d', value: '7d' },
  { label: '30d', value: '30d' },
]

// --- Data fetch ---
async function fetchData(r: string) {
  isLoading.value = true
  try {
    const [statsRes, trafficRes, pagesRes, countriesRes, referrersRes, eventsRes] = await Promise.all([
      fetch(`${props.baseUrl}/stats?range=${r}`),
      fetch(`${props.baseUrl}/traffic?range=${r}`),
      fetch(`${props.baseUrl}/pages?range=${r}`),
      fetch(`${props.baseUrl}/countries?range=${r}`),
      fetch(`${props.baseUrl}/referrers?range=${r}`),
      fetch(`${props.baseUrl}/events?range=${r}`),
      new Promise(resolve => setTimeout(resolve, 600)), // min skeleton duration
    ])
    if (statsRes.ok) stats.value = await statsRes.json()
    if (trafficRes.ok) traffic.value = await trafficRes.json()
    if (pagesRes.ok) pages.value = await pagesRes.json()
    if (countriesRes.ok) countries.value = await countriesRes.json()
    if (referrersRes.ok) referrers.value = await referrersRes.json()
    if (eventsRes.ok) events.value = await eventsRes.json()
  } catch (e) {
    console.error('Failed to load dashboard data', e)
  } finally {
    isLoading.value = false
  }
}

function setRange(r: '24h' | '7d' | '30d') {
  range.value = r
  fetchData(r)
  if (typeof window !== 'undefined' && (window as any).simble) {
    (window as any).simble('Changed_Time_Range_to_' + r)
  }
}

onMounted(() => fetchData(range.value))

// --- Chart.js config ---
const chartData = computed<ChartData<'line'>>(() => ({
  labels: traffic.value.map((p) => {
    const d = new Date(p.hour)
    if (range.value === '24h') {
      return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
    }
    return d.toLocaleDateString([], { month: 'short', day: 'numeric' })
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
  <!-- Header row -->
  <div class="flex items-start justify-between mb-8 gap-4">
    <div>
      <!-- Slot for back button or badge, rendered above the title -->
      <slot name="above-title" />
      <h1 class="text-2xl font-bold text-white">{{ title }}</h1>
      <p class="text-xs text-gray-500 mt-1 uppercase tracking-widest">
        Last {{ range === '24h' ? '24 hours' : range === '7d' ? '7 days' : '30 days' }}
      </p>
    </div>

    <!-- Range Toggle -->
    <div class="flex items-center gap-1 bg-gray-900 border border-gray-800 rounded-lg p-1 mt-6">
      <button
        v-for="r in RANGES"
        :key="r.value"
        @click="setRange(r.value)"
        :class="[
          'px-3 py-1.5 rounded-md text-xs font-semibold transition-all',
          range === r.value
            ? 'bg-emerald-600 text-white shadow'
            : 'text-gray-500 hover:text-gray-300'
        ]"
      >
        {{ r.label }}
      </button>
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

    <!-- Data Tables Row -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">

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
            <span class="text-gray-300 font-mono truncate mr-2">{{ page.path }}</span>
            <span class="text-right text-white font-medium">{{ formatNum(page.views) }}</span>
            <span class="text-right text-gray-400">{{ formatNum(page.unique_visitors) }}</span>
          </div>
        </div>
        <div v-else class="text-center text-gray-600 text-sm py-8">No page data yet for this period.</div>
      </div>

      <!-- Top Locations -->
      <div class="p-6 border border-gray-800 rounded-xl bg-gray-900/60">
        <h2 class="text-sm font-semibold uppercase tracking-widest text-gray-400 mb-6">Top Locations</h2>
        <div v-if="countries.length > 0">
          <div class="grid grid-cols-6 text-xs text-gray-500 uppercase tracking-widest pb-3 border-b border-gray-800 mb-2">
            <span class="col-span-4">Country</span>
            <span class="text-right">Views</span>
            <span class="text-right">Unique</span>
          </div>
          <div
            v-for="country in countries"
            :key="country.country_code"
            class="relative py-3 border-b border-gray-800/50 text-sm hover:bg-gray-800/30 px-1 rounded transition-colors"
          >
            <div
              class="absolute top-0 left-0 h-full bg-emerald-500/10 rounded pointer-events-none"
              :style="{ width: `${(country.unique_visitors / countries[0].unique_visitors) * 100}%` }"
            ></div>
            <div class="relative grid grid-cols-6 items-center">
              <div class="col-span-4 flex items-center gap-3">
                <span class="text-lg leading-none" v-if="country.country_code !== 'Unknown'">
                  {{ String.fromCodePoint(...country.country_code.toUpperCase().split('').map(c => 127397 + c.charCodeAt(0))) }}
                </span>
                <span class="text-gray-300">{{ country.country_code }}</span>
              </div>
              <span class="text-right text-white font-medium">{{ formatNum(country.views) }}</span>
              <span class="text-right text-gray-400">{{ formatNum(country.unique_visitors) }}</span>
            </div>
          </div>
        </div>
        <div v-else class="text-center text-gray-600 text-sm py-8">No location data yet for this period.</div>
      </div>

      <!-- Top Referrers -->
      <div class="p-6 border border-gray-800 rounded-xl bg-gray-900/60">
        <h2 class="text-sm font-semibold uppercase tracking-widest text-gray-400 mb-6">Top Referrers</h2>
        <div v-if="referrers.length > 0">
          <div class="grid grid-cols-3 text-xs text-gray-500 uppercase tracking-widest pb-3 border-b border-gray-800 mb-2">
            <span>Source</span>
            <span class="text-right">Views</span>
            <span class="text-right">Unique</span>
          </div>
          <div
            v-for="ref in referrers"
            :key="ref.referrer"
            class="grid grid-cols-3 py-3 border-b border-gray-800/50 text-sm hover:bg-gray-800/30 px-1 rounded transition-colors"
          >
            <span class="text-gray-300 truncate mr-2">{{ ref.referrer }}</span>
            <span class="text-right text-white font-medium">{{ formatNum(ref.views) }}</span>
            <span class="text-right text-gray-400">{{ formatNum(ref.unique_visitors) }}</span>
          </div>
        </div>
        <div v-else class="text-center text-gray-600 text-sm py-8">No referrer data yet for this period.</div>
      </div>

      <!-- Top Events -->
      <div class="p-6 border border-gray-800 rounded-xl bg-gray-900/60">
        <h2 class="text-sm font-semibold uppercase tracking-widest text-gray-400 mb-6">Top Events</h2>
        <div v-if="events.length > 0">
          <div class="grid grid-cols-3 text-xs text-gray-500 uppercase tracking-widest pb-3 border-b border-gray-800 mb-2">
            <span>Name</span>
            <span class="text-right">Events</span>
            <span class="text-right">Unique</span>
          </div>
          <div
            v-for="ev in events"
            :key="ev.name"
            class="grid grid-cols-3 py-3 border-b border-gray-800/50 text-sm hover:bg-gray-800/30 px-1 rounded transition-colors"
          >
            <span class="text-gray-300 font-mono truncate mr-2">{{ ev.name }}</span>
            <span class="text-right text-white font-medium">{{ formatNum(ev.events) }}</span>
            <span class="text-right text-gray-400">{{ formatNum(ev.unique_visitors) }}</span>
          </div>
        </div>
        <div v-else class="text-center text-gray-600 text-sm py-8">No custom events tracked yet.</div>
      </div>

    </div>
  </template>
</template>
