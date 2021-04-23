package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"paintbot-client/basebot"
	"paintbot-client/models"
	"paintbot-client/utilities/maputility"
)

func main() {
	basebot.Start("\x00Golor Bot", models.Training, desiredGameSettings, calculateMove)
}

const MaxInt = int(^uint(0) >> 1)

var (
	moves   = []models.Action{models.Right, models.Down, models.Left, models.Up} // models.Explode, models.Stay}
	lastDir = 0
	graph   maputility.Graph = nil
)

// Implement your paintbot here
func calculateMove(settings models.GameSettings, updateEvent models.MapUpdateEvent) models.Action {
	utility := maputility.New(updateEvent.Map, nil, *updateEvent.ReceivingPlayerID)
	me := utility.GetMe()
	if graph == nil {
		fmt.Println("making map")
		graph = maputility.GraphOfMap(*utility)
	}

	utility.SetGraph(graph)

	if me.StunnedForTicks() > 0 {
		return models.Stay
	}

	powerupCoordinates := utility.ListCoordinatesContainingPowerUps()

	if len(powerupCoordinates) == 0 {
		return models.Up
	}

	distanceToClosestPowerUpCoord := MaxInt
	var closestPowerUpCoord models.Coordinates

	for _, coord := range powerupCoordinates {
		dist, err := utility.DistanceTo(coord)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if dist <= distanceToClosestPowerUpCoord {
			//fmt.Println(dist)
			distanceToClosestPowerUpCoord = dist
			closestPowerUpCoord = coord
		}
	}

	path, err := utility.ShortestPathTo(closestPowerUpCoord)
	if err != nil {
		panic("ojoj")
	}
	//myPos := utility.ConvertCoordinatesToPosition(utility.GetMyCoordinates())
	//closestPos := utility.ConvertCoordinatesToPosition(closestPowerUpCoord)
	//bestPath := pathMatrix[myPos][closestPos]
	//path := bestPath.Path
	//fmt.Println(path)
	if distanceToClosestPowerUpCoord == 1 && me.HasPowerUp() {
		fmt.Println("don't waste")
		return models.Explode
		//fmt.Println(distanceToClosestPowerUpCoord)
	}
	if utility.IsAnyPlayerWithinExplosionRange() && me.HasPowerUp() {
		fmt.Println("get them!")
		return models.Explode
	}
	nextPos := path[1] // [0] is the current position
	move := utility.DirectionToPoint(nextPos)
	return move
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:            true,
		ForceQuote:             true,
		FullTimestamp:          true,
		TimestampFormat:        "15:04:05.999",
		DisableLevelTruncation: true,
		PadLevelText:           true,
		QuoteEmptyFields:       true,
	})

	log.SetOutput(os.Stdout)

	log.SetLevel(log.InfoLevel)
}

// desired game settings can be changed to nil to get default settings
var desiredGameSettings = &models.GameSettings{
	MaxNOOFPlayers:                 5,
	TimeInMSPerTick:                250,
	ObstaclesEnabled:               true,
	PowerUpsEnabled:                true,
	AddPowerUpLikelihood:           38,
	RemovePowerUpLikelihood:        5,
	TrainingGame:                   true,
	PointsPerTileOwned:             1,
	PointsPerCausedStun:            5,
	NOOFTicksInvulnerableAfterStun: 3,
	NOOFTicksStunned:               10,
	StartObstacles:                 40,
	StartPowerUps:                  41,
	GameDurationInSeconds:          15,
	ExplosionRange:                 4,
	PointsPerTick:                  false,
}
