# Rationalist

A math-physics-chemistry compiler (Preety Weird Right...!?!? )..Yes because i am working on it on that time when my mind is mostly weird..

In Rationalist, there are no polygons or pre-built models. Everything—the ground you walk on, the shapes you see, and the forces that move you—is defined by mathematical equations you type in real time.

## Quick Start

1. **Install the basics**: You'll need Go, Rust, and Python 3.11+.
2. **Setup**:
   ```bash
   cd MathCompiler && pip install fastapi uvicorn sympy
   ```
3. **Run**:
   ```bash
   ./start_universe.sh
   ```

## How it works

- **Math (Structure)**: You write Signed Distance Fields (SDFs) to build the world.
- **Physics (Motion)**: The universe uses Gaussian Potential Fields to calculate gravity and movement.
- **Client**: A Rust/WebGPU renderer draws the equations as you type them.
- **Server**: A Go hub keeps everyone in the same reality.
