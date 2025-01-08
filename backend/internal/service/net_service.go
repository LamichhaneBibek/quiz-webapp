package service

import "github.com/gofiber/contrib/websocket"

type NetService struct {
	quizeService *QuizService
}

func NewNetService(quizService *QuizService) *NetService {
	return &NetService{
		quizeService: quizService,
	}
}

type connectPacket struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type HostGamePacket struct {
	QuidId string `json:"quiz_id"`
}

func (ns *NetService) OnIncomingMessage(conn *websocket.Conn, mt int, msg []byte) {

}
