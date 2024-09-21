export class Stats {
  fps: number = 0
  ping: number = 0
  mouseCoordinate: { x: number; y: number } = { x: 0, y: 0 }
  headCoordinate: { x: number; y: number } = { x: innerWidth / 2, y: innerHeight / 2 }
  playerId: string = ""
  status: string = "offline"

  // internal
  times: number[] = []
  pingStart: number = 0

  renderStats(c: CanvasRenderingContext2D): void {
    if (!c) {
      throw new Error("Can't find canvas element")
    }
    c.fillStyle = "White";
    c.font = "normal 12px Arial";
    c.fillText(Math.floor(this.fps) + " fps", 5, 15);
    c.fillText(`Coordinates: ${this.headCoordinate.x}, ${this.headCoordinate.y}`, 5, 30)
    c.fillText(`Mouse coordinates: ${this.mouseCoordinate.x}, ${this.mouseCoordinate.y}`, 5, 45)
    c.fillText(`Player id: ${this.playerId}`, 5, 60)
    c.fillText(`Status: ${this.status}`, 5, 75)
    c.fillText(`Ping: ${Math.trunc(this.ping * 100) / 100} ms`, 5, 90)
  }

  calculateFps(): void {
    window.requestAnimationFrame((): void => {
      const now = performance.now();
      while (this.times.length > 0 && this.times[0] <= now - 1000) {
        this.times.shift();
      }
      this.times.push(now);
      this.fps = this.times.length;
      this.calculateFps();
    });
  }

  updatePingStart(): void {
    this.pingStart = Date.now();
  }

  updatePing(): void {
    this.ping = Date.now() - this.pingStart;
  }

  updatePlayerId(playerId: string): void {
    this.playerId = playerId;
  }

  updateStatus(status: string): void {
    this.status = status;
  }

  updateMouseCoordinate(x: number, y: number): void {
    this.mouseCoordinate.x = x;
    this.mouseCoordinate.y = y;
  }

  updateHeadCoordinate(x: number, y: number): void {
    this.headCoordinate.x = x;
    this.headCoordinate.y = y;
  }

  reset(): void {
    this.playerId = "";
    this.status = "offline"
    this.headCoordinate.x = 0;
    this.headCoordinate.y = 0;
  }
}
