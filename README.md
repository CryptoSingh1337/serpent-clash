# 🎮 Serpent Clash — Multiplayer Snake Game 
**Serpent Clash** is a **fast-paced**, **multiplayer snake game** where players compete to outmaneuver and 
outgrow their rivals in a dynamic online arena. Control your serpent with precision, dodge enemies, and strategically 
cut them off as you slither your way to the top of the leaderboard.

The game is built using Vue.js on the frontend and Golang on the backend, utilizing modern **WebSocket communication** and 
an **authoritative server** model to maintain a fair and synchronized game state. With **server reconciliation** and 
**interpolation**, players experience smooth gameplay and accurate movement.

🎥 **Demo** - [https://youtu.be/hHJq1ubGmuw](https://youtu.be/hHJq1ubGmuw)

[![wakatime](https://wakatime.com/badge/github/CryptoSingh1337/serpent-clash.svg)](https://wakatime.com/badge/github/CryptoSingh1337/serpent-clash)

### 🔑 Key Features
- ⚔️ **Multiplayer** with low-latency WebSocket communication
- 🛡️ **Authoritative Server Model** to ensure fairness and consistency
- 🔄 **Server Reconciliation** for accurate game state even under lag
- ✨ **Player Movement Interpolation** for smooth rendering of remote players
- 🐍 **Dynamic Snake Rendering** using multiple coordinates and mouse input
- 🌐 **Efficient Collision Detection** powered by Quad Tree structures
- 🧩 **Entity Component System (ECS)** Architecture for efficient resource management and maintainability

### 🚀 Future Enhancements
- 🍎 **Food Generation & Snake Growth**
<br>Introduce dynamic elements where snakes grow by consuming food, adding a new layer of strategy and progression.
- 🏆 **Leaderboard & UI Enhancements**
<br>Improve the overall player experience with a more interactive UI and detailed leaderboard statistics to showcase top players.

### 🧱 Backend ECS Architecture

#### 🧍‍♂️ Entities
- `Player`
- `Food`

#### 🧩 Components
- `Expiry`
- `Input`
- `Network`
- `PlayerInfo`
- `Position`
- `Snake`

#### ⚙️ Systems
- `Collision`
- `Food despawn`
- `Food spawn`
- `Movement`
- `Network`
- `Player spawn`
- `Quad tree`

#### 🔗 Entity-Component Relationships
- **Player** -> `Input`, `Network`, `PlayerInfo`, `Snake`
- **Food** -> `Expiry`, `Position`

### ✅ TODO Tracker
- [x] Serve Vue files from backend
- [x] POC: WebSocket server
- [x] Connect/disconnect player
- [x] Show statistics on client-side
- [x] POC: Render snake based on multiple coordinates (client-side)
- [x] POC: Player movement based on mouse coordinates (client-side)
- [x] Design authoritative server based on ticker
- [x] Server-side player movement
- [x] Adopt class-based design for client-side
- [x] Add support for ping calculation
- [x] ~~Add client-side prediction for smoother movement (client-side)~~ (removed)
- [x] Add server reconciliation for handling lag (client-side)
- [x] Add interpolation for smoother movement (client-side)
- [x] Make world finite with a larger dimension
- [x] Add camera logic
- [x] Add collision detection
- [x] Create debug menu (e.g., teleport feature)
- [x] Adjust player speed on `mousedown` / `mouseup`
- [x] Improve snake spawn logic (consider world radius)
- [x] Refactor into driver classes for better readability
- [x] Detect collisions with world boundaries
- [x] Re-architect backend using ECS
- [x] Optimize collision detection using Quad Tree
- [x] Improve collision detection logic (beyond head-to-head only)
- [x] Create spawn system for snake with world-awareness
- [x] Migrate from HTML 2D canvas to Pixi.js renderer engine
- [x] Create dashboard & API for server metrics
- [x] Implement food spawning functionality with world-awareness
- [x] Randomized food generation
- [ ] Food consumption and snake growth logic
- [ ] Compensate speed boosts by reducing snake length
- [ ] Implement leaderboard
- [ ] Chat system using SSE and channels
- [ ] Client UI design improvements

### ⚡ Optimizations
- [ ] Switch to binary data formats instead of JSON for faster network communication
- [x] Switch to gorilla websocket
- [x] Explore Pixi.js to improve performance

### 📚 Resources
- [https://www.gabrielgambetta.com/client-server-game-architecture.html](https://www.gabrielgambetta.com/client-server-game-architecture.html)

### 🔌 Dependencies
~~- **Websocket** - [https://github.com/lesismal/nbio](https://github.com/lesismal/nbio)~~
- **Websocket** - [https://github.com/gorilla/websocket](https://github.com/gorilla/websocket)
