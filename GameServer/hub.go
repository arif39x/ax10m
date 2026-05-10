package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		handleClientMessage(c, message)
	}
}

func (c *Client) WritePump() {
	for message := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
	c.conn.Close()
}

type IncomingMsg struct {
	Equation string `json:"equation"`
}

type CompileRequest struct {
	Equation string `json:"equation"`
}

func handleClientMessage(client *Client, message []byte) {
	var in IncomingMsg
	if err := json.Unmarshal(message, &in); err == nil && in.Equation != "" {
		log.Printf("Received equation: %s", in.Equation)
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
				client.hub.broadcast <- outBytes
			} else if detail, ok := compResp["detail"].(string); ok {
				log.Printf("Compilation failed: %s", detail)
			}
		}
	}
}
