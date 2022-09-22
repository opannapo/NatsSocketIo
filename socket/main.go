package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"socket/cmd"
	"socket/config"
)

func init() {
	config.InitConfig()

	isDebugMode := config.Config.AppMode == "debug"

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.With().Caller().Logger()
	if isDebugMode {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}
}

func main() {
	cmd.Execute()
}
