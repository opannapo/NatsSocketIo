package server

import (
	"fmt"
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"github.com/rs/zerolog/log"
	"net/http"
	"socket/api"
	"socket/config"
)

func StartServer() {
	startSocketServer()
}

func startSocketServer() {
	tr := transport.WebsocketTransport{
		PingInterval:   transport.WsDefaultPingInterval,
		PingTimeout:    transport.WsDefaultPingTimeout,
		ReceiveTimeout: transport.WsDefaultReceiveTimeout,
		SendTimeout:    transport.WsDefaultSendTimeout,
		BufferSize:     transport.WsDefaultBufferSize,
	}

	server := gosocketio.NewServer(&tr)
	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)

	server.On(gosocketio.OnConnection, api.SocketHandler.OnConnect)
	server.On(gosocketio.OnDisconnection, api.SocketHandler.OnDisconnect)

	log.Print("Starting socket server...")
	log.Print(http.ListenAndServe(fmt.Sprintf("%s:%d", config.Config.AppHost, config.Config.AppPort), serveMux))
}
