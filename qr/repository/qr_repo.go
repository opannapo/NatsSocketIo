package repository

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	ierr "qr/error"
	"time"
)

type Qr struct {
	ID        string    `json:"id" gorm:"column:id"`
	Status    int       `json:"status" gorm:"column:status"` //0 new, 1 scanned, 2 success, 3 failed
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
	Scan(ID string) (err error)
	GetByID(ID string) (result Qr, err error)
}

func newQrDao() IQrDao {
	return &qrDao{}
}

type qrDao struct{}

func (q qrDao) GetByID(ID string) (result Qr, err error) {
	result = Qr{ID: ID}
	err = Database.Mysql.First(&result).Error
	if err != nil {
		log.Err(err).Send()
		if err == gorm.ErrRecordNotFound {
			err = ierr.ErrQrNotFound
		}
		return
	}

	return
}

func (q qrDao) Scan(ID string) (err error) {
	query := `UPDATE qr SET status=1 WHERE id = ? and expired_at > ? and status=0`
	tx := Database.Mysql.Exec(query, ID, time.Now())
	if tx.Error != nil {
		log.Err(tx.Error).Send()
		err = ierr.ErrUpdateQrScan
		return
	}
	if tx.RowsAffected != 1 {
		err = ierr.ErrNoRowsAffected
		log.Err(err).Send()
		return
	}

	return
}

func (q qrDao) Create(data Qr) (err error) {
	err = Database.Mysql.Create(&data).Error
	return
}
