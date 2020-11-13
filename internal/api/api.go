package api

import "github.com/bnkamalesh/htmlhost/internal/pages"

type API struct {
	pagesSvc *pages.Pages
}

func New(pgs *pages.Pages) *API {
	return &API{
		pagesSvc: pgs,
	}
}
