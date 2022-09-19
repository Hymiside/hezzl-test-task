package config

import (
	"github.com/Hymiside/hezzl-test-task/pkg/rediscache"
	"os"
	"strconv"

	"github.com/Hymiside/hezzl-test-task/pkg/repository"
	"github.com/Hymiside/hezzl-test-task/pkg/server"
	"github.com/joho/godotenv"
)

func InitConfig() (server.ConfigServer, repository.ConfigRepository, rediscache.ConfigRedis) {
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
	DbString, _ := os.LookupEnv("Db")
	DbInt, _ := strconv.Atoi(DbString)

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
		DB:   DbInt,
	}
	return config, configDb, configRedis
}
