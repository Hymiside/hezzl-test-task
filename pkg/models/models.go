package models

import "time"

type NewItem struct {
	CampaignId  int    `db:"campaign_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Removed     bool
	CreatedAt   time.Time
}

type Item struct {
	ID          int       `db:"id"`
	CampaignId  int       `db:"campaign_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Priority    int       `db:"priority"`
	Removed     bool      `db:"removed"`
	CreatedAt   time.Time `db:"created_at"`
}

type DeleteItem struct {
	ID         int  `db:"id"`
	CampaignId int  `db:"campaign_id"`
	Removed    bool `db:"removed"`
}
