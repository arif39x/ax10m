package main

import (
	"encoding/json"
	"math"
	"time"
)

func StartPhysicsLoop(registry *SpatialRegistry, fields *PotentialFieldRegistry, hub *Hub) {
	ticker := time.NewTicker(time.Millisecond * 50)
	defer ticker.Stop()

	dt := float32(0.05)

	for {
		<-ticker.C

		registry.mu.Lock()
		fields.mu.RLock()

		// In-process Physics Calculation
		for i := 0; i < len(registry.X); i++ {
			netAX := float32(0.0)
			netAY := float32(0.0)
			netAZ := float32(-9.8)

			// THE VOID ZONE (Mystery Law)
			// Inverted gravity in a specific coordinate range
			if registry.X[i] > 50.0 && registry.X[i] < 100.0 &&
				registry.Y[i] > 50.0 && registry.Y[i] < 100.0 {
				netAZ = 5.0 // Upward pull
			}

			mi := registry.Mass[i]
			if mi == 0 {
				mi = 1.0
			}

			for j := 0; j < len(fields.X); j++ {
				dx := registry.X[i] - fields.X[j]
				dy := registry.Y[i] - fields.Y[j]
				dz := registry.Z[i] - fields.Z[j]

				distSq := dx*dx + dy*dy + dz*dz
				sigmaSq := fields.Sigma[j] * fields.Sigma[j]

				val := fields.Amplitude[j] * float32(math.Exp(float64(-distSq/(2.0*sigmaSq))))
				invMSigmaSq := 1.0 / (mi * sigmaSq)
				netAX += dx * val * invMSigmaSq
				netAY += dy * val * invMSigmaSq
				netAZ += dz * val * invMSigmaSq
			}

			registry.AX[i] = netAX
			registry.AY[i] = netAY
			registry.AZ[i] = netAZ
		}

		for i := 0; i < len(registry.X); i++ {
			registry.VX[i] += registry.AX[i] * dt
			registry.VY[i] += registry.AY[i] * dt
			registry.VZ[i] += registry.AZ[i] * dt

			registry.X[i] += registry.VX[i] * dt
			registry.Y[i] += registry.VY[i] * dt
			registry.Z[i] += registry.VZ[i] * dt

			if registry.Z[i] < -50.0 {
				registry.Z[i] = -50.0
				registry.VZ[i] *= -0.5
			}
		}

		state := map[string]interface{}{
			"x": registry.X,
			"y": registry.Y,
			"z": registry.Z,
		}

		fields.mu.RUnlock()
		registry.mu.Unlock()

		payload, err := json.Marshal(state)
		if err == nil {
			hub.broadcast <- payload
		}
	}
}
