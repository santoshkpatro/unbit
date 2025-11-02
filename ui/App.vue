<script setup>
import { RouterLink, RouterView } from 'vue-router'
import { Notivue, Notification } from 'notivue'

const theme = {
  token: {
    colorPrimary: '#51bc8f',
    borderRadius: 4,
    fontSize: 13,
  },
}
</script>

<template>
  <Notivue v-slot="item">
    <Notification :item="item" />
  </Notivue>

  <a-config-provider :theme="theme">
    <router-view></router-view>
  </a-config-provider>
</template>

<style scoped></style>
