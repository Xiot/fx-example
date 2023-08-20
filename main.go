package main

import (
	"context"
	"net"
	"net/http"
	"time"

	"usemotion.com/fx-example/handlers"
	"usemotion.com/fx-example/handlers/echo"
	"usemotion.com/fx-example/handlers/hello"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var startTime = time.Now()

func main() {
	fx.New(
		fx.Provide(
			NewHTTPServer,
			fx.Annotate(
				handlers.NewServeMux,
				fx.ParamTags(`group:"routes"`),
			),

			handlers.AsRoute(echo.NewEchoHandler),
			handlers.AsRoute(hello.NewHelloHandler),

			zap.NewExample,
		),
		fx.Invoke(func(*http.Server) {}),
		// fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
		// 	return &fxevent.ZapLogger{Logger: log}
		// }),
	).Run()
}

// NewHTTPServer builds an HTTP server that will begin serving requests
// when the Fx application starts.
func NewHTTPServer(lc fx.Lifecycle, mux *http.ServeMux, log *zap.Logger) *http.Server {
	srv := &http.Server{Addr: ":8080", Handler: mux}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}

			log.Info("Starting HTTP server", zap.String("addr", srv.Addr))
			go srv.Serve(ln)

			var afterStart = time.Now()
			log.Info("Started", zap.Duration("duration", afterStart.Sub(startTime)))

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}
