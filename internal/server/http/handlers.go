package http

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"text/template"
	"time"

	"github.com/bnkamalesh/errors"
	"github.com/bnkamalesh/htmlhost/internal/api"
	"github.com/bnkamalesh/htmlhost/internal/pages"
	"github.com/bnkamalesh/webgo/v4"
)

// Handler has all the dependencies initialized and stored, and made available to all the handler
// methods
type Handler struct {
	templates map[string]*template.Template
	api       *api.API
}

type PageResponse struct {
	pages.Page
	Link    string
	Message string
}

func (handler *Handler) recoverer(w http.ResponseWriter) {
	rec := recover()
	if rec == nil {
		return
	}

	log.Println(rec, string(debug.Stack()))
	webgo.R500(w, "sorry, unknown error occurred")
}

func httpRespondError(w http.ResponseWriter, err error) {
	status, msg, _ := errors.HTTPStatusCodeMessage(err)
	webgo.SendError(w, msg, status)
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	defer h.recoverer(w)

	buff := bytes.NewBuffer([]byte{})
	err := h.templates["home"].Execute(buff, PageResponse{})
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.Write(buff.Bytes())
}

func (h *Handler) CreatePage(w http.ResponseWriter, r *http.Request) {
	defer h.recoverer(w)

	pg := new(pages.Page)
	pg.Content = r.PostFormValue("body")
	pg, err := h.api.PageCreate(r.Context(), pg)

	pr := PageResponse{}
	if err != nil {
		status, msg, _ := errors.HTTPStatusCodeMessage(err)
		w.WriteHeader(status)
		pr.Message = msg
	} else {
		pr.Link = "https://htmlhost.live/p/" + pg.ID
		pr.Content = pg.Content
	}

	buff := bytes.NewBuffer([]byte{})
	err = h.templates["home"].Execute(buff, pr)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.Write(buff.Bytes())
}

func (h *Handler) ViewPage(w http.ResponseWriter, r *http.Request) {
	defer h.recoverer(w)

	wctx := webgo.Context(r)
	pageID := wctx.Params()["pageID"]
	pg, err := h.api.PageRead(r.Context(), pageID)
	if err != nil {
		status, msg, _ := errors.HTTPStatusCodeMessage(err)
		pg = &pages.Page{
			ID:      pageID,
			Content: msg,
		}
		w.WriteHeader(status)
	} else {
		w.Header().Add("Cache-Control", fmt.Sprintf("public,max-age=%v,must-revalidate", pg.Expiry.Sub(time.Now()).Seconds()))
		w.Header().Add("ETag", pg.ID)
		if r.Header.Get("If-None-Match") == pg.ID {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.Write([]byte(pg.Content))
}

func newHandler(a *api.API) (*Handler, error) {
	tpl, err := template.ParseFiles(
		"./internal/server/http/web/templates/home.html",
	)
	if err != nil {
		return nil, err
	}

	return &Handler{
		api: a,
		templates: map[string]*template.Template{
			"home": tpl,
		},
	}, nil
}
