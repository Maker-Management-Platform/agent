package events

import (
	"log"
	"sync"
)

type session struct {
	ID            string
	Out           chan *Message
	subscriptions map[string]struct{}
}

type stateType struct {
	sessions          map[string]*session
	sessionsLock      sync.Mutex
	subscriptions     map[string][]*session
	subscriptionsLock sync.Mutex
	publishers        map[string]Publisher
	publishersLock    sync.Mutex
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
		sessions:          make(map[string]*session),
		sessionsLock:      sync.Mutex{},
		subscriptions:     make(map[string][]*session),
		subscriptionsLock: sync.Mutex{},
		publishers:        make(map[string]Publisher),
		publishersLock:    sync.Mutex{},
	}
}

func RegisterSession(id string) (chan *Message, func()) {
	state.sessionsLock.Lock()
	defer state.sessionsLock.Unlock()
	state.sessions[id] = &session{
		ID:            id,
		Out:           make(chan *Message, 100),
		subscriptions: make(map[string]struct{}, 0),
	}

	return state.sessions[id].Out, func() {
		state.sessionsLock.Lock()
		state.subscriptionsLock.Lock()
		for topic := range state.sessions[id].subscriptions {
			for i, sub := range state.subscriptions[topic] {
				if sub.ID == id {
					state.subscriptions[topic] = append(state.subscriptions[topic][:i], state.subscriptions[topic][i+1:]...)
				}
				if len(state.subscriptions[topic]) == 0 {
					delete(state.subscriptions, topic)
					state.publishers[topic].Stop()
					state.publishersLock.Lock()
					delete(state.publishers, topic)
					state.publishersLock.Unlock()
				}
			}
		}
		state.subscriptionsLock.Unlock()
		delete(state.sessions, id)
		state.sessionsLock.Unlock()
	}
}

func Subscribe(sessionId string, topic string, publisher Publisher) {

	state.sessionsLock.Lock()
	sess, ok := state.sessions[sessionId]
	state.sessionsLock.Unlock()
	if !ok {
		log.Println("session not found :", sessionId)
		return
	}

	state.publishersLock.Lock()
	_, ok = state.publishers[topic]
	state.publishersLock.Unlock()
	if !ok {
		state.publishersLock.Lock()
		state.publishers[topic] = publisher
		state.publishersLock.Unlock()

		err := publisher.Start()
		if err != nil {
			log.Println("failed to start topic :", topic)
			return
		}
	}
	state.subscriptionsLock.Lock()
	defer state.subscriptionsLock.Unlock()
	if _, ok := state.subscriptions[topic]; !ok {
		state.subscriptions[topic] = make([]*session, 0)
		state.subscriptions[topic] = append(state.subscriptions[topic], sess)
		go runSubscription(topic)
	} else {
		state.subscriptions[topic] = append(state.subscriptions[topic], sess)
	}
	state.sessionsLock.Lock()
	state.sessions[sessionId].subscriptions[topic] = struct{}{}
	state.sessionsLock.Unlock()
	state.publishers[topic].OnNewSub()

}

func UnSubscribe(sessionId string, topic string) {
	state.subscriptionsLock.Lock()
	sessions, ok := state.subscriptions[topic]
	if !ok {
		state.subscriptionsLock.Unlock()
		return
	}
	for i, sess := range sessions {
		if sess.ID == sessionId {
			sessions = append(sessions[:i], sessions[i+1:]...)
			state.sessionsLock.Lock()
			delete(state.sessions[sessionId].subscriptions, topic)
			state.sessionsLock.Unlock()
			break
		}
	}
	state.subscriptions[topic] = sessions
	if len(sessions) == 0 {
		delete(state.subscriptions, topic)
		state.publishersLock.Lock()
		pub, ok := state.publishers[topic]
		if ok {
			if err := pub.Stop(); err != nil {
				log.Println("failed to stop publisher :", topic)
			}
			delete(state.publishers, topic)
		}
		state.publishersLock.Unlock()
	}
	state.subscriptionsLock.Unlock()
}

func runSubscription(topic string) {
	publisher, ok := state.publishers[topic]
	if !ok {
		log.Println("publisher not found :", topic)
		return
	}
	for msg := range publisher.Read() {
		subCount := 0
		state.sessionsLock.Lock()
		for _, sess := range state.subscriptions[topic] {
			subCount++
			sess.Out <- msg

		}
		state.sessionsLock.Unlock()
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
	state.publishersLock.Lock()
	delete(state.publishers, topic)
	state.publishersLock.Unlock()

	state.subscriptionsLock.Lock()

	state.sessionsLock.Lock()
	for _, sess := range state.subscriptions[topic] {
		delete(sess.subscriptions, topic)
	}
	state.sessionsLock.Unlock()

	delete(state.subscriptions, topic)
	state.subscriptionsLock.Unlock()
}
