package ws

import (
	"github.com/gofiber/websocket/v2"
	"sync"
)

type MsgType string
type MsgStatus string

const (
	TypeGet      MsgType = "get"
	TypeSub      MsgType = "sub"
	TypeUnsub    MsgType = "unsub"
	TypeResource MsgType = "resource"
)

const (
	StatusOK MsgStatus = "ok"
)

type ResourceType string

const (
	ResourceBoardStatus ResourceType = "board:status"
)

type Message struct {
	ID   int     `json:"id"`
	Type MsgType `json:"type"`
	// requested resource
	Resource ResourceType `json:"resource"`
	// board where resource will be taken from
	BoardID string `json:"board_id,omitempty"`
	// message body
	Body interface{} `json:"body,omitempty"`
}

type MessageResponse struct {
	Message
	Status MsgStatus `json:"status,omitempty"`
}

type MessageResponseError struct {
	*Message
	Code int `json:"code"`
}

type Conn struct {
	c  *websocket.Conn
	mu sync.Mutex

	jsonChan chan interface{}
	closed   bool
}

func (c *Conn) conn() *websocket.Conn {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.c
}

func (c *Conn) run() {
	c.jsonChan = make(chan interface{})
	for j := range c.jsonChan {
		c.c.WriteJSON(j)
	}
}

func (c *Conn) close() {
	c.closed = true
	close(c.jsonChan)
	c.conn().Close()
}

func (c *Conn) writeJson(v interface{}) {
	if c.closed {
		return
	}
	c.jsonChan <- v
}
