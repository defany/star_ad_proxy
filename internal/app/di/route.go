package di

import (
	"context"

	"github.com/defany/goblin/inject"
	redisroute "github.com/defany/star_ad_proxy/internal/repo/route/redis"
	"github.com/defany/star_ad_proxy/internal/service/route"
)

func (d *DI) RouteRepo(ctx context.Context) *redisroute.Repo {
	return inject.Once(ctx, func(ctx context.Context) *redisroute.Repo {
		return redisroute.New(d.Redis(ctx))
	})
}

func (d *DI) RouteService(ctx context.Context) *route.Service {
	return inject.Once(ctx, func(ctx context.Context) *route.Service {
		return route.New(d.RouteRepo(ctx))
	})
}
