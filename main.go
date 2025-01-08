package main

import (
	"time"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	opts := []fx.Option{}
	opts = append(opts, appFxOptions()...)
	opts = append(opts, startupFxOptions()...)
	fx.New(opts...).Run()
}

func appFxOptions() []fx.Option {
	opts := []fx.Option{
		fx.Provide(zap.NewExample), // provide logger
		fx.Provide(NewAPIv0),       // provide api handler
		fx.Provide(NewRouterv0),    // provide router
	}
	return opts
}

func startupFxOptions() []fx.Option {
	opts := []fx.Option{
		fx.StartTimeout(20 * time.Second),
		fx.StopTimeout(15 * time.Second),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{
				Logger: zap.New(log.Core()).
					With(
						zap.String("service", "weather-fx"),
						zap.String("version", "0.0.1"),
					),
			}
		}),
		fx.Invoke(NewHTTPServer), // start http server while starting the app
	}
	return opts
}

// server -> router -> handlers
