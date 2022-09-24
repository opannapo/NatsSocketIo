package logic

import (
	"github.com/rs/zerolog/log"
	"socket/storage"
)

var Logic = NewLogic()

type ILogic interface {
}

func NewLogic() ILogic {
	return &logic{}
}

type logic struct {
}

func (l logic) OpenRedisConnection() (err error) {
	err = storage.Storage.ConnectRedisCache()
	if err != nil {
		log.Err(err).Send()
		return
	}

	return
}

func (l logic) CloseRedisConnection() {
	err := storage.Database.Redis.Close()
	if err != nil {
		log.Err(err).Send()
	}
}
