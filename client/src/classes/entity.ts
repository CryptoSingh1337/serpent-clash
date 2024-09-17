export class Player {
  clientId: string
  sessionId: string
  positions: Position[]
  color: string

  constructor(clientId:string, sessionId: string, color: string, positions: Position[]) {
    this.clientId = clientId
    this.sessionId = sessionId
    this.color = color
    this.positions = positions
  }

  draw(c: CanvasRenderingContext2D) {
    c.beginPath()
    c.arc(this.positions[0].x, this.positions[0].y, 2, 0, Math.PI * 2, false)
    c.fillStyle = this.color
    c.fill()
  }
}
