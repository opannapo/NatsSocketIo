package service

import (
	"github.com/google/uuid"
	"qr/dto"
	"qr/repository"
	"time"
)

var QrService = newQrService()

type IQrService interface {
	Create(request dto.CreateQrRequest) (result dto.CreateQrResponse, err error)
}

func newQrService() IQrService {
	return &qrService{}
}

type qrService struct{}

func (q qrService) Create(request dto.CreateQrRequest) (result dto.CreateQrResponse, err error) {
	now := time.Now()

	newQr := repository.Qr{
		ID:        uuid.New().String(),
		Status:    0,
		CreatedAt: now,
		ExpiredAt: now.Add(time.Minute * 15),
		Amount:    request.Amount,
	}
	err = repository.QrDao.Create(newQr)

	result = dto.CreateQrResponse{
		ID:        newQr.ID,
		Status:    newQr.Status,
		CreatedAt: newQr.CreatedAt,
		ExpiredAt: newQr.ExpiredAt,
		Amount:    newQr.Amount,
	}
	return
}

func (q qrService) Scan(request dto.CreateQrRequest) (result dto.CreateQrResponse, err error) {
	now := time.Now()

	newQr := repository.Qr{
		ID:        uuid.New().String(),
		Status:    0,
		CreatedAt: now,
		ExpiredAt: now.Add(time.Minute * 15),
		Amount:    request.Amount,
	}
	err = repository.QrDao.Create(newQr)

	result = dto.CreateQrResponse{
		ID:        newQr.ID,
		Status:    newQr.Status,
		CreatedAt: newQr.CreatedAt,
		ExpiredAt: newQr.ExpiredAt,
		Amount:    newQr.Amount,
	}
	return
}
