package main

import (
	"bufio"
	"fmt"
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"log"
	"net/http"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Print("Masukan QRID : ")
	reader := bufio.NewReader(os.Stdin)
	for {
		data, _, _ := reader.ReadLine()
		qrID := string(data)
		connectANdListenQRStatus(qrID)
	}
}

func connectANdListenQRStatus(qrID string) {
	//valid tapi expired
	//jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiJhZGYwMTRiNTUwMGY0NWUwNjBmZTJkYjM4OGIzMDFkZSIsImlzcyI6InVzZXJfbG9jYWwiLCJzdWIiOiJhdCIsImV4cCI6MTY2MTk4ODkxMywiaWF0IjoxNjYxOTQ1NzEzfQ._UTaHI6sT2_PyppmFF-ULBc3Vj5Frw87kcS10QbHfH0"

	//jwt invalid
	//jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	//exp 21 april 203
	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiI5YzNjMGNkYmU5ZDEzYzM3OGMzZGU1ZDExOWFhNTRkNiIsImlzcyI6InVzZXJfbG9jYWwiLCJzdWIiOiJhdCIsImV4cCI6MTk2NjE0ODI5NywiaWF0IjoxNjYzNTU2Mjk3fQ.PxV4NEytbI_tLdzAhu9j9c3XfUfb6wrClmj_ggKdTr0"

	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+jwt)
	headers.Add("x-qrcodesId", qrID)

	tr := transport.WebsocketTransport{
		PingInterval:   transport.WsDefaultPingInterval,
		PingTimeout:    transport.WsDefaultPingTimeout,
		ReceiveTimeout: transport.WsDefaultReceiveTimeout,
		SendTimeout:    transport.WsDefaultSendTimeout,
		BufferSize:     transport.WsDefaultBufferSize,
		RequestHeader:  headers,
	}
	c, err := gosocketio.Dial(gosocketio.GetUrl("localhost:1111/socket/", 1111, false), &tr)
	if err != nil {
		log.Fatal(err)
	}

	err = c.On("/message", func(h *gosocketio.Channel, args string) {
		log.Println("On message : ", args)
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
		log.Fatal("On Disconnected")
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		log.Println("On Connected")
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On("/status", func(h *gosocketio.Channel, args string) {
		log.Println("On status : ", args)
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On("/error", func(h *gosocketio.Channel, args string) {
		log.Println("On Error : ", args)
	})
	if err != nil {
		log.Fatal(err)
	}
}
