package pingo_pongo

type MessageType string

const (
	MessageTypeConfig       MessageType = "msg-config"
	MessageTypeEndGame      MessageType = "msg-end-game"
	MessageTypeUUID         MessageType = "msg-uuid"
	MessageTypePlayerInfo   MessageType = "msg-player-info"
	MessageTypeGameStart    MessageType = "msg-game-start"
	MessageTypeClientAction MessageType = "msg-client-action"
)

var (
	MessageStartGame = &Message{Type: MessageTypeGameStart}
	MessageEndGame   = &Message{Type: MessageTypeEndGame}
	MessageConfig    = &Message{Type: MessageTypeConfig, Data: cfg}
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

type ClientAction struct {
	DT    float64 `json:"dt"`
	Left  bool    `json:"left"`
	Right bool    `json:"right"`
	Up    bool    `json:"up"`
	Down  bool    `json:"down"`
}
