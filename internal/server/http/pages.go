package http

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bnkamalesh/errors"
	"github.com/bnkamalesh/htmlhost/internal/pages"
	"github.com/bnkamalesh/webgo/v4"
)

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	defer h.recoverer(w)

	buff := bytes.NewBuffer([]byte{})
	err := h.templates["home"].Execute(buff, newPageResponse("", r.Host))
	if err != nil {
		return
	}
	bodybytes := buff.Bytes()
	expiry := time.Now().Add(time.Hour * 2)
	etag := "home" + startedAtHTTPFormatted + strconv.Itoa(len(bodybytes))
	if !cacheHeaders(w, r, etag, startedAtHTTPFormatted, &expiry) {
		return
	}

	h.minifiedHTML(w, bodybytes)
}

func (h *Handler) CreatePage(w http.ResponseWriter, r *http.Request) {
	defer h.recoverer(w)

	pg := new(pages.Page)
	pg.Content = r.PostFormValue("body")
	pg, err := h.api.PageCreate(r.Context(), pg)

	pr := newPageResponse("", r.Host)
	if err != nil {
		status, msg, _ := errors.HTTPStatusCodeMessage(err)
		w.WriteHeader(status)
		pr.Message = msg
	} else {
		base := r.Header.Get("Origin")
		if base == "" {
			base = r.Host
		}
		pr.Link = pg.URL(base)
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

	h.minifiedHTML(w, buff.Bytes())
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
			Expiry:    time.Now().Add(time.Hour * 1),
		}
		w.WriteHeader(status)
	}

	createdAtFormatted := pg.CreatedAt.Format(http.TimeFormat)
	if !cacheHeaders(w, r, pg.ID, createdAtFormatted, &pg.Expiry) {
		return
	}

	h.minifiedHTML(w, []byte(pg.Content))
}

type PageResponse struct {
	pages.Page
	Title   string
	BaseURL string
	Link    string
	Message string
}

func pageTitle(baseURL string) string {
	parts := strings.Split(baseURL, ".")
	if len(parts) == 1 {
		return strings.ReplaceAll(parts[0], "html", "HTML")
	}
	return strings.ReplaceAll(strings.Join(parts[:len(parts)-1], "."), "html", "HTML")
}

func newPageResponse(scheme string, host string) *PageResponse {
	baseURL := host
	if len(baseURL) > 5 {
		if baseURL[:4] != "http" {
			if scheme != "" {
				baseURL = scheme + "://" + baseURL
			} else {
				baseURL = "https://" + baseURL
			}
		}
	}
	return &PageResponse{
		Title:   pageTitle(host),
		BaseURL: baseURL,
	}
}
