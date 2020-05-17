package pingo_pongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"sync"
	"time"
)

var (
	gameIsFull       = errors.New("game is full")
	gameHasPlayer    = errors.New("player already in game")
	opponentNotFound = errors.New("could not find opponent")
)

type Game struct {
	ID      uuid.UUID
	Player1 *Player
	Player2 *Player

	mutex sync.Mutex
}

func NewGame() *Game {
	return &Game{
		ID: uuid.New(),
	}
}

func (g *Game) errorWrap(err error) error {
	if err != nil {
		return fmt.Errorf("[%v]: %w", g.ID, err)
	}
	return nil
}

func (g *Game) getPlayerIndex(p *Player) (int, error) {
	if g.Player1 != nil && g.Player1.ID == p.ID {
		return 0, nil
	}
	if g.Player2 != nil && g.Player2.ID == p.ID {
		return 1, nil
	}

	return 0, fmt.Errorf("unknown player %v", p.ID)
}

func (g *Game) updateOwner(p *Player) {
	log.Printf("update: %v %v %v", p.ID, g.Player1, g.Player2)
	if g.Player1 != nil && g.Player1.ID == p.ID {
		g.Player1.conn.Close()
		g.Player1 = p
	} else if g.Player2 != nil && g.Player2.ID == p.ID {
		g.Player2.conn.Close()
		g.Player2 = p
	}
}

func (g *Game) isWaiting() bool {
	return g.Player1 == nil || g.Player2 == nil
}

func (g *Game) isOwner(player *Player) bool {
	return (g.Player1 != nil && g.Player1.ID == player.ID) ||
		(g.Player2 != nil && g.Player2.ID == player.ID)
}

func (g *Game) addPlayer(player *Player) error {
	if g.Player1 == nil {
		g.Player1 = player
		return nil
	} else if g.Player2 == nil {
		g.Player2 = player
		return nil
	}

	if g.Player1 == player || g.Player2 == player {
		return gameHasPlayer
	}

	return gameIsFull
}

func (g *Game) IsWaiting() bool {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	return g.isWaiting()
}

func (g *Game) IsOwner(player *Player) bool {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	return g.isOwner(player)
}

func (g *Game) AddNewPlayer(player *Player) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	return g.addPlayer(player)
}

func (g *Game) UpdateOwner(player *Player) (isOwner bool) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	isOwner = g.isOwner(player)
	if isOwner {
		g.updateOwner(player)
	}

	return
}

func (g *Game) pingPlayer(ctx context.Context, cancel context.CancelFunc, p *Player) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(1 * time.Second):
			if p == nil {
				continue
			}
			if err := p.Ping(); err != nil {
				log.Println(g.errorWrap(err))
				cancel()
				return
			}
		}
	}
}

func (g *Game) broadcast(msg *Message) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	fu := func(p *Player) {
		defer wg.Done()
		if err := p.Send(msg); err != nil {
			log.Println(g.errorWrap(err))
		}
	}

	go fu(g.Player1)
	go fu(g.Player2)
	wg.Wait()
}

func (g *Game) Run() {
	for {
		ctx, cancel := context.WithCancel(context.Background())
		if err := g.init(ctx, cancel); err != nil {
			continue
		}
		g.process(ctx, cancel)
		g.finish()
	}
}

func (g *Game) process(ctx context.Context, cancel context.CancelFunc) {
	g.broadcast(MessageStartGame)
	log.Printf("[%v]: game start\n", g.ID)

	onReceive := func(p *Player, o *Player) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			var msg Message
			if err := p.Receive(&msg); err != nil {
				log.Println(g.errorWrap(err))
				cancel()
				return
			}

			if msg.Type == MessageTypeClientAction {
				var ca ClientAction
				data, ok := msg.Data.(map[string]interface{})
				if !ok {
					log.Printf("error parsing data: %#v\n", msg)
					cancel()
					return
				}

				ca.Down = data["down"].(bool)
				ca.Up = data["up"].(bool)
				ca.Left = data["left"].(bool)
				ca.Right = data["right"].(bool)
				ca.DT = data["dt"].(float64)

				dx := int(ca.DT * 200)
				dy := int(ca.DT * 200)

				if ca.Left {
					if p.Position.X > p.Radius {
						p.Position.X -= dx
					} else {
						p.Position.X = p.Radius
					}
				}
				if ca.Up {
					if p.Position.Y > p.Radius {
						p.Position.Y -= dy
					} else {
						p.Position.Y = p.Radius
					}
				}
				if ca.Right {
					if p.Position.X < cfg.Width - p.Radius {
						p.Position.X += dx
					} else {
						p.Position.X = cfg.Width - p.Radius
					}
				}
				if ca.Down {
					if p.Position.Y < cfg.Height - p.Radius {
						p.Position.Y += dy
					} else {
						p.Position.Y = cfg.Height - p.Radius
					}
				}

				g.broadcast(NewPlayerInfoMessage(p))
			} else {
				if err := o.Send(NewPlayerInfoMessage(p)); err != nil {
					log.Println(g.errorWrap(err))
					cancel()
					return
				}
			}
		}
	}
	go onReceive(g.Player1, g.Player2)
	go onReceive(g.Player2, g.Player1)

	<-ctx.Done()
	if err := ctx.Err(); err != nil {
		log.Println(g.errorWrap(err))
	}
}

func (g *Game) init(ctx context.Context, cancel context.CancelFunc) error {
	g.mutex.Lock()
	g.Player1 = nil
	g.Player2 = nil
	g.mutex.Unlock()

	go g.pingPlayer(ctx, cancel, g.Player1)
	go g.pingPlayer(ctx, cancel, g.Player2)

	for g.IsWaiting() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(1 * time.Second):
		}
	}

	var wg sync.WaitGroup
	wg.Add(2)

	sendOpponent := func(p *Player, op *Player) {
		defer wg.Done()
		if err := p.Send(NewPlayerInfoMessage(op)); err != nil {
			log.Println(g.errorWrap(err))
			cancel()
		}
	}

	go sendOpponent(g.Player1, g.Player2)
	go sendOpponent(g.Player2, g.Player1)

	wg.Wait()

	return ctx.Err()
}

func (g *Game) finish() {
	g.broadcast(MessageEndGame)
	g.Player1.Close()
	g.Player2.Close()
}
