package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"reflect"
)

var Config AppConfig

type AppConfig struct {
	AppName                string `mapstructure:"NAME" validate:"required"`
	AppEnv                 string `mapstructure:"ENV" validate:"required"`
	AppHost                string `mapstructure:"HOST" validate:"required"`
	AppPort                int    `mapstructure:"PORT" validate:"required"`
	AppMode                string `mapstructure:"MODE" validate:"required"`
	JwtSecret              string `mapstructure:"JWT_SECRET" validate:"required"`
	JwtTokenExpired        int64  `mapstructure:"JWT_TOKEN_EXPIRED" validate:"required"`
	JwtRefreshTokenExpired int64  `mapstructure:"JWT_REFRESH_TOKEN_EXPIRED" validate:"required"`
	RedisHost              string `mapstructure:"REDIS_HOST" validate:"required"`
	RedisPort              int    `mapstructure:"REDIS_PORT" validate:"required"`
	RedisPassword          string `mapstructure:"REDIS_PASSWORD"`
	NatsAddress            string `mapstructure:"NATS_ADDRESS" validate:"required"`
}

func InitConfig() {
	cfg := AppConfig{}
	v := reflect.ValueOf(cfg)
	for i := 0; i < v.NumField(); i++ {
		f := v.Type().Field(i)
		t := f.Tag.Get("mapstructure")
		viper.Set(getOsEnv(t))
	}

	err := viper.Unmarshal(&cfg)
	if err != nil {
		panic(fmt.Errorf("err Failed to load config: %+v", err))
	}

	validate := validator.New()
	err = validate.Struct(cfg)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		panic(fmt.Errorf("err invalid configuration %s", validationErrors[0].Error()))
	}

	Config = cfg
	log.Printf("Config %+v", cfg)
}

func getOsEnv(name string) (key string, value string) {
	key = name
	value = os.Getenv(name)
	return
}
