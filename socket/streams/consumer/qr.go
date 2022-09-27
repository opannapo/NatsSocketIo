package consumer

import (
	cdto "common/dto"
	"encoding/json"
	"fmt"
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

	fmt.Printf("(QrConsumer) NATS Received %s : Message in : %+v \n", msg.Subject, payload)
	logic.SocketService.HandleQrCodeUpdate(payload)
}
