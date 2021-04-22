package models

type Action string

const (
	Left    Action = "LEFT"
	Right   Action = "RIGHT"
	Up      Action = "UP"
	Down    Action = "DOWN"
	Stay    Action = "STAY"
	Explode Action = "EXPLODE"
)

type Tile string

const (
	Obstacle Tile = "OBSTACLE"
	PowerUp  Tile = "POWERUP"
	Player   Tile = "PLAYER"
	Open     Tile = "OPEN"
)

type GameMode string

const (
	Tournament GameMode = "/tournament"
	Training   GameMode = "/training"
)

type MessageType string

const (
	MessageTypeInvalidMessage       MessageType = "se.cygni.paintbot.api.exception.InvalidMessage"
	MessageTypePlayerRegistered     MessageType = "se.cygni.paintbot.api.response.PlayerRegistered"
	MessageTypeGameLinkEvent        MessageType = "se.cygni.paintbot.api.event.GameLinkEvent"
	MessageTypeGameStartingEvent    MessageType = "se.cygni.paintbot.api.event.GameStartingEvent"
	MessageTypeMapUpdateEvent       MessageType = "se.cygni.paintbot.api.event.MapUpdateEvent"
	MessageTypeGameResultEvent      MessageType = "se.cygni.paintbot.api.event.GameResultEvent"
	MessageTypeGameEndedEvent       MessageType = "se.cygni.paintbot.api.event.GameEndedEvent"
	MessageTypeTournamentEndedEvent MessageType = "se.cygni.paintbot.api.event.TournamentEndedEvent"
	MessageTypeHeartBeatResponse    MessageType = "se.cygni.paintbot.api.response.HeartBeatResponse"
)
