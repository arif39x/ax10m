# AxiomMP: World Building Rules

You can build an entire universe by compiling a single, massive visual equation and registering appropriate Potential Field emitters alongside it. Below is a step-by-step guide to generating a full world state with a ground layer, an arena wall, and a celestial sun.

## 1. The Visual World Payload
To render the world, you must submit a combined equation using `sympy` operators over the WebSocket to the `GameServer`.

```json
{
  "equation": "Min(z + 50, Min(x*x + y*y + (z-500)*(z-500) - 10000, Max(x*x + y*y - 400, -(x*x + y*y - 100))))"
}
```

### Breakdown of the Equation above:
1. **The Ground:** `z + 50` (A flat plane lying at `Z = -50`).
2. **The Sun:** `x*x + y*y + (z-500)*(z-500) - 10000` (A massive spherical body high in the sky at `z=500`).
3. **The Arena Ring:** `Max(x*x + y*y - 400, -(x*x + y*y - 100))` (A hollow cylinder on the Z-axis, creating interior and exterior structural walls).
4. **Combining the Universe:** They are all joined together sequentially using the CSG operator `Min()`.

*Send this JSON string to the GameServer WebSocket instance. The Python compiler will lock and distribute the translated WGSL to the Rust rendering pipeline.*

---

## 2. The Physics Configuration Rules
To make the mathematical world tangible and have physical consequences, you must map Potential Fields in the GameServer's `main.go` source logic to closely match your visual SDFs. 

| Object type | X, Y, Z Coordinates | Amplitude (Strength) | Sigma (Radius) |
| --- | --- | --- | --- |
| **Ground Hover Lock** | `0, 0, -20` | `50000.0` | `50.0` |
| **Sun Gravity Well** | `0, 0, 500` | `-100000.0` | `300.0` |
| **Arena Central Repulsion** | `0, 0, 0` | `1000.0` | `80.0` |

### Setting up the Server Initialization
Before starting the Go `GameServer`, edit the initialization logic in `GameServer/main.go` to inject the fields directly into the `PotentialFieldRegistry`:

```go
fields.AddEmitters(0, 0, -20, 50000.0, 50.0)      // The default Ground Lock
fields.AddEmitters(0, 0, 500, -100000.0, 300.0)   // The Sun Gravity Source
fields.AddEmitters(0, 0, 0, 1000.0, 80.0)         // Central Arena repellant
```

By ensuring that the Go backend physics completely mirror the mathematical visual inputs compiled by the Python compiler, the complete universe is effectively realized.
