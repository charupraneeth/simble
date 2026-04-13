<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuth } from '../composables/useAuth'
import SimbleLogo from './icons/SimbleLogo.vue'
import LogOut from './icons/LogOut.vue'
import GitHub from './icons/GitHub.vue'

const router = useRouter()
const { user, isLoggedIn, logout, loading } = useAuth()

const handleGithubLogin = () => {
  window.location.href = '/auth/github'
}

const handleSignOut = async () => {
  await logout()
}

// Derive initials for the avatar fallback
const initials = (username: string) => username.slice(0, 2).toUpperCase()
</script>

<template>
  <nav class="border-b border-gray-800/50 bg-[#02040a]/80 backdrop-blur-md sticky top-0 z-10">
    <div class="container mx-auto px-6 h-16 flex items-center justify-between">

      <!-- Logo -->
      <div class="flex items-center gap-2 cursor-pointer" @click="router.push('/')">
        <div class="w-8 h-8 rounded-lg bg-emerald-600 flex items-center justify-center">
          <SimbleLogo class="w-5 h-5 text-white" />
        </div>
        <span class="text-xl font-bold tracking-tight text-white">simble</span>
      </div>

      <!-- Right Side -->
      <div class="flex items-center gap-4">

        <!-- Loading State -->
        <template v-if="loading">
           <svg class="animate-spin h-5 w-5 text-emerald-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
        </template>

        <!-- Logged In: Avatar + Username + Sign Out -->
        <template v-else-if="isLoggedIn && user">
          <!-- Avatar -->
          <div class="flex items-center gap-3">
            <div
              class="w-8 h-8 rounded-full bg-emerald-600 flex items-center justify-center text-white text-xs font-bold ring-2 ring-emerald-500/30">
              {{ initials(user.username) }}
            </div>
            <span class="text-sm font-medium text-gray-300 hidden sm:block">{{ user.username }}</span>
          </div>

          <!-- Sign Out -->
          <button class="flex items-center gap-2 text-sm text-gray-400 hover:text-white transition-colors group"
            @click="handleSignOut">
            <LogOut class="w-4 h-4 group-hover:-translate-x-0.5 transition-transform" />
            <span class="hidden sm:block">Sign out</span>
          </button>
        </template>

        <!-- Logged Out: GitHub Login Button -->
        <template v-else>
          <a
            href="https://github.com/charupraneeth/simble"
            target="_blank"
            rel="noopener noreferrer"
            class="text-gray-400 hover:text-white transition-colors"
            title="Star on GitHub"
          >
            <GitHub class="w-5 h-5" />
          </a>
          <button
            class="inline-flex items-center justify-center gap-2 px-4 py-2 rounded-lg bg-emerald-500 hover:bg-emerald-400 text-emerald-950 font-bold text-sm transition-all"
            @click="handleGithubLogin">
            Start for free
          </button>
        </template>

      </div>
    </div>
  </nav>
</template>