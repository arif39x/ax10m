# AxiomMP: Beta Tester Math Guide

Welcome to the Sandbox! In this game, there are no predefined 3D models or hardcoded physical buttons. **Everything is built and controlled using simple math.** 

Don't panic! You don't need a PhD in calculus to build here. Think of math simply as a way to describe "distances" and "pushes". 

Here is everything you need to know to start shaping the world.

---

## 1. Shaping Matter (Geometry)

To create an object, you describe its surface by writing a **Distance Equation** (technically called an SDF). The rule is simple: **The surface of your object exists exactly where your equation equals zero.**

### The Basic Shapes
If `x`, `y`, and `z` are your coordinates in the world:

* **Sphere**: Describe the radius around a center point.
  * *Equation:* `x^2 + y^2 + z^2 - 100` (Creates a sphere of radius 10)
* **The Ground (Infinite Plane)**: Set a flat boundary on an axis.
  * *Equation:* `z - 0` (Creates a flat floor at height 0)
* **Infinite Cylinder**:
  * *Equation:* `x^2 + y^2 - 25` (A pillar of radius 5 zooming up the Z-axis)

### Moving Shapes Around
To move a shape, subtract the destination coordinates from `x`, `y`, and `z`:
* *Sphere at X=50, Y=20:* `(x - 50)^2 + (y - 20)^2 + z^2 - 100`

---

## 2. Building Complex Structures (Constructive Solid Geometry)

You don't just build single blobs. You can combine simple equations together using `min()` and `max()` to build anything you can imagine like digital Lego.

* **Glue Together (Union)**: Use `min(A, B)`
  * *Example:* `min(sphere1, sphere2)` joins two spheres into a peanut shape.
* **Carve Out (Subtraction)**: Use `max(A, -B)`
  * *Example:* `max(wall, -doorway_sphere)` carves a perfect spherical hole out of a wall.
* **Keep Only the Overlap (Intersection)**: Use `max(A, B)`
  * *Example:* `max(sphere, box)` creates a shape only where the box and sphere overlap (great for making dice or rounded gems).

---

## 3. Controlling Physics (Antigravity & Forces)

Normal games use rigid collision boxes to stop you from falling. AxiomMP doesn't have collisions. Instead, it uses **Potential Fields**—invisible bubbles of math that push or pull matter.

To make an object hover, stop, or fly away, you spawn an **Emitter**. You define three things for an emitter:

1. **Location** `(x, y, z)`: Where the force originates.
2. **Strength** `(Amplitude)`: How hard it pushes. 
   * *Positive Number = Repel (Forcefield)*
   * *Negative Number = Attract (Gravity Well/Black Hole)*
3. **Radius** `(Sigma)`: How far the force reaches before it fades out.

### Real Gameplay Example: Building a Hoverboard

1. **Draw the Board**: 
   `max(box_equation, rounded_edges_equation)`
2. **Make it Hover**:
   You place an **Emitter** underneath it with a massive positive strength `(Strength: 50000)` but a very short radius `(Radius: 2)`.
   * *Result:* As gravity pulls the board down, the emitter violently pushes it back up the moment it gets close to the ground, causing it to bounce and lock into a perfect mathematical hover.

---

## Quick Reference Cheat Sheet

You can type these directly into the terminal to spawn shapes:

| Object | Code to Type |
| :--- | :--- |
| **Basic Sphere** | `x*x + y*y + z*z - 400` |
| **Moving Floor** | `z - (sin(time) * 10)` |
| **Swiss Cheese** | `max(box_math, -(holes_math))` |
| **Gravity Well** | `Emitter(x=0, y=0, z=50, strength=-10000, radius=50)` |

### How the Engine Reads This
When you hit "Build", Python instantly translates these simple mathematical sentences into compiled GPU shaders and Julia tensor arrays. You just focus on the creativity, and the distributed engine handles the extreme math required to make it real-time.
