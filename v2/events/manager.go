package events

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// EventType represents the type of event
type EventType string

const (
	EventAssetCreated   EventType = "asset.created"
	EventAssetUpdated   EventType = "asset.updated"
	EventAssetDeleted   EventType = "asset.deleted"
	EventScanStarted    EventType = "scan.started"
	EventScanProgress   EventType = "scan.progress"
	EventScanCompleted  EventType = "scan.completed"
	EventProcessStarted EventType = "process.started"
	EventProcessProgress EventType = "process.progress"
	EventProcessCompleted EventType = "process.completed"
	EventPrinterStatus  EventType = "printer.status"
	EventError          EventType = "error"
)

// Event represents a WebSocket event
type Event struct {
	Type      EventType   `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// Client represents a connected WebSocket client
type Client struct {
	ID       string
	conn     *websocket.Conn
	send     chan *Event
	manager  *EventManager
	ctx      context.Context
	cancel   context.CancelFunc
}

// EventManager manages WebSocket connections and event broadcasting
type EventManager struct {
	clients    map[*Client]bool
	broadcast  chan *Event
	register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex
	upgrader   websocket.Upgrader
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewEventManager creates a new event manager
func NewEventManager() *EventManager {
	ctx, cancel := context.WithCancel(context.Background())

	return &EventManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *Event, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins in development
			},
		},
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start begins the event manager's main loop
func (m *EventManager) Start() {
	go m.run()
}

// Stop gracefully stops the event manager
func (m *EventManager) Stop() {
	m.cancel()

	// Close all client connections
	m.mutex.Lock()
	for client := range m.clients {
		client.cancel()
		client.conn.Close()
	}
	m.mutex.Unlock()
}

// run is the main event loop
func (m *EventManager) run() {
	slog.Info("Event manager started")

	for {
		select {
		case <-m.ctx.Done():
			slog.Info("Event manager stopping")
			return

		case client := <-m.register:
			m.mutex.Lock()
			m.clients[client] = true
			m.mutex.Unlock()
			slog.Info("Client registered", "client_id", client.ID)

		case client := <-m.unregister:
			m.mutex.Lock()
			if _, ok := m.clients[client]; ok {
				delete(m.clients, client)
				close(client.send)
			}
			m.mutex.Unlock()
			slog.Info("Client unregistered", "client_id", client.ID)

		case event := <-m.broadcast:
			m.mutex.RLock()
			for client := range m.clients {
				select {
				case client.send <- event:
					// Event sent successfully
				default:
					// Client's send channel is full, close connection
					m.mutex.RUnlock()
					m.unregister <- client
					m.mutex.RLock()
				}
			}
			m.mutex.RUnlock()
		}
	}
}

// HandleWebSocket handles WebSocket upgrade requests
func (m *EventManager) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := m.upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("WebSocket upgrade failed", "error", err)
		return
	}

	ctx, cancel := context.WithCancel(m.ctx)
	client := &Client{
		ID:      generateClientID(),
		conn:    conn,
		send:    make(chan *Event, 256),
		manager: m,
		ctx:     ctx,
		cancel:  cancel,
	}

	m.register <- client

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()
}

// Broadcast sends an event to all connected clients
func (m *EventManager) Broadcast(eventType EventType, data interface{}) {
	event := &Event{
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now(),
	}

	select {
	case m.broadcast <- event:
		// Event queued successfully
	default:
		// Broadcast channel is full, log warning
		slog.Warn("Broadcast channel full, event dropped", "type", eventType)
	}
}

// BroadcastAssetCreated broadcasts an asset creation event
func (m *EventManager) BroadcastAssetCreated(assetID string, assetData interface{}) {
	m.Broadcast(EventAssetCreated, map[string]interface{}{
		"id":   assetID,
		"data": assetData,
	})
}

// BroadcastAssetUpdated broadcasts an asset update event
func (m *EventManager) BroadcastAssetUpdated(assetID string, assetData interface{}) {
	m.Broadcast(EventAssetUpdated, map[string]interface{}{
		"id":   assetID,
		"data": assetData,
	})
}

// BroadcastAssetDeleted broadcasts an asset deletion event
func (m *EventManager) BroadcastAssetDeleted(assetID string) {
	m.Broadcast(EventAssetDeleted, map[string]interface{}{
		"id": assetID,
	})
}

// BroadcastScanProgress broadcasts scan progress
func (m *EventManager) BroadcastScanProgress(current, total int, message string) {
	m.Broadcast(EventScanProgress, map[string]interface{}{
		"current": current,
		"total":   total,
		"message": message,
	})
}

// ClientCount returns the number of connected clients
func (m *EventManager) ClientCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.clients)
}

// writePump pumps messages from the send channel to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case <-c.ctx.Done():
			return

		case event, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			if err := json.NewEncoder(w).Encode(event); err != nil {
				slog.Error("Failed to encode event", "error", err)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// readPump pumps messages from the WebSocket connection
func (c *Client) readPump() {
	defer func() {
		c.manager.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			_, _, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					slog.Error("WebSocket error", "error", err)
				}
				return
			}
			// For now, we ignore client messages (one-way communication)
		}
	}
}

// generateClientID generates a unique client ID
func generateClientID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
