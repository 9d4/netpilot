<script lang="ts" setup>
import DashBase from '@/layout/DashBase.vue'
import BoardList from '@/components/dashboard/BoardList.vue'
import BoardDetail from '@/components/dashboard/BoardDetail.vue'
import { useRoute } from 'vue-router'
import { ref, watch } from 'vue'

const route = useRoute()
const uuid = ref('')

const updateBoardDetail = () => {
  if (typeof route.params.uuid == 'string') {
    uuid.value = route.params.uuid!
  }
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
      <BoardList />
    </template>
    <template #modal>
      <BoardDetail v-if="route.params.uuid" :uuid="uuid" />
    </template>
  </DashBase>
</template>
