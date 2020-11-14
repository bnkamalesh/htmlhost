package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"text/template"
	"time"

	"github.com/bnkamalesh/errors"
	"github.com/bnkamalesh/htmlhost/internal/api"
	"github.com/bnkamalesh/htmlhost/internal/pages"
	"github.com/bnkamalesh/webgo/v4"
)

var (
	startedAt              = time.Now()
	startedAtHTTPFormatted = startedAt.Format(http.TimeFormat)
)

// Handler has all the dependencies initialized and stored, and made available to all the handler
// methods
type Handler struct {
	templates        map[string]*template.Template
	generatedBaseURL string
	api              *api.API
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

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	defer h.recoverer(w)

	buff := bytes.NewBuffer([]byte{})
	err := h.templates["home"].Execute(buff, PageResponse{})
	if err != nil {
		return
	}
	bodybytes := buff.Bytes()

	etag := "home" + startedAtHTTPFormatted + strconv.Itoa(len(bodybytes))
	if r.Header.Get("If-None-Match") == etag ||
		r.Header.Get("If-Modified-Since") == startedAtHTTPFormatted {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.Header().Set("Date", startedAtHTTPFormatted)
	w.Header().Set("Last-Modified", startedAtHTTPFormatted)
	w.Header().Add("Content-Length", strconv.Itoa(len(bodybytes)))
	w.Header().Add("ETag", etag)
	w.Write(bodybytes)
}

func (h *Handler) Static(w http.ResponseWriter, r *http.Request) {
	defer h.recoverer(w)

	wctx := webgo.Context(r)
	path := wctx.Params()["path"]

	etag := fmt.Sprintf("%s-%s", path, startedAt.String())
	if r.Header.Get("If-None-Match") == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	dat, err := ioutil.ReadFile(fmt.Sprintf("./internal/server/http/web/static/%s", path))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	}

	kind, err := detectFileType(dat)
	if err != nil {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte("not supported"))
		return
	}

	w.Header().Add("Content-Type", kind)
	w.Header().Set("Date", startedAtHTTPFormatted)
	w.Header().Set("Last-Modified", startedAtHTTPFormatted)
	w.Header().Add("Content-Length", strconv.Itoa(len(dat)))
	w.Header().Add("ETag", etag)
	w.Write(dat)
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
		pr.Link = pg.URL(h.generatedBaseURL)
		pr.Content = pg.Content
	}

	buff := bytes.NewBuffer([]byte{})
	err = h.templates["home"].Execute(buff, pr)
	if err != nil {
		status, msg, _ := errors.HTTPStatusCodeMessage(err)
		w.WriteHeader(status)
		w.Write([]byte(msg))
		return
	}

	nowFormatted := time.Now().Format(http.TimeFormat)
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.Header().Set("Date", nowFormatted)
	w.Header().Set("Last-Modified", nowFormatted)

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
			ID:        pageID,
			Content:   msg,
			CreatedAt: startedAt,
		}
		w.WriteHeader(status)
	}

	createdAtFormatted := pg.CreatedAt.Format(http.TimeFormat)

	if err == nil {
		w.Header().Add(
			"Cache-Control",
			fmt.Sprintf(
				"public,max-age=%v,must-revalidate",
				time.Until(pg.Expiry).Seconds(),
			),
		)
		w.Header().Set("Expires", pg.Expiry.Format(http.TimeFormat))
		w.Header().Add("ETag", pg.ID)

		if r.Header.Get("If-None-Match") == pg.ID ||
			r.Header.Get("If-Modified-Since") == createdAtFormatted {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.Header().Set("Date", createdAtFormatted)
	w.Header().Set("Last-Modified", createdAtFormatted)
	w.Write([]byte(pg.Content))
}

func newHandler(a *api.API, baseURL string) (*Handler, error) {
	tpl, err := template.ParseFiles(
		"./internal/server/http/web/templates/home.html",
	)
	if err != nil {
		return nil, err
	}

	return &Handler{
		api:              a,
		generatedBaseURL: baseURL,
		templates: map[string]*template.Template{
			"home": tpl,
		},
	}, nil
}
