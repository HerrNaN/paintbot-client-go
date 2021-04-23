package maputility

import (
	"fmt"
	"paintbot-client/models"
	"paintbot-client/utilities/arrays"

	"github.com/RyanCarrier/dijkstra"
)

type Graph interface {
	Shortest(from int, to int) (dijkstra.BestPath, error)
}

// Utility for getting information from the map object in a bit more developer friendly format
type MapUtility struct {
	mapp            models.Map
	graph           Graph
	currentPlayerID string
}

func (u *MapUtility) SetGraph(g Graph) {
	u.graph = g
}

func New(m models.Map, g Graph, currentPlayerID string) *MapUtility {
	return &MapUtility{
		mapp:            m,
		currentPlayerID: currentPlayerID,
		graph:           g,
	}
}

func GraphOfMap(u MapUtility) Graph {
	g := dijkstra.NewGraph()
	maxID := u.mapp.Height * u.mapp.Width
	for i := 0; i < maxID; i++ {
		g.AddVertex(i)
	}
	for i := 0; i < maxID; i++ {
		connectedNeighbours := u.inBoundsAccessibleNeighbours(i)
		for _, n := range connectedNeighbours {
			err := g.AddArc(i, n, 1)
			if err != nil {
				panic("cannot add arc to Graph")
			}
		}
	}
	return g
}

// returns true if the current player can perform the given action given no action for all other players
func (u *MapUtility) CanIMoveInDirection(action models.Action) bool {
	info := u.getMyCharacterInfo()

	if info.StunnedForGameTicks > 0 {
		return false
	}

	if action == models.Explode {
		return info.CarryingPowerUp
	}

	if action == models.Stay {
		return true
	}

	pos := u.GetMyCoordinates()
	pos = u.TranslateCoordinateByAction(action, pos)

	return u.IsTileAvailableForMovementTo(pos)
}

// Returns the coordinates given after an action has been performed successfully
func (u *MapUtility) TranslateCoordinateByAction(action models.Action, pos models.Coordinates) models.Coordinates {
	switch action {
	case models.Left:
		return models.Coordinates{X: pos.X - 1, Y: pos.Y}
	case models.Right:
		return models.Coordinates{X: pos.X + 1, Y: pos.Y}
	case models.Up:
		return models.Coordinates{X: pos.X, Y: pos.Y - 1}
	case models.Down:
		return models.Coordinates{X: pos.X, Y: pos.Y + 1}
	case models.Stay, models.Explode:
		return models.Coordinates{X: pos.X, Y: pos.Y}
	default:
		panic("Unknown Action: " + action)
	}
}

func (u *MapUtility) GetColouredBy(coordinates models.Coordinates) *Player {
	pos := u.ConvertCoordinatesToPosition(coordinates)

	for _, c := range u.mapp.CharacterInfos {
		if arrays.Contains(c.ColouredPosition, pos) {
			p := u.toPlayer(c)
			return &p
		}
	}

	return nil
}

// returns list of all the coordinates containing a power up
func (u *MapUtility) ListCoordinatesContainingPowerUps() []models.Coordinates {
	return u.ConvertPositionsToCoordinates(u.mapp.PowerUpPositions)
}

// returns list of all the coordinates containing a obstacle
func (u *MapUtility) GetObstacleCoordinates() []models.Coordinates {
	return u.ConvertPositionsToCoordinates(u.mapp.ObstacleUpPositions)
}

// returns List of all the coordinates coloured by the given player
func (u *MapUtility) ListCoordinatesColouredByPlayer(playerId string) []models.Coordinates {
	return u.ConvertPositionsToCoordinates(u.getCharacterInfo(playerId).ColouredPosition)
}

// returns true if tile is walkable
func (u *MapUtility) IsTileAvailableForMovementTo(coord models.Coordinates) bool {
	tile := u.GetTileAt(coord)

	return tile == models.Open || tile == models.PowerUp || tile == models.Player
}

// returns the coordinates of the current player
func (u *MapUtility) GetMyCoordinates() models.Coordinates {
	return u.ConvertPositionToCoordinates(u.getMyCharacterInfo().Position)
}

// returns information about the current player
func (u *MapUtility) getMyCharacterInfo() models.CharacterInfo {
	return u.getCharacterInfo(u.currentPlayerID)
}

// returns the current player
func (u *MapUtility) GetMe() Player {
	c := u.getCharacterInfo(u.currentPlayerID)
	return Player{info: &c, utility: u}
}

// returns information about the given player
func (u *MapUtility) getCharacterInfo(playerID string) models.CharacterInfo {
	for i := range u.mapp.CharacterInfos {
		if u.mapp.CharacterInfos[i].ID == playerID {
			return u.mapp.CharacterInfos[i]
		}
	}
	panic("Trying to find invalid playerID: " + playerID)
}

// Returns true if the coordinate is withing the game field
func (u *MapUtility) IsCoordinatesOutOfBounds(coord models.Coordinates) bool {
	w := u.mapp.Width
	h := u.mapp.Height
	return coord.X < 0 || coord.Y < 0 || coord.X >= w || coord.Y >= h
}

// returns the type of object at the given coordinates
// returns OBSTACLE if the coordinate is out of bounds
func (u *MapUtility) GetTileAt(coordinates models.Coordinates) models.Tile {
	if u.IsCoordinatesOutOfBounds(coordinates) {
		return models.Obstacle
	}

	return u.getTileAtPosition(u.ConvertCoordinatesToPosition(coordinates))
}

func (u *MapUtility) getTileAtPosition(position int) models.Tile {
	if arrays.Contains(u.mapp.ObstacleUpPositions, position) {
		return models.Obstacle
	}

	if arrays.Contains(u.mapp.PowerUpPositions, position) {
		return models.PowerUp
	}

	if arrays.Contains(u.getPlayerPositions(), position) {
		return models.Player
	}

	return models.Open
}

// Converts a position in the flattened single array representation
// of the Map to a Coordinates.
func (u *MapUtility) ConvertPositionToCoordinates(position int) models.Coordinates {
	w := u.mapp.Width
	return models.Coordinates{
		X: position % w,
		Y: position / w,
	}
}

// Converts a MapCoordinate to the same position in the flattened
// single array representation of the Map.
func (u *MapUtility) ConvertCoordinatesToPosition(coordinates models.Coordinates) int {
	w := u.mapp.Width
	return coordinates.Y*w + coordinates.X
}

// converts a list of positions to coordinates
func (u *MapUtility) ConvertPositionsToCoordinates(positions []int) []models.Coordinates {
	coords := make([]models.Coordinates, len(positions))
	for i := range positions {
		coords[i] = u.ConvertPositionToCoordinates(positions[i])
	}
	return coords
}

// converts a list of coordinates to positions
func (u *MapUtility) ConvertCoordinatesToPositions(coordinates []models.Coordinates) []int {
	positions := make([]int, len(coordinates))
	for i := range coordinates {
		positions[i] = u.ConvertCoordinatesToPosition(coordinates[i])
	}
	return positions
}

// DistanceTo returns the distance to a specified coordinate from the players current position
func (u *MapUtility) DistanceTo(destination models.Coordinates) (int, error) {
	if !u.IsTileAvailableForMovementTo(destination) || u.IsCoordinatesOutOfBounds(destination) {
		return 0, fmt.Errorf("coordinates are unreachable: %v", destination)
	}

	myPos := u.GetMe().info.Position

	destinationPos := u.ConvertCoordinatesToPosition(destination)
	bestPath, err := u.graph.Shortest(myPos, destinationPos)
	if err != nil {
		return 0, err
	}

	// Assuming Width * Height < int.max
	return int(bestPath.Distance), nil
}

func (u *MapUtility) DirectionToPoint(p int) models.Action {
	coord := u.ConvertPositionToCoordinates(p)
	myCoord := u.GetMyCoordinates()

	if myCoord.X+1 == coord.X &&
		myCoord.Y == coord.Y {
		return models.Right
	}

	if myCoord.X-1 == coord.X &&
		myCoord.Y == coord.Y {
		return models.Left
	}

	if myCoord.X == coord.X &&
		myCoord.Y-1 == coord.Y {
		return models.Up
	}

	if myCoord.X == coord.X &&
		myCoord.Y+1 == coord.Y {
		return models.Down
	}

	panic("p should be a neighbour")
}

// ShortestPathTo returns the shortest path to the given destination.
// If the destination is unreachable an error is returned.
func (u *MapUtility) ShortestPathTo(destination models.Coordinates) ([]int, error) {
	if !u.IsTileAvailableForMovementTo(destination) || u.IsCoordinatesOutOfBounds(destination) {
		return nil, fmt.Errorf("coordinates are unreachable: %v", destination)
	}

	myPos := u.GetMe().info.Position
	destinationPos := u.ConvertCoordinatesToPosition(destination)
	g := u.graph

	bestPath, err := g.Shortest(myPos, destinationPos)
	if err != nil {
		return nil, err
	}
	return bestPath.Path, nil
}

func (u *MapUtility) IsAnyPlayerWithinExplosionRange() bool {
	pps := u.getPlayerPositions()
	pcs := u.ConvertPositionsToCoordinates(pps)
	myCoord := u.GetMyCoordinates()
	for _, pc := range pcs {
		if pc == myCoord {
			continue
		}
		if d, _ := u.DistanceTo(pc); d <= 4 {
			return true
		}
	}
	return false
}

func (u *MapUtility) inBoundsAccessibleNeighbours(pos int) []int {
	coord := u.ConvertPositionToCoordinates(pos)

	var neighbours []int
	for _, a := range models.Movements {
		neighbour := u.TranslateCoordinateByAction(a, coord)
		if !u.IsCoordinatesOutOfBounds(neighbour) && u.IsTileAvailableForMovementTo(neighbour) {
			neighbours = append(neighbours, u.ConvertCoordinatesToPosition(neighbour))
		}
	}
	return neighbours
}

func (u *MapUtility) getPlayerPositions() []int {
	positions := make([]int, len(u.mapp.CharacterInfos))
	for i := range u.mapp.CharacterInfos {
		positions[i] = u.mapp.CharacterInfos[i].Position
	}
	return positions
}

func (u *MapUtility) toPlayer(info models.CharacterInfo) Player {
	return Player{
		info:    &info,
		utility: u,
	}
}
