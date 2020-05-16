package pingo_pongo

import "github.com/google/uuid"

type MessageType int

const (
	MessageTypePlayersInfo MessageType = iota
	MessageTypeEndGame
	MessageTypeUUID
	MessageTypeUUIDAndIndex
	MessageTypeOpponentUUID
)

type Message struct {
	Type MessageType `json:"type"`
	Data interface{} `json:"data"`
}

func NewUUIDAndIndexMessage(id uuid.UUID, index int) *Message {
	return &Message{
		Type: MessageTypeUUIDAndIndex,
		Data: struct {
			Index int       `json:"index"`
			UUID  uuid.UUID `json:"uuid"`
		}{index, id},
	}
}

func NewOpponentMessage(id uuid.UUID) *Message {
	return &Message{
		Type: MessageTypeOpponentUUID,
		Data: id,
	}
}

var (
	MessageEndGame = &Message{Type: MessageTypeEndGame}
)
