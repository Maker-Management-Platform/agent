package integration_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	v2events "github.com/eduardooliveira/stLib/v2/events"
	v2entities "github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEventManagerCreation tests event manager creation
func TestEventManagerCreation(t *testing.T) {
	mgr := v2events.NewEventManager()
	require.NotNil(t, mgr)

	// Start the event manager
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go mgr.Start(ctx)

	// Give it time to start
	time.Sleep(100 * time.Millisecond)
}

// TestEventBroadcasting tests basic event broadcasting
func TestEventBroadcasting(t *testing.T) {
	mgr := v2events.NewEventManager()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go mgr.Start(ctx)
	time.Sleep(100 * time.Millisecond)

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mgr.HandleWebSocket(w, r)
	}))
	defer server.Close()

	// Connect WebSocket client
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer ws.Close()

	// Broadcast an event
	go func() {
		time.Sleep(100 * time.Millisecond)
		mgr.Broadcast(v2events.EventAssetCreated, map[string]interface{}{
			"id":   "test-123",
			"name": "test.stl",
		})
	}()

	// Read event from WebSocket
	ws.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, message, err := ws.ReadMessage()
	require.NoError(t, err)

	// Parse event
	var event v2events.Event
	err = json.Unmarshal(message, &event)
	require.NoError(t, err)

	assert.Equal(t, v2events.EventAssetCreated, event.Type)
	assert.NotNil(t, event.Data)
}

// TestMultipleClients tests broadcasting to multiple clients
func TestMultipleClients(t *testing.T) {
	mgr := v2events.NewEventManager()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go mgr.Start(ctx)
	time.Sleep(100 * time.Millisecond)

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mgr.HandleWebSocket(w, r)
	}))
	defer server.Close()

	// Connect multiple clients
	numClients := 5
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	var clients []*websocket.Conn
	for i := 0; i < numClients; i++ {
		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		require.NoError(t, err)
		defer ws.Close()
		clients = append(clients, ws)
	}

	time.Sleep(200 * time.Millisecond)

	// Broadcast event
	testData := map[string]interface{}{
		"message": "test broadcast",
	}
	mgr.Broadcast(v2events.EventScanStarted, testData)

	// All clients should receive the event
	var wg sync.WaitGroup
	wg.Add(numClients)

	for i, ws := range clients {
		go func(clientNum int, client *websocket.Conn) {
			defer wg.Done()

			client.SetReadDeadline(time.Now().Add(2 * time.Second))
			_, message, err := client.ReadMessage()
			if err != nil {
				t.Errorf("Client %d failed to read: %v", clientNum, err)
				return
			}

			var event v2events.Event
			err = json.Unmarshal(message, &event)
			if err != nil {
				t.Errorf("Client %d failed to unmarshal: %v", clientNum, err)
				return
			}

			assert.Equal(t, v2events.EventScanStarted, event.Type)
		}(i, ws)
	}

	wg.Wait()
}

// TestEventTypes tests all event types
func TestEventTypes(t *testing.T) {
	mgr := v2events.NewEventManager()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go mgr.Start(ctx)
	time.Sleep(100 * time.Millisecond)

	eventTypes := []v2events.EventType{
		v2events.EventAssetCreated,
		v2events.EventAssetUpdated,
		v2events.EventAssetDeleted,
		v2events.EventScanStarted,
		v2events.EventScanProgress,
		v2events.EventScanCompleted,
		v2events.EventProcessStarted,
		v2events.EventProcessProgress,
		v2events.EventProcessCompleted,
		v2events.EventPrinterStatus,
		v2events.EventError,
	}

	for _, eventType := range eventTypes {
		t.Run(string(eventType), func(t *testing.T) {
			// Just verify the event type can be broadcast
			mgr.Broadcast(eventType, map[string]interface{}{
				"test": true,
			})
		})
	}
}

// TestAssetCreatedEvent tests asset created event helper
func TestAssetCreatedEvent(t *testing.T) {
	mgr := v2events.NewEventManager()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go mgr.Start(ctx)
	time.Sleep(100 * time.Millisecond)

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mgr.HandleWebSocket(w, r)
	}))
	defer server.Close()

	// Connect client
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer ws.Close()

	// Create test asset
	testAsset := &v2entities.Asset{
		ID: "test-asset-123",
	}
	label := "test.stl"
	testAsset.Label = &label

	// Broadcast asset created
	go func() {
		time.Sleep(100 * time.Millisecond)
		mgr.BroadcastAssetCreated(testAsset.ID, testAsset)
	}()

	// Read event
	ws.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, message, err := ws.ReadMessage()
	require.NoError(t, err)

	var event v2events.Event
	err = json.Unmarshal(message, &event)
	require.NoError(t, err)

	assert.Equal(t, v2events.EventAssetCreated, event.Type)
}

// TestAssetUpdatedEvent tests asset updated event helper
func TestAssetUpdatedEvent(t *testing.T) {
	mgr := v2events.NewEventManager()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go mgr.Start(ctx)
	time.Sleep(100 * time.Millisecond)

	testAsset := &v2entities.Asset{
		ID: "test-asset-456",
	}

	// Should not panic
	mgr.BroadcastAssetUpdated(testAsset.ID, testAsset)
}

// TestAssetDeletedEvent tests asset deleted event helper
func TestAssetDeletedEvent(t *testing.T) {
	mgr := v2events.NewEventManager()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go mgr.Start(ctx)
	time.Sleep(100 * time.Millisecond)

	// Should not panic
	mgr.BroadcastAssetDeleted("test-asset-789")
}

// TestConnectionHandling tests WebSocket connection lifecycle
func TestConnectionHandling(t *testing.T) {
	mgr := v2events.NewEventManager()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go mgr.Start(ctx)
	time.Sleep(100 * time.Millisecond)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mgr.HandleWebSocket(w, r)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Connect
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)

	// Disconnect
	err = ws.Close()
	assert.NoError(t, err)

	// Manager should handle disconnection gracefully
	time.Sleep(100 * time.Millisecond)
}

// TestEventOrdering tests that events maintain order
func TestEventOrdering(t *testing.T) {
	mgr := v2events.NewEventManager()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go mgr.Start(ctx)
	time.Sleep(100 * time.Millisecond)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mgr.HandleWebSocket(w, r)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer ws.Close()

	// Send multiple events
	numEvents := 5
	go func() {
		time.Sleep(100 * time.Millisecond)
		for i := 0; i < numEvents; i++ {
			mgr.Broadcast(v2events.EventScanProgress, map[string]interface{}{
				"step": i,
			})
			time.Sleep(10 * time.Millisecond)
		}
	}()

	// Read events and verify order
	receivedSteps := make([]int, 0, numEvents)
	ws.SetReadDeadline(time.Now().Add(3 * time.Second))

	for i := 0; i < numEvents; i++ {
		_, message, err := ws.ReadMessage()
		if err != nil {
			break
		}

		var event v2events.Event
		json.Unmarshal(message, &event)

		if data, ok := event.Data.(map[string]interface{}); ok {
			if step, ok := data["step"].(float64); ok {
				receivedSteps = append(receivedSteps, int(step))
			}
		}
	}

	// Should maintain order
	for i := 0; i < len(receivedSteps)-1; i++ {
		assert.LessOrEqual(t, receivedSteps[i], receivedSteps[i+1],
			"Events should be in order")
	}
}

// TestConcurrentBroadcasts tests concurrent broadcasting
func TestConcurrentBroadcasts(t *testing.T) {
	mgr := v2events.NewEventManager()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go mgr.Start(ctx)
	time.Sleep(100 * time.Millisecond)

	// Broadcast from multiple goroutines
	var wg sync.WaitGroup
	numGoroutines := 10
	eventsPerGoroutine := 10

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < eventsPerGoroutine; j++ {
				mgr.Broadcast(v2events.EventScanProgress, map[string]interface{}{
					"goroutine": id,
					"event":     j,
				})
			}
		}(i)
	}

	wg.Wait()
	// Should not panic or deadlock
}

// TestEventTimestamps tests event timestamps are set correctly
func TestEventTimestamps(t *testing.T) {
	mgr := v2events.NewEventManager()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go mgr.Start(ctx)
	time.Sleep(100 * time.Millisecond)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mgr.HandleWebSocket(w, r)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer ws.Close()

	beforeBroadcast := time.Now()

	go func() {
		time.Sleep(100 * time.Millisecond)
		mgr.Broadcast(v2events.EventAssetCreated, map[string]interface{}{})
	}()

	ws.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, message, err := ws.ReadMessage()
	require.NoError(t, err)

	afterReceive := time.Now()

	var event v2events.Event
	err = json.Unmarshal(message, &event)
	require.NoError(t, err)

	// Timestamp should be between before broadcast and after receive
	assert.True(t, event.Timestamp.After(beforeBroadcast) || event.Timestamp.Equal(beforeBroadcast))
	assert.True(t, event.Timestamp.Before(afterReceive) || event.Timestamp.Equal(afterReceive))
}

// BenchmarkEventBroadcasting benchmarks event broadcasting performance
func BenchmarkEventBroadcasting(b *testing.B) {
	mgr := v2events.NewEventManager()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go mgr.Start(ctx)
	time.Sleep(100 * time.Millisecond)

	testData := map[string]interface{}{
		"test": "data",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mgr.Broadcast(v2events.EventAssetCreated, testData)
	}
}

// BenchmarkMultiClientBroadcast benchmarks broadcasting to multiple clients
func BenchmarkMultiClientBroadcast(b *testing.B) {
	mgr := v2events.NewEventManager()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go mgr.Start(ctx)
	time.Sleep(100 * time.Millisecond)

	// This would require setting up actual WebSocket connections
	// For now, just benchmark the broadcast mechanism
	testData := map[string]interface{}{
		"test": "data",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mgr.Broadcast(v2events.EventAssetCreated, testData)
	}
}
