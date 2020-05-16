package pingo_pongo

import (
	"errors"
	"log"
)

type GamePool []*Game

var gamePoolFull = errors.New("all games are busy")

func NewGamePool() GamePool {
	gs := make(GamePool, 5)
	for i, _ := range gs {
		gs[i] = NewGame()
		go gs[i].Run()
	}
	return gs
}

func (gs GamePool) FindOwnedGame(player *Player) *Game {
	for _, g := range gs {
		if g.UpdateOwner(player) {
			return g
		}
	}
	return nil
}

func (gs GamePool) AddPlayerToWaitingGame(player *Player) (*Game, error) {
	for _, g := range gs {
		if err := g.AddNewPlayer(player); err == nil {
			return g, nil
		}
	}

	return nil, gamePoolFull
}

func (gs GamePool) EnterTheGame(player *Player) error {
	game := gs.FindOwnedGame(player)
	if game == nil {
		g, err := gs.AddPlayerToWaitingGame(player)
		if err != nil {
			return err
		}
		game = g
	}

	ind, err := game.getPlayerIndex(player)
	if err != nil {
		return err
	}

	log.Printf("[%v]: player[%v] enter the game", game.ID, player.ID)
	return player.Send(NewUUIDAndIndexMessage(player.ID, ind))
}
