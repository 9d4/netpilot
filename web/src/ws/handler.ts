import { useBoardStore } from '@/stores/boards'

export const handleResource = (event: MessageEvent, json: any) => {
  const boardStore = useBoardStore()

  switch(json.resource) {
    case 'board:status':
      boardStore.status = json.body
      break
  }
}

export const handleOnOpen = (event: Event) => {

}

