import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const useMessageCounter = defineStore('messageCounter', () => {
  const counter = ref<number>(0)

  const get = computed(() => ++counter.value)

  return { counter, get }
})
