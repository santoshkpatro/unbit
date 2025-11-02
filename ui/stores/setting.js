import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export const useSettingStore = defineStore('setting', () => {
  const setting = ref(null)

  const setSetting = (settingData) => {
    setting.value = settingData
  }

  const isInstalled = computed(() => !!setting.value)

  return { setting, isInstalled, setSetting }
})
