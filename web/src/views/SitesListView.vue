<script setup lang="ts">
// Sites List Dashboard
</script>

<template>
  <div class="p-8">
    <h1 class="text-3xl font-bold mb-6 text-gray-900 dark:text-white">Your Sites</h1>
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      
      <!-- Example static card -->
      <router-link to="/sites/123" class="block p-6 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg shadow hover:bg-gray-50 dark:hover:bg-gray-700 transition cursor-pointer">
        <h5 class="mb-2 text-xl font-bold tracking-tight text-gray-900 dark:text-white">charuvlogs.com</h5>
        <p class="font-normal text-gray-700 dark:text-gray-400">1,204 Visitors (7d)</p>
      </router-link>
      
      <!-- Add Site card -->
      <div class="flex flex-col items-center justify-center p-6 border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg hover:border-blue-500 dark:hover:border-blue-500 transition cursor-pointer text-gray-500 hover:text-blue-500">
        <span class="text-lg font-semibold">+ Add Site</span>
      </div>
      
    </div>
  </div>
</template>
