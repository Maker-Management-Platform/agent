package events

import (
	"errors"
	"log"
	"sync"
)

type session struct {
	ID            string
	Out           chan *Message
	subscriptions map[string]struct{}
}

type stateType struct {
	sessions      map[string]*session
	subscriptions map[string]map[string]*session
	publishers    map[string]Publisher
	eventLock     sync.Mutex
}

type Publisher interface {
	Start() error
	Stop() error
	OnNewSub() error
	Read() chan *Message
}

var state *stateType

func init() {
	state = &stateType{
		sessions:      make(map[string]*session),
		subscriptions: make(map[string]map[string]*session),
		publishers:    make(map[string]Publisher),
		eventLock:     sync.Mutex{},
	}
}

func RegisterSession(id string) (chan *Message, func()) {
	state.eventLock.Lock()
	defer state.eventLock.Unlock()
	state.sessions[id] = &session{
		ID:            id,
		Out:           make(chan *Message, 100),
		subscriptions: make(map[string]struct{}, 0),
	}
	return state.sessions[id].Out, func() {
		state.eventLock.Lock()
		defer state.eventLock.Unlock()

		for topic := range state.sessions[id].subscriptions {
			delete(state.subscriptions[topic], id)
			if len(state.subscriptions[topic]) == 0 {
				delete(state.subscriptions, topic)

				delete(state.publishers, topic)
			}
		}
		delete(state.sessions, id)
	}
}

func Subscribe(sessionId string, topic string, publisher Publisher) error {
	state.eventLock.Lock()
	defer state.eventLock.Unlock()

	sess, ok := state.sessions[sessionId]
	if !ok {
		log.Println("session not found :", sessionId)
		return errors.New("session not found")
	}

	_, ok = state.publishers[topic]
	if !ok {
		err := publisher.Start()
		if err != nil {
			log.Println("failed to start topic :", topic)
			return err
		}
		state.publishers[topic] = publisher
	}

	if _, ok := state.subscriptions[topic]; !ok {
		state.subscriptions[topic] = make(map[string]*session, 0)
		state.subscriptions[topic][sess.ID] = sess
		go runSubscription(topic)
	} else {
		state.subscriptions[topic][sess.ID] = sess
	}

	state.sessions[sessionId].subscriptions[topic] = struct{}{}

	state.publishers[topic].OnNewSub()
	return nil
}

func UnSubscribe(sessionId string, topic string) {
	state.eventLock.Lock()
	defer state.eventLock.Unlock()

	_, ok := state.subscriptions[topic]
	if !ok {

		return
	}
	delete(state.subscriptions[topic], sessionId)

	delete(state.sessions[sessionId].subscriptions, topic)

	if len(state.subscriptions[topic]) == 0 {
		delete(state.subscriptions, topic)
		pub, ok := state.publishers[topic]
		if ok {
			if err := pub.Stop(); err != nil {
				log.Println("failed to stop publisher :", topic)
			}
			delete(state.publishers, topic)
		}
	}
}

func runSubscription(topic string) {
	publisher, ok := state.publishers[topic]
	if !ok {
		log.Println("publisher not found :", topic)
		return
	}
	for msg := range publisher.Read() {
		// Copy sessions under lock to avoid race conditions
		state.eventLock.Lock()
		sessions := make([]*session, 0, len(state.subscriptions[topic]))
		for _, sess := range state.subscriptions[topic] {
			sessions = append(sessions, sess)
		}
		state.eventLock.Unlock()

		subCount := 0
		// Send to sessions without holding lock
		for _, sess := range sessions {
			subCount++
			select {
			case sess.Out <- msg:
				// Message sent successfully
			default:
				// Channel full, log and skip to prevent blocking
				log.Printf("Session %s channel full, dropping message for topic %s", sess.ID, topic)
			}
		}

		if subCount == 0 {
			log.Println("no subscribers found for topic :", topic)
			log.Println("publisher stopped :", topic)
			err := publisher.Stop()
			if err != nil {
				log.Println("failed to stop publisher :", topic)
			}
			return
		}
	}
	state.eventLock.Lock()
	defer state.eventLock.Unlock()
	delete(state.publishers, topic)

	for _, sess := range state.subscriptions[topic] {
		delete(sess.subscriptions, topic)
	}

	delete(state.subscriptions, topic)
}
