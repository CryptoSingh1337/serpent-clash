import { getServerBaseUrl } from "@/utils/helper.ts"

export class SocketDriver {
  socket: WebSocket

  constructor(
    username: string,
    onOpen: () => void,
    onClose: () => void,
    onError: (err: any) => void,
    onMessage: (data: any) => void
  ) {
    this.socket = this.init(username, onOpen, onClose, onError, onMessage)
  }

  init(
    username: string,
    onOpen: () => void,
    onClose: () => void,
    onError: (err: any) => void,
    onMessage: (data: any) => void
  ): WebSocket {
    const baseUrl = getServerBaseUrl(true)
    const socket = new WebSocket(`${baseUrl}/ws?username=${username}`)
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
