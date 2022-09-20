package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Hymiside/hezzl-test-task/pkg/models"
)

type ConfigRepository struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Repository struct {
	db *sqlx.DB
}

var (
	ItemAlreadyDelete = errors.New("item already delete")
	ItemNotFound      = errors.New("item not found")
)

// NewRepository инициализация работы с БД
func NewRepository(ctx context.Context, c ConfigRepository) (*Repository, error) {
	connect := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Name)
	db, err := sqlx.Connect("postgres", connect)
	if err != nil {
		return nil, fmt.Errorf("failed to connection: %w", err)
	}

	go func(ctx context.Context) {
		<-ctx.Done()
		db.Close()
	}(ctx)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("connection test error: %w", err)
	}

	return &Repository{db: db}, err
}

// CreateItem создает новый элемент в таблице
func (r *Repository) CreateItem(ctx context.Context, item models.NewItem) (int, int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	// Проверяем наличие кампании (можно это сделать на ключах в БД и без запросов, но в ТЗ описано явно проверять)
	q1 := `select id from campaigns where id = $1`
	row := r.db.QueryRowContext(ctx, q1, item.CampaignId)

	var campaignId int
	if err = row.Scan(&campaignId); err != nil {
		return 0, 0, fmt.Errorf("failed to query campaning: %w", err)
	}

	if err = row.Err(); err != nil {
		return 0, 0, fmt.Errorf("got row error of query campaning: %w", err)
	}

	// Создаем айтем
	q2 := `INSERT INTO items (campaign_id, name, description, removed, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, priority`
	row = r.db.QueryRowContext(ctx, q2, item.CampaignId, item.Name, item.Description, item.Removed, item.CreatedAt)

	var itemId, priority int
	if err = row.Scan(&itemId, &priority); err != nil {
		return 0, 0, fmt.Errorf("failed to query item id: %w", err)
	}

	if err = row.Err(); err != nil {
		return 0, 0, fmt.Errorf("got row error of query item: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, 0, fmt.Errorf("failed to commit: %w", err)
	}

	return itemId, priority, nil
}

// GetItem Получает айтем по идентификаторам кампании и айтема
func (r *Repository) GetItem(ctx context.Context, campaignId, itemId int) (models.Item, error) {
	q := `
		select 
			id, 
			campaign_id, 
			name,
			description,
			priority,
			removed,
			created_at
		from items  
		where id = $1 and campaign_id = $2 and removed = false
		`

	var item models.Item

	row := r.db.QueryRowContext(ctx, q, itemId, campaignId)
	if err := row.Scan(
		&item.ID,
		&item.CampaignId,
		&item.Name,
		&item.Description,
		&item.Priority,
		&item.Removed,
		&item.CreatedAt,
	); err != nil {
		return models.Item{}, fmt.Errorf("failed to query item: %w", err)
	}

	if err := row.Err(); err != nil {
		return models.Item{}, fmt.Errorf("got row error of query item: %w", err)
	}

	return item, nil
}

// GetAllItems Получает все айтемы из БД
func (r *Repository) GetAllItems(ctx context.Context) ([]models.Item, error) {
	q := `
		select
			id, 
			campaign_id, 
			name,
			description,
			priority,
			removed,
			created_at
		from items
		where removed=false
	`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}

	var items []models.Item
	for rows.Next() {
		var item models.Item
		if err = rows.Scan(
			&item.ID,
			&item.CampaignId,
			&item.Name,
			&item.Description,
			&item.Priority,
			&item.Removed,
			&item.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to query item id: %w", err)
		}

		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("got rows error: %w", err)
	}

	return items, nil
}

// UpdateItem Обновляет значение айтема
func (r *Repository) UpdateItem(ctx context.Context, campaignId, itemId int, name, description string) (models.Item, error) {
	q := `update items set name = $3, description = $4 where id = $1 and campaign_id = $2 and removed = false returning id, campaign_id, name, description, priority, removed, created_at`
	var item models.Item

	if description == "" {
		q1 := `select description from items where id = $1 and campaign_id = $2`

		row := r.db.QueryRowContext(ctx, q1, itemId, campaignId)
		if err := row.Scan(&description); err != nil {
			return models.Item{}, fmt.Errorf("failed to query description item: %w", err)
		}
	}

	row := r.db.QueryRowContext(ctx, q, itemId, campaignId, name, description)
	if err := row.Scan(
		&item.ID,
		&item.CampaignId,
		&item.Name,
		&item.Description,
		&item.Priority,
		&item.Removed,
		&item.CreatedAt,
	); err != nil {
		return models.Item{}, fmt.Errorf("failed to query item: %w", err)
	}

	if err := row.Err(); err != nil {
		return models.Item{}, fmt.Errorf("got row error of query item: %w", err)
	}

	return item, nil
}

// DeleteItem Удаляет айтем
func (r *Repository) DeleteItem(ctx context.Context, campaignId, itemId int) error {
	q := `select removed from items where id = $1 and campaign_id = $2`
	q1 := `update items set removed = true where id = $1 and campaign_id = $2`

	var status bool
	row := r.db.QueryRowContext(ctx, q, itemId, campaignId)
	if err := row.Scan(&status); err != nil {
		return fmt.Errorf("failed to query item: %w", err)
	}
	if status == true {
		return ItemAlreadyDelete
	}

	res, err := r.db.ExecContext(ctx, q1, itemId, campaignId)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}
	deletedLines, _ := res.RowsAffected()
	if deletedLines == 0 {
		return ItemNotFound
	}
	return nil
}