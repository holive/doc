package http

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/holive/doc/app/http/handler"
)

type RouterConfig struct {
	MiddlewareTimeout time.Duration
}

func NewRouter(cfg *RouterConfig, handler *handler.Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(cfg.MiddlewareTimeout))

	r.Get("/health", handler.Health)

	r.Route("/doc", func(r chi.Router) {
		r.Post("/{squad}/{projeto}/{versao}", handler.CreateDocApi)
		r.Get("/", handler.GetAllDocs)
		r.Get("/{squad}/{projeto}/{versao}", handler.GetDocApi)
		r.Delete("/{squad}/{projeto}/{versao}", handler.DeleteDocApi)
	})

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "documentos"))
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
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
