package server

import (
	"encoding/json"
	"fmt"
	"github.com/cristalhq/jwt/v3"
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"github.com/rs/zerolog/log"
	"mono/common/utils"
	"net/http"
	"socket/api"
	"socket/config"
	"strings"
	"time"
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

	//no middleware, handling on service layer, checking after soket connected
	//serveMux.Handle("/socket/", server)

	//handling request on middleware layer, checking before sokey connected
	serveMux.Handle("/socket/", ValidateRequest(server))

	server.On(gosocketio.OnConnection, api.SocketHandler.OnConnect)
	server.On(gosocketio.OnDisconnection, api.SocketHandler.OnDisconnect)

	log.Printf(fmt.Sprintf("Server running on %s:%d", config.Config.AppHost, config.Config.AppPort))
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.Config.AppHost, config.Config.AppPort), serveMux)
	if err != nil {
		log.Err(err).Send()
		return
	}
}

func ValidateRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header
		authHeader := strings.Split(header.Get("Authorization"), "Bearer ")
		qrID := header.Get("x-qrcodesId")

		if len(authHeader) != 2 {
			log.Err(fmt.Errorf("error invalid authorization header")).Send()
			w.WriteHeader(401)
			w.Write([]byte(`error invalid authorization header`))
			return
		}

		tokenStr := authHeader[1]
		jwtToken, err := jwt.ParseString(authHeader[1])
		if err != nil {
			log.Err(err).Send()
			w.WriteHeader(401)
			w.Write([]byte(err.Error()))
			return
		}

		var claims jwt.StandardClaims
		err = json.Unmarshal(jwtToken.RawClaims(), &claims)
		if err != nil {
			log.Err(err).Send()
			w.WriteHeader(401)
			w.Write([]byte(err.Error()))
			return
		}
		if claims.ExpiresAt == nil {
			err := fmt.Errorf("claims.ExpiresAt nil")
			log.Err(err).Send()
			w.WriteHeader(401)
			w.Write([]byte(err.Error()))
			return
		}

		if claims.ExpiresAt.Before(time.Now()) {
			err := fmt.Errorf("token expired")
			log.Err(err).Send()
			w.WriteHeader(401)
			w.Write([]byte(err.Error()))
			return
		}

		if err = utils.JwtVerify(tokenStr, config.Config.JwtSecret); err != nil {
			log.Err(err).Send()
			w.WriteHeader(401)
			w.Write([]byte(err.Error()))
			return
		}

		if len(qrID) == 0 {
			err := fmt.Errorf("invalid header x-qrcodesId")
			log.Err(err).Send()
			w.WriteHeader(401)
			w.Write([]byte(err.Error()))
			return
		}

		next.ServeHTTP(w, r)
	})
}
