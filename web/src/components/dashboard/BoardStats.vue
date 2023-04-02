<script lang="ts" setup>
import { useWebSocket } from '@vueuse/core'
import type { BoardInfo } from '@/types/board'
import { compile, computed, inject, onMounted, ref, unref, watch } from 'vue'
import { useBoardStore } from '@/stores/boards'
import { storeToRefs } from 'pinia'

const selectedBoard = storeToRefs(useBoardStore()).selected
const wsBaseUrl = inject('wsBaseUrl')
const wsUrl = ref('')
const data = ref<BoardInfo>({} as BoardInfo)

watch(selectedBoard, (board) => (wsUrl.value = `${wsBaseUrl}/boards/${board?.uuid}`), {
  immediate: true
})

const { data: wsData, open } = useWebSocket<string>(wsUrl, {
  autoReconnect: true,
  heartbeat: true
})

watch(wsData, (d) => {
  if (d !== null) {
    const json = JSON.parse(d)
    data.value = json
  }
})

const memUsage = computed(() => {
  const free: number = parseInt(data.value['free-memory'])
  const total: number = parseInt(data.value['total-memory'])
  const usage = ((total - free) / total) * 100
  return usage.toFixed(0)
})

const memFree = computed(() => {
  const free: number = parseInt(data.value['free-memory'])
  const total: number = parseInt(data.value['total-memory'])
  const usage = (free / total) * 100
  return usage.toFixed(0)
})
</script>

<template>
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-5">
    <div class="card p-4 shadow-lg">
      <div class="card-body">
        <h2 class="text-lg font-medium">Board Information</h2>
        <div class="flex flex-col">
          <div class="flex justify-between my-2">
            <span>Board Name</span>
            <span>{{ data['board-name'] }}</span>
          </div>
          <div class="flex justify-between my-2">
            <span>Uptime</span>
            <span>{{ data.uptime }}</span>
          </div>
          <div class="flex justify-between my-2">
            <span>Version</span>
            <span>{{ data.version }}</span>
          </div>
        </div>
      </div>
    </div>

    <div class="card p-4 shadow-lg">
      <div class="card-body">
        <h2 class="text-lg font-medium">CPU Usage</h2>
        <div class="flex flex-col">
          <div class="flex justify-between my-2">
            <span>CPU Load</span>
            <span>{{ data['cpu-load'] }}</span>
          </div>
          <div class="flex justify-between my-2">
            <span>CPU Count</span>
            <span>{{ data['cpu-count'] }}</span>
          </div>
          <div class="flex justify-between my-2">
            <span>CPU Frequency</span>
            <span>{{ data['cpu-frequency'] }}</span>
          </div>
        </div>
        <div class="mt-4">
          <progress class="progress progress-primary" :value="data['cpu-load']" max="100" />
        </div>
      </div>
    </div>

    <div class="card p-4 shadow-lg">
      <div class="card-body">
        <h2 class="text-lg font-medium">Memory Usage</h2>
        <div class="flex flex-col">
          <div class="flex justify-between my-2">
            <span>Total Memory</span>
            <span>{{ data['total-memory'] }}</span>
          </div>
          <div class="flex justify-between my-2">
            <span>Free Memory</span>
            <span>{{ data['free-memory'] }}</span>
          </div>
        </div>
        <div class="mt-4 flex justify-center gap-10">
          <div class="flex flex-col items-center">
            <div class="radial-progress text-primary" :style="{ '--value': memUsage }">
              {{ memUsage }}%
            </div>
            <span class="mt-2">Usage</span>
          </div>
          <div class="flex flex-col items-center">
            <div class="radial-progress text-primary" :style="{ '--value': memFree }">
              {{ memFree }}%
            </div>
            <span class="mt-2">Free</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
