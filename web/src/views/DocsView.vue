<script setup lang="ts">
import { ref } from 'vue'
import Copy from '../components/icons/Copy.vue'
import CheckCircle from '../components/icons/CheckCircle.vue'

interface Snippet {
  label: string
  code: string
  copied: boolean
}

const snippets = ref<Snippet[]>([
  {
    label: 'basic-install',
    code: '<script defer data-domain="yourdomain.com" src="https://simble.dev/script.js"><\/script>',
    copied: false,
  },
  {
    label: 'css-class',
    code: '<button class="simble-event-Checkout">Checkout<\/button>',
    copied: false,
  },
  {
    label: 'css-multiword',
    code: '<a href="/pricing" class="simble-event-Viewed_Pricing">Pricing<\/a>',
    copied: false,
  },
  {
    label: 'js-api-setup',
    code: `<script>\n  window.simble = window.simble || function() {\n    (window.simble.q = window.simble.q || []).push(arguments)\n  }\n<\/script>\n<script defer data-domain="yourdomain.com" src="https://simble.dev/script.js"><\/script>`,
    copied: false,
  },
  {
    label: 'js-api-call',
    code: "// Fire a custom event anywhere in your JS\nsimble('Signed Up')\n\n// Works great inside form submit handlers:\ndocument.querySelector('form').addEventListener('submit', () => {\n  simble('Form Submitted')\n})",
    copied: false,
  },
])

async function copy(snippet: Snippet) {
  try {
    await navigator.clipboard.writeText(snippet.code)
    snippet.copied = true
    if (typeof window !== 'undefined' && (window as any).simble) {
      (window as any).simble('Copied_Docs_Snippet_' + snippet.label)
    }
    setTimeout(() => (snippet.copied = false), 2000)
  } catch {}
}

function getSnippet(label: string) {
  return snippets.value.find((s) => s.label === label)!
}
</script>

<template>
  <div class="px-6 py-12 max-w-4xl mx-auto text-gray-300">
    <!-- Header -->
    <div class="mb-12 text-center">
      <span class="inline-block bg-emerald-900/50 text-emerald-400 text-xs font-semibold px-3 py-1 rounded-full uppercase tracking-widest mb-4">Integration Guide</span>
      <h1 class="text-4xl font-bold text-white mb-4">Documentation</h1>
      <p class="text-gray-400 text-lg max-w-xl mx-auto">
        Everything you need to take Simble from zero to tracking in under 2 minutes.
      </p>
    </div>

    <!-- Prerequisite banner -->
    <div class="mb-10 flex items-start gap-4 bg-amber-900/20 border border-amber-500/30 rounded-xl px-5 py-4">
      <span class="text-amber-400 text-lg mt-0.5">⚠</span>
      <p class="text-sm text-amber-300/80 leading-relaxed">
        <strong class="text-amber-300">Before you start</strong> — you need to register your website in Simble first.
        Once you've added your domain, the snippet below will start tracking automatically.
        <router-link to="/onboarding" class="underline underline-offset-2 hover:text-amber-200 transition-colors ml-1">Add your site →</router-link>
      </p>
    </div>

    <!-- Section 1: Basic Install -->
    <section class="mb-10 bg-gray-900/60 p-8 rounded-2xl border border-gray-800">
      <div class="flex items-center gap-3 mb-4">
        <span class="flex-shrink-0 w-7 h-7 rounded-full bg-emerald-500/20 text-emerald-400 text-xs font-bold flex items-center justify-center">1</span>
        <h2 class="text-xl font-semibold text-white">Basic Pageview Tracking</h2>
      </div>
      <p class="mb-5 text-gray-400 leading-relaxed">
        Paste this single line into the <code class="bg-gray-800 px-1.5 py-0.5 rounded text-gray-200 text-sm">&lt;head&gt;</code> of every page you want to track. There's no cookie banner needed — Simble is privacy-first by design.
      </p>

      <div class="relative group">
        <div class="bg-gray-950 p-4 pr-16 rounded-xl border border-gray-800/80 font-mono text-sm text-emerald-300 overflow-x-auto whitespace-pre">{{ getSnippet('basic-install').code }}</div>
        <button @click="copy(getSnippet('basic-install'))" class="absolute top-2 right-2 bg-gray-800/80 hover:bg-gray-700 inline-flex items-center justify-center h-8 px-3 rounded text-sm text-white transition-colors">
          <CheckCircle v-if="getSnippet('basic-install').copied" class="w-4 h-4 text-emerald-400" />
          <Copy v-else class="w-4 h-4" />
          <span class="ml-2 text-xs">{{ getSnippet('basic-install').copied ? 'Copied' : 'Copy' }}</span>
        </button>
      </div>

      <p class="mt-4 text-sm text-gray-500">
        Works out of the box with Vue, React, Next.js, and any SPA — Simble patches <code class="text-gray-400">pushState</code> and <code class="text-gray-400">replaceState</code> to track client-side navigation automatically.
      </p>
    </section>

    <!-- Section 2: CSS Class Events -->
    <section class="mb-10 bg-gray-900/60 p-8 rounded-2xl border border-gray-800">
      <div class="flex items-center gap-3 mb-4">
        <span class="flex-shrink-0 w-7 h-7 rounded-full bg-blue-500/20 text-blue-400 text-xs font-bold flex items-center justify-center">2</span>
        <h2 class="text-xl font-semibold text-white">Custom Events via CSS Classes</h2>
      </div>
      <p class="mb-6 text-gray-400 leading-relaxed">
        No Javascript needed. Add a class starting with <code class="bg-gray-800 px-1.5 py-0.5 rounded text-gray-200 text-sm">simble-event-</code> to any clickable element and we'll track it automatically.
      </p>

      <div class="space-y-6">
        <div>
          <p class="text-sm text-gray-400 mb-2">Simple single-word events:</p>
          <div class="relative group">
            <div class="bg-gray-950 p-4 pr-16 rounded-xl border border-gray-800/80 font-mono text-sm text-blue-300 overflow-x-auto whitespace-pre">{{ getSnippet('css-class').code }}</div>
            <button @click="copy(getSnippet('css-class'))" class="absolute top-2 right-2 bg-gray-800/80 hover:bg-gray-700 inline-flex items-center justify-center h-8 px-3 rounded text-sm text-white transition-colors">
              <CheckCircle v-if="getSnippet('css-class').copied" class="w-4 h-4 text-emerald-400" />
              <Copy v-else class="w-4 h-4" />
              <span class="ml-2 text-xs">{{ getSnippet('css-class').copied ? 'Copied' : 'Copy' }}</span>
            </button>
          </div>
        </div>

        <div>
          <p class="text-sm text-gray-400 mb-2">Multi-word events — use underscores, shown as spaces in the dashboard:</p>
          <div class="relative group">
            <div class="bg-gray-950 p-4 pr-16 rounded-xl border border-gray-800/80 font-mono text-sm text-blue-300 overflow-x-auto whitespace-pre">{{ getSnippet('css-multiword').code }}</div>
            <button @click="copy(getSnippet('css-multiword'))" class="absolute top-2 right-2 bg-gray-800/80 hover:bg-gray-700 inline-flex items-center justify-center h-8 px-3 rounded text-sm text-white transition-colors">
              <CheckCircle v-if="getSnippet('css-multiword').copied" class="w-4 h-4 text-emerald-400" />
              <Copy v-else class="w-4 h-4" />
              <span class="ml-2 text-xs">{{ getSnippet('css-multiword').copied ? 'Copied' : 'Copy' }}</span>
            </button>
          </div>
        </div>
      </div>
    </section>

    <!-- Section 3: JS API -->
    <section class="mb-10 bg-gray-900/60 p-8 rounded-2xl border border-gray-800">
      <div class="flex items-center gap-3 mb-4">
        <span class="flex-shrink-0 w-7 h-7 rounded-full bg-purple-500/20 text-purple-400 text-xs font-bold flex items-center justify-center">3</span>
        <h2 class="text-xl font-semibold text-white">Custom Events via JavaScript API</h2>
      </div>
      <p class="mb-2 text-gray-400 leading-relaxed">
        Need to track events from inside your code — like after a successful payment or a form submission? Use the <code class="bg-gray-800 px-1.5 py-0.5 rounded text-gray-200 text-sm">window.simble()</code> function.
      </p>
      <p class="mb-6 text-sm text-gray-500">
        Add this tiny queue snippet anywhere in your HTML <em>before</em> your Simble script tag. It ensures events fired before the script loads are not lost.
      </p>

      <div class="space-y-4">
        <div>
          <p class="text-sm text-gray-400 mb-2">Step 1 — Add the queue snippet and your script tag:</p>
          <div class="relative group">
            <div class="bg-gray-950 p-4 pr-16 rounded-xl border border-gray-800/80 font-mono text-sm text-purple-300 overflow-x-auto whitespace-pre">{{ getSnippet('js-api-setup').code }}</div>
            <button @click="copy(getSnippet('js-api-setup'))" class="absolute top-2 right-2 bg-gray-800/80 hover:bg-gray-700 inline-flex items-center justify-center h-8 px-3 rounded text-sm text-white transition-colors">
              <CheckCircle v-if="getSnippet('js-api-setup').copied" class="w-4 h-4 text-emerald-400" />
              <Copy v-else class="w-4 h-4" />
              <span class="ml-2 text-xs">{{ getSnippet('js-api-setup').copied ? 'Copied' : 'Copy' }}</span>
            </button>
          </div>
        </div>

        <div>
          <p class="text-sm text-gray-400 mb-2">Step 2 — Call <code class="text-gray-200">simble()</code> anywhere in your code:</p>
          <div class="relative group">
            <div class="bg-gray-950 p-4 pr-16 rounded-xl border border-gray-800/80 font-mono text-sm text-purple-300 overflow-x-auto whitespace-pre">{{ getSnippet('js-api-call').code }}</div>
            <button @click="copy(getSnippet('js-api-call'))" class="absolute top-2 right-2 bg-gray-800/80 hover:bg-gray-700 inline-flex items-center justify-center h-8 px-3 rounded text-sm text-white transition-colors">
              <CheckCircle v-if="getSnippet('js-api-call').copied" class="w-4 h-4 text-emerald-400" />
              <Copy v-else class="w-4 h-4" />
              <span class="ml-2 text-xs">{{ getSnippet('js-api-call').copied ? 'Copied' : 'Copy' }}</span>
            </button>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>
