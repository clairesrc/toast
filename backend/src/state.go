package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
)

const playerSpriteHeight = 48
const playerSpriteWidth = 48

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
	IsWalking   bool   `json:"isWalking"`
	IsDodging   bool   `json:"isDodging"`
	lastAttack  int64
	lastWalk    int64
	lastDodge   int64
	Skin        string `json:"skin"`
}

type gameEvent struct {
	Type string `json:"type"`
	Data player `json:"data"`
}

type attackHitbox struct {
	topLeftCornerX     int
	topLeftCornerY     int
	bottomRightCornerX int
	bottomRightCornerY int
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
		// handle refresh event
		gs.refresh()
	}
	if event.Type == "attack" {
		// handle attack event
		gs.playerAttack(event.Data.Name)
	}
	if event.Type == "walk" {
		// handle walk event
		gs.playerWalk(event.Data.Name, event.Data.Facing)
	}
	if event.Type == "dodge" {
		// handle dodge event
		gs.playerDodge(event.Data.Name)
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

func (gs *gameState) refresh() {
	// refresh the game state
	for i, p := range gs.Players {
		if p.IsAttacking && time.Now().UnixMilli()-p.lastAttack > 400 {
			fmt.Println("clearing attack for player", p.Name)
			p.IsAttacking = false
			gs.Players[i] = p
		}
		if p.IsWalking && time.Now().UnixMilli()-p.lastWalk > 250 {
			p.IsWalking = false
			gs.Players[i] = p
		}
		if p.IsDodging && time.Now().UnixMilli()-p.lastDodge > 300 {
			p.IsDodging = false
			gs.Players[i] = p
		}
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

func (gs *gameState) playerDodge(name string) {
	// set the player to be dodging
	p, err := gs.getPlayer(name)
	if err != nil {
		log.Println("cannot find dodging player")
		return
	}
	p.IsDodging = true
	p.lastDodge = time.Now().UnixMilli()
	gs.updatePlayer(p)

	// dodge roll should advance player 10 units in the direction they are facing
	switch p.Facing {
	case "up":
		gs.movePlayer(name, p.X, p.Y-10)
	case "down":
		gs.movePlayer(name, p.X, p.Y+10)
	case "left":
		gs.movePlayer(name, p.X-10, p.Y)
	case "right":
		gs.movePlayer(name, p.X+10, p.Y)
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
		if p.IsDodging {
			continue
		}

		hitbox := attackHitbox{}
		switch playerFacing {
		case "up":
			hitbox.topLeftCornerX = playerX
			hitbox.topLeftCornerY = playerY - 10
			hitbox.bottomRightCornerX = playerX + playerSpriteWidth
			hitbox.bottomRightCornerY = playerY
		case "down":
			hitbox.topLeftCornerX = playerX
			hitbox.topLeftCornerY = playerY + playerSpriteHeight
			hitbox.bottomRightCornerX = playerX + playerSpriteWidth
			hitbox.bottomRightCornerY = playerY + playerSpriteHeight + 10
		case "left":
			hitbox.topLeftCornerX = playerX - 10
			hitbox.topLeftCornerY = playerY
			hitbox.bottomRightCornerX = playerX
			hitbox.bottomRightCornerY = playerY + playerSpriteHeight
		case "right":
			hitbox.topLeftCornerX = playerX + playerSpriteWidth
			hitbox.topLeftCornerY = playerY
			hitbox.bottomRightCornerX = playerX + playerSpriteWidth + 10
			hitbox.bottomRightCornerY = playerY + playerSpriteHeight
		}

		if p.X+playerSpriteWidth >= hitbox.topLeftCornerX && p.X <= hitbox.bottomRightCornerX && p.Y+playerSpriteHeight >= hitbox.topLeftCornerY && p.Y <= hitbox.bottomRightCornerY {
			return true, p.Name
		}

	}

	return false, ""
}

func (gs *gameState) playerAttack(name string) {
	// set the attacking player to be attacking
	p, err := gs.getPlayer(name)
	if err != nil {
		log.Println("cannot find attacking player")
		return
	}
	p.IsAttacking = true
	p.lastAttack = time.Now().UnixMilli()
	gs.updatePlayer(p)

	// apply damage if another player was hit
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

	p.IsWalking = true
	p.lastWalk = time.Now().UnixMilli()
	gs.updatePlayer(p)

	switch direction {
	case "up":
		gs.movePlayer(name, p.X, p.Y-2)
	case "down":
		gs.movePlayer(name, p.X, p.Y+2)
	case "left":
		gs.movePlayer(name, p.X-2, p.Y)
	case "right":
		gs.movePlayer(name, p.X+2, p.Y)
	}
}

func (gs *gameState) movePlayer(name string, x, y int) {
	p, err := gs.getPlayer(name)
	if err != nil {
		log.Println("cannot find moving player")
		return
	}
	// don't allow player to collide with other players' bounding box (taking into account sprite dimensions)
	for _, player := range gs.Players {
		if player.Name == name {
			continue
		}
		if x+playerSpriteWidth >= player.X && x <= player.X+playerSpriteWidth && y+playerSpriteHeight >= player.Y && y <= player.Y+playerSpriteHeight {
			return
		}
	}

	// no collision, so move player
	p.X = x
	p.Y = y
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
