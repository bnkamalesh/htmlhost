package pages

import (
	"context"
	"time"

	"github.com/bnkamalesh/errors"
	"github.com/gomodule/redigo/redis"
)

const (
	pagesKeyPrefix = "pages/"
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
	Active(context.Context) (int, error)
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

	key := pagesKeyPrefix + pg.ID

	_, err = conn.Do("SET", key, payload)
	if err != nil {
		return errors.Wrap(err, "failed storing page data")
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

	key := pagesKeyPrefix + pageID

	payload, err := conn.Do("GET", key)
	if err != nil {
		return nil, errors.Wrap(err, "failed reading page data")
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

func (p *pagestore) Active(ctx context.Context) (int, error) {
	conn, err := p.rd.Conn(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	reply, err := conn.Do("KEYS", pagesKeyPrefix+"*")
	if err != nil {
		return 0, errors.Wrap(err, "failed getting keys")
	}

	bslices, err := redis.ByteSlices(reply, err)
	if err != nil {
		return 0, errors.Wrap(err, "failed getting bytes")
	}
	return len(bslices), nil
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
