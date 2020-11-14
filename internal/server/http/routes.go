package http

import (
	"net/http"

	"github.com/bnkamalesh/webgo/v4"
)

func routes(handlers *Handler) []*webgo.Route {
	allroutes := []*webgo.Route{
		&webgo.Route{
			Name:     "home",
			Pattern:  "/",
			Method:   http.MethodGet,
			Handlers: []http.HandlerFunc{handlers.Home},
		},
		&webgo.Route{
			Name:     "static",
			Pattern:  "/static/:path*",
			Method:   http.MethodGet,
			Handlers: []http.HandlerFunc{handlers.Static},
		},
		&webgo.Route{
			Name:     "submitpage",
			Pattern:  "/",
			Method:   http.MethodPost,
			Handlers: []http.HandlerFunc{handlers.CreatePage},
		},
		&webgo.Route{
			Name:     "renderpage",
			Pattern:  "/p/:pageID",
			Method:   http.MethodGet,
			Handlers: []http.HandlerFunc{handlers.ViewPage},
		},
	}
	return allroutes
}
