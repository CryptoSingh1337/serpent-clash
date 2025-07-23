# 🎮 Serpent Clash — Multiplayer Snake Game 
**Serpent Clash** is a **fast-paced**, **multiplayer snake game** where players compete to outmaneuver and 
outgrow their rivals in a dynamic online arena. Control your serpent with precision, dodge enemies, and strategically 
cut them off as you slither your way to the top of the leaderboard.

The game is built using Vue.js on the frontend and Golang on the backend, utilizing modern **WebSocket communication** and 
an **authoritative server** model to maintain a fair and synchronized game state. With **server reconciliation** and 
**interpolation**, players experience smooth gameplay and accurate movement.

🎥 **Demo** - [https://www.youtube.com/watch?v=qjxnNrzlXaQ](https://www.youtube.com/watch?v=qjxnNrzlXaQ?utm_source=github)

[![wakatime](https://wakatime.com/badge/github/CryptoSingh1337/serpent-clash.svg)](https://wakatime.com/badge/github/CryptoSingh1337/serpent-clash)

### 🔑 Key Features
- ⚔️ **Multiplayer** with low-latency WebSocket communication
- 🛡️ **Authoritative Server Model** to ensure fairness and consistency
- 🔄 **Server Reconciliation** for accurate game state even under lag
- ✨ **Player Movement Interpolation** for smooth rendering of remote players
- 🐍 **Dynamic Snake Rendering** using multiple coordinates and mouse input
- 🍔 **Growth and stamina mechanics** where snakes can consume food to grow longer and use stamina for speed boosts
- 🌐 **Efficient Collision Detection** powered by Quad Tree structures
- 🧩 **Entity Component System (ECS)** Architecture for efficient resource management and maintainability

### 🚀 Future Enhancements
- 🏆 **Leaderboard & UI Enhancements**
<br>Improve the overall player experience with a more interactive UI and detailed leaderboard statistics to showcase top players.

### 💪 Motivation & Journey
The goal of this project is to get hands-on experience with backend game development, focusing on low latency communication,
architecture design, and multiplayer game simulation. During the development, I also learned Golang and its ecosystem,
which has been a great experience. Initially, I was not familiar with the language and how to structure a game server, but I
managed to create a working prototype in a short time. Later, I got to know about the Entity Component System (ECS) architecture,
which I found to be a great fit for this project. It allows for better organization of game entities and their behaviors, 
making the codebase more maintainable and scalable.

While working on collision detection, I learned about the Quad Tree data structure, which significantly improved the 
performance of collision checks in the game. For rendering, I used the HTML 2D canvas, but as the complexity of the game increased,
it became unmanageable. Therefore, I switched to Pixi.js, which provided a more efficient rendering engine and better performance, 
basically I offloaded the rendering logic to Pixi.js, allowing me to focus on the game mechanics and logic.

Overall, this project has been a great learning experience and has helped me understand the complexities of game development
and real-time systems.

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
- `Collision` - Detect and handle collisions:
  - Between players (e.g., head-to-head, head-to-body)
  - With world boundaries
  - With food entities
- `Food despawn` - Remove food entities either when consumed or after a predefined number of ticks.
- `Food spawn` - Spawn food entities at random valid positions within the game world.
- `Movement` - Update player movement based on the last recorded mouse direction. If speed boost is active, decrease both snake length and stamina accordingly.
- `Network` - Broadcast the current game state to all clients. Send pong responses to clients for ping calculation.
- `Player despawn` - Handle player leave events by removing the player entity and all its associated components.
- `Player spawn` - Handle player join events by creating a player entity and initializing its associated components.
- `Quad tree` - Rebuild the quad tree each tick using all food entities and snake segments for optimized spatial queries (e.g., collision detection).

#### 🔗 Entity-Component Relationships
- **Player** -> `Input`, `Network`, `PlayerInfo`, `Snake`
- **Food** -> `Expiry`, `Position`

### 📸 Screenshots

##### Landing page
![Landing page](/assets/landing-page.png)

##### Game menu
![Game menu](/assets/game-menu.png)

##### Hud & Gameplay
![Hud and Gameplay](/assets/hud-plus-gameplay.png)

#### Dashboard

##### Server metrics
![Server metrics](/assets/server-metrics.png)

##### Game metrics
![Game metrics](/assets/game-metrics.png)

##### Quad Tree visualization
![Quad Tree visualization](/assets/quad-tree.png)

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
- [x] Food consumption and snake growth logic
- [x] Compensate speed boosts by reducing snake length
- [ ] Implement leaderboard
- [ ] Chat system using SSE and channels
- [ ] Client UI design improvements

### ⚡ Optimizations
- [ ] Store delta state in each system and send only the delta to the client
- [ ] Switch to binary data formats instead of JSON for faster network communication
- [x] Switch to gorilla websocket
- [x] Explore Pixi.js to improve performance

### 📚 Resources
- [https://www.gabrielgambetta.com/client-server-game-architecture.html](https://www.gabrielgambetta.com/client-server-game-architecture.html)

### 🔌 Dependencies
~~- **Websocket** - [https://github.com/lesismal/nbio](https://github.com/lesismal/nbio)~~
- **Websocket** - [https://github.com/gorilla/websocket](https://github.com/gorilla/websocket)
- **Pixijs** - [https://pixijs.com/](https://pixijs.com/)
- **Echo** - [https://echo.labstack.com/](https://echo.labstack.com/)
