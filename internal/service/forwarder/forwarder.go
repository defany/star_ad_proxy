package forwarder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/defany/goblin/errfmt"
	"github.com/defany/star_ad_proxy/internal/model"
	"github.com/defany/star_ad_proxy/internal/model/modelerr"
	httpwrap "github.com/defany/star_ad_proxy/pkg/fiber_http"
	"github.com/gofiber/fiber/v3/client"
	"github.com/sourcegraph/conc/pool"
)

type Webhook struct {
	WebhookID     int64       `json:"webhook_id"`
	AdsID         json.Number `json:"ads_id"`
	Link          string      `json:"link"`
	UserID        int64       `json:"user_id"`
	BotID         int64       `json:"bot_id"`
	Status        string      `json:"status"`
	SubscribeDate string      `json:"subscribe_date,omitempty"`
}

type Result struct {
	BotID  int64 `json:"bot_id"`
	Status int   `json:"status,omitempty"`
	Count  int   `json:"count"`
}

type RouteLookup interface {
	Get(ctx context.Context, botID int64) (model.Route, error)
}

type Forwarder struct {
	routes RouteLookup
	client *client.Client
	log    *slog.Logger
}

func New(routes RouteLookup, c *client.Client, log *slog.Logger) *Forwarder {
	return &Forwarder{routes: routes, client: c, log: log}
}

type job struct {
	route model.Route
	items []Webhook
}

// Forward резолвит per-bot маршруты и параллельно форвардит каждую группу.
// Незарегистрированные bot_id логируются и пропускаются — их форвардить некуда,
// ретрай subgram'у не поможет.
func (f *Forwarder) Forward(ctx context.Context, batch []Webhook) ([]Result, error) {
	groups := make(map[int64][]Webhook, 4)
	for _, w := range batch {
		groups[w.BotID] = append(groups[w.BotID], w)
	}

	botIDs := make([]int64, 0, len(groups))
	for botID := range groups {
		botIDs = append(botIDs, botID)
	}
	f.log.InfoContext(ctx, "received subgram batch",
		slog.Int("events", len(batch)),
		slog.Int("groups", len(groups)),
		slog.Any("bot_ids", botIDs),
	)

	jobs := make([]job, 0, len(groups))
	for botID, items := range groups {
		rt, err := f.routes.Get(ctx, botID)
		if err != nil {
			if errors.Is(err, modelerr.ErrRouteNotFound) {
				f.log.WarnContext(ctx, "skipping events for unregistered bot",
					slog.Int64("bot_id", botID),
					slog.Int("dropped", len(items)),
				)
				continue
			}
			return nil, errfmt.WithSource(fmt.Errorf("bot %d: route lookup: %w", botID, err))
		}
		jobs = append(jobs, job{route: rt, items: items})
	}

	if len(jobs) == 0 {
		return nil, nil
	}

	p := pool.NewWithResults[Result]().
		WithContext(ctx).
		WithCollectErrored().
		WithMaxGoroutines(8)

	for _, j := range jobs {
		j := j
		p.Go(func(ctx context.Context) (Result, error) {
			return f.forwardOne(ctx, j.route, j.items)
		})
	}

	results, err := p.Wait()
	return results, errfmt.WithSource(err)
}

func (f *Forwarder) forwardOne(ctx context.Context, rt model.Route, items []Webhook) (Result, error) {
	res := Result{BotID: rt.BotID, Count: len(items)}

	f.log.DebugContext(ctx, "forwarding events",
		slog.Int64("bot_id", rt.BotID),
		slog.Int("count", len(items)),
		slog.String("target", rt.TargetURL),
	)

	resp, err := httpwrap.Post[json.RawMessage](f.client, rt.TargetURL, client.Config{
		Ctx:  ctx,
		Body: map[string]any{"webhooks": items},
		Header: map[string]string{
			"Content-Type": "application/json",
			"Api-Key":      rt.APIKey,
		},
	})
	if err != nil {
		return res, errfmt.WithSource(fmt.Errorf("bot %d: forward: %w", rt.BotID, err))
	}

	res.Status = resp.StatusCode()
	if res.Status < 200 || res.Status >= 300 {
		return res, errfmt.WithSource(fmt.Errorf("bot %d: non-2xx status %d", rt.BotID, res.Status))
	}

	f.log.InfoContext(ctx, "forwarded events",
		slog.Int64("bot_id", rt.BotID),
		slog.Int("count", len(items)),
		slog.Int("status", res.Status),
		slog.String("target", rt.TargetURL),
	)

	return res, nil
}
