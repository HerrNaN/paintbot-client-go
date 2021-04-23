package maputility

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"paintbot-client/models"
)

func TestMapUtility_convertPositionToCoordinates(t *testing.T) {
	mu := MapUtility{
		mapp: models.Map{Width: 5},
	}

	c := mu.ConvertPositionToCoordinates(0)
	if c.X != 0 || c.Y != 0 {
		t.Fail()
	}

	c = mu.ConvertPositionToCoordinates(5)
	if c.X != 0 || c.Y != 1 {
		t.Fail()
	}

	c = mu.ConvertPositionToCoordinates(6)
	if c.X != 1 || c.Y != 1 {
		t.Fail()
	}
}

func TestMapUtility_convertCoordinatesToPosition(t *testing.T) {
	mu := MapUtility{
		mapp: models.Map{Width: 5},
	}
	pos := mu.ConvertCoordinatesToPosition(models.Coordinates{X: 0, Y: 0})
	if pos != 0 {
		t.Fail()
	}

	pos = mu.ConvertCoordinatesToPosition(models.Coordinates{X: 1, Y: 1})
	if pos != 6 {
		t.Fail()
	}
}

func TestMapUtility_coordOutOfBounds(t *testing.T) {
	mu := MapUtility{
		mapp: models.Map{Width: 5, Height: 5},
	}

	if mu.IsCoordinatesOutOfBounds(models.Coordinates{X: 1, Y: 1}) {
		t.Fail()
	}

	if !mu.IsCoordinatesOutOfBounds(models.Coordinates{X: -1, Y: 1}) {
		t.Fail()
	}

	if !mu.IsCoordinatesOutOfBounds(models.Coordinates{X: 6, Y: 1}) {
		t.Fail()
	}
}

func TestMapUtility_CanIMove(t *testing.T) {
	mu := MapUtility{
		mapp: models.Map{
			Width:  3,
			Height: 3,
			CharacterInfos: []models.CharacterInfo{{
				Position: 0,
				ID:       "myId",
			}},
		},
		currentPlayerID: "myId",
	}

	if !mu.CanIMoveInDirection(models.Right) {
		t.Fail()
	}

	if mu.CanIMoveInDirection(models.Left) {
		t.Fail()
	}

	mu = MapUtility{
		mapp: models.Map{
			Width:  3,
			Height: 3,
			CharacterInfos: []models.CharacterInfo{{
				Position: 8,
				ID:       "myId",
			}},
		},
		currentPlayerID: "myId",
	}
	if !mu.CanIMoveInDirection(models.Left) {
		t.Fail()
	}
	if mu.CanIMoveInDirection(models.Down) {
		t.Fail()
	}
	if mu.CanIMoveInDirection(models.Right) {
		t.Fail()
	}
	if mu.CanIMoveInDirection(models.Explode) {
		t.Fail()
	}

	mu = MapUtility{
		mapp: models.Map{
			Width:  3,
			Height: 3,
			CharacterInfos: []models.CharacterInfo{{
				Position:        5,
				CarryingPowerUp: true,
				ID:              "myId",
			}},
		},
		currentPlayerID: "myId",
	}
	if !mu.CanIMoveInDirection(models.Explode) {
		t.Fail()
	}

	mu = MapUtility{
		mapp: models.Map{
			Width:  3,
			Height: 3,
			CharacterInfos: []models.CharacterInfo{{
				Position: 5,
				ID:       "myId",
			}},
			ObstacleUpPositions: []int{4},
		},
		currentPlayerID: "myId",
	}

	if mu.CanIMoveInDirection(models.Left) {
		t.Fail()
	}
	if !mu.CanIMoveInDirection(models.Up) {
		t.Fail()
	}
	if mu.CanIMoveInDirection(models.Right) {
		t.Fail()
	}
}

func Test_GraphOfMap_emptyMapFindsPathAcrossTheMapWithoutErrors(t *testing.T) {
	mu := MapUtility{
		mapp: models.Map{Width: 5, Height: 5},
	}
	g := GraphOfMap(mu)
	_, err := g.Shortest(0,24)
	assert.NoError(t, err)
}

func Test_GraphOfMap_returnsErrorWhenNoPathExists(t *testing.T) {
	mu := MapUtility{
		mapp: models.Map{Width: 5, Height: 5, ObstacleUpPositions: []int{10,11,12,13,14}},
	}
	g := GraphOfMap(mu)
	_, err := g.Shortest(0,24)
	assert.Error(t, err)
}