package service

import (
	"common/dto"
	"fmt"
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
	fmt.Printf("Notification Service Send Email with data : %+v \n", template)
	return nil
}
