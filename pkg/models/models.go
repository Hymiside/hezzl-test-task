package models

import (
	"github.com/jmoiron/sqlx"
	"net/http"
)

type ConfigServer struct {
	host string
	port string
}

type ConfigRepository struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Server struct {
	server *http.Server
}

type Repository struct {
	db *sqlx.DB
}
