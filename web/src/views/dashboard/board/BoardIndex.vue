<script lang="ts" setup>
import DashBase from '@/layout/DashBase.vue'
import BoardList from '@/components/dashboard/BoardList.vue'
import BoardDetail from '@/components/dashboard/BoardDetail.vue'
import { useRoute } from 'vue-router'
import { ref, watch } from 'vue'
import { useBoardsApi } from '@/composables/api'
import { useBoardStore } from '@/stores/boards'

const boardStore = useBoardStore()
const route = useRoute()
const uuid = ref('')

const updateBoardDetail = () => {
  if (typeof route.params.uuid == 'string') {
    uuid.value = route.params.uuid!
  }
}

const addBoard = () => {
  useBoardsApi({ params: { quick: 1 } })
    .post()
    .then((res) => {
      if (res.status == 201) {
        boardStore.refreshBoards()
      }
    })
    .catch(() => {})
}

watch(
  route,
  () => {
    updateBoardDetail()
  },
  { immediate: true }
)
</script>
<template>
  <DashBase>
    <template #default>
      <button class="btn btn-sm mb-3 normal-case" @click="addBoard">Add</button>
      <BoardList />
    </template>
    <template #modal>
      <BoardDetail v-if="route.params.uuid" :uuid="uuid" />
    </template>
  </DashBase>
</template>
