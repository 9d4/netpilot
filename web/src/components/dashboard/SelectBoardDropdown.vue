<script lang="ts" setup>
import { useBoardStore } from '@/stores/boards'
import { storeToRefs } from 'pinia'
import { ref, watch } from 'vue'

const boardStore = useBoardStore()
const boardStoreRef = storeToRefs(boardStore)

const boards = boardStoreRef.boards
const selected = boardStoreRef.selected
const selectedLabel = ref('')
const opened = ref(false)

const toggle = () => {
  opened.value = !opened.value
}

const updateLabel = () => {
  if (selected.value === null) {
    selectedLabel.value = 'Select Board'
    return
  }

  selectedLabel.value = selected.value.name
}

watch(
  selected,
  () => {
    updateLabel()
  },
  { immediate: true }
)

const onItemClick = (uuid: string) => {
  toggle()
  boardStore.updateSelected(uuid)
}
</script>

<template>
  <div class="dropdown dropdown-bottom" :class="{ 'dropdown-open': opened }">
    <label class="btn btn-ghost rounded-btn normal-case" @click="toggle">
      <span>{{ selectedLabel }}</span>
      <span>
        <svg
          width="12px"
          height="12px"
          class="ml-1 h-3 w-3 fill-current opacity-60"
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 2048 2048"
        >
          <path d="M1799 349l242 241-1017 1017L7 590l242-241 775 775 775-775z" />
        </svg>
      </span>
    </label>
    <ul class="menu dropdown-content p-2 shadow bg-base-100 rounded-box w-52 mt-4">
      <li v-for="board in boards" :key="board.uuid">
        <span @click="onItemClick(board.uuid)">{{ board.name }}</span>
      </li>
    </ul>
  </div>
</template>
