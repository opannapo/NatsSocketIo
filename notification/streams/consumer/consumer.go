package consumer

import (
	"common"
	"fmt"
	"github.com/rs/zerolog/log"
	"notification/streams"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartConsumer() (err error) {
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	emailConsumer := newEmailConsumer()
	_, err = streams.MessageBroker.Nats.Subscribe(common.SubjectSendEmail, emailConsumer.Send)
	if err != nil {
		log.Err(err).Send()
		return
	}

	<-interrupt
	time.Sleep(time.Second)
	fmt.Println("Notification Consumer Stopped.")
	return
}
