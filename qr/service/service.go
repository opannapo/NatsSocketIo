package service

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
)

var Service = NewService()

type IService interface{}

func NewService() IService {
	return &service{}
}

type service struct{}

func getRequestID(ctx context.Context) string {
	reId, ok := ctx.Value("request_id").(string)
	if !ok {
		log.Err(errors.New("error no context request_id")).Send()
		reId = ""
	}

	return reId
}
