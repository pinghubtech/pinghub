package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/antgubarev/pingbot/internal"
	"github.com/antgubarev/pingbot/internal/handler"
	"github.com/antgubarev/pingbot/internal/storage/redis"
)

func main() {

	logger := internal.NewLogger()
	defer logger.Info("Buy")

	gfShutdown := make(chan os.Signal, 1)
	signal.Notify(gfShutdown, syscall.SIGTERM, syscall.SIGINT)

	redisOpt, err := redis.ParseRedisOpt()
	if err != nil {
		logger.Fatalf("Parse redis config: %v", err)
		return
	}
	responseStorage := redis.NewRedisResponseStorage(redisOpt)

	httpServerOpt, err := handler.ParseHTTPServerOptions()
	if err != nil {
		logger.Fatalf("Parse HTTP server config: %v", err)
	}
	statusHandler := handler.NewStatusHandler(responseStorage, logger)
	logger.Infof("Starting http server on %s \n", httpServerOpt.ListenAddr)
	go func() {
		if err := http.ListenAndServe(httpServerOpt.ListenAddr, statusHandler); err != nil {
			logger.Fatalf("Http server didn't start: %v", err)
		}
	}()

	<-gfShutdown
}
