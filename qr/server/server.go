package server

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"qr/config"
	"qr/handler"
	"syscall"
	"time"
)

func StartServer() {
	startInternalServer()
}

func startInternalServer() {
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Config.AppHost, config.Config.AppPort),
		Handler: internalRouting(),
	}

	go func() {
		log.Printf("Server running at %s:%d", config.Config.AppHost, config.Config.AppPort)
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal().Err(err).Send()
			}
		}
	}()

	<-interrupt
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Print("Shutting down ...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Send()
	}
	log.Print("Http Server Stopped.")
}

func internalRouting() *mux.Router {
	r := mux.NewRouter()
	r.Use(BaseMiddleware)
	v1 := r.PathPrefix("/v1").Subrouter().SkipClean(true)

	handler.Post(v1, "/qr/create", handler.Qr.Create)
	handler.Get(v1, "/qr/scan/{id}", handler.Qr.Scan)
	_ = r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, _ := route.GetPathTemplate()
		met, _ := route.GetMethods()
		fmt.Println(met, tpl)
		return nil
	})

	return r
}

func BaseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := uuid.New().String()
		ctx := context.WithValue(r.Context(), "time_request", time.Now())
		ctx = context.WithValue(ctx, "request_id", reqId)
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
