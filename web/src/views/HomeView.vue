<script setup lang="ts">
import { useRouter } from 'vue-router'
import GithubIcon from '../assets/icons/github.svg'

const router = useRouter()

const handleGithubLogin = () => {
  window.location.href = `/auth/github`
}
</script>

<template>


  <!-- Main Content -->
  <main class="container mx-auto px-6 py-24 flex flex-col items-center text-center">

    <!-- Privacy Badge -->
    <div
      class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-gray-800/50 border border-gray-700/50 text-emerald-400 text-sm font-medium mb-8">
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
      <!-- Primary Action -->
      <button
        class="inline-flex items-center justify-center gap-2 px-8 py-3 w-full sm:w-auto rounded-lg bg-emerald-500 hover:bg-emerald-400 text-emerald-950 font-bold transition-all shadow-lg shadow-emerald-500/25 group text-base"
        @click="handleGithubLogin">
        <img :src="GithubIcon" />
        Signin with GitHub
      </button>
    </div>

    <!-- Feature Cards Grid -->
    <div class="mt-24 grid md:grid-cols-3 gap-8 max-w-5xl mx-auto w-full text-left">

      <!-- Card 1 -->
      <div class="p-6 rounded-2xl bg-gray-900/50 border border-gray-800 flex flex-col gap-4">
        <div class="w-12 h-12 rounded-xl bg-gray-800 flex items-center justify-center border border-gray-700/50">
          <!-- Zap Icon -->
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
            stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
            class="text-yellow-400">
            <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"></polygon>
          </svg>
        </div>
        <h3 class="text-xl font-semibold text-white">Lightning Fast</h3>
        <p class="text-gray-400 leading-relaxed">Lightweight script &lt; 1KB. Won't slow down your shiny new frontend.
        </p>
      </div>

      <!-- Card 2 -->
      <div class="p-6 rounded-2xl bg-gray-900/50 border border-gray-800 flex flex-col gap-4">
        <div class="w-12 h-12 rounded-xl bg-gray-800 flex items-center justify-center border border-gray-700/50">
          <!-- Shield Icon -->
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
            stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
            class="text-emerald-400">
            <path
              d="M20 13c0 5-3.5 7.5-7.66 8.95a1 1 0 0 1-.67-.01C7.5 20.5 4 18 4 13V6a1 1 0 0 1 1-1c2-1 4-2 7-2c3 0 5 1 7 2a1 1 0 0 1 1 1v7z">
            </path>
          </svg>
        </div>
        <h3 class="text-xl font-semibold text-white">Privacy First</h3>
        <p class="text-gray-400 leading-relaxed">No cookies, no GDPR banners needed. We don't track personal data.</p>
      </div>

      <!-- Card 3 -->
      <div class="p-6 rounded-2xl bg-gray-900/50 border border-gray-800 flex flex-col gap-4">
        <div class="w-12 h-12 rounded-xl bg-gray-800 flex items-center justify-center border border-gray-700/50">
          <!-- Activity Icon -->
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
            stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
            class="text-emerald-400">
            <polyline points="22 12 18 12 15 21 9 3 6 12 2 12"></polyline>
          </svg>
        </div>
        <h3 class="text-xl font-semibold text-white">Beautiful Dashboards</h3>
        <p class="text-gray-400 leading-relaxed">Everything you need to know on a single page. No training required.
        </p>
      </div>

    </div>
  </main>
</template>
