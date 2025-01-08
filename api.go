package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

type APIv0 struct {
	version string
}

// APIResult struct for the output of NewAPI.
type APIv0Result struct {
	fx.Out

	API APIHandler `group:"v0apis"`
}

// NewAPIv0 returns a new API handler.
func NewAPIv0() APIv0Result {
	return APIv0Result{
		API: &APIv0{
			version: "v0",
		},
	}
}

// Router returns a new router.
func (a *APIv0) Router() http.Handler {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello World!"))
		})

		r.Post("/echo", func(w http.ResponseWriter, r *http.Request) {
			if _, err := io.Copy(w, r.Body); err != nil {
				fmt.Fprintln(os.Stderr, "Failed to handle request:", err)
			}
		})
	})

	return r
}

// Path returns the path for the API.
func (a *APIv0) Path() string {
	return "/"
}
