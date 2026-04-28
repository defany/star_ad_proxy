package redisroute

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/defany/goblin/errfmt"
	"github.com/defany/star_ad_proxy/internal/model"
	"github.com/defany/star_ad_proxy/internal/model/modelerr"
	"github.com/redis/go-redis/v9"
)

const (
	keyPrefix = "proxy:routes"
	indexKey  = "proxy:routes:index"

	fieldTargetURL = "target_url"
	fieldAPIKey    = "api_key"
	fieldCreatedAt = "created_at"
	fieldUpdatedAt = "updated_at"
)

type Repo struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) *Repo {
	return &Repo{rdb: rdb}
}

func key(botID int64) string {
	return keyPrefix + ":" + strconv.FormatInt(botID, 10)
}

func (r *Repo) Get(ctx context.Context, botID int64) (model.Route, error) {
	m, err := r.rdb.HGetAll(ctx, key(botID)).Result()
	if err != nil {
		return model.Route{}, errfmt.WithSource(err)
	}
	if len(m) == 0 {
		return model.Route{}, modelerr.ErrRouteNotFound
	}

	createdAt, err := time.Parse(time.RFC3339Nano, m[fieldCreatedAt])
	if err != nil {
		return model.Route{}, errfmt.WithSource(err)
	}
	updatedAt, err := time.Parse(time.RFC3339Nano, m[fieldUpdatedAt])
	if err != nil {
		return model.Route{}, errfmt.WithSource(err)
	}

	return model.Route{
		BotID:     botID,
		TargetURL: m[fieldTargetURL],
		APIKey:    m[fieldAPIKey],
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func (r *Repo) Set(ctx context.Context, route model.Route) error {
	_, err := r.rdb.TxPipelined(ctx, func(p redis.Pipeliner) error {
		p.HSet(ctx, key(route.BotID), map[string]any{
			fieldTargetURL: route.TargetURL,
			fieldAPIKey:    route.APIKey,
			fieldCreatedAt: route.CreatedAt.UTC().Format(time.RFC3339Nano),
			fieldUpdatedAt: route.UpdatedAt.UTC().Format(time.RFC3339Nano),
		})
		p.SAdd(ctx, indexKey, route.BotID)
		return nil
	})
	return errfmt.WithSource(err)
}

func (r *Repo) Delete(ctx context.Context, botID int64) error {
	n, err := r.rdb.Del(ctx, key(botID)).Result()
	if err != nil {
		return errfmt.WithSource(err)
	}
	if n == 0 {
		return modelerr.ErrRouteNotFound
	}
	return errfmt.WithSource(r.rdb.SRem(ctx, indexKey, botID).Err())
}

func (r *Repo) List(ctx context.Context) ([]model.Route, error) {
	ids, err := r.rdb.SMembers(ctx, indexKey).Result()
	if err != nil {
		return nil, errfmt.WithSource(err)
	}

	out := make([]model.Route, 0, len(ids))
	for _, id := range ids {
		botID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nil, errfmt.WithSource(err)
		}

		rt, err := r.Get(ctx, botID)
		if err != nil {
			if errors.Is(err, modelerr.ErrRouteNotFound) {
				continue
			}
			return nil, err
		}

		out = append(out, rt)
	}

	return out, nil
}
