package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/bnkamalesh/errors"
	"github.com/bnkamalesh/webgo/v5"
	"github.com/h2non/filetype"
	svg "github.com/h2non/go-is-svg"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
)

func (h *Handler) readStaticFile(path string, fileDir string, w http.ResponseWriter, r *http.Request) {
	defer h.recoverer(w)

	expiry := time.Now().Add(time.Hour * 2)
	etag := fmt.Sprintf("%s-%s", path, startedAt.String())
	if !cacheHeaders(w, r, etag, startedAtHTTPFormatted, &expiry) {
		return
	}

	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", fileDir, path))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	}

	ftype := ""
	if webgo.Context(r).Route.Name == "site-manifest" {
		ftype = "application/json"
	} else {
		ftype, err = detectFileType(path, data)
		if err != nil {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte("not supported"))
			return
		}
	}

	w.Header().Set("Content-Type", ftype)
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Write(data)
}

func (h *Handler) Static(w http.ResponseWriter, r *http.Request) {
	defer h.recoverer(w)

	wctx := webgo.Context(r)
	path := wctx.Params()["path"]
	h.readStaticFile(path, "./internal/server/http/web/static", w, r)
}

func (h *Handler) MetaStatic(w http.ResponseWriter, r *http.Request) {
	defer h.recoverer(w)

	path := r.RequestURI[1:]
	h.readStaticFile(path, "./internal/server/http/web/static/meta", w, r)
}

func nativeFtypecheck(fname string, content []byte) (string, error) {
	ftype := http.DetectContentType(content)
	switch ftype {
	case "":
		{
			return "", errors.Validation("unknown file type")
		}
	case "text/plain; charset=utf-8":
		{
			fnamelen := len(fname)
			if fnamelen > 3 && fname[fnamelen-3:] == "css" {
				return "text/css", nil
			}
			if fnamelen > 2 && fname[fnamelen-2:] == "js" {
				return "text/javascript; charset=UTF-8", nil
			}
		}
	default:
		return ftype, nil
	}

	return ftype, nil
}

func detectFileType(fname string, content []byte) (string, error) {
	if svg.Is(content) {
		return "image/svg+xml", nil
	}

	ftype, err := filetype.Match(content)
	if ftype == filetype.Unknown || err != nil {
		return nativeFtypecheck(fname, content)
	}

	return fmt.Sprintf("%s/%s", ftype.MIME.Type, ftype.MIME.Subtype), nil
}

func newMinifier() *minify.M {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	return m
}
