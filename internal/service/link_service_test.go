package service

import (
	"testing"
	"time"

	"go-url-shortener/internal/model"
)

type fakeLinkRepo struct {
	link model.ShortLink
}

func (f *fakeLinkRepo) Init() error { return nil }
func (f *fakeLinkRepo) Create(link model.ShortLink) error {
	f.link = link
	return nil
}
func (f *fakeLinkRepo) GetByID(id string) (*model.ShortLink, error)   { return &f.link, nil }
func (f *fakeLinkRepo) GetByCode(code string) (*model.ShortLink, error) { return &f.link, nil }
func (f *fakeLinkRepo) List() ([]model.ShortLink, error) { return []model.ShortLink{f.link}, nil }
func (f *fakeLinkRepo) Deactivate(id string) error        { return nil }
func (f *fakeLinkRepo) CleanupExpired(now time.Time) error { return nil }

type fakeEventRepo struct{}

func (f *fakeEventRepo) Init() error                                              { return nil }
func (f *fakeEventRepo) Create(event model.VisitEvent) error                      { return nil }
func (f *fakeEventRepo) ListByLinkID(linkID string, limit int) ([]model.VisitEvent, error) { return nil, nil }
func (f *fakeEventRepo) CountByLinkID(linkID string) (int, error)                 { return 1, nil }
func (f *fakeEventRepo) LastClickAt(linkID string) (*time.Time, error)            { return nil, nil }
func (f *fakeEventRepo) TopUserAgents(linkID string, limit int) ([]string, error) { return []string{}, nil }

func TestCreateValidatesURL(t *testing.T) {
	svc := NewLinkService(&fakeLinkRepo{}, &fakeEventRepo{}, "http://localhost:8080")
	if _, err := svc.Create("bad-url", 1); err == nil {
		t.Fatal("expected invalid url error")
	}
}

func TestCreateSetsDefaults(t *testing.T) {
	repo := &fakeLinkRepo{}
	svc := NewLinkService(repo, &fakeEventRepo{}, "http://localhost:8080")
	link, err := svc.Create("https://example.com", 0)
	if err != nil {
		t.Fatal(err)
	}
	if link.ShortCode == "" {
		t.Fatal("short code must be generated")
	}
}
