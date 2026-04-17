package service

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"go-url-shortener/internal/model"
	"go-url-shortener/internal/repository"
)

var ErrInvalidURL = errors.New("invalid url")
var ErrLinkExpired = errors.New("link expired")
var ErrLinkInactive = errors.New("link inactive")

type LinkService struct {
	links   repository.LinkRepository
	events  repository.EventRepository
	mu      sync.Mutex
	clicks  map[string]int
	baseURL string
}

func NewLinkService(links repository.LinkRepository, events repository.EventRepository, baseURL string) *LinkService {
	return &LinkService{
		links:   links,
		events:  events,
		clicks:  make(map[string]int),
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (s *LinkService) Create(originalURL string, ttlHours int) (*model.ShortLink, error) {
	if !isValidURL(originalURL) {
		return nil, ErrInvalidURL
	}
	if ttlHours <= 0 {
		ttlHours = 24
	}
	code := generateShortCode(originalURL + time.Now().UTC().String())
	link := model.ShortLink{
		ID:          uuid.NewString(),
		OriginalURL: originalURL,
		ShortCode:   code,
		CreatedAt:   time.Now().UTC(),
		ExpiresAt:   time.Now().UTC().Add(time.Duration(ttlHours) * time.Hour),
		IsActive:    true,
	}
	if err := s.links.Create(link); err != nil {
		return nil, err
	}
	return &link, nil
}

func (s *LinkService) ResolveByCode(code string, userAgent string, ip string, referrer string) (*model.ShortLink, error) {
	link, err := s.links.GetByCode(code)
	if err != nil {
		return nil, err
	}
	if !link.IsActive {
		return nil, ErrLinkInactive
	}
	if time.Now().UTC().After(link.ExpiresAt) {
		return nil, ErrLinkExpired
	}

	s.mu.Lock()
	s.clicks[link.ID]++
	s.mu.Unlock()

	_ = s.events.Create(model.VisitEvent{
		ID:          uuid.NewString(),
		ShortLinkID: link.ID,
		VisitedAt:   time.Now().UTC(),
		UserAgent:   userAgent,
		IPAddress:   ip,
		Referrer:    referrer,
	})
	return link, nil
}

func (s *LinkService) Analytics(id string) (*model.LinkAnalytics, error) {
	link, err := s.links.GetByID(id)
	if err != nil {
		return nil, err
	}
	total, err := s.events.CountByLinkID(id)
	if err != nil {
		return nil, err
	}
	lastClickAt, err := s.events.LastClickAt(id)
	if err != nil {
		return nil, err
	}
	topUA, err := s.events.TopUserAgents(id, 3)
	if err != nil {
		return nil, err
	}
	recent, err := s.events.ListByLinkID(id, 10)
	if err != nil {
		return nil, err
	}
	return &model.LinkAnalytics{
		LinkID:        id,
		OriginalURL:   link.OriginalURL,
		ShortCode:     link.ShortCode,
		TotalClicks:   total,
		LastClickAt:   lastClickAt,
		TopUserAgents: topUA,
		RecentEvents:  recent,
	}, nil
}

func (s *LinkService) List() ([]model.ShortLink, error) {
	return s.links.List()
}

func (s *LinkService) Get(id string) (*model.ShortLink, error) {
	return s.links.GetByID(id)
}

func (s *LinkService) Deactivate(id string) error {
	return s.links.Deactivate(id)
}

func (s *LinkService) CleanupExpired() error {
	return s.links.CleanupExpired(time.Now().UTC())
}

func (s *LinkService) ShortURL(code string) string {
	return s.baseURL + "/r/" + code
}

func generateShortCode(source string) string {
	hash := sha1.Sum([]byte(source))
	return hex.EncodeToString(hash[:])[:8]
}

func isValidURL(raw string) bool {
	u, err := url.ParseRequestURI(raw)
	return err == nil && u.Scheme != "" && u.Host != ""
}
