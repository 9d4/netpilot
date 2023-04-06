<script lang="ts" setup>
import { useBoardStore } from '@/stores/boards'
import type { BoardStatus } from '@/types/board'
import ws from '@/ws'
import { storeToRefs } from 'pinia'
import { onUnmounted, ref, watch } from 'vue'

const boardStore = useBoardStore()
const boardStoreRef = storeToRefs(boardStore)
const online = ref<BoardStatus>({} as BoardStatus)
const timestamp = ref<Date>()

watch(
  boardStoreRef.status,
  () => {
    if (boardStoreRef.status.value != null) {
      online.value = boardStoreRef.status.value
      timestamp.value = new Date(online.value.timestamp)
    }
  },
  { immediate: 1 == 1 }
)

const interval = setInterval(() => {
  ws.fetchBoardStatus()
}, 1000)

onUnmounted(() => {
  clearInterval(interval)
})
</script>

<template>
  <span
    class="rounded-full w-2 h-2 mx-2 hover:cursor-pointer"
    :class="{ 'bg-green-400': online.status, 'bg-red-400': !online.status }"
    :title="timestamp?.toLocaleDateString() + ' ' + timestamp?.toLocaleTimeString()"
  ></span>
</template>
