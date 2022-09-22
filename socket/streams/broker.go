package streams

import (
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"socket/config"
)

var MessageBroker *messageBroker

type messageBroker struct {
	Nats *nats.Conn
}

func NewMessageBroker() (broker *messageBroker, err error) {
	nats, err := natsClient()
	if err != nil {
		log.Err(err).Send()
		return
	}

	broker = &messageBroker{Nats: nats}
	return
}

func natsClient() (clientConn *nats.Conn, err error) {
	clientConn, err = nats.Connect(config.Config.NatsAddress)
	if err != nil {
		log.Err(err).Send()
		return
	}

	return
}

func ConnectMessageBroker() (err error) {
	broker, err := NewMessageBroker()
	if err != nil {
		log.Err(err).Send()
		return
	}

	MessageBroker = broker
	return
}
