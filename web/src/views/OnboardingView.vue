<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import Copy from '../components/icons/Copy.vue'
import CheckCircle from '../components/icons/CheckCircle.vue'
import SimbleLogo from '../components/icons/SimbleLogo.vue'

const router = useRouter()

const domain = ref('')
const step = ref(1)
const copied = ref(false)
const isSubmitting = ref(false)
const submitError = ref('')

const origin = window.location.origin

const snippet = computed(() => {
  return `<script defer data-domain="${domain.value || 'yourdomain.com'}" src="${origin}/script.js"><\/script>`
})

const isFormValid = computed(() => {
  return domain.value.trim().length > 3 && domain.value.includes('.')
})

const handleCopy = async () => {
  try {
    await navigator.clipboard.writeText(snippet.value)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  } catch (err) {
    console.error('Failed to copy', err)
  }
}

const submitSite = async () => {
  isSubmitting.value = true
  submitError.value = ''
  try {
    const response = await fetch('/api/sites', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ domain: domain.value.trim() })
    })

    if (response.ok) {
      step.value = 2
    } else if (response.status === 409) {
      submitError.value = 'This domain is already registered. Use a different domain or go to your dashboard.'
    } else {
      submitError.value = 'Something went wrong. Please try again.'
    }
  } catch (error) {
    submitError.value = 'Network error. Check your connection and try again.'
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="min-h-screen text-gray-200 flex flex-col items-center py-24 px-6">
    <!-- Header -->
    <div class="mb-12 flex flex-col items-center">
      <div class="w-8 h-8 rounded-lg bg-emerald-600 flex items-center justify-center mb-5">
        <SimbleLogo class="w-5 h-5 text-white" />
      </div>
      <h1 class="text-3xl font-bold text-white mb-2">Welcome to simble!</h1>
      <p class="text-gray-400">Let's get your first site set up.</p>
    </div>

    <!-- Onboarding Card -->
    <div class="w-full max-w-lg p-8 border border-gray-800 rounded-2xl bg-gray-900/80 backdrop-blur-sm">

      <!-- Step 1: Input Domain -->
      <div v-if="step === 1" class="flex flex-col gap-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
        <div>
          <h2 class="text-xl font-semibold text-white mb-1">Add your website</h2>
          <p class="text-sm text-gray-400">Enter the domain of the site you want to track.</p>
        </div>

        <div class="flex flex-col gap-2">
          <label for="domain" class="text-sm font-medium text-gray-300">Domain name</label>
          <input id="domain" v-model="domain" placeholder="e.g. example.com" autofocus
            :class="['bg-gray-950 border rounded-lg h-12 px-4 text-base text-white focus:ring-2 focus:outline-none transition-all', submitError ? 'border-red-500 focus:ring-red-500' : 'border-gray-700 focus:ring-emerald-500']"
            @keydown.enter="isFormValid && !isSubmitting ? submitSite() : null" />
          <p v-if="submitError" class="text-xs text-red-400 mt-1">{{ submitError }}</p>
        </div>

        <button @click="submitSite" :disabled="!isFormValid || isSubmitting"
          class="mt-4 w-full flex items-center justify-center gap-2 h-12 bg-emerald-500 hover:bg-emerald-400 disabled:opacity-50 disabled:cursor-not-allowed text-emerald-950 font-bold rounded-lg transition-all">
          <template v-if="isSubmitting">
            <svg class="animate-spin h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
              </path>
            </svg>
            Registering...
          </template>
          <template v-else>
            Continue
          </template>
        </button>
      </div>

      <!-- Step 2: Snippet -->
      <div v-else class="flex flex-col gap-6 animate-in fade-in slide-in-from-right-8 duration-500">
        <div>
          <h2 class="text-xl font-semibold text-white mb-1">Install snippet</h2>
          <p class="text-sm text-gray-400">Paste this code in the <code
              class="bg-gray-800 px-1 py-0.5 rounded text-gray-300">&lt;head&gt;</code> of your website.</p>
        </div>

        <div class="relative group">
          <div class="absolute inset-0 bg-emerald-500/20 rounded-lg blur group-hover:bg-emerald-500/30 transition-all">
          </div>
          <pre v-text="snippet"
            class="relative bg-gray-950 border border-gray-700/50 p-4 rounded-lg overflow-x-auto text-sm text-emerald-300 font-mono whitespace-pre-wrap" />
          <button @click="handleCopy"
            class="absolute top-2 right-2 bg-gray-800/80 hover:bg-gray-700 inline-flex items-center justify-center h-8 px-3 rounded text-sm text-white transition-colors">
            <CheckCircle v-if="copied" class="w-4 h-4 text-emerald-400" />
            <Copy v-else class="w-4 h-4" />
            <span class="ml-2">{{ copied ? 'Copied' : 'Copy' }}</span>
          </button>
        </div>

        <button @click="router.push('/sites')"
          class="mt-4 w-full h-12 bg-emerald-500 hover:bg-emerald-400 text-emerald-950 font-bold rounded-lg shadow-emerald-500/20 shadow-lg transition-all">
          Go to Dashboard
        </button>
      </div>
    </div>

    <!-- Stepper indicator -->
    <div class="flex items-center gap-2 mt-8">
      <div :class="['w-2 h-2 rounded-full transition-colors', step === 1 ? 'bg-emerald-500' : 'bg-gray-700']"></div>
      <div :class="['w-2 h-2 rounded-full transition-colors', step === 2 ? 'bg-emerald-500' : 'bg-gray-700']"></div>
    </div>
  </div>
</template>
