<script setup lang="ts">
const handleGetMe = async ()=>{
  const response = await fetch("/api/me")
  const json = await response.json()
  console.log(json)
}

</script>

<template>
  <button @click="handleGetMe">get me</button>
</template>
