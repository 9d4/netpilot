<script lang="ts" setup>
import type { BoardOfBoards } from '@/types/board'
import type { PropType } from 'vue'

const props = defineProps({
  board: { type: Object as PropType<BoardOfBoards>, required: true }
})

defineEmits(['click'])
</script>
<template>
  <div class="card card-body bg-base-300 border-base-300 prose hover:cursor-pointer hover:bg-base-200 hover:animate-pulse" @click="$event => $emit('click', $event, props.board.uuid)">
    <p class="text-lg font-bold m-0">{{ props.board.name }}</p>
    <p class="text-sm text-gray-500 m-0">{{ props.board.host }}:{{ props.board.port }}</p>
  </div>
</template>
