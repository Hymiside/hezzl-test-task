package repository

import (
	"fmt"
	"github.com/Hymiside/hezzl-test-task/pkg/models"
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

// CreateItem создает новый элемент в таблице
func (r *Repository) CreateItem(ni models.NewItem) (models.Item, error) {
	var (
		rows   *sqlx.Rows
		i      models.Item
		itemId int
	)

	tx, _ := r.Db.Begin()
	err := r.Db.QueryRowx(`INSERT INTO items (campaign_id, name, description, removed, created_at) 
							   VALUES ($1, $2, $3, $4, $5) RETURNING id;`, ni.CampaignId, ni.Name, ni.Description, ni.Removed, ni.CreatedAt).Scan(&itemId)
	if err != nil {
		return models.Item{}, fmt.Errorf("error create item: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return models.Item{}, fmt.Errorf("error commit item: %w", err)
	}

	//

	rows, err = r.Db.Queryx(`SELECT * FROM items WHERE id=$1 AND campaign_id=$2;`, itemId, ni.CampaignId)
	if err != nil {
		return models.Item{}, fmt.Errorf("error return item: %w", err)
	}

	for rows.Next() {
		err = rows.StructScan(&i)
		if err != nil {
			return models.Item{}, fmt.Errorf("error return item: %w", err)
		}
	}
	return i, nil
}

// GetCampaign ищет кампанию по ее id
func (r *Repository) GetCampaign(ni models.NewItem) error {
	var id int

	tx := r.Db.MustBegin()
	if err := tx.QueryRowx(`SELECT id FROM campaigns WHERE id = $1`, ni.CampaignId).Scan(&id); err != nil {
		return fmt.Errorf("error item not found: %w", err)
	}
	return nil
}

func (r *Repository) GetItems() ([]models.Item, error) {
	var i []models.Item

	if err := r.Db.Select(&i, `SELECT * FROM items;`); err != nil {
		return nil, fmt.Errorf("error return item: %w", err)
	}
	return i, nil
}
