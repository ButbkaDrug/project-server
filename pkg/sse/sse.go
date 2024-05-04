package sse

import (
	"log"
	"sync"
)

type SSE struct {
	clients map[string]*chan interface{}
	mu      sync.Mutex
}

func NewSSE() *SSE {
	return &SSE{
		clients: make(map[string]*chan interface{}),
	}
}

func (sse *SSE) AddClient(id string) *chan interface{} {
	sse.mu.Lock()
	defer sse.mu.Unlock()
	if ch, ok := sse.clients[id]; ok {
		return ch
	}

	ch := make(chan interface{})

	sse.clients[id] = &ch

	return &ch
}

func (sse *SSE) RemoveClient(id string) {
	sse.mu.Lock()
	defer sse.mu.Unlock()
	ch, ok := sse.clients[id]

	if !ok {
		return
	}

	delete(sse.clients, id)
	close(*ch)
}

func (sse *SSE) Broadcast(data interface{}) {
    sse.mu.Lock()
	for _, ch := range sse.clients {
		*ch <- data
	}
    sse.mu.Unlock()

	defer func() {
        if err := recover(); err != nil {
            log.Println("panic recovered")
            log.Println(err)
        }
	}()
}

func (sse *SSE) GetClientsId() []string {
	sse.mu.Lock()
	result := make([]string, len(sse.clients))
	i := 0
	for client := range sse.clients {
		result[i] = client
		i++
	}
	sse.mu.Unlock()

	return result
}
