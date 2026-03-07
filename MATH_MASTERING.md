# AxiomMP: Math Mastering Guide

Welcome to the absolute foundation of the AxiomMP sandbox. In this engine, you don't place graphical meshes or draw polygons. You declare mathematics. Everything you see and interact with is the result of continuous functions evaluated across 3D space.

## 1. Signed Distance Fields (SDFs)
An SDF is a mathematical function that takes a point in space `(x, y, z)` and returns the shortest distance to the surface of an object.
*   If the result is **positive**, you are outside the object.
*   If the result is **negative**, you are inside the object.
*   If the result is **exactly 0**, you are strictly on the surface.

### Primitive Geometric Shapes
*   **Infinite Floor (Plane):** `z - height`
    *   *Example:* `z + 50` creates an infinite floor at `z = -50`.
*   **Sphere:** `x*x + y*y + z*z - (radius*radius)`
    *   *Example:* `x*x + y*y + z*z - 100` creates a sphere of radius 10 at the origin `(0,0,0)`.
*   **Cylinder (Infinite):** `x*x + y*y - (radius*radius)`
    *   *Example:* `x*x + y*y - 400` creates a vertical cylindrical pillar with radius 20.

## 2. Constructive Solid Geometry (CSG)
To build complex worlds out of simple primitives, you combine SDFs using minimum and maximum functions. Because AxiomMP's internal `MathCompiler` uses Python's `sympy` library, you must use capitalization for operators: `Min` and `Max`.

*   **Union (`Min(A, B)`)**: Combines two shapes into one. The surface is whichever shape is closer to the evaluation point.
*   **Intersection (`Max(A, B)`)**: Keeps only the overlapping region where both shape A and shape B exist.
*   **Subtraction (`Max(A, -B)`)**: Evaluates the negative of shape B, completely carving it out of shape A.

## 3. Potential Fields (The Physics Engine)
Visuals alone don't prevent you from falling through the floor. The Go `GameServer` and Julia `PhysicsEngine` compute gravitational and repulsive forces using Gaussian Potential Fields, completely bypassing traditional rigid-body collision detection.

**The Underlying Equation for Field Emitters is:**
`Strength * exp(-(distance^2) / (2 * Radius^2))`

Instead of hitting a physical "floor box" and calculating a bouncy normal force, you place a high-strength Potential Field Emitter slightly beneath your visual floor. As you fall closer to the emitter's center, the exponential resistive force increases infinitely, resulting in a perfectly stable hover lock.

### Archetypes of Physics Fields:
*   **Platform Hover Lock:** Extremely high positive strength (`50000.0`), tight radius (`50.0`). Places you in a stable float just above the visual asset.
*   **Repulsion Field:** Medium positive strength (`5000.0`), large radius. Pushes objects away gently as they approach.
*   **Black Hole / Gravity Well:** Massive negative strength (`-100000.0`), massive radius. Pulls objects forcefully towards its center `(x,y,z)`.
