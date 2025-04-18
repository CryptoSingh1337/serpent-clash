export const Constants = {
  debug: true,
  serverPort: 8080,
  tickRate: 60,
  playerSpeed: 5,
  playerSpeedBoost: 3,
  maxTurnRate: 0.03,
  snakeSegmentDistance: 15,
  snakeSegmentDiameter: 50,
  pingCooldown: 75,
  worldBoundary: {
    radius: 2750,
    minX: -3000,
    maxX: 3000,
    minY: -3000,
    maxY: 3000
  }
}

export enum WsMessageType {
  GameState = "game_state",
  Ping = "ping",
  Pong = "pong",
  hello = "hello",
  Movement = "movement"
}
