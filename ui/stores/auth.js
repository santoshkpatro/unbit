import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', () => {
  const loggedInUser = ref(null)

  const setLoggedInUser = (user) => {
    loggedInUser.value = user
  }

  const isLoggedIn = computed(() => !!loggedInUser.value)

  return { loggedInUser, setLoggedInUser, isLoggedIn }
})
