package basebot

import (
	"encoding/json"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"paintbot-client/models"
	"paintbot-client/utilities/timeHelper"
)

var mux sync.Mutex

type GameState struct {
	gameMode models.GameMode
	settings models.GameSettings
}

func Start(
	playerName string,
	gameMode models.GameMode,
	desiredGameSettings *models.GameSettings,
	calculateMove func(settings models.GameSettings, event models.MapUpdateEvent) models.Action,
) {
	state := GameState{
		gameMode: gameMode,
	}
	conn := getWebsocketConnection(gameMode)
	defer conn.Close()

	registerPlayer(conn, playerName, desiredGameSettings)

	handleMapUpdate := func(conn *websocket.Conn, event models.MapUpdateEvent) {
		s := time.Now()
		action := calculateMove(state.settings, event)
		e := time.Now()
		decisionTime := e.Sub(s)
		fmt.Printf("[%-3dms] Action: %s\n", decisionTime.Milliseconds(), action)
		sendMove(conn, event, action)
	}

	for !state.recv(conn, handleMapUpdate) {
	}
}

func (s *GameState) recv(conn *websocket.Conn, handleMapUpdate func(*websocket.Conn, models.MapUpdateEvent)) (done bool) {
	var msg []byte
	var err error
	if _, msg, err = conn.ReadMessage(); err != nil {
		panic(err)
	}

	log.Debugf("Received: %s\n", msg)

	gameMSG := models.GameMessage{}
	if err := json.Unmarshal(msg, &gameMSG); err != nil {
		panic(err)
	}

	switch models.MessageType(gameMSG.Type) {
	case models.MessageTypeInvalidMessage:
		panic("invalid message: " + string(msg))
	case models.MessageTypePlayerRegistered:
		playerRegisteredEvent := models.PlayerRegisteredEvent{}
		if err = json.Unmarshal(msg, &playerRegisteredEvent); err != nil {
			panic(err)
		}

		s.settings = playerRegisteredEvent.GameSettings
		log.Infof("Player registered")
		sendClientInfo(conn, gameMSG)
		go heartbeat(conn, gameMSG.ReceivingPlayerID)
		StartGame(conn)
	case models.MessageTypeGameLinkEvent:
		gameLinkEvent := &models.GameLinkEvent{}
		if err := json.Unmarshal(msg, gameLinkEvent); err != nil {
			panic(err)
		}
		log.Infof("Game can be viewed at: %s\n", gameLinkEvent.URL)
	case models.MessageTypeGameStartingEvent:
		log.Infof("Game started\n")
	case models.MessageTypeMapUpdateEvent:
		updateEvent := models.MapUpdateEvent{}
		if err := json.Unmarshal(msg, &updateEvent); err != nil {
			panic(err)
		}
		if updateEvent.GameTick%10 == 0 {
			log.Infof("Game tick: %d/%d\n", updateEvent.GameTick, s.settings.TotalTicks())
		}
		handleMapUpdate(conn, updateEvent)
	case models.MessageTypeGameResultEvent:
		event := models.GameResultEvent{}
		if err := json.Unmarshal(msg, &event); err != nil {
			panic(err)
		}

		log.Infof("### Game Results ###\n")
		for _, player := range event.PlayerRanks {
			log.Infof("%d: %s - %d\n", player.Rank, player.PlayerName, player.Points)
		}
	case models.MessageTypeGameEndedEvent:
		event := models.GameEndedEvent{}
		if err := json.Unmarshal(msg, &event); err != nil {
			panic(err)
		}

		if event.PlayerWinnerID == *event.ReceivingPlayerID {
			log.Info("You won the game")
		}

		if s.gameMode == models.Training {
			return true
		}
	case models.MessageTypeTournamentEndedEvent:
		event := models.TournamentEndedEvent{}
		if err := json.Unmarshal(msg, &event); err != nil {
			panic(err)
		}

		log.Infof("### Tournament Ended ###")
		for _, player := range event.GameResult {
			log.Infof("%s - %d\n", player.Name, player.Points)
		}
		return true
	case models.MessageTypeHeartBeatResponse:
	default:
		panic(fmt.Sprintf("unknown message: %s\n", msg))
	}
	return false
}

func registerPlayer(conn *websocket.Conn, playerName string, desiredGameSettings *models.GameSettings) {
	registerMSG := &models.RegisterPlayerEvent{
		Type:              "se.cygni.paintbot.api.request.RegisterPlayer",
		PlayerName:        playerName,
		GameSettings:      desiredGameSettings,
		ReceivingPlayerID: nil,
		Timestamp:         timeHelper.Now(),
	}

	log.Debugf("Registering player: %v\n", registerMSG)
	send(conn, registerMSG)
}

func sendClientInfo(conn *websocket.Conn, msg models.GameMessage) {
	clientInfoMSG := &models.ClientInfoMSG{
		Type:                   "se.cygni.paintbot.api.event.GameStartingEvent",
		Language:               "Go",
		LanguageVersion:        runtime.Version(),
		OperatingSystem:        runtime.GOOS,
		OperatingSystemVersion: "",
		ClientVersion:          "0.3",
		ReceivingPlayerID:      msg.ReceivingPlayerID,
		Timestamp:              timeHelper.Now(),
	}
	send(conn, clientInfoMSG)
}

func StartGame(conn *websocket.Conn) {
	startGame := &models.StartGameEvent{
		Type:              "se.cygni.paintbot.api.request.StartGame",
		ReceivingPlayerID: nil,
		Timestamp:         timeHelper.Now(),
	}

	send(conn, startGame)
}

func sendMove(conn *websocket.Conn, updateEvent models.MapUpdateEvent, action models.Action) {
	moveEvent := &models.RegisterMoveEvent{
		Type:              "se.cygni.paintbot.api.request.RegisterMove",
		GameID:            updateEvent.GameID,
		GameTick:          updateEvent.GameTick,
		Action:            string(action),
		ReceivingPlayerID: updateEvent.ReceivingPlayerID,
		Timestamp:         timeHelper.Now(),
	}
	if marshal, err := json.Marshal(moveEvent); err != nil {
		panic(err)
	} else {
		log.Debugf("send action: %s\n", marshal)
	}

	send(conn, moveEvent)
}
