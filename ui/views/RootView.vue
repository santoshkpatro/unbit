<script setup>
import { useRouter } from 'vue-router'
import { onMounted } from 'vue'

const router = useRouter()

onMounted(() => {
  router.push({ name: 'issue-list' })
})
</script>

<template>
  <div>Root</div>
</template>
