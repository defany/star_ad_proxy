package di

import (
	"context"

	"github.com/defany/goblin/errfmt"
	"github.com/defany/goblin/inject"
	"github.com/defany/goblin/lc"
	"github.com/defany/star_ad_proxy/internal/config"
	"github.com/redis/go-redis/v9"
)

func (d *DI) Redis(ctx context.Context) *redis.Client {
	return inject.Once(ctx, func(ctx context.Context) *redis.Client {
		client := redis.NewClient(&redis.Options{
			Addr:     config.DragonflyAddr(),
			Password: config.DragonflyPassword(),
			DB:       config.DragonflyDB(),
		})

		if err := client.Ping(ctx).Err(); err != nil {
			d.Log(ctx).Error("redis was not reached")
			d.mustExit(err)
		}

		d.Log(ctx).Info("connected to redis")

		lc.Defer(ctx, func(ctx context.Context) error {
			d.Log(ctx).Debug("closing redis")
			err := client.Close()
			d.Log(ctx).Debug("redis closed")
			return errfmt.WithSource(err)
		})

		return client
	})
}
