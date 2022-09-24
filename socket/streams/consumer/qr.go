package consumer

import (
	cdto "common/dto"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	logic "socket/service"
)

type IQrConsumer interface {
	Update(msg *nats.Msg)
}

type qrConsumer struct{}

func newQrConsumer() IQrConsumer {
	return &qrConsumer{}
}

func (c qrConsumer) Update(msg *nats.Msg) {
	payload := cdto.QrCodesMessage{}
	if err := json.Unmarshal(msg.Data, &payload); err != nil {
		log.Err(err).Send()
	}

	log.Printf("Message in : %+v", payload)
	logic.SocketService.HandleQrCodeUpdate(payload)
}
