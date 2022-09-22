package storage

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/rs/zerolog/log"
	"socket/config"
)

var Database *Db

type Db struct {
	Redis *redis.Pool
}

var Storage = NewStorage()

type IStorage interface {
	InitDatabase() (err error)
	ConnectRedisCache() (err error)
}

type storage struct{}

func NewStorage() IStorage {
	return &storage{}
}

func (s storage) InitDatabase() (err error) {
	log.Info().Msg("InitDatabase")

	_redis := s.RedisClient()
	cacheConn, err := _redis.Dial()
	if err != nil {
		log.Err(err).Send()
		return
	}
	defer cacheConn.Close()

	Database = &Db{
		Redis: _redis,
	}
	return
}

func (s storage) RedisClient() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(
				"tcp",
				fmt.Sprintf("%s:%d", config.Config.RedisHost, config.Config.RedisPort),
				redis.DialPassword(config.Config.RedisPassword),
			)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

func (s storage) ConnectRedisCache() (err error) {
	_redis := s.RedisClient()
	_, err = _redis.Dial()
	if err != nil {
		log.Err(err).Send()
		return
	}

	Database.Redis = _redis
	return
}
