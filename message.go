package pingo_pongo

type MessageType string

const (
	MessageTypeConfig MessageType = "msg-config"
	MessageTypeEndGame MessageType = "msg-end-game"
	MessageTypeUUID MessageType = "msg-uuid"
	MessageTypePlayerInfo MessageType = "msg-player-info"
	MessageTypeGameStart MessageType = "msg-game-start"
)

type Message struct {
	Type MessageType `json:"type"`
	Data interface{} `json:"data"`
}

func NewPlayerInfoMessage(player *Player) *Message {
	return &Message{
		Type: MessageTypePlayerInfo,
		Data: player,
	}
}

var (
	MessageStartGame = &Message{Type: MessageTypeGameStart}
	MessageEndGame   = &Message{Type: MessageTypeEndGame}
	MessageConfig    = &Message{Type: MessageTypeConfig, Data: cfg}
)
