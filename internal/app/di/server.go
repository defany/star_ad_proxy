package di

import (
	"context"

	"github.com/defany/goblin/inject"
	"github.com/gofiber/fiber/v3"
)

func (d *DI) HttpServer(ctx context.Context) *fiber.App {
	return inject.Once(ctx, func(_ context.Context) *fiber.App {
		return fiber.New()
	})
}
