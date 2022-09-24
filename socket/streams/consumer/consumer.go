package consumer

import (
	"common"
	"github.com/rs/zerolog/log"
	"socket/streams"
)

func StartConsumer() (err error) {
	qrConsumer := newQrConsumer()
	_, err = streams.MessageBroker.Nats.Subscribe(common.SubjectQrcodeUpdate, qrConsumer.Update)
	if err != nil {
		log.Err(err).Send()
		return
	}

	return
}
