<script setup lang="ts">
import { useRouter } from 'vue-router'
import AnalyticsDashboard from '../components/AnalyticsDashboard.vue'

const router = useRouter()
</script>

<template>
  <div class="min-h-screen text-gray-200">
    <!-- Demo Banner -->
    <div class="bg-emerald-600/10 border-b border-emerald-500/20 py-3 px-6 text-center flex flex-col sm:flex-row items-center justify-center gap-3">
      <div class="flex items-center gap-2 text-sm text-emerald-300">
        <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse inline-block"></span>
        <span><strong class="font-semibold">Live demo</strong> — You're watching real analytics for <span class="font-mono">simble.dev</span></span>
      </div>
      <button
        @click="router.push('/')"
        class="text-xs px-4 py-1.5 bg-emerald-500 hover:bg-emerald-400 text-emerald-950 font-bold rounded-full transition-colors"
      >
        Track your site →
      </button>
    </div>

    <main class="container mx-auto px-6 py-12 max-w-6xl">
      <AnalyticsDashboard base-url="/api/demo" title="simble.dev" />

      <!-- Bottom CTA -->
      <div class="mt-12 text-center p-8 border border-emerald-500/20 rounded-2xl bg-emerald-500/5">
        <h2 class="text-xl font-bold text-white mb-2">Like what you see?</h2>
        <p class="text-gray-400 text-sm mb-6">Add simble to your site in 30 seconds. No cookies, no GDPR banners.</p>
        <button
          @click="router.push('/')"
          class="inline-flex items-center gap-2 px-6 py-2.5 bg-emerald-500 hover:bg-emerald-400 text-emerald-950 font-bold rounded-lg transition-colors"
        >
          Start tracking for free →
        </button>
      </div>
    </main>
  </div>
</template>
