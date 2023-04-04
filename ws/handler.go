package ws

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/9d4/netpilot/store"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var clients = make(map[string]*websocket.Conn)

var Handler = websocket.New(func(conn *websocket.Conn) {

	cID := uuid.New()
	clients[cID.String()] = conn

	defer conn.Close()
	defer delete(clients, cID.String())

	var (
		mt  int
		msg []byte
		err error
	)

	for {
		if mt, msg, err = conn.ReadMessage(); err != nil {
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
			wsHandleGet(conn, &parsedMsg)
		}
	}
})

func wsWriteError(conn *websocket.Conn, msg *Message, err *fiber.Error) {
	msg.Body = err.Error()

	e := &MessageResponseError{
		Code:    err.Code,
		Message: msg,
	}

	conn.WriteJSON(e)
}

func wsHandleGet(conn *websocket.Conn, msg *Message) {
	if msg.BoardID == "" || msg.Resource == "" {
		wsWriteError(conn, msg, fiber.ErrBadRequest)
		return
	}

	switch msg.Resource {
	case ResourceBoardStatus:
		board, err := store.Board.FindByUUID(msg.BoardID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				wsWriteError(conn, msg, fiber.ErrNotFound)
			}
			return
		}

		resp := MessageResponse{
			Message: Message{
				ID:       msg.ID,
				Type:     "resource",
				Resource: ResourceBoardStatus,
				BoardID:  msg.BoardID,
				Body:     board.Status(),
			},
		}

		fmt.Printf("%+v\n", resp)
		fmt.Printf("%+v\n", resp.Message.Body)

		conn.WriteJSON(resp)
	}

}
