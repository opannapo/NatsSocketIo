package publisher

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	ierr "qr/error"
	"qr/streams"
)

var Nats = NewNatsPublisher()

type INatsPublisher interface {
	Publish(subject string, payload interface{}) (err error)
}

type natsPublisher struct {
}

func NewNatsPublisher() INatsPublisher {
	return &natsPublisher{}
}

func (n natsPublisher) Publish(subject string, payload interface{}) (err error) {
	data, _ := json.Marshal(payload)
	if streams.MessageBroker == nil {
		err = ierr.ErrMessageBrokerInstance
		log.Err(err).Send()
		return
	}

	if streams.MessageBroker.Nats == nil {
		err = ierr.ErrNatsInstance
		log.Err(err).Send()
		return
	}

	if streams.MessageBroker.Nats.IsConnected() {
		err = streams.MessageBroker.Nats.Publish(subject, data)
		if err != nil {
			log.Err(err).Send()
			return err
		}
		log.Info().Msg("Published message on subject " + subject)
	} else {
		err = nats.ErrConnectionClosed
		log.Err(err).Send()
	}

	return
}
