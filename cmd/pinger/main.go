package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/antgubarev/pingbot/internal"
	"github.com/antgubarev/pingbot/internal/pinger"
	"github.com/antgubarev/pingbot/internal/storage"
	"github.com/antgubarev/pingbot/internal/storage/redis"
)

func main() {
	logger := internal.NewLogger()
	defer logger.Info("Buy")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gfShutdown := make(chan os.Signal, 1)
	signal.Notify(gfShutdown, syscall.SIGTERM, syscall.SIGINT)

	redisOpt, err := redis.ParseRedisOpt()
	if err != nil {
		logger.Fatalf("Parse redis config: %v", err)
		return
	}
	responseStorage := redis.NewRedisResponseStorage(redisOpt)

	targetStorage := storage.NewTargetFileProvider()
	pinger := pinger.NewPinger(logger, targetStorage, responseStorage)
	go pinger.Run(ctx)

	<-gfShutdown
}
