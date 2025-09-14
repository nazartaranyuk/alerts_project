package ws

import "sync"

type Payload struct {
	UserID   string  `json:"user_id"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Accuracy float64 `json:"accuracy"`
	TS       string  `json:"ts"`
}

type Hub struct {
	mu    sync.RWMutex
	rooms map[string]map[*Conn]struct{}
}

func NewHub() *Hub {
	return &Hub{rooms: make(map[string]map[*Conn]struct{})}
}

func (h *Hub) Join(room string, c *Conn) {
	h.mu.Lock()
	if _, ok := h.rooms[room]; !ok {
		h.rooms[room] = make(map[*Conn]struct{})
	}
	h.rooms[room][c] = struct{}{}
	h.mu.Unlock()
}

func (h *Hub) Leave(room string, c *Conn) {
	h.mu.Lock()
	if conns, ok := h.rooms[room]; ok {
		delete(conns, c)
		if len(conns) == 0 {
			delete(h.rooms, room)
		}
	}
	h.mu.Unlock()
}

func (h *Hub) Broadcast(room string, p Payload) {
	h.mu.RLock()
	conns := h.rooms[room]
	for c := range conns {
		c.Send(p)
	}
	h.mu.RUnlock()
}
