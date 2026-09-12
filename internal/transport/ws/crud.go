package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// ---- generic CRUD dispatch helpers (typed per entity) ----

type creator[T any] interface {
	Create(context.Context, json.RawMessage) (*T, error)
}
type getter[ID any, T any] interface {
	Get(context.Context, ID) (*T, error)
}
type lister[T any] interface {
	List(context.Context, int, int) ([]T, error)
}
type updater[ID any, T any] interface {
	Update(context.Context, ID, json.RawMessage) (*T, error)
}
type deleter[ID any] interface {
	Delete(context.Context, ID) error
}

func msgCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func errBadPayload() error { return fmt.Errorf("invalid payload") }
func errUnknownType(t string) error {
	return fmt.Errorf("unknown message type %q", t)
}

func handleCreate[T any](c *ClientConn, m Message, svc creator[T]) {
	ctx, cancel := msgCtx()
	defer cancel()
	e, err := svc.Create(ctx, m.Payload)
	if err != nil {
		c.Send(ErrReply(m.ID, m.Type, "CREATE", err))
		return
	}
	c.Send(OKReply(m.ID, m.Type, e))
}

func handleGet[ID any, T any](c *ClientConn, m Message, svc getter[ID, T]) {
	var p struct {
		ID ID `json:"id"`
	}
	if json.Unmarshal(m.Payload, &p) != nil {
		c.Send(ErrReply(m.ID, m.Type, "VALIDATION", errBadPayload()))
		return
	}
	ctx, cancel := msgCtx()
	defer cancel()
	e, err := svc.Get(ctx, p.ID)
	if err != nil {
		c.Send(ErrReply(m.ID, m.Type, "GET", err))
		return
	}
	c.Send(OKReply(m.ID, m.Type, e))
}

func handleList[T any](c *ClientConn, m Message, svc lister[T]) {
	var p struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	}
	_ = json.Unmarshal(m.Payload, &p)
	ctx, cancel := msgCtx()
	defer cancel()
	items, err := svc.List(ctx, p.Limit, p.Offset)
	if err != nil {
		c.Send(ErrReply(m.ID, m.Type, "LIST", err))
		return
	}
	c.Send(OKReply(m.ID, m.Type, items))
}

func handleUpdate[ID any, T any](c *ClientConn, m Message, svc updater[ID, T]) {
	var p struct {
		ID   ID              `json:"id"`
		Data json.RawMessage `json:"data"`
	}
	if json.Unmarshal(m.Payload, &p) != nil || len(p.Data) == 0 {
		c.Send(ErrReply(m.ID, m.Type, "VALIDATION", errBadPayload()))
		return
	}
	ctx, cancel := msgCtx()
	defer cancel()
	e, err := svc.Update(ctx, p.ID, p.Data)
	if err != nil {
		c.Send(ErrReply(m.ID, m.Type, "UPDATE", err))
		return
	}
	c.Send(OKReply(m.ID, m.Type, e))
}

func handleDelete[ID any](c *ClientConn, m Message, svc deleter[ID]) {
	var p struct {
		ID ID `json:"id"`
	}
	if json.Unmarshal(m.Payload, &p) != nil {
		c.Send(ErrReply(m.ID, m.Type, "VALIDATION", errBadPayload()))
		return
	}
	ctx, cancel := msgCtx()
	defer cancel()
	if err := svc.Delete(ctx, p.ID); err != nil {
		c.Send(ErrReply(m.ID, m.Type, "DELETE", err))
		return
	}
	c.Send(OKReply(m.ID, m.Type, map[string]bool{"done": true}))
}