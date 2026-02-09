package account

import (
	"net/http"

	"paybridge-transaction-service/internal/infra/logger"
	"paybridge-transaction-service/internal/server/middleware"
	"paybridge-transaction-service/pkg/response"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Handler struct {
	service Service
	log     *logger.Logger
}

func NewHandler(svc Service, log *logger.Logger) *Handler {
	return &Handler{service: svc, log: log}
}

func (h *Handler) RegisterInternalRoutes(g *echo.Group) {
	g.GET("/account/:owner_id", h.GetAccount, middleware.ValidateInternalToken)
}

// GetAccount godoc
//
// @Summary     Get account by owner ID
// @Description Internal account lookup
// @Tags        Internal Account
// @Produce     json
// @Param       owner_id path string true "Owner UUID"
// @Security    BearerAuth
// @Router      /account/{owner_id} [get]
func (h *Handler) GetAccount(c echo.Context) error {
	ctx := c.Request().Context()

	ownerIDStr := c.Param("owner_id")
	ownerID, err := uuid.Parse(ownerIDStr)
	h.log.Info(ctx, "get account request received",
		zap.String("owner_id", ownerIDStr),
	)
	if err != nil {
		h.log.Warn(ctx, "invalid owner_id format",
			zap.String("owner_id", ownerIDStr),
		)
		return c.JSON(http.StatusBadRequest,
			response.Error("invalid owner_id", http.StatusBadRequest),
		)
	}

	resp, err := h.service.GetAccount(ctx, ownerID)

	if err != nil {
		return mapError(c, h.log, err)
	}

	return c.JSON(http.StatusOK,
		response.Success("Account is found",
			AccountResponse{
				resp.OwnerID,
				resp.Status,
			},
			http.StatusOK),
	)
}
