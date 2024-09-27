package main

import (
	"encoding/json"
	"errors"
	"log"
)

type gameState struct {
	Players []player `json:"players"`
}

type player struct {
	X           int    `json:"x"`
	Y           int    `json:"y"`
	Name        string `json:"name"`
	Health      int    `json:"health"`
	Facing      string `json:"facing"`
	IsAttacking bool   `json:"isAttacking"`
	Skin        string `json:"skin"`
}

type gameEvent struct {
	Type string `json:"type"`
	Data player `json:"data"`
}

func (gs *gameState) toJSON() []byte {
	// convert the gameState to a json byte slice
	json, err := json.Marshal(gs)
	if err != nil {
		log.Println("json marshal:", err)
	}

	return json
}

func (gs *gameState) handleEvent(event gameEvent) {
	if event.Type == "refresh" {
		// do nothing. this is just the client asking for the latest state
	}
	if event.Type == "attack" {
		// handle attack event
		gs.playerAttack(event.Data.Name)
	}
	if event.Type == "walk" {
		// handle walk event
		gs.playerWalk(event.Data.Name, event.Data.Facing)
	}
	if event.Type == "join" {
		// handle join event
		// data should be the name of the player joining
		gs.addPlayer(event.Data)
	}
	if event.Type == "leave" {
		// handle leave event
		// data should be the name of the player leaving
		gs.removePlayer(event.Data.Name)
	}
}

func (gs *gameState) removePlayer(name string) {
	// remove the player with the given name
	for i, p := range gs.Players {
		if p.Name == name {
			gs.Players = append(gs.Players[:i], gs.Players[i+1:]...)
			return
		}
	}
}

func (gs *gameState) playerAttackHit(name string) (bool, string) {
	// check if player is facing another player within 10 units of them
	// if so, return true and the name of the player they hit
	// otherwise, return false and an empty string

	player, err := gs.getPlayer(name)
	if err != nil {
		log.Println("cannot find attacking player:", err)
		return false, ""
	}

	playerFacing := player.Facing
	playerX := player.X
	playerY := player.Y

	for _, p := range gs.Players {
		if p.Name == name {
			continue
		}

		// is player within 10 units to the left, right, above or below this player?
		if p.X >= playerX-10 && p.X <= playerX+10 && p.Y >= playerY-10 && p.Y <= playerY+10 {
			// is player facing this player?
			if playerFacing == "up" && p.Y < playerY {
				return true, p.Name
			}
			if playerFacing == "down" && p.Y > playerY {
				return true, p.Name
			}
			if playerFacing == "left" && p.X < playerX {
				return true, p.Name
			}
			if playerFacing == "right" && p.X > playerX {
				return true, p.Name
			}
		}
	}

	return false, ""
}

func (gs *gameState) playerAttack(name string) {
	hit, hitName := gs.playerAttackHit(name)
	if hit {
		hitPlayer, err := gs.getPlayer(hitName)
		if err != nil {
			log.Println("cannot find player to attack")
			return
		}
		hitPlayer.Health -= 10
		gs.updatePlayer(hitPlayer)
	}
}

func (gs *gameState) playerWalk(name, direction string) {
	p, err := gs.getPlayer(name)
	if err != nil {
		log.Println("cannot find walking player")
		return
	}
	p.Facing = direction
	switch direction {
	case "up":
		p.Y--
	case "down":
		p.Y++
	case "left":
		p.X--
	case "right":
		p.X++
	}
	gs.updatePlayer(p)
}

func (gs *gameState) updatePlayer(p player) {
	for i, player := range gs.Players {
		if player.Name == p.Name {
			gs.Players[i] = p
			return
		}
	}
}

func (gs *gameState) addPlayer(p player) {
	gs.Players = append(gs.Players, p)
}

func (gs *gameState) getPlayer(name string) (player, error) {
	// get player with the given name
	for _, p := range gs.Players {
		if p.Name == name {
			return p, nil
		}
	}
	return player{}, errors.New("player not found")
}

func (gs *gameState) getPlayers() []player {
	players := []player{}
	for _, p := range gs.Players {
		players = append(players, p)
	}
	return players
}

func newGameState() gameState {
	return gameState{
		Players: []player{},
	}
}
