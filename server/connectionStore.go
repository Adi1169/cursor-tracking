package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

// ConnectionStore is a struct that stores all the websocket connections
// based on the channel they belong to in a map
// concurrently safe operations are supported
type ConnectionStore struct {
	sync.RWMutex
	clients map[string]map[*websocket.Conn]bool
}

func NewConnectionStore() *ConnectionStore {
	return &ConnectionStore{sync.RWMutex{}, make(map[string]map[*websocket.Conn]bool)}
}

func (cs *ConnectionStore) Set(key string, conn *websocket.Conn) bool {
	cs.Lock()
	var exists bool
	if _, exists = cs.clients[key]; !exists {
		cs.clients[key] = make(map[*websocket.Conn]bool)
	}
	cs.clients[key][conn] = true
	cs.Unlock()
	return exists
}

func (cs *ConnectionStore) Get(key string) (map[*websocket.Conn]bool, bool) {
	cs.RLock()
	value, exists := cs.clients[key]
	cs.RUnlock()
	return value, exists
}

func (cs *ConnectionStore) Delete(key string, conn *websocket.Conn) {
	cs.Lock()
	if _, ok := cs.clients[key]; ok {
		delete(cs.clients[key], conn)
		if len(cs.clients[key]) == 0 {
			delete(cs.clients, key)
		}
	}
	cs.Unlock()
}

func (cs *ConnectionStore) Iterate(f func(key string, value map[*websocket.Conn]bool)) {
	cs.RLock()
	for key, value := range cs.clients {
		f(key, value)
	}
	cs.RUnlock()
}
