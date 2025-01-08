package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/LamichhaneBibek/quiz-webapp/internal/entity"
	"github.com/LamichhaneBibek/quiz-webapp/internal/game"
	"github.com/gofiber/contrib/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NetService struct {
	quizeService *QuizService
	games        []*game.Game
}

func NewNetService(quizService *QuizService) *NetService {
	return &NetService{
		quizeService: quizService,
		games:        []*game.Game{},
	}
}

type ConnectPacket struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type HostGamePacket struct {
	QuidId string `json:"quiz_id"`
}

type QestionShowPacket struct {
	Question entity.QuizQuestion `json:"question"`
}

type ChangeGameStatePacket struct {
	State game.GameState `json:"state"`
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
	}

	return nil
}

func (ns *NetService) packetToPacketId(packet any) (uint8, error) {
	switch packet.(type) {
	case QestionShowPacket:
		{
			return 2, nil
		}
	case ChangeGameStatePacket:
		{
			return 3, nil
		}
	}
	return 0, errors.New("invalid packet")
}

func (ns *NetService) getGamesByCode(code string) *game.Game {
	for _, game := range ns.games {
		if game.Code == code {
			return game
		}
	}
	return nil
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
			newGame := game.NewGame(*quiz, conn)
			fmt.Println("new game created", newGame.Code)
			ns.games = append(ns.games, &newGame)

			ns.SendPacket(conn, ChangeGameStatePacket{
				State: game.LobbyState,
			})
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
