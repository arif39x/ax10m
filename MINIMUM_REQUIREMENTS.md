# AxiomMP: Minimum Requirement Logic

Before you can start playing this mathematical sandbox, you must ensure both your **System Requirements** and **Mathematical Requirements** meet the baseline logic. If either of these is missing, the universe will fail to materialize or your avatar will fall infinitely into the void.

---

## 1. System Requirements (The 4 Pillars)
To sustain the universe, you need a multi-language distributed environment. All 4 microservices must be active simultaneously:

1. **Rust (Visual Conduit)**: Cargo & Rust 1.70+ 
2. **Go (State Authority)**: Go 1.20+
3. **Python (Math Compiler)**: Python 3.10+ (with sympy, fastapi, uvicorn)
4. **Julia (Physics Oracle)**: Julia 1.9+ 

*Without all 4 running, the GameClient will either show a black screen, or objects will not have physics computed.*

## 2. In-Game Minimum Logic (The Floor and The Force)

Because AxiomMP has no predefined levels, starting the game requires you to manually define the mathematical logic to stand on. You need two things:

### A. The Visual Floor (SDF)
You need to compile an equation so the GPU draws a floor. The simplest flat plane equation below your starting coordinate is:
```json
{
  "equation": "z + 50"
}
```
*Submit this payload to the GameServer over the WebSocket. The Python MathCompiler will translate `z + 50` so the Rust client renders an infinite floor at `z = -50`.*

### B. The Anti-Gravity Platform (Potential Field)
You need a physical force to stop your `-9.8` gravitational descent. The Go server by default registers an emitter under you:
```go
// Location: GameServer/main.go
fields.AddEmitters(0, 0, -20, 50000.0, 50.0)
```
*This emitter uses a Gaussian potential field. The massive amplitude (`50000.0`) repels your mass exponentially as you approach `z = -20`, creating a perfectly locked hover state right above the visual floor.*

### Summary of Minimum Play State:
To not fall into the void, ensure the GameServer's `AddEmitters` logic aligns with your submitted SDF equation, so the visual floor and the physical hover-field exist in the same spatial region.
