package di

import (
	"context"

	"github.com/defany/goblin/inject"
	"github.com/defany/star_ad_proxy/internal/config"
	"github.com/defany/star_ad_proxy/internal/service/forwarder"
	"github.com/gofiber/fiber/v3/client"
)

func (d *DI) FiberClient(ctx context.Context) *client.Client {
	return inject.Once(ctx, func(_ context.Context) *client.Client {
		c := client.New()
		c.SetTimeout(config.ForwarderTimeout())

		return c
	})
}

func (d *DI) Forwarder(ctx context.Context) *forwarder.Forwarder {
	return inject.Once(ctx, func(ctx context.Context) *forwarder.Forwarder {
		return forwarder.New(d.RouteService(ctx), d.FiberClient(ctx), d.Log(ctx))
	})
}
