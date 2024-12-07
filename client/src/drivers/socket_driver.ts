import { Constants } from "@/utils/constants.ts"

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
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:"
    const url = window.location.hostname.includes("localhost")
      ? `${protocol}//${window.location.hostname}:${Constants.serverPort}/ws`
      : `${protocol}//${window.location.hostname}/ws`
    const socket = new WebSocket(url)
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
