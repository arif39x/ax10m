package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type SpatialRegistry struct {
	mu   sync.RWMutex
	X    []float32
	Y    []float32
	Z    []float32
	VX   []float32
	VY   []float32
	VZ   []float32
	AX   []float32
	AY   []float32
	AZ   []float32
	Mass []float32
}

type PotentialFieldRegistry struct {
	mu        sync.RWMutex
	X         []float32
	Y         []float32
	Z         []float32
	Amplitude []float32
	Sigma     []float32
}

func NewSpatialRegistry() *SpatialRegistry {
	return &SpatialRegistry{
		X:    make([]float32, 0),
		Y:    make([]float32, 0),
		Z:    make([]float32, 0),
		VX:   make([]float32, 0),
		VY:   make([]float32, 0),
		VZ:   make([]float32, 0),
		AX:   make([]float32, 0),
		AY:   make([]float32, 0),
		AZ:   make([]float32, 0),
		Mass: make([]float32, 0),
	}
}

func NewPotentialFieldRegistry() *PotentialFieldRegistry {
	return &PotentialFieldRegistry{
		X:         make([]float32, 0),
		Y:         make([]float32, 0),
		Z:         make([]float32, 0),
		Amplitude: make([]float32, 0),
		Sigma:     make([]float32, 0),
	}
}

func (sr *SpatialRegistry) AddEntity(x, y, z, mass float32) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	sr.X = append(sr.X, x)
	sr.Y = append(sr.Y, y)
	sr.Z = append(sr.Z, z)
	sr.VX = append(sr.VX, 0)
	sr.VY = append(sr.VY, 0)
	sr.VZ = append(sr.VZ, 0)
	sr.AX = append(sr.AX, 0)
	sr.AY = append(sr.AY, 0)
	sr.AZ = append(sr.AZ, -9.8)
	sr.Mass = append(sr.Mass, mass)
}

func (pfr *PotentialFieldRegistry) AddEmitters(x, y, z, amplitude, sigma float32) {
	pfr.mu.Lock()
	defer pfr.mu.Unlock()
	pfr.X = append(pfr.X, x)
	pfr.Y = append(pfr.Y, y)
	pfr.Z = append(pfr.Z, z)
	pfr.Amplitude = append(pfr.Amplitude, amplitude)
	pfr.Sigma = append(pfr.Sigma, sigma)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

var (
	clients    = make(map[*Client]bool)
	clientsMu  sync.Mutex
	broadcast  = make(chan []byte)
	register   = make(chan *Client)
	unregister = make(chan *Client)
)

func runHub() {
	for {
		select {
		case client := <-register:
			clientsMu.Lock()
			clients[client] = true
			clientsMu.Unlock()
		case client := <-unregister:
			clientsMu.Lock()
			if _, ok := clients[client]; ok {
				delete(clients, client)
				close(client.send)
			}
			clientsMu.Unlock()
		case message := <-broadcast:
			clientsMu.Lock()
			for client := range clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(clients, client)
				}
			}
			clientsMu.Unlock()
		}
	}
}

func clientWriter(client *Client) {
	for message := range client.send {
		err := client.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
	client.conn.Close()
}

type IncomingMsg struct {
	Equation string `json:"equation"`
}

type CompileRequest struct {
	Equation string `json:"equation"`
}

type CompileResponse struct {
	Status string `json:"status"`
	Wgsl   string `json:"wgsl"`
}

func handleClientMessage(client *Client, message []byte) {
	var in IncomingMsg
	if err := json.Unmarshal(message, &in); err == nil && in.Equation != "" {
		reqBody, _ := json.Marshal(CompileRequest{Equation: in.Equation})
		resp, err := http.Post("http://127.0.0.1:8081/compile_sdf", "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			log.Println("Error calling python:", err)
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var compResp map[string]interface{}
		if err := json.Unmarshal(body, &compResp); err == nil {
			if wgsl, ok := compResp["wgsl"].(string); ok {
				outMsg := map[string]string{
					"wgsl": wgsl,
				}
				outBytes, _ := json.Marshal(outMsg)
				broadcast <- outBytes
			}
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	client := &Client{conn: conn, send: make(chan []byte, 256)}
	register <- client
	go clientWriter(client)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			unregister <- client
			break
		}
		handleClientMessage(client, message)
	}
}

func main() {
	registry := NewSpatialRegistry()
	fields := NewPotentialFieldRegistry()

	registry.AddEntity(0, 0, 100, 1.0)

	fields.AddEmitters(0, 0, -20, 50000.0, 50.0)

	go runHub()

	http.HandleFunc("/ws", serveWs)
	go func() {
		log.Println("Go Spatial Core listening on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal("ListenAndServe:", err)
		}
	}()

	ticker := time.NewTicker(time.Millisecond * 50)
	defer ticker.Stop()

	dt := float32(0.05)

	type JuliaRequest struct {
		X              []float32 `json:"x"`
		Y              []float32 `json:"y"`
		Z              []float32 `json:"z"`
		FieldX         []float32 `json:"field_x"`
		FieldY         []float32 `json:"field_y"`
		FieldZ         []float32 `json:"field_z"`
		FieldAmplitude []float32 `json:"field_amplitude"`
		FieldSigma     []float32 `json:"field_sigma"`
		Mass           []float32 `json:"mass"`
	}

	type JuliaResponse struct {
		AX []float32 `json:"ax"`
		AY []float32 `json:"ay"`
		AZ []float32 `json:"az"`
	}

	for {
		<-ticker.C

		registry.mu.Lock()
		fields.mu.RLock()

		reqData := JuliaRequest{
			X:              registry.X,
			Y:              registry.Y,
			Z:              registry.Z,
			FieldX:         fields.X,
			FieldY:         fields.Y,
			FieldZ:         fields.Z,
			FieldAmplitude: fields.Amplitude,
			FieldSigma:     fields.Sigma,
			Mass:           registry.Mass,
		}

		reqBytes, _ := json.Marshal(reqData)
		resp, err := http.Post("http://127.0.0.1:50051", "application/json", bytes.NewBuffer(reqBytes))

		if err == nil {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			var juliaResp JuliaResponse
			if err := json.Unmarshal(body, &juliaResp); err == nil {
				if len(juliaResp.AX) == len(registry.AX) {
					copy(registry.AX, juliaResp.AX)
					copy(registry.AY, juliaResp.AY)
					copy(registry.AZ, juliaResp.AZ)
				}
			}
		} else {
			for i := 0; i < len(registry.AX); i++ {
				registry.AX[i] = 0
				registry.AY[i] = 0
				registry.AZ[i] = -9.8
			}
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
			broadcast <- payload
		}
	}
}
