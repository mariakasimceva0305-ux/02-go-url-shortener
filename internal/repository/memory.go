package repository

import (
	"database/sql"
	"sort"
	"sync"
	"time"

	"go-url-shortener/internal/model"
)

type MemoryLinkRepository struct {
	mu      sync.RWMutex
	byID    map[string]model.ShortLink
	byCode  map[string]string
	ordered []string
}

func NewMemoryLinkRepository() *MemoryLinkRepository {
	return &MemoryLinkRepository{
		byID:   make(map[string]model.ShortLink),
		byCode: make(map[string]string),
	}
}

func (r *MemoryLinkRepository) Init() error { return nil }

func (r *MemoryLinkRepository) Create(link model.ShortLink) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byID[link.ID] = link
	r.byCode[link.ShortCode] = link.ID
	r.ordered = append(r.ordered, link.ID)
	return nil
}

func (r *MemoryLinkRepository) GetByID(id string) (*model.ShortLink, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	link, ok := r.byID[id]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return &link, nil
}

func (r *MemoryLinkRepository) GetByCode(code string) (*model.ShortLink, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.byCode[code]
	if !ok {
		return nil, sql.ErrNoRows
	}
	link := r.byID[id]
	return &link, nil
}

func (r *MemoryLinkRepository) List() ([]model.ShortLink, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]model.ShortLink, 0, len(r.byID))
	for _, item := range r.byID {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})
	return items, nil
}

func (r *MemoryLinkRepository) Deactivate(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	link, ok := r.byID[id]
	if !ok {
		return sql.ErrNoRows
	}
	link.IsActive = false
	r.byID[id] = link
	return nil
}

func (r *MemoryLinkRepository) CleanupExpired(now time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for id, link := range r.byID {
		if link.ExpiresAt.Before(now) {
			link.IsActive = false
			r.byID[id] = link
		}
	}
	return nil
}

type MemoryEventRepository struct {
	mu     sync.RWMutex
	events []model.VisitEvent
}

func NewMemoryEventRepository() *MemoryEventRepository {
	return &MemoryEventRepository{events: make([]model.VisitEvent, 0)}
}

func (r *MemoryEventRepository) Init() error { return nil }

func (r *MemoryEventRepository) Create(event model.VisitEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, event)
	return nil
}

func (r *MemoryEventRepository) ListByLinkID(linkID string, limit int) ([]model.VisitEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]model.VisitEvent, 0)
	for i := len(r.events) - 1; i >= 0; i-- {
		if r.events[i].ShortLinkID == linkID {
			out = append(out, r.events[i])
			if len(out) == limit {
				break
			}
		}
	}
	return out, nil
}

func (r *MemoryEventRepository) CountByLinkID(linkID string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	count := 0
	for _, event := range r.events {
		if event.ShortLinkID == linkID {
			count++
		}
	}
	return count, nil
}

func (r *MemoryEventRepository) LastClickAt(linkID string) (*time.Time, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := len(r.events) - 1; i >= 0; i-- {
		if r.events[i].ShortLinkID == linkID {
			t := r.events[i].VisitedAt
			return &t, nil
		}
	}
	return nil, nil
}

func (r *MemoryEventRepository) TopUserAgents(linkID string, limit int) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	counts := make(map[string]int)
	for _, event := range r.events {
		if event.ShortLinkID == linkID {
			counts[event.UserAgent]++
		}
	}
	type pair struct {
		ua    string
		count int
	}
	pairs := make([]pair, 0, len(counts))
	for ua, count := range counts {
		pairs = append(pairs, pair{ua: ua, count: count})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})
	out := make([]string, 0, limit)
	for _, item := range pairs {
		if item.ua == "" {
			continue
		}
		out = append(out, item.ua)
		if len(out) == limit {
			break
		}
	}
	return out, nil
}
