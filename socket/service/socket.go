package logic

import (
	cdto "common/dto"
	"common/utils"
	"encoding/json"
	"fmt"
	"github.com/cristalhq/jwt/v3"
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/rs/zerolog/log"
	"socket/config"
	"socket/dto"
	"socket/errors"
	"socket/storage"
	"strings"
	"time"
)

var SocketService = NewSocketService()

type ISocketService interface {
	AddUserToPrivateRoom(c *gosocketio.Channel)
	ValidateRequest(c *gosocketio.Channel) (err error)
	HandleQrCodeUpdate(payload cdto.WalletTransactionQrcodesMessage)
}

func NewSocketService() ISocketService {
	return &socketService{
		channels: map[string]*gosocketio.Channel{},
	}
}

type socketService struct {
	channels map[string]*gosocketio.Channel
}

func (s socketService) ValidateRequest(c *gosocketio.Channel) (err error) {
	header := c.RequestHeader()
	authHeader := strings.Split(header.Get("Authorization"), "Bearer ")
	qrID := header.Get("x-qrcodesId")

	if len(authHeader) != 2 {
		s.closeAndEmitError(c, errors.SocketErrorAuthorization(fmt.Errorf("error invalid authorization header")))
		return
	}

	tokenStr := authHeader[1]
	jwtToken, err := jwt.ParseString(authHeader[1])
	if err != nil {
		log.Err(err).Send()
		s.closeAndEmitError(c, errors.SocketErrorAuthorization(err))
		return
	}

	var claims jwt.StandardClaims
	err = json.Unmarshal(jwtToken.RawClaims(), &claims)
	if err != nil {
		log.Err(err).Send()
		s.closeAndEmitError(c, errors.SocketErrorAuthorization(err))
		return err
	}
	if claims.ExpiresAt == nil {
		err := fmt.Errorf("claims.ExpiresAt nil")
		log.Err(err).Send()
		s.closeAndEmitError(c, errors.SocketErrorAuthorization(err))
		return err
	}

	if claims.ExpiresAt.Before(time.Now()) {
		err := fmt.Errorf("token expired")
		log.Err(err).Send()
		s.closeAndEmitError(c, errors.SocketErrorAuthorization(err))
		return err
	}

	if err = utils.JwtVerify(tokenStr, config.Config.JwtSecret); err != nil {
		log.Err(err).Send()
		s.closeAndEmitError(c, errors.SocketErrorAuthorization(err))
		return
	}

	if len(qrID) == 0 {
		err := fmt.Errorf("invalid header x-qrcodesId")
		log.Err(err).Send()
		s.closeAndEmitError(c, errors.SocketErrorInvalidHeaderXQrCodesId)
		return err
	}

	return
}

func (s socketService) closeAndEmitError(c *gosocketio.Channel, se dto.SocketError) {
	dataError, _ := json.Marshal(se)
	err := c.Emit("/error", string(dataError))
	if err != nil {
		log.Err(err).Send()
		time.Sleep(1 * time.Second)
		delete(s.channels, c.Id())
		c.Close()
		return
	}

	time.Sleep(1 * time.Second)
	delete(s.channels, c.Id())
	c.Close()
}

func (s socketService) AddUserToPrivateRoom(c *gosocketio.Channel) {
	header := c.RequestHeader()
	qrcodesId := header.Get("x-qrcodesId")

	if len(qrcodesId) == 0 {
		err := fmt.Errorf("error invalid header x-qrcodesId")
		log.Err(err).Send()
		s.closeAndEmitError(c, errors.SocketErrorInvalidHeaderXQrCodesId)
		return
	}

	s.channels[c.Id()] = c
	newRoom := fmt.Sprintf("room-%s", c.Id())
	c.Join(newRoom)
	c.BroadcastTo(newRoom, "/message", fmt.Sprintf("User %s join to room %s", c.Id(), newRoom))

	go func() {
		time.Sleep(time.Minute * 5)
		s.closeAndEmitError(c, errors.SocketErrorTimeToLive)
	}()

	pool := storage.Database.Redis
	exp := 1 * 60 * 60 * 24 //1day
	utils.CacheRedis.SetExpired(pool, config.Config.AppName, config.Config.AppEnv, qrcodesId, newRoom, exp)
}

func (s socketService) HandleQrCodeUpdate(payload cdto.WalletTransactionQrcodesMessage) {
	log.Printf("check %+v", payload)

	//ignore status != success / failed
	if strings.ToLower(payload.QrcodesStatus) != "s" && strings.ToLower(payload.QrcodesStatus) != "f" {
		return
	}

	//Check private room by qrcode id
	pool := storage.Database.Redis
	val, err := utils.CacheRedis.Get(pool, config.Config.AppName, config.Config.AppEnv, payload.QrcodesID)
	if err != nil {
		log.Err(err).Send()
		return
	}

	log.Printf("result for key %s is %s", payload.QrcodesID, val)
	if val != "" {
		data, _ := json.Marshal(payload)
		socketId := strings.Replace(val.(string), "room-", "", 1)
		c := s.channels[socketId]
		c.BroadcastTo(val.(string), "/status", string(data))

		//Delay 3 second to close client connection
		time.Sleep(3 * time.Second)
		c.Close()
		delete(s.channels, c.Id())
	}
}
