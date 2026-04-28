package app

import (
	"time"

	"github.com/defany/goblin/errfmt"
	"github.com/defany/goblin/lc"
	"github.com/defany/star_ad_proxy/internal/app/di"
)

type App struct {
	di *di.DI
}

func New() *App {
	return &App{di: di.New()}
}

func (a *App) Run() error {
	l := lc.New(lc.WithShutdownTimeout(time.Second * 30))

	ctx := l.Context()

	log := a.di.Log(ctx)

	l.Go(a.runHttp)

	log.Info("app started")

	return errfmt.WithSource(l.Run())
}
