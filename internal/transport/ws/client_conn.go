package ws

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 30 * time.Second * 2
	maxMsgSize = 1 << 20 // 1 MiB
)

// ClientConn wraps a websocket connection with an outbound queue.
type ClientConn struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan Message
	rooms  map[string]bool
	userID string
}

func NewClientConn(hub *Hub, conn *websocket.Conn, userID string) *ClientConn {
	return &ClientConn{
		hub:    hub,
		conn:   conn,
		send:   make(chan Message, 64),
		rooms:  make(map[string]bool),
		userID: userID,
	}
}

// ReadPump pumps messages from the socket to the handler dispatch loop.
func (c *ClientConn) ReadPump(dispatch func(*ClientConn, Message)) {
	defer func() { c.hub.Unregister(c); _ = c.conn.Close() }()
	c.conn.SetReadLimit(maxMsgSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})
	for {
		var m Message
		if err := c.conn.ReadJSON(&m); err != nil {
			return
		}
		dispatch(c, m)
	}
}

// WritePump drains the outbound queue and sends pings.
func (c *ClientConn) WritePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() { ticker.Stop(); _ = c.conn.Close() }()
	for {
		select {
		case msg, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteJSON(msg); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *ClientConn) Send(m Message) {
	select {
	case c.send <- m:
	default:
	}
}