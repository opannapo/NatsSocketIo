package handler

import (
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/rs/zerolog/log"
	"socket/service"
)

var SocketHandler = NewSocketHandler()

type socketHandler struct {
}

func NewSocketHandler() ISocketHandler {
	return &socketHandler{}
}

type ISocketHandler interface {
	OnConnect(c *gosocketio.Channel) error
	OnDisconnect(c *gosocketio.Channel)
}

func (sh socketHandler) OnConnect(c *gosocketio.Channel) error {
	log.Print("OnConnect ", c.Id())
	//Handling diganti ke middleware
	/*err := logic.SocketService.ValidateRequest(c)
	if err != nil {
		c.Close()
		log.Err(err).Send()
		return err
	}*/

	logic.SocketService.AddUserToPrivateRoom(c)
	return nil
}

func (sh socketHandler) OnDisconnect(c *gosocketio.Channel) {
	log.Print("User:", c.Id(), " closed")
}
