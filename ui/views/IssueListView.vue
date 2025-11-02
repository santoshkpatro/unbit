<script setup>
import { onMounted, ref } from 'vue'
import { projectListAPI } from '@/api/projects'
import { issueListAPI } from '@/api/issues'

const projects = ref([])
const issues = ref([])

const loadProjects = async () => {
  projects.value = await projectListAPI()
}

const loadIssues = async () => {
  issues.value = await issueListAPI()
}

onMounted(async () => {
  await loadProjects()
  await loadIssues()
})
</script>

<template>
  <main>
    <a-typography-title :level="3">Issues</a-typography-title>
  </main>
</template>
