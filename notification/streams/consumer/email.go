package consumer

import (
	cdto "common/dto"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"notification/service"
)

type IEmailConsumer interface {
	Send(msg *nats.Msg)
}

type emailConsumer struct{}

func newEmailConsumer() IEmailConsumer {
	return &emailConsumer{}
}

func (e emailConsumer) Send(msg *nats.Msg) {
	payload := cdto.MasterSendEmailMessage{}
	if err := json.Unmarshal(msg.Data, &payload); err != nil {
		log.Err(err).Send()
	}

	fmt.Printf("(EmailConsumer) NATS Received %s : Message in : %+v \n", msg.Subject, payload)
	err := service.EmailService.SendEmail(payload)
	if err != nil {
		log.Err(err).Send()
		return
	}
}
