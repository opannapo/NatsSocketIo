package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"qr/dto"
	"qr/repository"
	"time"
)

var QrService = newQrService()

type IQrService interface {
	Create(ctx context.Context, request dto.CreateQrRequest) (result dto.CreateQrResponse, err error)
	Scan(ctx context.Context, qrID string) (result dto.ScanQrResponse, err error)
}

func newQrService() IQrService {
	return &qrService{}
}

type qrService struct{}

func (q qrService) Create(ctx context.Context, request dto.CreateQrRequest) (result dto.CreateQrResponse, err error) {
	now := time.Now()
	reqId := getRequestID(ctx)

	newQr := repository.Qr{
		ID:        uuid.New().String(),
		Status:    0,
		CreatedAt: now,
		ExpiredAt: now.Add(time.Minute * 15),
		Amount:    request.Amount,
	}
	err = repository.QrDao.Create(newQr)
	if err != nil {
		log.Err(err).Interface("request_id", reqId).Send()
		return
	}

	result = dto.CreateQrResponse{
		ID:        newQr.ID,
		Status:    newQr.Status,
		CreatedAt: newQr.CreatedAt,
		ExpiredAt: newQr.ExpiredAt,
		Amount:    newQr.Amount,
	}
	return
}

func (q qrService) Scan(ctx context.Context, qrID string) (result dto.ScanQrResponse, err error) {
	reqId := getRequestID(ctx)

	qr, err := repository.QrDao.GetByID(qrID)
	if err != nil {
		log.Err(err).Interface("request_id", reqId).Send()
		return
	}

	err = repository.QrDao.Scan(qrID)
	if err != nil {
		log.Err(err).Interface("request_id", reqId).Send()
		return
	}

	ttlInSeconds := qr.ExpiredAt.Sub(time.Now()).String()
	result = dto.ScanQrResponse{
		ID:        qr.ID,
		Status:    qr.Status,
		CreatedAt: qr.CreatedAt,
		ExpiredAt: qr.ExpiredAt,
		Amount:    qr.Amount,
		TTL:       ttlInSeconds,
	}
	return
}
