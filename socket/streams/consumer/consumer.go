package consumer

import (
	"github.com/rs/zerolog/log"
	"mono/common"
	"socket/streams"
)

func StartConsumer() (err error) {
	err = streams.ConnectMessageBroker()
	if err != nil {
		log.Err(err).Send()
		return
	}

	walletConsumer := newWalletConsumer()
	_, err = streams.MessageBroker.Nats.Subscribe(common.SubjectWalletQrcodeUpdate, walletConsumer.Update)
	if err != nil {
		log.Err(err).Send()
		return
	}

	return
}
