package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type HTTPServerParams struct {
	fx.In
	Lifecycle fx.Lifecycle

	Logger *zap.Logger
	Router http.Handler
}

// NewHTTPServer returns a new HTTP server.
func NewHTTPServer(p HTTPServerParams) *http.Server {
	server := &http.Server{
		Addr:              ":9090",
		Handler:           p.Router,
		ReadHeaderTimeout: 1 * time.Second,
	}

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			ln, err := net.Listen("tcp", server.Addr)
			if err != nil {
				return err
			}
			p.Logger.Info("Starting HTTP server at", zap.String("addr", server.Addr))
			go func() {
				if err := server.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
					p.Logger.Error("error at runtime in http server", zap.Error(err))
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			p.Logger.Info("Stopping HTTP server at", zap.String("addr", server.Addr))
			return server.Shutdown(ctx)
		},
	})

	return server
}
