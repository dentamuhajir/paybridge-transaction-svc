package account

import (
	"net/http"
	"sync/atomic"

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
	counter int32
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

	count := atomic.AddInt32(&h.counter, 1)

	h.log.Info(ctx, "retry test attempt",
		zap.Int32("attempt", count),
	)

	// First 2 attempts -> 500
	if count <= 2 {
		return c.JSON(
			http.StatusInternalServerError,
			response.Error("temporary failure", http.StatusInternalServerError),
		)
	}

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
