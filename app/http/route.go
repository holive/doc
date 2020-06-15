package http

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/holive/doc/app/docApi"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/holive/doc/app/http/handler"
)

type RouterConfig struct {
	MiddlewareTimeout time.Duration
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

var gzPool = sync.Pool{
	New: func() interface{} {
		w := gzip.NewWriter(ioutil.Discard)
		return w
	},
}

func NewRouter(cfg *RouterConfig, handler *handler.Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(cfg.MiddlewareTimeout))

	r.Get("/health", handler.Health)

	r.Route("/", func(r chi.Router) {
		r.Get("/", handler.GetAllDocs)
		r.Get("/{squad}", handler.ListBySquad)
		r.Get("/search/{projeto}", handler.SearchByProject)
		r.Get("/{squad}/{projeto}/{versao}", handler.GetDoc)
		r.Post("/{squad}/{projeto}/{versao}", handler.CreateDoc)
		r.Delete("/{squad}/{projeto}/{versao}", handler.DeleteDoc)
	})

	r.Route("/squad", func(r chi.Router) {
		r.Post("/", handler.CreateSquad)
	})

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, docApi.FilesFolder))
	fileServer(r, "/files", filesDir)

	return r
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("fileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())

		if !validFilePath(rctx) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		handler := http.StripPrefix(pathPrefix, http.FileServer(root))

		w.Header().Set("Cache-Control", "max-age=2592000")
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			handler.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")

		gz := gzPool.Get().(*gzip.Writer)
		defer gzPool.Put(gz)

		gz.Reset(w)
		defer gz.Close()

		handler.ServeHTTP(&gzipResponseWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func validFilePath(rctx *chi.Context) bool {
	for _, el := range docApi.FileTypes {
		if strings.HasSuffix(rctx.URLParams.Values[len(rctx.URLParams.Values)-1], el) {
			return true
		}
	}

	return false
}

func (w *gzipResponseWriter) WriteHeader(status int) {
	w.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(status)
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
