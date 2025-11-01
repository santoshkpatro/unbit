<script setup>
import { RouterLink, RouterView } from 'vue-router'

const theme = {
  token: {
    colorPrimary: '#1890ff',
    borderRadius: 4,
  },
}
</script>

<template>
  <a-config-provider :theme="theme">
    <router-view></router-view>
  </a-config-provider>
</template>

<style scoped></style>
