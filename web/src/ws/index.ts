import { useBoardStore } from '@/stores/boards'
import { useMessageCounter } from '@/stores/msgCounter'

import { handleOnOpen, handleResource } from './handler'
import { ref } from 'vue'

let socket: WebSocket
let wsUrl: string
const state = ref<number>(0)

const socketMessageListener = (event: MessageEvent) => {
  state.value = socket.readyState
  const json = JSON.parse(event.data)

  switch (json.type) {
    case 'resource':
      handleResource(event, json)
      break
  }
}

const socketOpenListener = (event: Event) => {
  state.value = socket.readyState
  console.log('Connected')
  handleOnOpen(event)
}

const socketCloseListener = (event?: CloseEvent) => {
  if (socket) {
    console.error('Disconnected.')
  }
  socket = new WebSocket(wsUrl)
  socket.addEventListener('open', socketOpenListener)
  socket.addEventListener('message', socketMessageListener)
  socket.addEventListener('close', socketCloseListener)
  socket.onclose

  state.value = socket.readyState
}

const init = (url: string) => {
  wsUrl = url
  socketCloseListener()
  setTimeout(() => {
    socket.close()
  }, 200)
}

const wait = async () =>
  new Promise((resolve) => {
    const ready = () => socket.readyState === 1

    const a = setTimeout(() => {
      if (ready()) {
        resolve(null)
        clearTimeout(a)
      }
    }, 500)
  })

const fetchBoardStatus = async () => {
  await wait()
  const counter = useMessageCounter()
  const boardStore = useBoardStore()

  socket.send(
    JSON.stringify({
      id: counter.get,
      type: 'get',
      resource: 'board:status',
      board_id: boardStore.selected?.uuid
    })
  )
}

const subscribeBoardStatus = async (uuid:string) => {
  await wait()
  const counter = useMessageCounter()

  socket.send(
    JSON.stringify({
      id: counter.get,
      type: 'sub',
      resource: 'board:status',
      board_id: uuid
    })
  )
}

const unsubscribeBoardStatus = async (uuid:string) => {
  await wait()
  const counter = useMessageCounter()

  socket.send(
    JSON.stringify({
      id: counter.get,
      type: 'unsub',
      resource: 'board:status',
      board_id: uuid
    })
  )
}


const getSocket = () => socket

export default {
  init,
  state,
  getSocket,
  fetchBoardStatus,
  subscribeBoardStatus,
  unsubscribeBoardStatus
}
