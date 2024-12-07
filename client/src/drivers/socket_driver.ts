import { getServerBaseUrl } from "@/utils/helper.ts"

export class SocketDriver {
  socket: WebSocket

  constructor(
    onOpen: () => void,
    onClose: () => void,
    onError: (err: any) => void,
    onMessage: (data: any) => void
  ) {
    this.socket = this.init(onOpen, onClose, onError, onMessage)
  }

  init(
    onOpen: () => void,
    onClose: () => void,
    onError: (err: any) => void,
    onMessage: (data: any) => void
  ): WebSocket {
    const baseUrl = getServerBaseUrl(true)
    const socket = new WebSocket(`${baseUrl}/ws`)
    socket.onopen = onOpen
    socket.onclose = onClose
    socket.onerror = onError
    socket.onmessage = onMessage
    return socket
  }

  getReadyState(): number {
    return this.socket.readyState
  }

  send(data: string): void {
    this.socket.send(data)
  }

  close(): void {
    this.socket.close()
  }
}
