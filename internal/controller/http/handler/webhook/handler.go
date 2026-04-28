package webhookhandler

import (
	"context"
	"log/slog"

	"github.com/defany/star_ad_proxy/internal/service/forwarder"
	"github.com/gofiber/fiber/v3"
)

type Forwarder interface {
	Forward(ctx context.Context, batch []forwarder.Webhook) ([]forwarder.Result, error)
}

type Handler struct {
	fwd Forwarder
	log *slog.Logger
}

func New(fwd Forwarder, log *slog.Logger) *Handler {
	return &Handler{fwd: fwd, log: log}
}

func (h *Handler) Setup(r fiber.Router) {
	r.Post("/webhook/subgram", h.subgram)
}
