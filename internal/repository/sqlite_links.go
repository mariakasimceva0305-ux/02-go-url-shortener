package repository

import (
	"database/sql"
	"time"

	"go-url-shortener/internal/model"
)

type SQLiteLinkRepository struct {
	DB *sql.DB
}

func (r *SQLiteLinkRepository) Init() error {
	_, err := r.DB.Exec(`
		CREATE TABLE IF NOT EXISTS short_links (
			id TEXT PRIMARY KEY,
			original_url TEXT NOT NULL,
			short_code TEXT NOT NULL UNIQUE,
			created_at DATETIME NOT NULL,
			expires_at DATETIME NOT NULL,
			is_active BOOLEAN NOT NULL
		);
	`)
	return err
}

func (r *SQLiteLinkRepository) Create(link model.ShortLink) error {
	_, err := r.DB.Exec(
		`INSERT INTO short_links (id, original_url, short_code, created_at, expires_at, is_active) VALUES (?, ?, ?, ?, ?, ?)`,
		link.ID, link.OriginalURL, link.ShortCode, link.CreatedAt, link.ExpiresAt, link.IsActive,
	)
	return err
}

func (r *SQLiteLinkRepository) GetByID(id string) (*model.ShortLink, error) {
	row := r.DB.QueryRow(`SELECT id, original_url, short_code, created_at, expires_at, is_active FROM short_links WHERE id = ?`, id)
	return scanLink(row)
}

func (r *SQLiteLinkRepository) GetByCode(code string) (*model.ShortLink, error) {
	row := r.DB.QueryRow(`SELECT id, original_url, short_code, created_at, expires_at, is_active FROM short_links WHERE short_code = ?`, code)
	return scanLink(row)
}

func (r *SQLiteLinkRepository) List() ([]model.ShortLink, error) {
	rows, err := r.DB.Query(`SELECT id, original_url, short_code, created_at, expires_at, is_active FROM short_links ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]model.ShortLink, 0)
	for rows.Next() {
		var link model.ShortLink
		if err := rows.Scan(&link.ID, &link.OriginalURL, &link.ShortCode, &link.CreatedAt, &link.ExpiresAt, &link.IsActive); err != nil {
			return nil, err
		}
		items = append(items, link)
	}
	return items, nil
}

func (r *SQLiteLinkRepository) Deactivate(id string) error {
	_, err := r.DB.Exec(`UPDATE short_links SET is_active = 0 WHERE id = ?`, id)
	return err
}

func (r *SQLiteLinkRepository) CleanupExpired(now time.Time) error {
	_, err := r.DB.Exec(`UPDATE short_links SET is_active = 0 WHERE expires_at < ?`, now)
	return err
}

func scanLink(row *sql.Row) (*model.ShortLink, error) {
	var link model.ShortLink
	if err := row.Scan(&link.ID, &link.OriginalURL, &link.ShortCode, &link.CreatedAt, &link.ExpiresAt, &link.IsActive); err != nil {
		return nil, err
	}
	return &link, nil
}
