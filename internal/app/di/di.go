package di

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/defany/goblin/inject"
	"github.com/defany/goblin/slogx"
)

type DI struct{}

func New() *DI {
	return &DI{}
}

func (d *DI) Log(ctx context.Context) *slog.Logger {
	return inject.Once(ctx, func(_ context.Context) *slog.Logger {
		log := slog.New(slogx.Pretty(
			slogx.WithLevel(slog.LevelDebug),
			slogx.WithAddSource(true),
		))

		slog.SetDefault(log)

		return log
	})
}

func (d *DI) mustExit(err error) {
	panic(fmt.Errorf("cannot init dep: %w", err))
}
