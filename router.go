package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/fx"
)

type APIHandler interface {
	Router() http.Handler
	Path() string
}

type RouterParams struct {
	fx.In

	V0Handlers []APIHandler `group:"v0apis"`
}

func NewRouterv0(p RouterParams) (http.Handler, error) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.CleanPath)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.RealIP)

	r.Mount("/debug", middleware.Profiler())

	for _, h := range p.V0Handlers {
		r.Mount(h.Path(), h.Router())
	}

	return r, nil
}
