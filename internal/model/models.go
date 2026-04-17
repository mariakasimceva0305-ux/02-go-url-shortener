package model

import "time"

type ShortLink struct {
	ID          string    `json:"id"`
	OriginalURL string    `json:"original_url"`
	ShortCode   string    `json:"short_code"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	IsActive    bool      `json:"is_active"`
}

type VisitEvent struct {
	ID          string    `json:"id"`
	ShortLinkID string    `json:"short_link_id"`
	VisitedAt   time.Time `json:"visited_at"`
	UserAgent   string    `json:"user_agent"`
	IPAddress   string    `json:"ip_address"`
	Referrer    string    `json:"referrer"`
}

type LinkAnalytics struct {
	LinkID        string       `json:"link_id"`
	OriginalURL   string       `json:"original_url"`
	ShortCode     string       `json:"short_code"`
	TotalClicks   int          `json:"total_clicks"`
	LastClickAt   *time.Time   `json:"last_click_at,omitempty"`
	TopUserAgents []string     `json:"top_user_agents"`
	RecentEvents  []VisitEvent `json:"recent_events"`
}
