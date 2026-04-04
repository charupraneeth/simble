<script setup lang="ts">
import { useRoute } from 'vue-router'

const route = useRoute()
// Access route.params.id to get the specific site ID to fetch data
</script>

<template>
  <div class="p-8">
    <div class="mb-6 flex items-center gap-4">
      <router-link to="/sites" class="text-gray-500 hover:text-gray-900 dark:hover:text-white transition">
        &larr; Back to Sites
      </router-link>
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Site Dashboard</h1>
    </div>
    
    <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6 shadow-sm">
      <h2 class="text-lg font-semibold mb-4 text-gray-800 dark:text-gray-200">Overview for Site ID: {{ route.params.id }}</h2>
      <p class="text-gray-600 dark:text-gray-400">
        Analytics graphs and metrics goes here.
      </p>
    </div>
  </div>
</template>
