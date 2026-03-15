# ax10m — Player Guide

> Welcome to ax10m, a sandbox where **you write reality**. The terrain, platforms, and forces in this universe are not built by artists — they are defined by the math equations you type in real time.

---

## Table of Contents

1. [Understanding the World](#1-understanding-the-world)
2. [The Live Math Editor](#2-the-live-math-editor)
3. [Writing SDF Equations — Shapes](#3-writing-sdf-equations--shapes)
4. [Gaussian Fields — Forces and Gravity](#4-gaussian-potential-fields--forces-and-gravity)
5. [Combining Shapes with Boolean CSG](#5-combining-shapes-with-boolean-csg)
6. [The AST Firewall — What You Can and Cannot Write](#6-the-ast-firewall--what-you-can-and-cannot-write)
7. [Physics Rules](#7-physics-rules)
8. [Multiplayer](#8-multiplayer)
9. [Tips & Tricks](#9-tips--tricks)
10. [Common Errors & How to Fix Them](#10-common-errors--how-to-fix-them)

---

## 1. Understanding the World

In ax10m, the world has **no meshes and no polygons**. Every surface is defined by a **Signed Distance Field (SDF)** — a math function `f(x, y, z)` that returns:

- **Negative** → you are _inside_ the shape
- **Zero** → you are _on the surface_
- **Positive** → you are _outside_ the shape

The renderer (WebGPU raymarching) steps along each ray until it finds a zero crossing, then shades that point. You never see polygons — you see the solution set of equations.

---

## 2. The Live Math Editor

When the game client window opens, you will see an overlay panel with three areas:

```
┌─────────────────────────────────────────┐
│  ax10m  ·  Live Math Editor             │
│                                         │
│  ┌───────────────────────────────────┐  │
│  │ sin(x) + cos(y) - z               │  │  ← multi-line input buffer
│  │                                   │  │
│  └───────────────────────────────────┘  │
│                                         │
│  [  Compile & Inject ]                  │  ← send button
│─────────────────────────────────────────│
│ ─── Firewall Output ───                 │
│  [ ax10m firewall ready ]               │  ← error / status log
│                                         │
└─────────────────────────────────────────┘
```

**Workflow:**

1. Type your SDF equation in the input buffer.
2. Press **⚡ Compile & Inject**.
3. The equation is sent to the Go server → forwarded to the Python AST firewall → compiled to WGSL shader code → injected into the renderer.
4. The surface updates live in your viewport. All other connected players see the change too.

If there is a syntax error or the equation contains disallowed code, the **Firewall Output** panel will show the rejection reason.

---

## 3. Writing SDF Equations — Shapes

All equations must be written in terms of `x`, `y`, and `z` only.

### Sphere

A sphere of radius `r` centred at the origin:

```
x**2 + y**2 + z**2 - r**2
```

Example (radius 10):

```
x**2 + y**2 + z**2 - 100
```

### Flat Plane (floor)

Horizontal plane at height `h`:

```
z - h
```

### Infinite Vertical Cylinder (along Z-axis)

```
x**2 + y**2 - r**2
```

### Torus (donut)

Major radius `R`, tube radius `r`:

```
(x**2 + y**2 + z**2 + R**2 - r**2)**2 - 4 * R**2 * (x**2 + y**2)
```

### Undulating Terrain

A rippled ground plane driven by sine waves:

```
z - sin(x) * cos(y)
```

### Twisted Tower

```
x**2 + (y - sin(z))**2 - 4
```

---

## 4. Gaussian Potential Fields — Forces and Gravity

Gaussian Potential Fields are invisible force emitters placed in the world by the server. They push or pull entities (players, objects) based on their position. The force follows the gradient of a Gaussian bell curve:

```
F = -A · exp(−‖r − r₀‖² / (2σ²))  ·  (r − r₀) / σ²
```

| Parameter       | Meaning                                                            |
| --------------- | ------------------------------------------------------------------ |
| `A` (Amplitude) | Strength of the field. Positive = repulsion, negative = attraction |
| `σ` (Sigma)     | Spread radius of the influence                                     |
| `r₀`            | Centre of the emitter (x, y, z)                                    |

The default world has one emitter at `(0, 0, -20)` with amplitude `50000` and sigma `50`, which acts as an upward repulsion floor — preventing you from falling through the world indefinitely.

---

## 5. Combining Shapes with Boolean CSG

Constructive Solid Geometry (CSG) operations are done with `Min` and `Max`:

| Operation                            | Formula        |
| ------------------------------------ | -------------- |
| **Union** (merge two shapes)         | `Min(f1, f2)`  |
| **Intersection** (keep overlap only) | `Max(f1, f2)`  |
| **Subtraction** (carve B from A)     | `Max(f1, -f2)` |

### Example — Sphere with a hole cut through it

```
Max(x**2 + y**2 + z**2 - 100, -(x**2 + y**2 - 4))
```

This creates a sphere of radius 10 with a cylindrical tunnel of radius 2 removed.

### Example — Multi-platform arena

```
Min(Min(z - sin(x) * cos(y), x**2 + z**2 - 9), y**2 + z**2 - 9)
```

---

## 6. The AST Firewall — What You Can and Cannot Write

Before your equation reaches the physics engine, it is parsed by a strict **Abstract Syntax Tree (AST) firewall**. This prevents Remote Code Execution (RCE). Only the following are allowed:

| Category      | Allowed                               |
| ------------- | ------------------------------------- |
| **Variables** | `x`, `y`, `z`                         |
| **Operators** | `+` `-` `*` `/` `**` (power)          |
| **Functions** | `sin`, `cos`, `Min`, `Max`            |
| **Literals**  | Any number (e.g. `3.14`, `100`, `-5`) |

**Everything else is rejected.** The firewall log will show exactly which node was rejected.

### ✅ Valid

```
sin(x) * cos(y) - z**2 + 10
Min(x**2 + y**2 - 4, z - 1)
Max(x**2 + z**2 - 9, -(y**2 - 1))
```

### ❌ Rejected

```python
__import__('os').system('rm -rf /')   # foreign function → rejected
x + a                                  # unknown variable 'a' → rejected
print(x)                               # disallowed function → rejected
x if x > 0 else -x                    # conditional → rejected
```

---

## 7. Physics Rules

The physics engine runs at **20 Hz** using Semi-Implicit Euler integration:

1. **Velocity update first:**  
   `V_new = V + A · dt`

2. **Position update using new velocity:**  
   `X_new = X + V_new · dt`

This order guarantees energy stability (no runaway jitter).

**Friction / Damping:** A velocity-dependent damping force `F = −kv` bleeds kinetic energy. Objects naturally come to rest rather than oscillating forever.

**Floor boundary:** If `z < -50`, the engine clamps `z = -50` and reverses `Vz` with 50% energy loss (inelastic bounce).

---

## 8. Multiplayer

- All players share the **same authoritative world state** broadcast by the Go server at 20 Hz.
- When you **Compile & Inject** an equation, the compiled WGSL shader is broadcast to **every connected client**. The surface change is instant and global.
- The server registers each player as an entity at spawn position `(0, 0, 100)` with mass `1.0`.

---

## 9. Tips & Tricks

- **Start simple.** Begin with a flat plane (`z - 0`) then gradually compose shapes with `Min`/`Max`.
- **Use `sin` and `cos` for organic surfaces.** `z - sin(x) * cos(y)` creates an infinite rippled terrain.
- **Scale with literals.** Multiply sigma expressions to scale a shape: `x**2/4 + y**2/9 + z**2 - 1` is an ellipsoid.
- **Power-of-two exponents are your friends.** `x**2 + y**4 + z**2 - 25` creates a pinched squircle-like shape.
- **Negate to flip.** Wrapping any SDF `f` in `-f` inverts inside/outside, which is required for CSG subtraction.
- **Chain `Min` for complex worlds:**  
  `Min(Min(f1, f2), f3)` unions three surfaces together.

---

## 10. Common Errors & How to Fix Them

| Firewall output                      | Cause                            | Fix                                                            |
| ------------------------------------ | -------------------------------- | -------------------------------------------------------------- |
| `Unauthorized variable: 'a'`         | Unknown symbol in equation       | Only use `x`, `y`, `z`                                         |
| `Unauthorized function call: 'sqrt'` | Disallowed function              | Use `(expr)**0.5` instead                                      |
| `Forbidden AST node type: IfExp`     | Ternary / conditional expression | Rewrqite with `Min`/`Max`                                      |
| `Syntax error: invalid syntax`       | Typo or Python syntax error      | Check parentheses and operators                                |
| `Compilation error: ...`             | Sympy could not parse expression | Check operator precedence; use explicit `*` for multiplication |
| `[WARN] WebSocket not connected`     | Server not running               | Start `GameServer` first (see README)                          |
