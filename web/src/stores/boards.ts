import { useApi } from '@/composables/api'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

import type { Board, BoardOfBoards, BoardStatus } from '@/types/board'

const saveSelectedBoard = (uuid: string) => {
  return window.localStorage.setItem('selected-board', uuid)
}

const getSelectedBoard = () => {
  return window.localStorage.getItem('selected-board')
}

export const useBoardStore = defineStore('boards', () => {
  const boards = ref<BoardOfBoards[]>([])
  const selected = ref<Board | null>(null)
  const status = ref<BoardStatus | null>(null)

  const getBoardById = computed(
    (state) => (uuid: string) => state.boards.find((b: any) => b.uuid === uuid)
  )

  function refreshBoards() {
    useApi('boards')
      .get()
      .then(({ data }) => (boards.value = data))
      .catch(() => null)
  }

  function updateSelected(uuid: string) {
    useApi(`boards/${uuid}`)
      .get()
      .then(({ data }) => {
        selected.value = data
        saveSelectedBoard(uuid)
      })
      .catch(({ message, response }) => {
        if (typeof response == 'undefined') {
          alert(message)
          return
        }

        refreshBoards()
      })
  }

  // sync with local saved value
  function syncSelected() {
    const saved = getSelectedBoard()
    if (saved !== null) {
      updateSelected(saved)
    }
  }

  function refreshStatus() {
  }

  refreshBoards()
  syncSelected()
  refreshStatus()

  return { boards, selected, refreshBoards, getBoardById, updateSelected, status, refreshStatus }
})
