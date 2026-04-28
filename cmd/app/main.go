package main

import (
	"github.com/defany/star_ad_proxy/internal/app"
	"github.com/defany/star_ad_proxy/internal/config"
)

func main() {
	config.MustLoad()

	a := app.New()

	if err := a.Run(); err != nil {
		panic(err)
	}
}
