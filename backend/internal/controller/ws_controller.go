package controller

import (
	"github.com/LamichhaneBibek/quiz-webapp/internal/service"
	"github.com/gofiber/contrib/websocket"
)

type WebsocketController struct {
	netService *service.NetService
}

func NewWebsocketController(netService *service.NetService) *WebsocketController {
	return &WebsocketController{
		netService: netService,
	}
}

func (wc *WebsocketController) HandleWS(conn *websocket.Conn) {
	var (
		mt  int
		msg []byte
		err error
	)

	for {
		if mt, msg, err = conn.ReadMessage(); err != nil {
			wc.netService.OnDisconnect(conn)
			break
		}
		wc.netService.OnIncomingMessage(conn, mt, msg)
	}
}
