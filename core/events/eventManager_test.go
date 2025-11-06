package events

import (
	"sync"
	"testing"
	"time"
)

// Mock Publisher for testing
type mockPublisher struct {
	started    bool
	stopped    bool
	messages   chan *Message
	onNewSubCh chan struct{}
	mu         sync.Mutex
}

func newMockPublisher() *mockPublisher {
	return &mockPublisher{
		messages:   make(chan *Message, 10),
		onNewSubCh: make(chan struct{}, 10),
	}
}

func (m *mockPublisher) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.started = true
	return nil
}

func (m *mockPublisher) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stopped = true
	close(m.messages)
	return nil
}

func (m *mockPublisher) OnNewSub() error {
	m.onNewSubCh <- struct{}{}
	return nil
}

func (m *mockPublisher) Read() chan *Message {
	return m.messages
}

func (m *mockPublisher) SendMessage(msg *Message) {
	m.messages <- msg
}

func (m *mockPublisher) IsStarted() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.started
}

func (m *mockPublisher) IsStopped() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.stopped
}

// Helper to reset state between tests
func resetState() {
	state.eventLock.Lock()
	defer state.eventLock.Unlock()
	state.sessions = make(map[string]*session)
	state.subscriptions = make(map[string]map[string]*session)
	state.publishers = make(map[string]Publisher)
}

func TestRegisterSession(t *testing.T) {
	resetState()

	sessionID := "test-session-1"
	ch, unregister := RegisterSession(sessionID)

	if ch == nil {
		t.Fatal("Expected channel to be non-nil")
	}

	state.eventLock.Lock()
	sess, exists := state.sessions[sessionID]
	state.eventLock.Unlock()

	if !exists {
		t.Fatal("Session not registered")
	}

	if sess.ID != sessionID {
		t.Errorf("Expected session ID %s, got %s", sessionID, sess.ID)
	}

	if cap(sess.Out) != 100 {
		t.Errorf("Expected channel buffer size 100, got %d", cap(sess.Out))
	}

	// Test unregister
	unregister()

	state.eventLock.Lock()
	_, exists = state.sessions[sessionID]
	state.eventLock.Unlock()

	if exists {
		t.Error("Session should be unregistered")
	}
}

func TestSubscribe(t *testing.T) {
	resetState()

	sessionID := "test-session-2"
	topic := "test-topic"

	_, unregister := RegisterSession(sessionID)
	defer unregister()

	mockPub := newMockPublisher()

	err := Subscribe(sessionID, topic, mockPub)
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	if !mockPub.IsStarted() {
		t.Error("Publisher should be started")
	}

	state.eventLock.Lock()
	subs, exists := state.subscriptions[topic]
	state.eventLock.Unlock()

	if !exists {
		t.Fatal("Topic not in subscriptions")
	}

	if _, exists := subs[sessionID]; !exists {
		t.Error("Session not subscribed to topic")
	}

	// Wait for OnNewSub to be called
	select {
	case <-mockPub.onNewSubCh:
		// Success
	case <-time.After(100 * time.Millisecond):
		t.Error("OnNewSub not called")
	}
}

func TestSubscribeNonExistentSession(t *testing.T) {
	resetState()

	mockPub := newMockPublisher()
	err := Subscribe("non-existent", "test-topic", mockPub)

	if err == nil {
		t.Error("Expected error when subscribing non-existent session")
	}

	if mockPub.IsStarted() {
		t.Error("Publisher should not be started for non-existent session")
	}
}

func TestUnSubscribe(t *testing.T) {
	resetState()

	sessionID := "test-session-3"
	topic := "test-topic"

	_, unregister := RegisterSession(sessionID)
	defer unregister()

	mockPub := newMockPublisher()
	Subscribe(sessionID, topic, mockPub)

	UnSubscribe(sessionID, topic)

	state.eventLock.Lock()
	subs, exists := state.subscriptions[topic]
	state.eventLock.Unlock()

	if exists && len(subs) > 0 {
		t.Error("Session should be unsubscribed from topic")
	}

	if !mockPub.IsStopped() {
		t.Error("Publisher should be stopped when no subscribers")
	}
}

func TestMessageDelivery(t *testing.T) {
	resetState()

	sessionID := "test-session-4"
	topic := "test-topic"

	ch, unregister := RegisterSession(sessionID)
	defer unregister()

	mockPub := newMockPublisher()
	Subscribe(sessionID, topic, mockPub)

	// Give runSubscription goroutine time to start
	time.Sleep(50 * time.Millisecond)

	testMsg := &Message{
		Event: "test-event",
		Data:  "test-data",
	}

	mockPub.SendMessage(testMsg)

	select {
	case receivedMsg := <-ch:
		if receivedMsg.Event != testMsg.Event {
			t.Errorf("Expected event %s, got %s", testMsg.Event, receivedMsg.Event)
		}
		if receivedMsg.Data != testMsg.Data {
			t.Errorf("Expected data %v, got %v", testMsg.Data, receivedMsg.Data)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Message not received")
	}
}

func TestNonBlockingMessageDelivery(t *testing.T) {
	resetState()

	sessionID := "test-session-5"
	topic := "test-topic"

	ch, unregister := RegisterSession(sessionID)
	defer unregister()

	mockPub := newMockPublisher()
	Subscribe(sessionID, topic, mockPub)

	// Give runSubscription goroutine time to start
	time.Sleep(50 * time.Millisecond)

	// Fill the channel buffer (capacity 100)
	for i := 0; i < 100; i++ {
		mockPub.SendMessage(&Message{
			Event: "fill",
			Data:  i,
		})
	}

	// Wait for messages to be sent
	time.Sleep(100 * time.Millisecond)

	// Don't read from the channel, so it's full

	// Send more messages - these should be dropped without blocking
	for i := 0; i < 10; i++ {
		mockPub.SendMessage(&Message{
			Event: "overflow",
			Data:  i,
		})
	}

	// If we get here without hanging, the test passes
	// The publisher should still be running

	// Drain the channel
	timeout := time.After(1 * time.Second)
	for i := 0; i < 100; i++ {
		select {
		case <-ch:
			// Message received
		case <-timeout:
			t.Fatal("Timeout draining channel")
		}
	}
}

func TestMultipleSubscribers(t *testing.T) {
	resetState()

	topic := "test-topic"
	sessionIDs := []string{"session-1", "session-2", "session-3"}
	channels := make([]chan *Message, 3)
	unregisters := make([]func(), 3)

	// Register multiple sessions
	for i, sessionID := range sessionIDs {
		ch, unreg := RegisterSession(sessionID)
		channels[i] = ch
		unregisters[i] = unreg
		defer unreg()
	}

	mockPub := newMockPublisher()

	// Subscribe all sessions to the same topic
	for _, sessionID := range sessionIDs {
		err := Subscribe(sessionID, topic, mockPub)
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}
	}

	// Give runSubscription goroutine time to start
	time.Sleep(50 * time.Millisecond)

	testMsg := &Message{
		Event: "broadcast",
		Data:  "test-data",
	}

	mockPub.SendMessage(testMsg)

	// All sessions should receive the message
	for i, ch := range channels {
		select {
		case receivedMsg := <-ch:
			if receivedMsg.Event != testMsg.Event {
				t.Errorf("Session %d: Expected event %s, got %s", i, testMsg.Event, receivedMsg.Event)
			}
		case <-time.After(500 * time.Millisecond):
			t.Errorf("Session %d: Message not received", i)
		}
	}
}

func TestConcurrentSubscribeUnsubscribe(t *testing.T) {
	resetState()

	const numGoroutines = 10
	const numIterations = 20

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Run concurrent subscribe/unsubscribe operations
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			for j := 0; j < numIterations; j++ {
				sessionID := "concurrent-session"
				topic := "concurrent-topic"

				ch, unregister := RegisterSession(sessionID)
				mockPub := newMockPublisher()

				Subscribe(sessionID, topic, mockPub)

				// Send a message
				testMsg := &Message{
					Event: "test",
					Data:  id,
				}
				mockPub.SendMessage(testMsg)

				// Try to receive (may or may not get it due to timing)
				select {
				case <-ch:
				case <-time.After(10 * time.Millisecond):
				}

				UnSubscribe(sessionID, topic)
				unregister()
			}
		}(i)
	}

	// Wait for all goroutines to complete
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Success - no deadlock or panic
	case <-time.After(10 * time.Second):
		t.Fatal("Test timed out - possible deadlock")
	}
}

func TestPublisherLifecycle(t *testing.T) {
	resetState()

	topic := "test-topic"
	sessionID1 := "session-1"
	sessionID2 := "session-2"

	_, unreg1 := RegisterSession(sessionID1)
	defer unreg1()
	_, unreg2 := RegisterSession(sessionID2)
	defer unreg2()

	mockPub := newMockPublisher()

	// First subscription should start publisher
	Subscribe(sessionID1, topic, mockPub)
	if !mockPub.IsStarted() {
		t.Error("Publisher should be started on first subscription")
	}

	// Second subscription should reuse publisher
	Subscribe(sessionID2, topic, mockPub)

	state.eventLock.Lock()
	pubCount := len(state.publishers)
	state.eventLock.Unlock()

	if pubCount != 1 {
		t.Errorf("Expected 1 publisher, got %d", pubCount)
	}

	// Unsubscribe one session - publisher should still be running
	UnSubscribe(sessionID1, topic)
	if mockPub.IsStopped() {
		t.Error("Publisher should not be stopped while subscribers remain")
	}

	// Unsubscribe last session - publisher should stop
	UnSubscribe(sessionID2, topic)
	if !mockPub.IsStopped() {
		t.Error("Publisher should be stopped when no subscribers remain")
	}
}

func TestSessionUnregisterCleansUpSubscriptions(t *testing.T) {
	resetState()

	sessionID := "test-session-cleanup"
	topic := "test-topic"

	_, unregister := RegisterSession(sessionID)

	mockPub := newMockPublisher()
	Subscribe(sessionID, topic, mockPub)

	// Unregister session (should clean up subscriptions)
	unregister()

	state.eventLock.Lock()
	_, sessExists := state.sessions[sessionID]
	subs, topicExists := state.subscriptions[topic]
	state.eventLock.Unlock()

	if sessExists {
		t.Error("Session should be removed")
	}

	if topicExists && len(subs) > 0 {
		t.Error("Subscriptions should be cleaned up")
	}
}

func TestMultipleTopicsPerSession(t *testing.T) {
	resetState()

	sessionID := "multi-topic-session"
	topics := []string{"topic-1", "topic-2", "topic-3"}

	ch, unregister := RegisterSession(sessionID)
	defer unregister()

	mockPubs := make([]*mockPublisher, len(topics))

	// Subscribe to multiple topics
	for i, topic := range topics {
		mockPubs[i] = newMockPublisher()
		err := Subscribe(sessionID, topic, mockPubs[i])
		if err != nil {
			t.Fatalf("Subscribe to %s failed: %v", topic, err)
		}
	}

	// Give goroutines time to start
	time.Sleep(50 * time.Millisecond)

	// Send message to each topic
	for i, mockPub := range mockPubs {
		mockPub.SendMessage(&Message{
			Event: topics[i],
			Data:  i,
		})
	}

	// Receive messages from all topics
	received := make(map[string]bool)
	timeout := time.After(1 * time.Second)

	for i := 0; i < len(topics); i++ {
		select {
		case msg := <-ch:
			received[msg.Event] = true
		case <-timeout:
			t.Fatal("Timeout receiving messages")
		}
	}

	// Verify all topics received
	for _, topic := range topics {
		if !received[topic] {
			t.Errorf("Did not receive message for topic %s", topic)
		}
	}
}
