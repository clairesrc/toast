package main

import (
	"encoding/json"
	"errors"
	"log"
	"time"
)

const playerSpriteHeight = 48
const playerSpriteWidth = 48

const playerHitboxHeight = 12
const playerHitboxWidth = 24

const playerDodgeDistance = 24

type gameState struct {
	Players []player `json:"players"`
}

type player struct {
	X           int    `json:"x"`
	Y           int    `json:"y"`
	Name        string `json:"name"`
	Health      int    `json:"health"`
	Stamina     int    `json:"stamina"`
	Facing      string `json:"facing"`
	IsAttacking bool   `json:"isAttacking"`
	IsWalking   bool   `json:"isWalking"`
	IsDodging   bool   `json:"isDodging"`
	lastAttack  int64
	lastWalk    int64
	lastDodge   int64
	Skin        string `json:"skin"`
}

type playerBoundingBox struct {
	hitbox boundingBox `json:"hitbox"`
	sprite boundingBox `json:"sprite"`
}

type boundingBox struct {
	x      int `json:"x"`
	y      int `json:"y"`
	width  int `json:"width"`
	height int `json:"height"`
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

		// restore stamina gradually if below 100
		if p.Stamina < 100 {
			p.Stamina += 1
			gs.Players[i] = p
		}
	}
}

func (gs *gameState) consumePlayerStamina(p player, staminaAmount int) {
	// consume stamina from the player
	p.Stamina -= staminaAmount
	if p.Stamina < 0 {
		p.Stamina = 0
	}
	gs.updatePlayer(p)
}

func (gs *gameState) playerHasStamina(p player, staminaAmount int) bool {
	// check if the player has enough stamina
	return p.Stamina >= staminaAmount
}

func (gs *gameState) removePlayer(name string) {
	// remove the player with the given name
	newPlayers := []player{}
	for _, p := range gs.Players {
		if p.Name != name {
			newPlayers = append(newPlayers, p)
		}
	}
	gs.Players = newPlayers
}

func (gs *gameState) playerDodge(name string) {
	p, err := gs.getPlayer(name)
	if err != nil {
		log.Println("cannot find dodging player")
		return
	}

	// check if the player has enough stamina to dodge
	if !gs.playerHasStamina(p, 30) {
		return
	}

	p.IsDodging = true
	p.lastDodge = time.Now().UnixMilli()
	gs.updatePlayer(p)

	// consume stamina
	gs.consumePlayerStamina(p, 30)

	// dodge roll should advance player 10 units in the direction they are facing
	switch p.Facing {
	case "up":
		gs.movePlayer(name, p.X, p.Y-playerDodgeDistance)
	case "down":
		gs.movePlayer(name, p.X, p.Y+playerDodgeDistance)
	case "left":
		gs.movePlayer(name, p.X-playerDodgeDistance, p.Y)
	case "right":
		gs.movePlayer(name, p.X+playerDodgeDistance, p.Y)
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
	p, err := gs.getPlayer(name)
	if err != nil {
		log.Println("cannot find attacking player")
		return
	}

	// check if the player has enough stamina to attack
	if !gs.playerHasStamina(p, 25) {
		return
	}

	// set the attacking player to be attacking
	p.IsAttacking = true
	p.lastAttack = time.Now().UnixMilli()
	gs.updatePlayer(p)

	// consume stamina
	gs.consumePlayerStamina(p, 25)

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
		otherPlayerBoundingBox := getPlayerBoundingBox(player)
		otherPlayerHitboxX := otherPlayerBoundingBox.hitbox.x
		otherPlayerHitboxY := otherPlayerBoundingBox.hitbox.y

		if player.Name == name {
			continue
		}
		if x+playerHitboxWidth >= otherPlayerHitboxX && x <= otherPlayerHitboxX+playerHitboxWidth && y+playerHitboxHeight >= otherPlayerHitboxY && y <= otherPlayerHitboxY+playerHitboxHeight {
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

func getPlayerBoundingBox(p player) playerBoundingBox {
	// figure out hitbox coordinates based on sprite position and dimensions
	// hitbox is a rectangle with the same center as the sprite
	const hitboxWidth = playerHitboxWidth
	const hitboxHeight = playerHitboxHeight
	const hitboxOffsetX = (playerSpriteWidth - hitboxWidth) / 2
	const hitboxOffsetY = (playerSpriteHeight - hitboxHeight) / 2
	hitboxTopLeftX := p.X + hitboxOffsetX
	hitboxTopLeftY := p.Y + hitboxOffsetY

	return playerBoundingBox{
		hitbox: boundingBox{
			x:      hitboxTopLeftX,
			y:      hitboxTopLeftY,
			width:  hitboxWidth,
			height: hitboxHeight,
		},
		sprite: boundingBox{
			x:      p.X,
			y:      p.Y,
			width:  playerSpriteWidth,
			height: playerSpriteHeight,
		},
	}
}
