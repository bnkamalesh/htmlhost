package http

import (
	"net/http"

	"github.com/bnkamalesh/webgo/v6"
)

func routesMetaStatic(handlers *Handler) []*webgo.Route {
	return []*webgo.Route{
		{
			Name:     "favicon.ico",
			Pattern:  "/favicon.ico",
			Method:   http.MethodGet,
			Handlers: []http.HandlerFunc{handlers.MetaStatic},
		},
		{
			Name:     "favicon32",
			Pattern:  "/favicon-32x32.png",
			Method:   http.MethodGet,
			Handlers: []http.HandlerFunc{handlers.MetaStatic},
		},
		{
			Name:     "favicon16",
			Pattern:  "/favicon-16x16.png",
			Method:   http.MethodGet,
			Handlers: []http.HandlerFunc{handlers.MetaStatic},
		},
		{
			Name:     "apple-touch-icon",
			Pattern:  "/apple-touch-icon.png",
			Method:   http.MethodGet,
			Handlers: []http.HandlerFunc{handlers.MetaStatic},
		},
		{
			Name:     "mstile-150x150",
			Pattern:  "/mstile-150x150.png",
			Method:   http.MethodGet,
			Handlers: []http.HandlerFunc{handlers.MetaStatic},
		},
		{
			Name:     "android-chrome-192x192",
			Pattern:  "/android-chrome-192x192.png",
			Method:   http.MethodGet,
			Handlers: []http.HandlerFunc{handlers.MetaStatic},
		},
		{
			Name:     "android-chrome-256x256",
			Pattern:  "/android-chrome-256x256.png",
			Method:   http.MethodGet,
			Handlers: []http.HandlerFunc{handlers.MetaStatic},
		},
		{
			Name:     "safari-pinned",
			Pattern:  "/safari-pinned-tab.svg",
			Method:   http.MethodGet,
			Handlers: []http.HandlerFunc{handlers.MetaStatic},
		},
		{
			Name:     "site-manifest",
			Pattern:  "/site.webmanifest",
			Method:   http.MethodGet,
			Handlers: []http.HandlerFunc{handlers.MetaStatic},
		},
	}
}

func routes(handlers *Handler) []*webgo.Route {
	allroutes := routesMetaStatic(handlers)
	allroutes = append(
		allroutes,
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
	)
	return allroutes
}
