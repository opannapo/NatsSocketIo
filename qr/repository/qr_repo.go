package repository

import (
	"time"
)

type Qr struct {
	ID        string    `json:"id" gorm:"column:id"`
	Status    int       `json:"status" gorm:"column:status"` // 0 pending, 1 success, 2 failed
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	ExpiredAt time.Time `json:"expired_at" gorm:"column:expired_at"`
	Amount    int64     `json:"amount" gorm:"column:amount"`
}

func (m *Qr) TableName() string {
	return "qr"
}

var QrDao = newQrDao()

type IQrDao interface {
	Create(data Qr) (err error)
}

func newQrDao() IQrDao {
	return &qrDao{}
}

type qrDao struct{}

func (q qrDao) Create(data Qr) (err error) {
	err = Database.Mysql.Debug().Create(&data).Error
	return
}
