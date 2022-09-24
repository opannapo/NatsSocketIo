package cmd

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"qr/server"
	"qr/storage"
	"qr/streams"
	"qr/streams/consumer"
)

var rootCmd = &cobra.Command{
	Use:   "report",
	Short: "A brief description of your application",
	Long:  `A longer description that spans multiple lines and likely contains examples and usage of using your application.`,
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func start() {
	//Init DB
	err := storage.Storage.InitDatabase()
	if err != nil {
		log.Err(err).Send()
		return
	}

	//Init Consumer
	err = consumer.StartConsumer()
	if err != nil {
		log.Err(err).Send()
		return
	}
	defer streams.MessageBroker.Nats.Close()

	server.StartServer()
}
