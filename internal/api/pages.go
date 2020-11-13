package api

import (
	"context"
	"log"

	"github.com/bnkamalesh/htmlhost/internal/pages"
)

func (a *API) PageCreate(ctx context.Context, pg *pages.Page) (*pages.Page, error) {
	pg, err := a.pagesSvc.Create(ctx, pg)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return pg, nil
}

func (a *API) PageRead(ctx context.Context, pageID string) (*pages.Page, error) {
	pg, err := a.pagesSvc.Read(ctx, pageID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return pg, nil
}
