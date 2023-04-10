<script lang="ts" setup>
import { useBoardStore } from '@/stores/boards'
import BoardListItem from './BoardListItem.vue'
import BoardDetail from './BoardDetail.vue'
import { useRoute, useRouter } from 'vue-router'

const boardStore = useBoardStore()
const router = useRouter()
const route = useRoute()

const onItemClick = (event: Event, uuid: string) => {
  router.push({ params: { uuid: uuid } })
}
</script>
<template>
  <div class="grid gap-2 sm:grid-cols-2">
    <BoardListItem v-for="b in boardStore.boards" :board="b" :key="b.uuid" @click="onItemClick" />
  </div>
</template>
