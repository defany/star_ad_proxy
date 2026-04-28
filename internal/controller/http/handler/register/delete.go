package registerhandler

import (
	"errors"
	"strconv"

	"github.com/defany/goblin/errfmt"
	"github.com/defany/star_ad_proxy/internal/model/modelerr"
	"github.com/gofiber/fiber/v3"
)

func (h *Handler) delete(c fiber.Ctx) error {
	botID, err := strconv.ParseInt(c.Params("bot_id"), 10, 64)
	if err != nil {
		return errfmt.WithSource(c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bot_id must be int64"}))
	}

	if err := h.routes.Delete(c.Context(), botID); err != nil {
		if errors.Is(err, modelerr.ErrRouteNotFound) {
			return errfmt.WithSource(c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "route not found"}))
		}
		return errfmt.WithSource(err)
	}

	return errfmt.WithSource(c.SendStatus(fiber.StatusNoContent))
}
