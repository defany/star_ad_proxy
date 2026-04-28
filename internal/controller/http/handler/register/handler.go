package registerhandler

import (
	"context"
	"crypto/subtle"
	"log/slog"
	"strings"

	"github.com/defany/goblin/errfmt"
	"github.com/defany/star_ad_proxy/internal/model"
	"github.com/gofiber/fiber/v3"
)

type RouteService interface {
	Get(ctx context.Context, botID int64) (model.Route, error)
	Upsert(ctx context.Context, in model.Route) (model.Route, error)
	Delete(ctx context.Context, botID int64) error
	List(ctx context.Context) ([]model.Route, error)
}

type Handler struct {
	routes     RouteService
	adminToken string
	log        *slog.Logger
}

func New(routes RouteService, adminToken string, log *slog.Logger) *Handler {
	return &Handler{routes: routes, adminToken: adminToken, log: log}
}

func (h *Handler) Setup(r fiber.Router) {
	g := r.Group("/register", h.authMiddleware)
	g.Post("", h.upsert)
	g.Get("", h.list)
	g.Delete("/:bot_id", h.delete)
}

func (h *Handler) authMiddleware(c fiber.Ctx) error {
	header := c.Get("Authorization")
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return errfmt.WithSource(c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"}))
	}
	got := strings.TrimPrefix(header, prefix)
	if subtle.ConstantTimeCompare([]byte(got), []byte(h.adminToken)) != 1 {
		return errfmt.WithSource(c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"}))
	}
	return errfmt.WithSource(c.Next())
}

type routeResponse struct {
	BotID     int64  `json:"bot_id"`
	TargetURL string `json:"target_url"`
	APIKey    string `json:"api_key"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func toResponse(r model.Route) routeResponse {
	return routeResponse{
		BotID:     r.BotID,
		TargetURL: r.TargetURL,
		APIKey:    maskKey(r.APIKey),
		CreatedAt: r.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: r.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func maskKey(k string) string {
	if len(k) <= 8 {
		return strings.Repeat("*", len(k))
	}
	return k[:4] + strings.Repeat("*", len(k)-8) + k[len(k)-4:]
}
