package pingo_pongo

import (
	"fmt"
	"github.com/google/uuid"
)

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}
type Player struct {
	conn *serverWS

	ID       uuid.UUID `json:"uuid"`
	Index    int       `json:"index"`
	Position Point     `json:"pos"`
	Radius   int       `json:"radius"`
}

func NewPlayer(conn *serverWS) (*Player, error) {
	p := &Player{conn: conn, Radius: 20}
	m := &Message{}

	if err := p.Receive(m); err != nil {
		return nil, err
	}

	if m.Type == MessageTypeUUID {
		if m.Data != nil {
			id, err := uuid.Parse(fmt.Sprintf("%v", m.Data))
			if err != nil {
				return nil, fmt.Errorf("data: %v, %v", m.Data, err)
			}
			p.ID = id
		} else {
			p.ID = uuid.New()
		}
	}

	return p, nil
}

func (p *Player) Init(ind int) {
	p.Index = ind
	if ind == 0 {
		p.Position = positionFirst
	} else {
		p.Position = positionSecond
	}
}

func (p *Player) errorWrap(err error) error {
	if err != nil {
		return fmt.Errorf("[%v]: %w", p.ID, err)
	}
	return nil
}

func (p *Player) Send(msg *Message) error {
	return p.errorWrap(p.conn.send(msg))
}

func (p *Player) Receive(msg *Message) error {
	return p.errorWrap(p.conn.receive(msg))
}

func (p *Player) Ping() error {
	return p.errorWrap(p.conn.ping())
}

func (p *Player) Close() {
	if err := p.conn.Close(); err != nil {
	}
}
