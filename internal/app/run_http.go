package app

import (
	"context"
	"log/slog"

	"github.com/defany/goblin/errfmt"
	"github.com/defany/goblin/lc"
	"github.com/defany/star_ad_proxy/internal/config"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
)

func (a *App) runHttp(ctx context.Context) error {
	server := a.di.HttpServer(ctx)
	server.Use(cors.New())

	router := server.Group("/")

	router.Get(healthcheck.LivenessEndpoint, healthcheck.New())
	router.Get(healthcheck.ReadinessEndpoint, healthcheck.New())
	router.Get(healthcheck.StartupEndpoint, healthcheck.New())

	a.di.HttpWebhookHandler(ctx).Setup(router)
	a.di.HttpRegisterHandler(ctx).Setup(router)

	lc.OnShutdown(ctx, func(ctx context.Context) error {
		a.di.Log(ctx).Debug("shutting down http server")
		err := server.ShutdownWithContext(ctx)
		a.di.Log(ctx).Debug("http server stopped")
		return errfmt.WithSource(err)
	})

	a.di.Log(ctx).Info("go http server!", slog.String("addr", config.HttpAddr()))

	if err := server.Listen(config.HttpAddr(), fiber.ListenConfig{
		DisableStartupMessage: true,
	}); err != nil {
		return errfmt.WithSource(err)
	}

	return nil
}
