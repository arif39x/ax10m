# AxiomMP: Official User Manual

Welcome to **AxiomMP**, the 4-language distributed math-driven sandbox. There are no levels, no predefined models, and no UI buttons. You are given a terminal and a rendering window. The universe is built entirely from mathematical equations.

This manual covers the technical operation of the game engine, how to start it, and the fundamentals of gameplay.

**New Extracurricular Guides:**
*   See `MINIMUM_REQUIREMENTS.md` for the explicit minimum logic required to spawn into the game without falling infinitely.
*   See `MATH_MASTERING.md` for a comprehensive guide on SDFs, CSG compilation, and defining physics equations.
*   See `WORLD_BUILD_RULES.md` for concrete examples and JSON payloads that build complex multi-object universes.

---

## Part 1: System Requirements & Architecture

AxiomMP abandons traditional monolith engines (like Unity or Unreal) in favor of a 4-part microservice architecture. You must run all four components simultaneously for the universe to calculate physics and render.

*   **GameServer (Go)**: The authoritative state spire. It runs the 20Hz physics tick rate and bounces data between the other languages.
*   **GameClient (Rust)**: The visual conduit. It does zero math. It simply maps the spatial float arrays from the Go server directly to the GPU using WebGPU raymarching.
*   **MathCompiler (Python)**: The translator. It takes the text equations you type and compiles them into safe WGSL shader code that the GPU can understand.
*   **PhysicsEngine (Julia)**: The oracle. It performs the brutal partial differential equations to calculate antigravity gradient fields across thousands of objects.

### Prerequisites
*   **Go** (1.20+)
*   **Rust** (1.70+ with Cargo)
*   **Python** (3.10+ with `pip` and `venv`)
*   **Julia** (1.9+)

---

## Part 2: Starting the Universe

To boot the sandbox, you must open four separate terminal windows and start each microservice. The startup order does not matter, but no visuals will appear until the Rust client and Go server connect successfully.

### 1. Start the Translation Layer (Terminal 1)
```bash
cd MathCompiler
source venv/bin/activate
uvicorn main:app --port 8081
```

### 2. Start the Physics Oracle (Terminal 2)
```bash
julia --project=PhysicsEngine PhysicsEngine/server.jl
```

### 3. Start the State Authority (Terminal 3)
```bash
cd GameServer
go run main.go
```

### 4. Boot the Visual Conduit (Terminal 4)
```bash
cd GameClient
cargo run --release
```
*(Note: Compiling the Rust client for the first time may take a few minutes. Subsequent launches will be instantaneous).*

---

## Part 3: Gameplay & Math Fundamentals

Once the servers are online, the rendering window will open. By default, you will see a sphere dropping towards the screen under standard `g = -9.8` gravity, hovering just as it approaches an invisible antigravity field.

To build new structures and forces, you will submit JSON payloads containing your geometric equations to the GameServer.

### Creating Matter (Signed Distance Fields)
Everything is defined by a distance formula. The surface of a shape exists exactly where the math evaluates to `0`. 

*   **Sphere**: `x*x + y*y + z*z - (radius*radius)`
*   **Infinite Floor**: `z - height`
*   **Cylinder**: `x*x + y*y - (radius*radius)`

You can combine these shapes using CSG (Constructive Solid Geometry) functions:
*   `min(A, B)` merges two shapes together.
*   `max(A, B)` finds the exact intersection between two shapes.
*   `max(A, -B)` carves shape B completely out of shape A.

### Creating Forces (Potential Fields)
There are no physical collision boxes to stop you from falling. Instead, you deploy mathematical Potential Field Emitters. 
An emitter uses a Gaussian field: `Strength * exp(-distance^2 / 2*(Radius)^2)`.

To make a platform, you don't build a solid floor block. You build a visual floor, and then you place an Emitter underneath it with a massive positive `Strength` and a tight `Radius`. When falling objects get too close, the exponential gradient pushes them with infinite acceleration, perfectly locking them into a hover state.

To create a black hole, place an emitter with a negative `Strength` and a massive `Radius`.

---

## Part 4: Troubleshooting

*   **"WebSocket connection failed" spam in Rust Client**: The Go `GameServer` is not running or failed to bind to port `8080`.
*   **Shapes immediately disappear**: The Julia `PhysicsEngine` is not running, so the `GameServer` is falling back to default gravity `(-9.8)` with no ground collision logic.
*   **"Compilation error" or "Unauthorized variables used"**: You submitted a math equation to the `MathCompiler` that contained illegal symbols. You may only use `x`, `y`, `z`, and standard math operations (`+`, `-`, `*`, `/`, `sin`, `cos`).
