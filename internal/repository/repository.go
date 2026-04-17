package repository

import (
	"database/sql"
	"time"

	"go-url-shortener/internal/model"
)

type LinkRepository interface {
	Init() error
	Create(link model.ShortLink) error
	GetByID(id string) (*model.ShortLink, error)
	GetByCode(code string) (*model.ShortLink, error)
	List() ([]model.ShortLink, error)
	Deactivate(id string) error
	CleanupExpired(now time.Time) error
}

type EventRepository interface {
	Init() error
	Create(event model.VisitEvent) error
	ListByLinkID(linkID string, limit int) ([]model.VisitEvent, error)
	CountByLinkID(linkID string) (int, error)
	LastClickAt(linkID string) (*time.Time, error)
	TopUserAgents(linkID string, limit int) ([]string, error)
}

type SQLiteRepository struct {
	DB *sql.DB
}
