package route

import (
	"context"
	"errors"
	"time"

	"github.com/defany/goblin/errfmt"
	"github.com/defany/star_ad_proxy/internal/model"
	"github.com/defany/star_ad_proxy/internal/model/modelerr"
)

type Repo interface {
	Get(ctx context.Context, botID int64) (model.Route, error)
	Set(ctx context.Context, r model.Route) error
	Delete(ctx context.Context, botID int64) error
	List(ctx context.Context) ([]model.Route, error)
}

type Service struct {
	repo Repo
}

func New(repo Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(ctx context.Context, botID int64) (model.Route, error) {
	rt, err := s.repo.Get(ctx, botID)
	return rt, errfmt.WithSource(err)
}

func (s *Service) Upsert(ctx context.Context, in model.Route) (model.Route, error) {
	now := time.Now().UTC()

	existing, err := s.repo.Get(ctx, in.BotID)
	switch {
	case err == nil:
		in.CreatedAt = existing.CreatedAt
	case errors.Is(err, modelerr.ErrRouteNotFound):
		in.CreatedAt = now
	default:
		return model.Route{}, errfmt.WithSource(err)
	}

	in.UpdatedAt = now

	if err := s.repo.Set(ctx, in); err != nil {
		return model.Route{}, errfmt.WithSource(err)
	}

	return in, nil
}

func (s *Service) Delete(ctx context.Context, botID int64) error {
	return errfmt.WithSource(s.repo.Delete(ctx, botID))
}

func (s *Service) List(ctx context.Context) ([]model.Route, error) {
	out, err := s.repo.List(ctx)
	return out, errfmt.WithSource(err)
}
