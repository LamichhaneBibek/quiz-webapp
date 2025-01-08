package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/LamichhaneBibek/quiz-webapp/internal/entity"
	"github.com/gofiber/contrib/websocket"
)

type NetService struct {
	quizeService *QuizService
}

func NewNetService(quizService *QuizService) *NetService {
	return &NetService{
		quizeService: quizService,
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
	}
	return 0, errors.New("invalid packet")
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
			fmt.Println(data.Name, "wants to join the game", data.Code)
			break
		}
	case *HostGamePacket:
		{
			fmt.Println("User wants to host a quiz", data.QuidId)
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
