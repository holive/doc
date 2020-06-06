package main

import (
	"fmt"
	"os"

	"github.com/holive/doc/app/http"

	"github.com/holive/doc/app/doc"
	"github.com/pkg/errors"
)

func main() {
	app, err := doc.New()
	if err != nil {
		fmt.Println(errors.Wrap(err, "could not run doc").Error())
		os.Exit(1)
	}

	server, err := http.NewServer(&http.ServerConfig{
		Addr:              app.Cfg.HTTPServer.Addr,
		MaxHeaderBytes:    app.Cfg.HTTPServer.MaxHeaderBytes,
		IdleTimeout:       app.Cfg.HTTPServer.IdleTimeout,
		ReadHeaderTimeout: app.Cfg.HTTPServer.ReadHeaderTimeout,
		ReadTimeout:       app.Cfg.HTTPServer.ReadTimeout,
		WriteTimeout:      app.Cfg.HTTPServer.WriteTimeout,
		Router:            &http.RouterConfig{MiddlewareTimeout: app.Cfg.HTTPServer.Router.MiddlewareTimeout},
	}, app.Services)
	if err != nil {
		fmt.Println(errors.Wrap(err, "could not run Doc").Error())
		os.Exit(1)
	}

	if err := server.Start(); err != nil {
		fmt.Println(errors.Wrap(err, "could not run Doc").Error())
		os.Exit(1)
	}
}
