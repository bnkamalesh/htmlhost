package http

import (
	"bytes"
	"net/http"
	"strconv"
	"time"

	"github.com/bnkamalesh/errors"
	"github.com/bnkamalesh/htmlhost/internal/pages"
	"github.com/bnkamalesh/webgo/v4"
)

func (h *Handler) minifiedHTML(w http.ResponseWriter, payload []byte) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	minified, err := h.minifier.Bytes("text/html", payload)
	if err != nil {
		status, msg, _ := errors.HTTPStatusCodeMessage(err)
		w.Header().Set("Content-Length", strconv.Itoa(len(msg)))
		w.WriteHeader(status)
		w.Write([]byte(msg))
		return
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(minified)))
	w.Write(minified)
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	defer h.recoverer(w)

	buff := bytes.NewBuffer([]byte{})
	err := h.templates["home"].Execute(buff, PageResponse{})
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
