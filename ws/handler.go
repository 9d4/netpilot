package ws

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/9d4/netpilot/database"
	p "github.com/9d4/netpilot/prefix"
	"github.com/9d4/netpilot/ros/board"
	"github.com/9d4/netpilot/worker"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

var clients = make(map[string]*websocket.Conn)

var Handler = websocket.New(func(conn *websocket.Conn) {

	cID := uuid.New()
	clients[cID.String()] = conn

	c := &Conn{
		c: conn,
	}
	go c.run()

	defer delete(clients, cID.String())

	var (
		mt  int
		msg []byte
		err error
	)

	for {
		if mt, msg, err = c.conn().ReadMessage(); err != nil {
			c.close()
			break
		}

		if mt != websocket.TextMessage {
			continue
		}

		var parsedMsg Message
		err := json.NewDecoder(bytes.NewReader(msg)).Decode(&parsedMsg)
		if err != nil {
			continue
		}

		switch parsedMsg.Type {
		case TypeGet:
			wsHandleGet(c, &parsedMsg)
		case TypeSub:
			handleSub(c, &parsedMsg)
		}
	}
})

func wsWriteError(conn *Conn, msg *Message, err *fiber.Error) {
	msg.Body = err.Error()

	e := &MessageResponseError{
		Code:    err.Code,
		Message: msg,
	}

	conn.writeJson(e)
}

func wsHandleGet(conn *Conn, msg *Message) {
	if msg.BoardID == "" || msg.Resource == "" {
		wsWriteError(conn, msg, fiber.ErrBadRequest)
		return
	}

	switch msg.Resource {
	case ResourceBoardStatus:
		board, err := worker.Boards.GetByUUID(msg.BoardID)
		if err != nil {
			wsWriteError(conn, msg, fiber.ErrNotFound)
			return
		}

		resp := MessageResponse{
			Message: Message{
				ID:       msg.ID,
				Type:     TypeResource,
				Resource: ResourceBoardStatus,
				BoardID:  msg.BoardID,
				Body:     board.Status(),
			},
		}

		conn.writeJson(resp)
	}

}

func handleSub(conn *Conn, msg *Message) {
	if msg.BoardID == "" || msg.Resource == "" {
		wsWriteError(conn, msg, fiber.ErrBadRequest)
		return
	}

	switch msg.Resource {
	case ResourceBoardStatus:
		handleSubBoardStatus(conn, msg)
	default:
		wsWriteError(conn, msg, fiber.ErrBadRequest)
	}
}

func handleSubBoardStatus(conn *Conn, msg *Message) {
	b, err := worker.Boards.GetByUUID(msg.BoardID)
	if err != nil {
		wsWriteError(conn, msg, fiber.ErrNotFound)
		return
	}

	resp := MessageResponse{
		Message: Message{
			ID:       msg.ID,
			Type:     TypeSub,
			Resource: ResourceBoardStatus,
			BoardID:  b.UUID,
		},
		Status: StatusOK,
	}
	ready := make(chan byte, 1)
	go func() {
		channel := p.BoardChannelPrefix.Status(b.UUID)
		pubsub := database.RedisCli().Subscribe(context.Background(), channel)
		defer pubsub.Close()
		ready <- 1

		for m := range pubsub.Channel() {
			if conn.closed {
				return
			}
			var s board.Status
			json.NewDecoder(bytes.NewBufferString(m.Payload)).Decode(&s)

			resp = MessageResponse{
				Message: Message{
					ID:       msg.ID,
					Type:     "resource",
					Resource: ResourceBoardStatus,
					BoardID:  b.UUID,
					Body:     s,
				},
			}
			conn.writeJson(resp)
		}
	}()

	<-ready
	conn.writeJson(resp)
}
