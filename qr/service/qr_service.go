package service

import (
	"common"
	cdto "common/dto"
	"context"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"qr/dto"
	ierr "qr/error"
	"qr/repository"
	"qr/streams/publisher"
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
		ExpiredAt: now.Add(time.Minute * 2),
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
		Status:    3,
		CreatedAt: qr.CreatedAt,
		ExpiredAt: qr.ExpiredAt,
		Amount:    qr.Amount,
		TTL:       ttlInSeconds,
	}

	err = publisher.Nats.Publish(common.SubjectQrcodeUpdate, cdto.QrCodesMessage{
		ID:        result.ID,
		Status:    3,
		CreatedAt: result.CreatedAt,
		ExpiredAt: result.ExpiredAt,
		Amount:    result.Amount,
		TTL:       result.TTL,
	})
	if err != nil {
		log.Err(err).Send()
		err = ierr.ErrNatsPublish
		return
	}

	emailEemplate := cdto.MasterSendEmailMessage{
		Email:        "opannapo@email.com",
		TemplateType: "qr_scan_success",
		TemplateData: cdto.TemplateTransactionQrSuccess{
			AppUrl:            "test.com",
			AccountFirstName:  "opan",
			AccountMiddleName: "napo",
			AccountLastName:   "opannapo",
			Amount:            result.Amount,
		},
	}
	err = publisher.Nats.Publish(common.SubjectSendEmail, emailEemplate)
	if err != nil {
		log.Err(err).Send()
		err = ierr.ErrNatsPublish
		return
	}

	return
}
