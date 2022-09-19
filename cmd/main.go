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
	"github.com/Hymiside/hezzl-test-task/pkg/rediscache"
	"github.com/Hymiside/hezzl-test-task/pkg/repository"
	"github.com/Hymiside/hezzl-test-task/pkg/server"
	"github.com/Hymiside/hezzl-test-task/pkg/service"
)

func main() {
	cfgSrv, cfgDb, cfgRd := config.InitConfig()

	srv := &server.Server{}
	h := &handler.Handler{}

	rdb, err := rediscache.NewRedis(cfgRd)
	if err != nil {
		log.Fatalf(err.Error())
	}

	repo, err := repository.NewRepository(cfgDb)
	if err != nil {
		log.Fatalf(err.Error())
	}
	services, _ := service.NewService(*repo, *rdb)

	go func() {
		if err = srv.RunServer(h.InitHandler(*services), cfgSrv); err != nil {
			log.Fatalf(err.Error())
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

	if err = srv.ShutdownServer(ctx); err != nil {
		log.Fatalf(err.Error())
	}
	if err = repo.CloseRepository(); err != nil {
		log.Fatalf(err.Error())
	}
	if err = rdb.CloseRedis(); err != nil {
		log.Fatalf(err.Error())
	}
}
