<script lang="ts" setup>
const layout = useDashLayout();

useHead({
  title: "Boards"
})

const { data, error } = useApi("boards");
</script>

<template>
  <NuxtLayout :name="layout">
    <div class="prose mb-4">
      <h1>Boards</h1>
    </div>
    <div v-if="error">
      {{ error }}
    </div>
    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
      <div
        v-for="board in data?.boards"
        :key="board.uuid"
        class="bg-base-300 shadow rounded-lg p-6"
      >
        <h2 class="text-xl font-semibold mb-2">
          {{ board.name == "" ? "unnamed" : board.name }}
        </h2>
        <p>
          <span class="text-gray-500">Host:</span> {{ board.host }}
          <span class="text-gray-500 ml-1">Port:</span> {{ board.port }}
        </p>
      </div>
    </div>
  </NuxtLayout>
</template>
