package dto

import "time"

type WalletTransactionQrcodesMessage struct {
	QrcodesID        string    `json:"qrcodes_id"`
	QrcodesCode      string    `json:"qrcodes_code"`
	QrcodesValue     string    `json:"qrcodes_value"`
	QrcodesStatus    string    `json:"qrcodes_status"`
	QrcodesExpired   time.Time `json:"qrcodes_expired"`
	QrcodesTimestamp time.Time `json:"qrcodes_timestamp"`
}
