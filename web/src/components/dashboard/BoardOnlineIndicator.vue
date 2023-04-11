<script lang="ts" setup>
import { useBoardStore } from '@/stores/boards'
import type { Board, BoardStatus } from '@/types/board'
import ws from '@/ws'
import { storeToRefs } from 'pinia'
import { ref, watch } from 'vue'

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

const updateSubscription = (old?:Board) => {
  if (old !== undefined && old !== null) {
    ws.unsubscribeBoardStatus(old.uuid)
  }

  ws.subscribeBoardStatus(boardStoreRef.selected.value?.uuid!)
}

const wsState = ref(ws.state)
watch(wsState, () => updateSubscription(), { immediate: 1 == 1 })
watch(
  boardStoreRef.selected,
  (_, old) => {
    updateSubscription(old!)
  },
  { immediate: true }
)
</script>

<template>
  <span
    class="rounded-full w-2 h-2 mx-2 hover:cursor-pointer"
    :class="{'bg-warning': online.status == 2, 'bg-green-400': online.status == 1, 'bg-red-400': online.status == 0 }"
    :title="timestamp?.toLocaleDateString() + ' ' + timestamp?.toLocaleTimeString()"
  ></span>
</template>
