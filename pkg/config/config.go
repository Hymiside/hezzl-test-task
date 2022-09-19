package config

import (
	"github.com/Hymiside/hezzl-test-task/pkg/natsqueue"
	"os"

	"github.com/Hymiside/hezzl-test-task/pkg/rediscache"
	"github.com/Hymiside/hezzl-test-task/pkg/repository"
	"github.com/Hymiside/hezzl-test-task/pkg/server"
	"github.com/joho/godotenv"
)

func InitConfig() (server.ConfigServer, repository.ConfigRepository, rediscache.ConfigRedis, natsqueue.ConfigNats) {
	_ = godotenv.Load()

	host, _ := os.LookupEnv("SERVICE_HOST")
	port, _ := os.LookupEnv("SERVICE_PORT")

	hostDb, _ := os.LookupEnv("DB_HOST")
	portDb, _ := os.LookupEnv("DB_PORT")
	userDb, _ := os.LookupEnv("DB_USER")
	passwordDb, _ := os.LookupEnv("DB_PASSWORD")
	nameDb, _ := os.LookupEnv("DB_NAME")

	hostRd, _ := os.LookupEnv("RD_HOST")
	portRd, _ := os.LookupEnv("RD_PORT")

	hostN, _ := os.LookupEnv("N_HOST")
	portN, _ := os.LookupEnv("N_PORT")

	configDb := repository.ConfigRepository{
		Host:     hostDb,
		Port:     portDb,
		User:     userDb,
		Password: passwordDb,
		Name:     nameDb,
	}

	config := server.ConfigServer{
		Host: host,
		Port: port,
	}

	configRedis := rediscache.ConfigRedis{
		Host: hostRd,
		Port: portRd,
	}

	configNats := natsqueue.ConfigNats{
		Host: hostN,
		Port: portN,
	}

	return config, configDb, configRedis, configNats
}
