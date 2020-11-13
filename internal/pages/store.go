package pages

import (
	"context"
	"time"

	"github.com/bnkamalesh/errors"
)

type Config struct {
	Host      string
	Port      string
	StoreName string
	Password  string
	PoolSize  int

	IdleTimeoutSecs  time.Duration
	ReadTimeoutSecs  time.Duration
	WriteTimeoutSecs time.Duration
	DialTimeoutSecs  time.Duration
}

type store interface {
	Create(context.Context, *Page) error
	Read(context.Context, string) (*Page, error)
}

type pagestore struct {
	rd *RedisDriver
}

func (p *pagestore) Create(ctx context.Context, pg *Page) error {
	conn, err := p.rd.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	payload, err := CacheSerialize(pg)
	if err != nil {
		return err
	}

	key := "pages/" + pg.ID

	_, err = conn.Do("SET", key, payload)
	if err != nil {
		return err
	}

	duration := pg.Expiry.Sub(pg.CreatedAt)
	_, err = conn.Do("EXPIRE", key, duration.Seconds())
	if err != nil {
		defer conn.Do("DELETE", key)
		return errors.Internal(err.Error())
	}

	return nil
}

func (p *pagestore) Read(ctx context.Context, pageID string) (*Page, error) {
	conn, err := p.rd.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	key := "pages/" + pageID

	payload, err := conn.Do("GET", key)
	if err != nil {
		return nil, err
	}

	pbytes, ok := payload.([]byte)
	if !ok {
		return nil, errors.New("page not found")
	}

	pg := new(Page)
	err = CacheDeserialize(pbytes, pg)
	if err != nil {
		return nil, err
	}

	return pg, nil
}

func newStore(cfg *Config) (*pagestore, error) {
	rd, err := newRedisDriver(cfg)
	if err != nil {
		return nil, err
	}

	return &pagestore{
		rd: rd,
	}, nil
}
