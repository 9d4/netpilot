<script lang="ts" setup>
import { useWebSocket } from '@vueuse/core'
import type { BoardInfo } from '@/types/board'
import { compile, computed, inject, onMounted, ref, unref, watch } from 'vue'
import { useBoardStore } from '@/stores/boards'
import { storeToRefs } from 'pinia'

const boardStore = useBoardStore()
const boardStoreRef = storeToRefs(boardStore)

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

const ram = computed(() => {
  // convert to MiB
  const free: number = parseInt(data.value['free-memory']) / 1.049e6
  const total: number = parseInt(data.value['total-memory']) / 1.049e6
  const usage = total - free

  return {
    free: free.toFixed(0),
    total: total.toFixed(0),
    usage: usage.toFixed(0),
    usePer: ((usage / total) * 100).toFixed(0)
  }
})

const hdd = computed(() => {
  const free: number = parseInt(data.value['free-hdd-space']) / 1.049e6
  const total: number = parseInt(data.value['total-hdd-space']) / 1.049e6
  const usage = total - free
  return {
    free: free.toFixed(0),
    total: total.toFixed(0),
    usage: usage.toFixed(0),
    usePer: ((usage / total) * 100).toFixed(0)
  }
})
</script>

<template>
  <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-4 mb-5">
    <div class="card card-bordered card-compact shadow-lg bg-base-300 rounded-sm">
      <div class="card-body">
        <div class="card-title text-lg">Board Information</div>
        <div class="flex flex-col leading-3">
          <div class="flex justify-between my-1">
            <span>Board Name</span>
            <span>{{ data['board-name'] }}</span>
          </div>
          <div class="flex justify-between my-1">
            <span>Uptime</span>
            <span>{{ data.uptime }}</span>
          </div>
          <div class="flex justify-between my-1">
            <span>Version</span>
            <span>{{ data.version }}</span>
          </div>
        </div>
      </div>
    </div>

    <div class="card card-bordered card-compact shadow-lg bg-base-300 rounded-sm">
      <div class="card-body">
        <h2 class="card-title text-lg">CPU Usage</h2>
        <div class="flex items-center">
          <div class="flex-1 flex flex-col leading-3">
            <div class="flex justify-between my-1">
              <span>Load</span>
              <span>{{ data['cpu-load'] }}%</span>
            </div>
            <div class="flex justify-between my-1">
              <span>Count</span>
              <span>{{ data['cpu-count'] }}</span>
            </div>
            <div class="flex justify-between my-1">
              <span>Freq.</span>
              <span>{{ data['cpu-frequency'] }} <span class="text-secondary">MHz</span></span>
            </div>
          </div>
          <div class="ml-5">
            <div class="radial-progress" :style="{ '--value': data['cpu-load'] }">
              {{ data['cpu-load'] }}%
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="card card-bordered card-compact shadow-lg bg-base-300 rounded-sm">
      <div class="card-body">
        <div class="card-title text-lg">Memory Usage</div>
        <div class="flex items-center">
          <div class="flex-1 flex flex-col leading-3">
            <div class="flex justify-between my-1">
              <span>Free</span>
              <span>{{ ram.free }} <span class="text-secondary">MiB</span></span>
            </div>
            <div class="flex justify-between my-1">
              <span>Used</span>
              <span>{{ ram.usage }} <span class="text-secondary">MiB</span></span>
            </div>
            <div class="flex justify-between my-1">
              <span>Total</span>
              <span>{{ ram.total }} <span class="text-secondary">MiB</span></span>
            </div>
          </div>
          <div class="ml-5">
            <div class="radial-progress" :style="{ '--value': ram.usePer }">{{ ram.usePer }}%</div>
          </div>
        </div>
      </div>
    </div>

    <div class="card card-bordered card-compact shadow-lg bg-base-300 rounded-sm">
      <div class="card-body">
        <div class="card-title text-lg">Storage Usage</div>
        <div class="flex items-center">
          <div class="flex-1 flex flex-col leading-3">
            <div class="flex justify-between my-1">
              <span>Free</span>
              <span>{{ hdd.free }} <span class="text-secondary">MiB</span></span>
            </div>
            <div class="flex justify-between my-1">
              <span>Used</span>
              <span>{{ hdd.usage }} <span class="text-secondary">MiB</span></span>
            </div>
            <div class="flex justify-between my-1">
              <span>Total</span>
              <span>{{ hdd.total }} <span class="text-secondary">MiB</span></span>
            </div>
          </div>
          <div class="ml-5">
            <div class="radial-progress" :style="{ '--value': hdd.usePer }">{{ hdd.usePer }}%</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
