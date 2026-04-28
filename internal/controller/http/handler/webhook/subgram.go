package webhookhandler

import (
	"log/slog"

	"github.com/defany/goblin/errfmt"
	"github.com/defany/star_ad_proxy/internal/service/forwarder"
	"github.com/gofiber/fiber/v3"
)

type subgramRequest struct {
	Webhooks []forwarder.Webhook `json:"webhooks"`
}

func (h *Handler) subgram(c fiber.Ctx) error {
	var req subgramRequest
	if err := c.Bind().JSON(&req); err != nil {
		return errfmt.WithSource(c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad json"}))
	}

	if len(req.Webhooks) == 0 {
		return errfmt.WithSource(c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok", "forwarded": 0}))
	}

	results, err := h.fwd.Forward(c.Context(), req.Webhooks)
	if err != nil {
		h.log.WarnContext(c.Context(), "forward partial failure",
			slog.String("err", err.Error()),
			slog.Int("groups", len(results)),
		)
		return errfmt.WithSource(c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status":  "partial",
			"error":   err.Error(),
			"results": results,
		}))
	}

	return errfmt.WithSource(c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "ok",
		"forwarded": len(results),
	}))
}
