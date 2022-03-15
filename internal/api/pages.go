package api

import (
	"context"
	"log"

	"github.com/bnkamalesh/errors"
	"github.com/bnkamalesh/htmlhost/internal/pages"
)

func (a *API) PageCreate(ctx context.Context, pg *pages.Page) (*pages.Page, error) {
	pg, err := a.pagesSvc.Create(ctx, pg)
	if err != nil {
		log.Println(errors.Stacktrace(err))
		return nil, err
	}
	return pg, nil
}

func (a *API) PageRead(ctx context.Context, pageID string) (*pages.Page, error) {
	pg, err := a.pagesSvc.Read(ctx, pageID)
	if err != nil {
		log.Println(errors.Stacktrace(err))
		return nil, err
	}
	return pg, nil
}

func (a *API) ActivePages(ctx context.Context) (int, error) {
	count, err := a.pagesSvc.Active(ctx)
	if err != nil {
		log.Println(errors.Stacktrace(err))
		return 0, err
	}
	return count, nil
}
