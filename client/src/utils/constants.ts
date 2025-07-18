export const Constants = {
  debug: true,
  serverPort: 8080,
  tickRate: 60,
  playerSpeed: 6,
  playerSpeedBoost: 3,
  maxTurnRate: 0.03,
  snakeSegmentDistance: 20,
  snakeSegmentRadius: 25,
  foodRadius: 5,
  pingCooldown: 75,
  worldBoundary: {
    radius: 2850,
    padding: 20,
    minX: -3000,
    maxX: 3000,
    minY: -3000,
    maxY: 3000
  }
}

export enum WsMessageType {
  Ping = "ping",
  Pong = "pong",
  hello = "hello",
  Movement = "movement",
  PlayerState = "player_state",
  FoodState = "food_state"
}
