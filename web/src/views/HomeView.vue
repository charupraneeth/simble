<script setup lang="ts">
import GithubIcon from '../assets/icons/github.svg'
import Zap from '../components/icons/Zap.vue'
import Shield from '../components/icons/Shield.vue'
import Activity from '../components/icons/Activity.vue'

const handleGithubLogin = () => {
  window.location.href = `/auth/github`
}
</script>

<template>
  <!-- Main Content -->
  <main class="container mx-auto px-6 py-24 flex flex-col items-center text-center">

    <!-- Privacy Badge -->
    <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-gray-800/50 border border-gray-700/50 text-emerald-400 text-sm font-medium mb-8">
      <span class="flex h-2 w-2 rounded-full bg-emerald-500 animate-pulse"></span>
      Privacy-friendly analytics
    </div>

    <!-- Headline -->
    <h1 class="text-5xl md:text-7xl font-extrabold tracking-tight text-white mb-6 max-w-4xl">
      Analytics so simple, <br />
      <span class="text-transparent bg-clip-text bg-gradient-to-r from-emerald-400 to-cyan-400">it's simble.</span>
    </h1>

    <!-- Subheadline -->
    <p class="text-lg md:text-xl text-gray-400 mb-10 max-w-2xl leading-relaxed">
      Stop digging through complex dashboards. Get the privacy-focused insights you need without the bloat.
      Keep your tracking simple, sweet, and to the point.
    </p>

    <!-- Buttons -->
    <div class="flex flex-col sm:flex-row items-center gap-4 w-full justify-center">
      <button
        class="inline-flex items-center justify-center gap-2 px-8 py-3 w-full sm:w-auto rounded-lg bg-emerald-500 hover:bg-emerald-400 text-emerald-950 font-bold transition-all shadow-lg shadow-emerald-500/25 group text-base"
        @click="handleGithubLogin">
        <img :src="GithubIcon" class="w-5 h-5" />
        Signin with GitHub
      </button>
    </div>

    <!-- Feature Cards Grid -->
    <div class="mt-24 grid md:grid-cols-3 gap-8 max-w-5xl mx-auto w-full text-left">

      <!-- Card 1 -->
      <div class="p-6 rounded-2xl bg-gray-900/50 border border-gray-800 flex flex-col gap-4">
        <div class="w-12 h-12 rounded-xl bg-gray-800 flex items-center justify-center border border-gray-700/50">
          <Zap class="w-6 h-6 text-yellow-400" />
        </div>
        <h3 class="text-xl font-semibold text-white">Lightning Fast</h3>
        <p class="text-gray-400 leading-relaxed">Lightweight script &lt; 1KB. Won't slow down your shiny new frontend.</p>
      </div>

      <!-- Card 2 -->
      <div class="p-6 rounded-2xl bg-gray-900/50 border border-gray-800 flex flex-col gap-4">
        <div class="w-12 h-12 rounded-xl bg-gray-800 flex items-center justify-center border border-gray-700/50">
          <Shield class="w-6 h-6 text-emerald-400" />
        </div>
        <h3 class="text-xl font-semibold text-white">Privacy First</h3>
        <p class="text-gray-400 leading-relaxed">No cookies, no GDPR banners needed. We don't track personal data.</p>
      </div>

      <!-- Card 3 -->
      <div class="p-6 rounded-2xl bg-gray-900/50 border border-gray-800 flex flex-col gap-4">
        <div class="w-12 h-12 rounded-xl bg-gray-800 flex items-center justify-center border border-gray-700/50">
          <Activity class="w-6 h-6 text-emerald-400" />
        </div>
        <h3 class="text-xl font-semibold text-white">Beautiful Dashboards</h3>
        <p class="text-gray-400 leading-relaxed">Everything you need to know on a single page. No training required.</p>
      </div>

    </div>
  </main>
</template>
