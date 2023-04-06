<script lang="ts" setup>
import BoardResource from '@/components/dashboard/BoardResource.vue'
import DashBase from '@/layout/DashBase.vue'
import { useBoardStore } from '@/stores/boards'

const boardStore = useBoardStore()
</script>

<template>
  <DashBase>
    <BoardResource v-if="boardStore.selected !== null" />
    <div class="prose" v-else>
      <p class="text-warning">No Board Selected</p>
    </div>
  </DashBase>
</template>
