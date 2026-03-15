# ax10m — Deterministic Math-Physics Sandbox

> A zero-polygon, multiplayer sandbox universe built entirely on **Signed Distance Fields** and **Gaussian Potential Fields**.  
> No traditional graphics. No rigid-body physics. No AI. 100% math.

---

## What is ax10m?

ax10m is a distributed, multiplayer game engine where the world has **no polygons, no meshes, and no collisions** in the traditional sense. Everything — terrain, objects, player boundaries — is defined by **mathematical equations** that players write in real time.

| Concept     | How ax10m does it                                              |
| ----------- | -------------------------------------------------------------- |
| Geometry    | Signed Distance Fields (SDFs) expressed as math equations      |
| Physics     | Gaussian Potential Field repulsion + damped kinematics (Julia) |
| Multiplayer | Authoritative Go server broadcasting state at 20 Hz            |
| Rendering   | WebGPU raymarching in a Rust client (no polygon rasterisation) |
| Security    | Python AST firewall — only safe math ever reaches the engine   |

### Service Map

```
GameClient  (Rust / wgpu + egui)
    │  WebSocket  ws://localhost:8080
GameServer  (Go)
    │  HTTP POST  :8081
MathCompiler (Python / FastAPI)       ← AST firewall lives here
    │  HTTP POST  :50051
PhysicsEngine (Julia)                 ← Gaussian gradient oracle
```

---

## Prerequisites

| Tool             | Minimum version         |
| ---------------- | ----------------------- |
| **Go**           | 1.21                    |
| **Rust + Cargo** | 1.77 (stable)           |
| **Julia**        | 1.10                    |
| **Python**       | 3.11                    |
| **pip packages** | `fastapi uvicorn sympy` |

Install Python deps:

```bash
cd MathCompiler
python3 -m venv venv
source venv/bin/activate
pip install fastapi uvicorn sympy
```

Install Julia packages (first run only):

```bash
julia --project=PhysicsEngine -e 'using Pkg; Pkg.instantiate()'
```

---

## Clone & Run

```bash
# 1. Clone the repo
git clone https://github.com/<your-org>/ax10m.git
cd ax10m

# 2. Start the entire universe with one script
chmod +x start_universe.sh
./start_universe.sh
```

`start_universe.sh` boots all four services in order and shuts them down cleanly when the Rust client window is closed.

### Manual start (individual services)

```bash
# Terminal 1 — Math Compiler
cd MathCompiler && source venv/bin/activate
uvicorn main:app --port 8081

# Terminal 2 — Physics Engine
julia --project=PhysicsEngine PhysicsEngine/server.jl

# Terminal 3 — Game Server
cd GameServer && go run main.go

# Terminal 4 — Game Client
cd GameClient && cargo run --release
```

## License

MIT
