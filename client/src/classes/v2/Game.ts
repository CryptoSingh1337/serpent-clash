import type {Ref, ShallowRef} from "vue"
import {DisplayDriver} from "@/classes/v2/DisplayDriver.ts"
import {Player} from "@/classes/v2/Player.ts"
import {Snake} from "@/classes/v2/Snake.ts"
import type {Coordinate} from "@/utils/types"
import {Constants} from "@/utils/constants.ts"
import {InputManager} from "@/classes/v2/InputManager.ts"
import {NetworkManager} from "@/classes/v2/NetworkManager.ts"
import Stats from "stats.js"
import {StatsManager} from "@/classes/v2/StatsManager.ts"

const debugMode: boolean = import.meta.env.VITE_DEBUG_MODE === "true"

export class Game {
  div: HTMLDivElement|null
  statsContainer: Readonly<ShallowRef<HTMLDivElement | null>> | null
  displayDriver: DisplayDriver
  inputManager: InputManager
  networkManager: NetworkManager|null
  statsManager: StatsManager
  clientStatusRef: Ref
  stats: Stats | null = null
  username: string
  player: Player|null

  constructor(
    div: HTMLDivElement|null,
    statsContainer: Readonly<ShallowRef<HTMLDivElement | null>> | null,
    clientStatusRef: Ref,
    username: string) {
    this.div = div
    this.statsContainer = statsContainer
    this.clientStatusRef = clientStatusRef
    this.displayDriver = new DisplayDriver(this)
    this.inputManager = new InputManager(this)
    this.networkManager = null
    this.statsManager = new StatsManager(this)
    this.player = null
    this.username = username
    if (debugMode && statsContainer && statsContainer.value) {
      this.stats = new Stats()
      this.stats.showPanel(0)
      statsContainer.value.appendChild(this.stats.dom)
    }
  }

  async init(): Promise<void> {
    await this.displayDriver.init()
    const theta = Math.atan2(0, 0)
    const x = 0
    const y = 0
    const segments: Coordinate[] = []
    segments.push({x: x, y: y})
    for (let i = 1; i < 10; i++) {
      segments.push({
        x: x - i*Constants.snakeSegmentDistance*Math.cos(theta),
        y: y - i*Constants.snakeSegmentDistance*Math.sin(theta),
      })
    }
    this.player = new Player(this, "1", new Snake(segments, 0xffffff))
    console.log(this.player.snake.segments)
    this.displayDriver.renderer.addEntity(this.player.sprite)
    this.connect()
  }

  removeEntity(entityId: string, entityType: string): void {
    //TODO: remove player entity
  }

  start(): void {
    if (!this.player) {
      return
    }
    this.displayDriver.camera.follow(this.player)
    this.displayDriver.renderer.app.ticker.add(() => {
      if (this.stats != null) {
        this.stats.begin()
      }
      this.displayDriver.render()
      this.displayDriver.update()
      if (this.stats != null) {
        this.stats.end()
      }
    })
  }

  stop(): void {
    this.player = null
    this.displayDriver.renderer.removeEntity()
    this.displayDriver.stop()
    if (this.networkManager) {
      this.networkManager.close()
    }
  }

  connect(): void {
    if (!this.networkManager ||
      this.networkManager.socketState() == WebSocket.CLOSED) {
      this.networkManager = new NetworkManager(
        this,
        this.clientStatusRef,
        this.username
      )
    }
  }

  disconnect(): void {
    if (
      this.networkManager &&
      this.networkManager.socketState() === WebSocket.OPEN
    ) {
      this.networkManager.close()
    }
  }
}
