package repository

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"qr/config"
)

var Database *db

type db struct {
	Mysql *gorm.DB
	Redis *redis.Pool
}

var Repository = newRepository()

type IRepository interface {
	InitDatabase() (err error)
}

type repository struct{}

func newRepository() IRepository {
	return &repository{}
}

func (r *repository) InitDatabase() (err error) {
	//init mysql
	_mysql, err := r.mysqlClient()
	if err != nil {
		log.Err(err).Send()
		return
	}

	//init redis
	_redis := r.redisClient()
	_, err = _redis.Dial()
	if err != nil {
		log.Err(err).Send()
		return
	}

	Database = &db{
		Redis: _redis,
		Mysql: _mysql,
	}
	return
}

func (r *repository) mysqlClient() (db *gorm.DB, err error) {
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

func (r *repository) redisClient() *redis.Pool {
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
