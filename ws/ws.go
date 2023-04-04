package ws

type MsgType string

const (
	TypeGet MsgType = "get"
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
	BoardID string `json:"board_id"`
	// message body
	Body interface{} `json:"body"`
}

type MessageResponse struct {
	Message
	Status string `json:"status"`
}

type MessageResponseError struct {
	*Message
	Code int `json:"code"`
}
