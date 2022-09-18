package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type ConfigRepository struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Repository struct {
	Db *sqlx.DB
}

// NewRepository инициализация работы с БД
func NewRepository(c ConfigRepository) (*Repository, error) {
	connect := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Name)
	db, err := sqlx.Connect("postgres", connect)
	if err != nil {
		return nil, fmt.Errorf("failed to connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("connection test error: %w", err)
	}
	return &Repository{Db: db}, err
}

// CloseRepository закрытие подключения к БД
func (r *Repository) CloseRepository() error {
	err := r.Db.Close()
	if err != nil {
		return fmt.Errorf("failed close connection: %w", err)
	}
	return nil
}
