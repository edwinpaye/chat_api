package ws

import "sync"

// Hub maintains active connections and rooms (pub/sub messaging).
type Hub struct {
	mu      sync.RWMutex
	clients map[*ClientConn]bool
	rooms   map[string]map[*ClientConn]bool

	register   chan *ClientConn
	unregister chan *ClientConn
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*ClientConn]bool),
		rooms:      make(map[string]map[*ClientConn]bool),
		register:   make(chan *ClientConn),
		unregister: make(chan *ClientConn),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			h.clients[c] = true
			h.mu.Unlock()
		case c := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				for room := range c.rooms {
					h.removeFromRoom(room, c)
				}
				close(c.send)
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) Register(c *ClientConn)   { h.register <- c }
func (h *Hub) Unregister(c *ClientConn) { h.unregister <- c }

func (h *Hub) JoinRoom(room string, c *ClientConn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.rooms[room] == nil {
		h.rooms[room] = make(map[*ClientConn]bool)
	}
	h.rooms[room][c] = true
	c.rooms[room] = true
}

func (h *Hub) LeaveRoom(room string, c *ClientConn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.removeFromRoom(room, c)
}

func (h *Hub) removeFromRoom(room string, c *ClientConn) {
	if members, ok := h.rooms[room]; ok {
		delete(members, c)
		if len(members) == 0 {
			delete(h.rooms, room)
		}
	}
	delete(c.rooms, room)
}

// BroadcastTo sends a message to every client in a room.
func (h *Hub) BroadcastTo(room string, msg Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.rooms[room] {
		select {
		case c.send <- msg:
		default: // drop for slow consumers; their pump will time out
		}
	}
}

// Count returns the number of connected clients.
func (h *Hub) Count() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}