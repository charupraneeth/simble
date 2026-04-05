<script setup lang="tsx">
import GithubIcon from '../assets/icons/github.svg'

const handleGithubLogin = () => {
    window.location.href = `/auth/github`
}
</script>

<template>
    <nav class="container mx-auto px-6 py-6 flex items-center justify-between border-b border-gray-800/50">
        <div class="flex items-center gap-2">
            <div class="w-8 h-8 rounded-lg bg-emerald-600 flex items-center justify-center">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"
                    stroke-linejoin="round" class="w-5 h-5 text-white">
                    <circle cx="4" cy="17" r="2"></circle>
                    <circle cx="20" cy="7" r="2"></circle>
                    <path d="M6 17 C 12 17 12 7 18 7"></path>
                </svg>
            </div>
            <span class="text-xl font-bold tracking-tight text-white">simble</span>
        </div>
        <div class="flex items-center gap-4">
            <!-- Start for free Button -->
            <button
                class="inline-flex items-center justify-center gap-2 px-4 py-2 rounded-lg bg-emerald-500 hover:bg-emerald-400 text-emerald-950 font-bold transition-all"
                @click="handleGithubLogin">
                <img :src="GithubIcon" class="w-5 h-5" />
                Start for free
            </button>
        </div>
    </nav>
</template>