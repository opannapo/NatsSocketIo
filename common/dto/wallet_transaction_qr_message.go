package dto

import "time"

type QrCodesMessage struct {
	ID        string    `json:"id"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
	Amount    int64     `json:"amount"`
	TTL       string    `json:"ttl"`
}
