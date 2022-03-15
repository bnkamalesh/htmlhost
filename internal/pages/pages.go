package pages

import (
	"context"
	"strings"
	"time"

	"github.com/bnkamalesh/errors"
)

type Pages struct {
	store store
}

func (ps *Pages) Create(ctx context.Context, pg *Page) (*Page, error) {
	pg = newPage(pg.Content)
	err := pg.Validate()
	if err != nil {
		return nil, err
	}

	err = ps.store.Create(ctx, pg)
	if err != nil {
		return nil, err
	}

	return pg, nil
}

func (ps *Pages) Read(ctx context.Context, pageID string) (*Page, error) {
	return ps.store.Read(ctx, pageID)
}

func (ps *Pages) Active(ctx context.Context) (int, error) {
	return ps.store.Active(ctx)
}

type Page struct {
	ID        string    `json:"id,omitempty"`
	Content   string    `json:"content,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	Expiry    time.Time `json:"expiry,omitempty"`
}

func (p *Page) Sanitize() {
	p.Content = strings.TrimSpace(p.Content)
}

func (p *Page) Validate() error {
	if p.Content == "" {
		return errors.Validation("content cannot be empty")
	}
	return nil
}

func (p *Page) URL(baseURL string) string {
	return baseURL + "/p/" + p.ID
}

func newPage(content string) *Page {
	now := time.Now()
	pg := &Page{
		ID:        Random(5),
		Content:   content,
		CreatedAt: now,
		Expiry:    now.Add(time.Hour),
	}

	pg.Sanitize()

	return pg
}

func NewService(cfg *Config) (*Pages, error) {
	pgs, err := newStore(cfg)
	if err != nil {
		return nil, err
	}

	return &Pages{
		store: pgs,
	}, nil
}
