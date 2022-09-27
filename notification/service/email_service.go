package service

import (
	"common/dto"
	"github.com/rs/zerolog/log"
)

var EmailService = newEmailService()

type IEmailService interface {
	SendEmail(template dto.MasterSendEmailMessage) (err error)
}

func newEmailService() IEmailService {
	return &emailService{}
}

type emailService struct{}

func (e emailService) SendEmail(template dto.MasterSendEmailMessage) (err error) {
	log.Printf("data email %+v ", template)
	return nil
}
