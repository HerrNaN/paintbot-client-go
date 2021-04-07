package basebot

import (
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"paintbot-client/models"
	"paintbot-client/utilities/timeHelper"
)

const timeBetweenHeartbeats = 30 * time.Second

func heartbeat(conn *websocket.Conn, playerID *string) {

	for {
		rq := &models.HearbeatMessage{
			Type:              "se.cygni.paintbot.api.request.HeartBeatRequest",
			ReceivingPlayerID: playerID,
			Timestamp:         timeHelper.Now(),
		}
		log.Debug("sending heartbeat")
		send(conn, rq)
		time.Sleep(timeBetweenHeartbeats)
	}
}
