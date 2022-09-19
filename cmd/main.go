package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Hymiside/hezzl-test-task/pkg/config"
	"github.com/Hymiside/hezzl-test-task/pkg/handler"
	"github.com/Hymiside/hezzl-test-task/pkg/natsqueue"
	"github.com/Hymiside/hezzl-test-task/pkg/rediscache"
	"github.com/Hymiside/hezzl-test-task/pkg/repository"
	"github.com/Hymiside/hezzl-test-task/pkg/server"
	"github.com/Hymiside/hezzl-test-task/pkg/service"
)

func main() {
	cfgSrv, cfgDb, cfgRd, cfgN := config.InitConfig()

	srv := &server.Server{}
	h := &handler.Handler{}

	rdb, err := rediscache.NewRedis(cfgRd)
	if err != nil {
		log.Panicf(err.Error())
	}

	repo, err := repository.NewRepository(cfgDb)
	if err != nil {
		log.Panicf(err.Error())
	}

	nc, err := natsqueue.NewNats(cfgN)
	if err != nil {
		log.Panicf(err.Error())
	}

	services := service.NewService(*repo, *rdb, *nc)

	go func() {
		if err = srv.RunServer(h.InitHandler(*services), cfgSrv); err != nil {
			log.Panicf(err.Error())
		}
	}()
	log.Printf("authentication microservice launched on http://%s:%s/", cfgSrv.Host, cfgSrv.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	nc.CloseNats()
	if err = srv.ShutdownServer(ctx); err != nil {
		log.Panicf(err.Error())
	}
	if err = repo.CloseRepository(); err != nil {
		log.Panicf(err.Error())
	}
	if err = rdb.CloseRedis(); err != nil {
		log.Panicf(err.Error())
	}
}
