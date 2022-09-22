package utils

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/rs/zerolog/log"
)

var CacheRedis = NewCacheRedis()

type ICacheRedis interface {
	Get(redisPool *redis.Pool, appName string, env string, key string) (interface{}, error)
	Set(redisPool *redis.Pool, appName string, appEnv string, key string, value string)
	SetExpired(redisPool *redis.Pool, appName string, appEnv string, key string, value string, expired int)
	Delete(redisPool *redis.Pool, appName string, appEnv string, key string)
}

func NewCacheRedis() ICacheRedis {
	return &cacheRedis{}
}

type cacheRedis struct{}

func (cr *cacheRedis) Get(redisPool *redis.Pool, appName string, appEnv string, key string) (interface{}, error) {
	if redisPool.Get() == nil {
		return nil, errors.New("no redis connection")
	}

	key = appName + "-" + appEnv + "-" + key
	log.Printf("Get cache %v ", key)

	val, err := redisPool.Get().Do("GET", key)
	if err == nil {
		val, err := redis.String(val, err)
		if err == nil {
			return val, nil
		}
		return nil, err
	}
	return nil, err
}

func (cr *cacheRedis) Set(redisPool *redis.Pool, appName string, appEnv string, key string, value string) {
	if redisPool.Get() != nil {
		key = appName + "-" + appEnv + "-" + key
		log.Printf("Set cache %v ", key)
		_, err := redisPool.Get().Do("SET", key, value, "")
		if err != nil {
			log.Err(err).Send()
		}

		_, err = redisPool.Get().Do("EXPIRE", key, 1*60*60)
		if err != nil {
			if err != nil {
				log.Err(err).Send()
			}
		}
	}

}

func (cr *cacheRedis) SetExpired(redisPool *redis.Pool, appName string, appEnv string, key string, value string, expired int) {
	if redisPool.Get() != nil {
		key = appName + "-" + appEnv + "-" + key
		log.Printf("Set cache %v ", key)
		_, err := redisPool.Get().Do("SET", key, value)
		if err != nil {
			return
		}
		if err != nil {
			log.Err(err).Send()
		}

		_, err = redisPool.Get().Do("EXPIRE", key, expired)
		if err != nil {
			return
		}
		if err != nil {
			log.Err(err).Send()
		}
	}
}

func (cr *cacheRedis) Delete(redisPool *redis.Pool, appName string, appEnv string, key string) {
	if redisPool.Get() != nil {
		key = appName + "-" + appEnv + "-" + key
		log.Printf("Delete cache %v ", key)
		_, err := redisPool.Get().Do("DEL", key)
		if err != nil {
			log.Err(err).Send()
		}
	}
}
