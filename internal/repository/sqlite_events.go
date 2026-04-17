package repository

import (
	"database/sql"
	"time"

	"go-url-shortener/internal/model"
)

type SQLiteEventRepository struct {
	DB *sql.DB
}

func (r *SQLiteEventRepository) Init() error {
	_, err := r.DB.Exec(`
		CREATE TABLE IF NOT EXISTS visit_events (
			id TEXT PRIMARY KEY,
			short_link_id TEXT NOT NULL,
			visited_at DATETIME NOT NULL,
			user_agent TEXT,
			ip_address TEXT,
			referrer TEXT
		);
	`)
	return err
}

func (r *SQLiteEventRepository) Create(event model.VisitEvent) error {
	_, err := r.DB.Exec(
		`INSERT INTO visit_events (id, short_link_id, visited_at, user_agent, ip_address, referrer) VALUES (?, ?, ?, ?, ?, ?)`,
		event.ID, event.ShortLinkID, event.VisitedAt, event.UserAgent, event.IPAddress, event.Referrer,
	)
	return err
}

func (r *SQLiteEventRepository) ListByLinkID(linkID string, limit int) ([]model.VisitEvent, error) {
	rows, err := r.DB.Query(`SELECT id, short_link_id, visited_at, user_agent, ip_address, referrer FROM visit_events WHERE short_link_id = ? ORDER BY visited_at DESC LIMIT ?`, linkID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]model.VisitEvent, 0)
	for rows.Next() {
		var event model.VisitEvent
		if err := rows.Scan(&event.ID, &event.ShortLinkID, &event.VisitedAt, &event.UserAgent, &event.IPAddress, &event.Referrer); err != nil {
			return nil, err
		}
		out = append(out, event)
	}
	return out, nil
}

func (r *SQLiteEventRepository) CountByLinkID(linkID string) (int, error) {
	row := r.DB.QueryRow(`SELECT COUNT(*) FROM visit_events WHERE short_link_id = ?`, linkID)
	var count int
	return count, row.Scan(&count)
}

func (r *SQLiteEventRepository) LastClickAt(linkID string) (*time.Time, error) {
	row := r.DB.QueryRow(`SELECT visited_at FROM visit_events WHERE short_link_id = ? ORDER BY visited_at DESC LIMIT 1`, linkID)
	var t time.Time
	if err := row.Scan(&t); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *SQLiteEventRepository) TopUserAgents(linkID string, limit int) ([]string, error) {
	rows, err := r.DB.Query(`SELECT user_agent, COUNT(*) AS c FROM visit_events WHERE short_link_id = ? GROUP BY user_agent ORDER BY c DESC LIMIT ?`, linkID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]string, 0)
	for rows.Next() {
		var ua string
		var count int
		if err := rows.Scan(&ua, &count); err != nil {
			return nil, err
		}
		result = append(result, ua)
	}
	return result, nil
}
