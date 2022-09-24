package consumer

import (
	cdto "common/dto"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	logic "socket/service"
)

type IWalletConsumer interface {
	Update(msg *nats.Msg)
}

type walletConsumer struct{}

func newWalletConsumer() IWalletConsumer {
	return &walletConsumer{}
}

func (c walletConsumer) Update(msg *nats.Msg) {
	payload := cdto.WalletTransactionQrcodesMessage{}
	if err := json.Unmarshal(msg.Data, &payload); err != nil {
		log.Err(err).Send()
	}

	log.Printf("Message in : %+v", payload)
	logic.SocketService.HandleQrCodeUpdate(payload)
}
