package registerhandler

import (
	"net/url"

	"github.com/defany/goblin/errfmt"
	"github.com/defany/star_ad_proxy/internal/model"
	"github.com/gofiber/fiber/v3"
)

type upsertRequest struct {
	BotID     int64  `json:"bot_id"`
	APIKey    string `json:"api_key"`
	TargetURL string `json:"target_url"`
}

func (r *upsertRequest) validate() error {
	if r.BotID <= 0 {
		return errBadField("bot_id must be > 0")
	}
	if r.APIKey == "" {
		return errBadField("api_key is required")
	}
	if len(r.APIKey) > 256 {
		return errBadField("api_key too long (max 256)")
	}
	u, err := url.Parse(r.TargetURL)
	if err != nil {
		return errBadField("target_url is not a valid url")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return errBadField("target_url scheme must be http or https")
	}
	if u.Host == "" {
		return errBadField("target_url host is empty")
	}
	return nil
}

type validationError struct{ msg string }

func (e *validationError) Error() string { return e.msg }

func errBadField(msg string) error { return &validationError{msg: msg} }

func (h *Handler) upsert(c fiber.Ctx) error {
	var req upsertRequest
	if err := c.Bind().JSON(&req); err != nil {
		return errfmt.WithSource(c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad json"}))
	}

	if err := req.validate(); err != nil {
		return errfmt.WithSource(c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()}))
	}

	stored, err := h.routes.Upsert(c.Context(), model.Route{
		BotID:     req.BotID,
		TargetURL: req.TargetURL,
		APIKey:    req.APIKey,
	})
	if err != nil {
		return errfmt.WithSource(err)
	}

	return errfmt.WithSource(c.Status(fiber.StatusOK).JSON(toResponse(stored)))
}
