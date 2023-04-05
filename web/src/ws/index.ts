import { useBoardStore } from '@/stores/boards'
import { useMessageCounter } from '@/stores/msgCounter'

import { handleResource }from './handler'

let socket: WebSocket
let wsUrl: string

const socketMessageListener = (event: MessageEvent) => {
  const json = JSON.parse(event.data)

  switch (json.type) {
    case 'resource':
      handleResource(event, json)
      break
  }
}

const socketOpenListener = () => {
  console.log('Connected')
  socket.send(JSON.stringify({}))
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
}

const init = (url: string) => {
  wsUrl = url
  socketCloseListener()
  setTimeout(() => {
    socket.close()
  }, 1000)
}

const fetchBoardStatus = () => {
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

export default {
  init,
  fetchBoardStatus
}
