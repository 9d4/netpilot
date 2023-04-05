<script lang="ts" setup>
import { inject, onMounted } from 'vue'
import ws from './ws'

const wsBaseUrl = inject<string>('wsBaseUrl')

onMounted(() => {
  ws.init(wsBaseUrl!)
})
</script>

<template>
  <div>
    <RouterView />
  </div>
</template>
