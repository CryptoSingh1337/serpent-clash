import type { Ref, ShallowRef } from "vue"
import Stats from "stats.js"
import { DisplayDriver } from "@/classes/DisplayDriver.ts"
import { InputManager } from "@/classes/InputManager.ts"
import { NetworkManager } from "@/classes/NetworkManager.ts"
import { StatsManager } from "@/classes/StatsManager.ts"
import { Player } from "@/classes/Player.ts"
import type { Players } from "@/utils/types"

const debugMode: boolean = import.meta.env.VITE_DEBUG_MODE === "true"

export class Game {
  div: HTMLDivElement | null
  statsContainer: Readonly<ShallowRef<HTMLDivElement | null>> | null
  displayDriver: DisplayDriver
  inputManager: InputManager
  networkManager: NetworkManager | null = null
  statsManager: StatsManager
  clientStatusRef: Ref
  stats: Stats | null = null
  pingPanel: Stats.Panel = new Stats.Panel("Ping", "#ff8", "#221")
  username: string
  player: Player | null = null
  playerEntities: Players = {}

  constructor(
    div: HTMLDivElement | null,
    statsContainer: Readonly<ShallowRef<HTMLDivElement | null>> | null,
    clientStatusRef: Ref,
    username: string
  ) {
    this.div = div
    this.statsContainer = statsContainer
    this.clientStatusRef = clientStatusRef
    this.displayDriver = new DisplayDriver(this)
    this.inputManager = new InputManager(this)
    this.statsManager = new StatsManager(this)
    this.username = username
    if (debugMode && statsContainer && statsContainer.value) {
      this.stats = new Stats()
      this.stats.addPanel(this.pingPanel)
      this.stats.showPanel(0)
      statsContainer.value.appendChild(this.stats.dom)
    }
  }

  async init(): Promise<void> {
    await this.displayDriver.init()
    this.connect()
  }

  start(): void {
    if (this.player) {
      this.displayDriver.camera.follow(this.player)
    }
    this.displayDriver.renderer.app.ticker.add(() => {
      if (!this.displayDriver.camera.target && this.player) {
        this.displayDriver.camera.follow(this.player)
      }
      if (this.stats != null) {
        this.stats.begin()
      }
      this.displayDriver.update()
      this.statsManager.update()
      this.statsManager.reducePingCooldown()
      if (this.statsManager.getPingCooldown() <= 0 && this.networkManager) {
        this.networkManager.sendPingPayload()
      }
      this.pingPanel.update(this.statsManager._stats.ping, 300)
      if (this.stats != null) {
        this.stats.end()
      }
    })
  }

  stop(): void {
    this.player = null
    this.displayDriver.camera.target = null
    this.displayDriver.renderer.removeEntity()
    this.displayDriver.stop()
    if (this.networkManager) {
      this.networkManager.close()
    }
  }

  connect(): void {
    if (
      !this.networkManager ||
      this.networkManager.socketState() === WebSocket.CLOSED
    ) {
      this.networkManager = new NetworkManager(this, this.username)
      if (this.networkManager.socketState() === WebSocket.OPEN) {
        this.statsManager.updateStatus("Online")
      }
    }
  }

  disconnect(): void {
    if (
      this.networkManager &&
      this.networkManager.socketState() === WebSocket.OPEN
    ) {
      this.networkManager.close()
      if (this.networkManager.socketState() === WebSocket.CLOSED) {
        this.statsManager.reset()
      }
    }
  }
}
