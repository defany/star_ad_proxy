package di

import (
	"context"

	"github.com/defany/goblin/inject"
	"github.com/defany/star_ad_proxy/internal/config"
	registerhandler "github.com/defany/star_ad_proxy/internal/controller/http/handler/register"
	webhookhandler "github.com/defany/star_ad_proxy/internal/controller/http/handler/webhook"
)

func (d *DI) HttpWebhookHandler(ctx context.Context) *webhookhandler.Handler {
	return inject.Once(ctx, func(ctx context.Context) *webhookhandler.Handler {
		return webhookhandler.New(d.Forwarder(ctx), d.Log(ctx))
	})
}

func (d *DI) HttpRegisterHandler(ctx context.Context) *registerhandler.Handler {
	return inject.Once(ctx, func(ctx context.Context) *registerhandler.Handler {
		return registerhandler.New(d.RouteService(ctx), config.AdminToken(), d.Log(ctx))
	})
}
