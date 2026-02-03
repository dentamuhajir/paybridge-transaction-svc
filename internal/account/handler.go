package account

import (
	"net/http"

	"paybridge-transaction-service/internal/server/middleware"
	"paybridge-transaction-service/pkg/response"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Handler struct {
	service Service
	log     *zap.Logger
}

func NewHandler(svc Service, log *zap.Logger) *Handler {
	return &Handler{service: svc, log: log}
}

func (h *Handler) RegisterInternalRoutes(g *echo.Group) {
	g.GET("/account/:owner_id", h.GetAccount, middleware.ValidateInternalToken)
}

func (h *Handler) GetAccount(c echo.Context) error {
	ctx := c.Request().Context()

	ownerIDStr := c.Param("owner_id")
	ownerID, err := uuid.Parse(ownerIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			response.Error("invalid owner_id", http.StatusBadRequest),
		)
	}

	resp, _ := h.service.GetAccount(ctx, ownerID)

	return c.JSON(http.StatusOK,
		response.Success("Account is found", resp, http.StatusOK),
	)
}
