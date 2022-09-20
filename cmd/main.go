package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Hymiside/hezzl-test-task/pkg/natsqueue"

	"github.com/Hymiside/hezzl-test-task/pkg/config"
	"github.com/Hymiside/hezzl-test-task/pkg/handler"
	"github.com/Hymiside/hezzl-test-task/pkg/repository/postgres"
	"github.com/Hymiside/hezzl-test-task/pkg/repository/redis"
	"github.com/Hymiside/hezzl-test-task/pkg/server"
	"github.com/Hymiside/hezzl-test-task/pkg/service"
)

// main инициализирует работу всего сервиса
func main() {
	cfgSrv, cfgDb, cfgRd, cfgN := config.InitConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := &server.Server{}
	h := &handler.Handler{}

	rdb, err := redis.NewRepository(ctx, cfgRd)
	if err != nil {
		log.Panicf("falied to create redis reository:%v", err)
	}

	repo, err := postgres.NewRepository(ctx, cfgDb)
	if err != nil {
		log.Panicf("failed to create postgres repository: %v", err)
	}

	nc, err := natsqueue.NewNats(ctx, cfgN, repo)
	if err != nil {
		log.Panicf("failed to create nats client: %v", err)
	}

	services := service.NewService(repo, rdb, nc)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
		select {
		case <-quit:
			cancel()
		case <-ctx.Done():
			return
		}
	}()

	if err = srv.RunServer(ctx, h.InitHandler(*services), cfgSrv); err != nil {
		log.Panicf("failed to run server: %v", err)
	}
	log.Printf("server was stopped")
}
