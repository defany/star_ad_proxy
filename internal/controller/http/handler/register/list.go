package registerhandler

import (
	"github.com/defany/goblin/errfmt"
	"github.com/gofiber/fiber/v3"
)

func (h *Handler) list(c fiber.Ctx) error {
	routes, err := h.routes.List(c.Context())
	if err != nil {
		return errfmt.WithSource(err)
	}

	out := make([]routeResponse, 0, len(routes))
	for _, r := range routes {
		out = append(out, toResponse(r))
	}

	return errfmt.WithSource(c.Status(fiber.StatusOK).JSON(fiber.Map{"routes": out}))
}
