package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var testPlayer1FacingRight = player{
	X:           0,
	Y:           0,
	Name:        "player1",
	Health:      100,
	Facing:      "right",
	IsAttacking: false,
	Skin:        "skin1",
}
var testPlayer1FacingLeft = player{
	X:           0,
	Y:           0,
	Name:        "player1",
	Health:      100,
	Facing:      "left",
	IsAttacking: false,
	Skin:        "skin1",
}
var testPlayer1FacingUp = player{
	X:           0,
	Y:           0,
	Name:        "player1",
	Health:      100,
	Facing:      "up",
	IsAttacking: false,
	Skin:        "skin1",
}
var testPlayer1FacingDown = player{
	X:           0,
	Y:           0,
	Name:        "player1",
	Health:      100,
	Facing:      "down",
	IsAttacking: false,
	Skin:        "skin1",
}

func TestGameState(t *testing.T) {
	// Create a new gameState
	gs := gameState{
		Players: make(map[string]player),
	}

	// Add a player
	gs.Players["player1"] = testPlayer1FacingRight

	// Add another player
	gs.Players["player2"] = player{
		X:           10,
		Y:           0,
		Name:        "player2",
		Health:      100,
		Facing:      "down",
		IsAttacking: false,
		Skin:        "skin2",
	}

	// Test playerAttackHit
	hit, hitName := gs.playerAttackHit("player1")
	require.Equal(t, true, hit)
	require.Equal(t, "player2", hitName)

	// Test playerAttack
	gs.playerAttack("player1")
	require.Equal(t, 90, gs.Players["player2"].Health)
}

var playerAttackHitTestCases = []struct {
	Name     string
	Players  map[string]player
	expected bool
}{
	{
		Name: "player1HitPlayerRight10UnitsAway",
		Players: map[string]player{
			"player1": testPlayer1FacingRight,
			"player2": {
				X:           10,
				Y:           0,
				Name:        "player2",
				Health:      100,
				Facing:      "down",
				IsAttacking: false,
				Skin:        "skin2",
			},
		},
		expected: true,
	},
	{
		Name: "player1HitPlayerRight11UnitsAway",
		Players: map[string]player{
			"player1": testPlayer1FacingRight,
			"player2": {
				X:           11,
				Y:           0,
				Name:        "player2",
				Health:      100,
				Facing:      "down",
				IsAttacking: false,
				Skin:        "skin2",
			},
		},
		expected: false,
	},
	{
		Name: "player1HitPlayerLeft10UnitsAway",
		Players: map[string]player{
			"player1": testPlayer1FacingLeft,
			"player2": {
				X:           -10,
				Y:           0,
				Name:        "player2",
				Health:      100,
				Facing:      "down",
				IsAttacking: false,
				Skin:        "skin2",
			},
		},
		expected: true,
	},
	{
		Name: "player1HitPlayerLeft11UnitsAway",
		Players: map[string]player{
			"player1": testPlayer1FacingLeft,
			"player2": {
				X:           -11,
				Y:           0,
				Name:        "player2",
				Health:      100,
				Facing:      "down",
				IsAttacking: false,
				Skin:        "skin2",
			},
		},
		expected: false,
	},
	{
		Name: "player1HitPlayerUp10UnitsAway",
		Players: map[string]player{
			"player1": testPlayer1FacingUp,
			"player2": {
				X:           0,
				Y:           -10,
				Name:        "player2",
				Health:      100,
				Facing:      "down",
				IsAttacking: false,
				Skin:        "skin2",
			},
		},
		expected: true,
	},
	{
		Name: "player1HitPlayerUp11UnitsAway",
		Players: map[string]player{
			"player1": testPlayer1FacingUp,
			"player2": {
				X:           0,
				Y:           -11,
				Name:        "player2",
				Health:      100,
				Facing:      "down",
				IsAttacking: false,
				Skin:        "skin2",
			},
		},
		expected: false,
	},
	{
		Name: "player1HitPlayerDown10UnitsAway",
		Players: map[string]player{
			"player1": testPlayer1FacingDown,
			"player2": {
				X:           0,
				Y:           10,
				Name:        "player2",
				Health:      100,
				Facing:      "down",
				IsAttacking: false,
				Skin:        "skin2",
			},
		},
		expected: true,
	},
	{
		Name: "player1HitPlayerDown11UnitsAway",
		Players: map[string]player{
			"player1": testPlayer1FacingDown,
			"player2": {
				X:           0,
				Y:           11,
				Name:        "player2",
				Health:      100,
				Facing:      "down",
				IsAttacking: false,
				Skin:        "skin2",
			},
		},
		expected: false,
	},
	{
		Name: "player1HitPlayerRight10UnitsAwayToTheLeft",
		Players: map[string]player{
			"player1": testPlayer1FacingRight,
			"player2": {
				X:           -10,
				Y:           0,
				Name:        "player2",
				Health:      100,
				Facing:      "left",
				IsAttacking: false,
				Skin:        "skin2",
			},
		},
		expected: false,
	},
}

// loop through test cases
// for each test case, create a new gameState with the players
// call playerAttackHit with the player1 name
// check if the result is equal to the expected result
func TestPlayerAttackHit(t *testing.T) {
	for _, tc := range playerAttackHitTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			gs := gameState{
				Players: tc.Players,
			}

			hit, _ := gs.playerAttackHit("player1")
			require.Equal(t, tc.expected, hit)
		})
	}
}
