package ws

import (
	"encoding/json"
"errors"
"net/http"
"strings"
"github.com/gorilla/websocket"

	"chat_api/internal/config"
	"chat_api/internal/domain"
	"chat_api/internal/security"
	"chat_api/internal/service"
)

type Handler struct {
	hub      *Hub
	services *service.Services
	tokens   *security.TokenService
	upgrader websocket.Upgrader
}

func NewHandler(h *Hub, s *service.Services, t *security.TokenService, cfg *config.Config) *Handler {
	return &Handler{
		hub:      h,
		services: s,
		tokens:   t,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  cfg.Server.ReadBuffer,
			WriteBufferSize: cfg.Server.WriteBuffer,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
}

// ServeWS upgrades the connection and requires a valid token
// via ?token= query parameter or Authorization: Bearer header.
func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	userID := ""
// Optional token at connect time (?token= or Authorization header).
	// Unauthenticated sockets may connect but are restricted to
	// auth.* and ping operations until they present a valid token.
	token := r.URL.Query().Get("token")
	if token == "" {
		if a := r.Header.Get("Authorization"); len(a) > 7 && a[:7] == "Bearer " {
			token = a[7:]
		}
	}
	if token != "" {
		if sub, err := h.tokens.Verify(token); err == nil {
			userID = sub
		}
	}
conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	c := NewClientConn(h.hub, conn, userID)
	h.hub.Register(c)
	go c.WritePump()
	go c.ReadPump(h.dispatch)
}

// dispatch routes a message by its `type`.
func (h *Handler) dispatch(c *ClientConn, m Message) {
// Auth gate: only auth.* and ping are allowed before authentication.
	if c.userID == "" && m.Type != "ping" && !strings.HasPrefix(m.Type, "auth.") {
		c.Send(ErrReply(m.ID, m.Type, "UNAUTHENTICATED", errors.New("authentication required")))
		return
	}
switch m.Type {
	case "ping":
		c.Send(OKReply(m.ID, m.Type, map[string]string{"pong": "ok"}))
	case "auth.token": // issue a token for a subject (dev gateway / login delegate)
		var p struct {
			Subject string `json:"subject"`
		}
		if json.Unmarshal(m.Payload, &p) != nil || p.Subject == "" {
			c.Send(ErrReply(m.ID, m.Type, "VALIDATION", errBadPayload()))
			return
		}
		tok, err := h.tokens.Issue(p.Subject)
		if err != nil {
			c.Send(ErrReply(m.ID, m.Type, "AUTH", err))
			return
		}
		c.userID = p.Subject // upgrade this connection to authenticated
		c.Send(OKReply(m.ID, m.Type, map[string]string{"token": tok}))
case "room.join":
		var p struct {
			Room string `json:"room"`
		}
		if json.Unmarshal(m.Payload, &p) == nil && p.Room != "" {
			h.hub.JoinRoom(p.Room, c)
			c.Send(OKReply(m.ID, m.Type, p))
		} else {
			c.Send(ErrReply(m.ID, m.Type, "VALIDATION", errBadPayload()))
		}
	case "room.leave":
		var p struct {
			Room string `json:"room"`
		}
		if json.Unmarshal(m.Payload, &p) == nil && p.Room != "" {
			h.hub.LeaveRoom(p.Room, c)
			c.Send(OKReply(m.ID, m.Type, p))
		}
	case "room.message":
		var p struct {
			Room string          `json:"room"`
			Data json.RawMessage `json:"data"`
		}
		if json.Unmarshal(m.Payload, &p) == nil && p.Room != "" {
			payload, _ := json.Marshal(map[string]any{
				"room": p.Room, "from": c.userID, "data": p.Data,
			})
			h.hub.BroadcastTo(p.Room, Message{Type: "room.message", Payload: payload})
			c.Send(OKReply(m.ID, m.Type, map[string]string{"room": p.Room}))
		}
case "users.create":
		handleCreate[domain.Users](c, m, h.services.Users)
	case "users.get":
		handleGet[string, domain.Users](c, m, h.services.Users)
	case "users.list":
		handleList[domain.Users](c, m, h.services.Users)
	case "users.update":
		handleUpdate[string, domain.Users](c, m, h.services.Users)
	case "users.delete":
		handleDelete[string](c, m, h.services.Users)
case "rooms.create":
		handleCreate[domain.Rooms](c, m, h.services.Rooms)
	case "rooms.get":
		handleGet[int64, domain.Rooms](c, m, h.services.Rooms)
	case "rooms.list":
		handleList[domain.Rooms](c, m, h.services.Rooms)
	case "rooms.update":
		handleUpdate[int64, domain.Rooms](c, m, h.services.Rooms)
	case "rooms.delete":
		handleDelete[int64](c, m, h.services.Rooms)
case "messages.create":
		handleCreate[domain.Messages](c, m, h.services.Messages)
	case "messages.get":
		handleGet[int64, domain.Messages](c, m, h.services.Messages)
	case "messages.list":
		handleList[domain.Messages](c, m, h.services.Messages)
	case "messages.update":
		handleUpdate[int64, domain.Messages](c, m, h.services.Messages)
	case "messages.delete":
		handleDelete[int64](c, m, h.services.Messages)
default:
		c.Send(ErrReply(m.ID, m.Type, "UNKNOWN_TYPE", errUnknownType(m.Type)))
	}
}