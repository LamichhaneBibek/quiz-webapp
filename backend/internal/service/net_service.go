package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/LamichhaneBibek/quiz-webapp/internal/entity"
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NetService struct {
	quizeService *QuizService
	games        []*Game
}

func NewNetService(quizService *QuizService) *NetService {
	return &NetService{
		quizeService: quizService,
		games:        []*Game{},
	}
}

type ConnectPacket struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type HostGamePacket struct {
	QuidId string `json:"quiz_id"`
}

type QuestionShowPacket struct {
	Question entity.QuizQuestion `json:"question"`
}

type ChangeGameStatePacket struct {
	State GameState `json:"state"`
}

type PlayerJoinPacket struct {
	Player Player `json:"player"`
}

type PlayerDisconnectPacket struct {
	PlayerId uuid.UUID `json:"playerId"`
}

type StartGamePacket struct {
}

type TickPacket struct {
	Tick int `json:"tick"`
}

type QuestionAnswerPacket struct {
	Question int `json:"question"`
}

type PlayerRevealPacket struct {
	Points int `json:"points"`
}

type LeaderboardPacket struct {
	Points []LeaderboardEntry `json:"points"`
}

func (ns *NetService) packetIdToPacket(packetId uint8) any {
	switch packetId {
	case 0:
		{
			return &ConnectPacket{}
		}
	case 1:
		{
			return &HostGamePacket{}
		}
	case 5:
		{
			return &StartGamePacket{}
		}
	case 7:
		{
			return &QuestionAnswerPacket{}
		}
	}

	return nil
}

func (ns *NetService) packetToPacketId(packet any) (uint8, error) {
	switch packet.(type) {
	case QuestionShowPacket:
		{
			return 2, nil
		}
	case HostGamePacket:
		{
			return 1, nil
		}
	case ChangeGameStatePacket:
		{
			return 3, nil
		}
	case PlayerJoinPacket:
		{
			return 4, nil
		}
	case TickPacket:
		{
			return 6, nil
		}
	case PlayerRevealPacket:
		{
			return 8, nil
		}
	case LeaderboardPacket:
		{
			return 9, nil
		}
	case PlayerDisconnectPacket:
		{
			return 10, nil
		}
	}
	return 0, errors.New("invalid packet")
}

func (ns *NetService) getGamesByCode(code string) *Game {
	for _, game := range ns.games {
		if game.Code == code {
			return game
		}
	}
	return nil
}

func (ns *NetService) getGameByHost(host *websocket.Conn) *Game {
	for _, game := range ns.games {
		if game.Host == host {
			return game
		}
	}
	return nil
}

func (c *NetService) getGameByPlayer(con *websocket.Conn) (*Game, *Player) {
	for _, game := range c.games {
		for _, player := range game.Players {
			if player.Connection == con {
				return game, player
			}
		}
	}

	return nil, nil
}

func (c *NetService) OnDisconnect(con *websocket.Conn) {
	game, player := c.getGameByPlayer(con)
	if game == nil {
		return
	}

	game.OnPlayerDisconnect(player)
}

func (ns *NetService) OnIncomingMessage(conn *websocket.Conn, mt int, msg []byte) {
	if len(msg) == 0 {
		return
	}

	packetId := msg[0]
	data := msg[1:]

	packet := ns.packetIdToPacket(packetId)
	if packet == nil {
		return
	}

	err := json.Unmarshal(data, packet)
	if err != nil {
		fmt.Println("onIncomingMessage: ", err)
		return
	}

	fmt.Println(packet)
	switch data := packet.(type) {
	case *ConnectPacket:
		{
			game := ns.getGamesByCode(data.Code)
			if game == nil {
				fmt.Println("game not found")
				return
			}

			game.OnPlayerJoin(data.Name, conn)
			break
		}
	case *HostGamePacket:
		{
			quizId, err := primitive.ObjectIDFromHex(data.QuidId)
			if err != nil {
				fmt.Println("invalid quiz id")
				return
			}

			quiz, err := ns.quizeService.quizCollection.GetQuizByID(quizId)
			if err != nil {
				fmt.Println("quiz not found")
				return
			}

			if quiz == nil {
				fmt.Println("quiz not found")
				return
			}
			game := newGame(*quiz, conn, ns)
			ns.games = append(ns.games, &game)

			ns.SendPacket(conn, HostGamePacket{
				QuidId: game.Code,
			})

			ns.SendPacket(conn, ChangeGameStatePacket{
				State: game.State,
			})
			break
		}
	case *StartGamePacket:
		{
			game := ns.getGameByHost(conn)
			if game == nil {
				fmt.Println("game not found")
				return
			}
			game.StartOrSkip()
			break
		}
	case *QuestionAnswerPacket:
		{
			game, player := ns.getGameByPlayer(conn)
			if game == nil {
				return
			}

			game.OnPlayerAnswer(data.Question, player)
			break
		}
	}
}

func (ns *NetService) PacketToBytes(packet any) ([]byte, error) {
	packetId, err := ns.packetToPacketId(packet)
	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(packet)
	if err != nil {
		return nil, err
	}

	return append([]byte{packetId}, bytes...), nil
}

func (ns *NetService) SendPacket(conn *websocket.Conn, packet any) error {
	bytes, err := ns.PacketToBytes(packet)
	if err != nil {
		return err
	}

	return conn.WriteMessage(websocket.BinaryMessage, bytes)
}
