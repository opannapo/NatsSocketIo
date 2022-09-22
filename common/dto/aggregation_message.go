package dto

import "time"

const (
	AggregationMessageTypeCreate = 1
	AggregationMessageTypeUpdate = 2
	AggregationMessageTypeDelete = 3
)

type AggregationMessage struct {
	Domain             string    `json:"domain"` //application,service name :: company, user, wallet, etc...
	Type               int       `json:"type"`
	ReferenceID        string    `json:"reference_id"`
	ReferenceValue     *[]byte   `json:"reference_value"`
	ReferenceTimestamp time.Time `json:"reference_timestamp"`
}
