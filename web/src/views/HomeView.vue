<script setup lang="ts">
import GitHub from '../components/icons/GitHub.vue'
import Zap from '../components/icons/Zap.vue'
import Shield from '../components/icons/Shield.vue'
import Activity from '../components/icons/Activity.vue'
import DashboardPreview from '../assets/dashboard-preview.png'

const handleGithubLogin = () => {
  window.location.href = `/auth/github`
}
</script>

<template>
  <!-- Main Content -->
  <main class="container mx-auto px-6 py-24 flex flex-col items-center text-center">

    <!-- Privacy Badge -->
    <!-- Top Badge -->
    <div
      class="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 text-sm font-medium mb-8">
      <span class="flex h-2 w-2 rounded-full bg-emerald-500 animate-pulse"></span>
      Privacy-friendly analytics, built in Go
    </div>

    <!-- Headline -->
    <h1 class="text-5xl md:text-7xl font-extrabold tracking-tight text-white mb-6 max-w-4xl leading-tight">
      Analytics so simple, <br class="hidden sm:block" />
      <span class="inline-block mt-2 md:mt-4 text-transparent bg-clip-text bg-gradient-to-r from-emerald-400 to-cyan-400 drop-shadow-[0_0_15px_rgba(52,211,153,0.3)]">it's simble.</span>
    </h1>

    <!-- Subheadline -->
    <p class="text-lg md:text-xl text-gray-400 mb-10 max-w-2xl leading-relaxed">
      A radically minimal, ultra-fast alternative to Google Analytics.
      No cookies, no GDPR banners, and a tracking script under 1KB.
    </p>

    <!-- Buttons -->
    <div class="flex flex-col sm:flex-row items-center gap-4 w-full justify-center mb-16">
      <button
        class="inline-flex items-center justify-center gap-2 px-8 py-3 w-full sm:w-auto rounded-lg bg-emerald-500 hover:bg-emerald-400 text-emerald-950 font-bold transition-all shadow-lg shadow-emerald-500/25 group text-base"
        @click="handleGithubLogin">
        <GitHub class="w-5 h-5" />
        Sign in with GitHub
      </button>
      <a href="/demo"
        class="inline-flex items-center justify-center gap-2 px-8 py-3 w-full sm:w-auto rounded-lg bg-gray-800 text-white font-medium transition-all duration-300 text-base border border-gray-700 hover:bg-gray-800/80 hover:border-emerald-500/40 hover:text-emerald-400 hover:shadow-[0_0_20px_rgba(16,185,129,0.1)]">
        <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse inline-block group-hover:bg-emerald-400"></span>
        Live Demo
      </a>
    </div>

    <!-- Hero Screenshot -->
    <div class="relative w-full max-w-5xl mx-auto mb-24">
      <!-- Glow Blob -->
      <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[80%] h-[60%] bg-emerald-500/15 blur-[100px] rounded-full pointer-events-none"></div>

      <!-- Image container -->
      <div class="relative rounded-xl lg:rounded-2xl border border-gray-800/60 bg-gray-900/50 p-2 lg:p-4 shadow-2xl overflow-hidden group">
        <div class="absolute inset-0 bg-gradient-to-b from-emerald-500/5 to-transparent pointer-events-none"></div>
        <img :src="DashboardPreview" alt="Simble Dashboard Preview"
          class="w-full h-auto rounded flex items-center justify-center bg-gray-950 min-h-[400px] text-gray-600 text-sm object-cover relative z-10"
          onerror="this.outerHTML='<div class=\'w-full h-[400px] rounded bg-gray-950 flex items-center justify-center border border-dashed border-gray-800 text-gray-500 relative z-10\'>Dashboard preview</div>'" />
      </div>
    </div>

    <!-- Feature Cards Grid -->
    <div class="grid md:grid-cols-3 gap-8 max-w-5xl mx-auto w-full text-left pb-24 border-b border-gray-800/50">

      <!-- Card 1 -->
      <div
        class="group p-6 rounded-2xl bg-gray-900/40 border border-gray-800 hover:border-gray-700 hover:-translate-y-1 hover:shadow-xl hover:shadow-emerald-900/10 transition-all duration-300 flex flex-col gap-4">
        <div class="w-12 h-12 rounded-xl bg-gray-800 flex items-center justify-center border border-gray-700/50 group-hover:border-yellow-500/30 group-hover:bg-gray-800/80 transition-colors">
          <Zap class="w-6 h-6 text-yellow-400" />
        </div>
        <h3 class="text-xl font-semibold text-white">Insanely Fast</h3>
        <p class="text-gray-400 leading-relaxed text-sm lg:text-base">Written in Go. The tracking script is tiny and the
          ingestion API handles concurrent requests instantly.</p>
      </div>

      <!-- Card 2 -->
      <div
        class="group p-6 rounded-2xl bg-gray-900/40 border border-gray-800 hover:border-gray-700 hover:-translate-y-1 hover:shadow-xl hover:shadow-emerald-900/10 transition-all duration-300 flex flex-col gap-4">
        <div class="w-12 h-12 rounded-xl bg-gray-800 flex items-center justify-center border border-gray-700/50 group-hover:border-emerald-500/30 group-hover:bg-gray-800/80 transition-colors">
          <Shield class="w-6 h-6 text-emerald-400" />
        </div>
        <h3 class="text-xl font-semibold text-white">No Cookies</h3>
        <p class="text-gray-400 leading-relaxed text-sm lg:text-base">We don't use cookies or local storage. Just
          privacy-first hashing for unique visitors. Say goodbye to cookie banners.</p>
      </div>

      <!-- Card 3 -->
      <div
        class="group p-6 rounded-2xl bg-gray-900/40 border border-gray-800 hover:border-gray-700 hover:-translate-y-1 hover:shadow-xl hover:shadow-emerald-900/10 transition-all duration-300 flex flex-col gap-4">
        <div class="w-12 h-12 rounded-xl bg-gray-800 flex items-center justify-center border border-gray-700/50 group-hover:border-cyan-500/30 group-hover:bg-gray-800/80 transition-colors">
          <Activity class="w-6 h-6 text-cyan-400" />
        </div>
        <h3 class="text-xl font-semibold text-white">Actionable Insights</h3>
        <p class="text-gray-400 leading-relaxed text-sm lg:text-base">Top pages, referrers, and locations. Exactly what
          you need to know about your traffic, all on one simple dashboard.</p>
      </div>

    </div>

    <footer class="w-full flex flex-col items-center justify-center gap-4 py-12 text-gray-500 text-sm">
      <div
        class="px-2.5 py-0.5 rounded bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 text-xs font-semibold uppercase tracking-wider">
        Public Beta
      </div>
      <p>Simble Analytics. Open source and self-hostable. More features coming soon!</p>
      <a href="https://github.com/charupraneeth/simble" target="_blank" rel="noopener noreferrer"
        class="inline-flex items-center gap-2 text-gray-500 hover:text-white transition-colors text-xs">
        <GitHub class="w-4 h-4" />
        Star us on GitHub
      </a>
    </footer>
  </main>
</template>
