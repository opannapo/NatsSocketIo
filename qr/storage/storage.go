package storage

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"qr/config"
)

var Database *Db

type Db struct {
	Mysql *gorm.DB
	Redis *redis.Pool
}

var Storage = NewStorage()

type IStorage interface {
	InitDatabase() (err error)
}

type storage struct{}

func NewStorage() IStorage {
	return &storage{}
}

func (s storage) InitDatabase() (err error) {
	log.Info().Msg("InitDatabase")
	_mysql, err := s.mysqlClient()
	if err != nil {
		log.Err(err).Send()
		return
	}
	mysqlDb, err := _mysql.DB()
	defer mysqlDb.Close()
	if err != nil {
		log.Err(err).Send()
	}

	_redis := s.redisClient()
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

func (s storage) mysqlClient() (db *gorm.DB, err error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.Config.MysqlUsername,
		config.Config.MysqlPassword,
		config.Config.MysqlHost,
		config.Config.MysqlPort,
		config.Config.MysqlDatabase,
	)
	db, err = gorm.Open(mysql.New(mysql.Config{DSN: dns}), &gorm.Config{SkipDefaultTransaction: false, DisableAutomaticPing: false})
	if err != nil {
		log.Err(err).Send()
		return nil, err
	}

	return db, err
}

func (s storage) redisClient() *redis.Pool {
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
