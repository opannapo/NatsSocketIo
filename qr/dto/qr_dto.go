package dto

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
	ierr "qr/error"
	"time"
)

type CreateQrRequest struct {
	Amount int64 `json:"amount"`
}

func (cqr *CreateQrRequest) Validate(r *http.Request) (err error) {
	err = json.NewDecoder(r.Body).Decode(&cqr)
	if err != nil {
		log.Err(err).Interface("context", r.Context()).Send()
		return ierr.ErrInvalidRequestPayload
	}

	if cqr.Amount < 500 {
		return ierr.ErrInvalidCreateQrAmount
	}

	return
}

type CreateQrResponse struct {
	ID        string    `json:"id"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
	Amount    int64     `json:"amount"`
}

type ScanQrResponse struct {
	ID        string    `json:"id"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
	Amount    int64     `json:"amount"`
	TTL       string    `json:"ttl"`
}
