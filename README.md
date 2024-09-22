# Serpent Clash
Serpent Clash is a fast-paced, real-time multiplayer snake game where players compete to outmaneuver and outgrow their opponents. Players control serpentine avatars, navigating a dynamic battlefield while dodging opponents and obstacles. The game leverages modern WebSocket communication for smooth, real-time player interaction and incorporates client-side prediction to ensure responsive gameplay even in high-latency situations.

The game is built with a combination of Vue.js for the client-side and Go for the backend, following an authoritative server model to ensure fairness and consistency in gameplay. With real-time movement tracking and competitive elements, Serpent Clash challenges players' reflexes and strategic thinking as they slither to the top of the leaderboard.

Key Features:
- Real-time multiplayer with WebSocket communication
- Smooth snake movement based on client-side prediction
- Dynamic player statistics and leaderboard
- Optimized for low-latency, fast-paced gameplay
- Scalable design with plans for further world expansion and advanced mechanics like collision detection and camera logic

### TODO:
- [x] Serve vue files from backend
- [x] POC - web socket server
- [x] Connect/disconnect player
- [x] Show statistics on client-side
- [x] POC - Render snake based on multiple coordinates (client-side)
- [x] POC - Player movement based on mouse coordinates (client-side)
- [x] Design authoritative server based on ticker
- [x] Server-side player movement
- [x] Adopt class-based design for client-side
- [x] Add support for ping calculation
- [x] Add client-side prediction for smooth and snappy player movement (client-side)
- [ ] Add server reconciliation to fix any server-side lag (client-side)
- [ ] Add interpolation for smooth other player movement (client-side)
- [ ] Make world as finite and with bigger dimension than current one
- [ ] Add camera logic
- [ ] Add collision detection logic
- [ ] Generate food on random coordinates
- [ ] Food consumption logic and snake growth logic
- [ ] Change player speed on right mouse button click
- [ ] Leaderboard
- [ ] Client UI design

### Optimizations
- [ ] Use of binary format instead of json
- [ ] Stop rendering when tab is switched

### Resources
- [https://www.gabrielgambetta.com/client-server-game-architecture.html](https://www.gabrielgambetta.com/client-server-game-architecture.html)

### Dependencies
- **Websocket** - [https://github.com/lesismal/nbio](https://github.com/lesismal/nbio)
