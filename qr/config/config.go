package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
)

var Config AppConfig

type AppConfig struct {
	AppName                string `mapstructure:"NAME" validate:"required"`
	AppEnv                 string `mapstructure:"ENV" validate:"required"`
	AppHost                string `mapstructure:"HOST" validate:"required"`
	AppPort                int    `mapstructure:"PORT" validate:"required"`
	AppSecret              string `mapstructure:"SECRET" validate:"required"`
	AppBasicAuthUser       string `mapstructure:"BASIC_AUTH_USER" validate:"required"`
	AppBasicAuthPassword   string `mapstructure:"BASIC_AUTH_PASSWORD" validate:"required"`
	AppMode                string `mapstructure:"MODE" validate:"required"`
	MysqlHost              string `mapstructure:"DB_HOST" validate:"required"`
	MysqlPort              int    `mapstructure:"DB_PORT" validate:"required"`
	MysqlDatabase          string `mapstructure:"DB_NAME" validate:"required"`
	MysqlUsername          string `mapstructure:"DB_USERNAME" validate:"required"`
	MysqlPassword          string `mapstructure:"DB_PASSWORD" validate:"required"`
	NatsAddress            string `mapstructure:"NATS_ADDRESS" validate:"required"`
	JwtSecret              string `mapstructure:"JWT_SECRET" validate:"required"`
	JwtTokenExpired        int64  `mapstructure:"JWT_TOKEN_EXPIRED" validate:"required"`
	JwtRefreshTokenExpired int64  `mapstructure:"JWT_REFRESH_TOKEN_EXPIRED" validate:"required"`
	SmtpHost               string `mapstructure:"SMTP_HOST" validate:"required"`
	SmtpMail               string `mapstructure:"SMTP_MAIL" validate:"required,email"`
	SmtpPass               string `mapstructure:"SMTP_PASS" validate:"required"`
	SmtpPort               string `mapstructure:"SMTP_PORT" validate:"required"`
	SmtpReplyTo            string `mapstructure:"SMTP_REPLY_TO" validate:"required,email"`
	RedisHost              string `mapstructure:"REDIS_HOST" validate:"required"`
	RedisPort              int    `mapstructure:"REDIS_PORT" validate:"required"`
	RedisPassword          string `mapstructure:"REDIS_PASSWORD"`
}

func InitConfig() {
	cfg := AppConfig{}
	viper.Set(getOsEnv("NAME"))
	viper.Set(getOsEnv("ENV"))
	viper.Set(getOsEnv("HOST"))
	viper.Set(getOsEnv("PORT"))
	viper.Set(getOsEnv("SECRET"))
	viper.Set(getOsEnv("BASIC_AUTH_USER"))
	viper.Set(getOsEnv("BASIC_AUTH_PASSWORD"))
	viper.Set(getOsEnv("MODE"))
	viper.Set(getOsEnv("DB_HOST"))
	viper.Set(getOsEnv("DB_PORT"))
	viper.Set(getOsEnv("DB_NAME"))
	viper.Set(getOsEnv("DB_USERNAME"))
	viper.Set(getOsEnv("DB_PASSWORD"))
	viper.Set(getOsEnv("NATS_ADDRESS"))
	viper.Set(getOsEnv("JWT_SECRET"))
	viper.Set(getOsEnv("JWT_TOKEN_EXPIRED"))
	viper.Set(getOsEnv("JWT_REFRESH_TOKEN_EXPIRED"))
	viper.Set(getOsEnv("SMTP_HOST"))
	viper.Set(getOsEnv("SMTP_MAIL"))
	viper.Set(getOsEnv("SMTP_PASS"))
	viper.Set(getOsEnv("SMTP_PORT"))
	viper.Set(getOsEnv("SMTP_REPLY_TO"))
	viper.Set(getOsEnv("REDIS_HOST"))
	viper.Set(getOsEnv("REDIS_PORT"))
	viper.Set(getOsEnv("REDIS_PASSWORD"))

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
