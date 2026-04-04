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
	counter int32
}

func NewHandler(svc Service, log *logger.Logger) *Handler {
	return &Handler{service: svc, log: log}
}

func (h *Handler) RegisterInternalRoutes(g *echo.Group) {
	g.GET("/account/:owner_id/balance", h.GetAccountBalance, middleware.ValidateInternalToken)
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

// GetAccountBalance godoc
//
// @Summary     Get account balance by owner ID
// @Description Retrieve the current balance for an account by owner ID
// @Tags        Internal Account
// @Produce     json
// @Param       owner_id path string true "Owner UUID"
// @Success     200 {object} map[string]interface{}
// @Failure     400 {object} map[string]interface{} "Invalid owner_id format"
// @Failure     404 {object} map[string]interface{} "Account not found"
// @Failure     500 {object} map[string]interface{} "Internal server error"
// @Security    BearerAuth
// @Router      /account/{owner_id}/balance [get]
func (h *Handler) GetAccountBalance(c echo.Context) error {

	ctx := c.Request().Context()
	ownerIDStr := c.Param("owner_id")
	ownerID, err := uuid.Parse(ownerIDStr)
	h.log.Info(ctx, "get account balance request received",
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

	balance, err := h.service.GetAccountBalance(ctx, ownerID)

	if err != nil {
		return mapError(c, h.log, err)
	}

	return c.JSON(http.StatusOK,
		response.Success("Account balance is found",
			AccountBalanceResponse{
				Balance: balance,
			},
			http.StatusOK),
	)

}
