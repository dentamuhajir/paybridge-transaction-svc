package wallet

import (
	"log"
	"net/http"
	"paybridge-transaction-service/internal/logger"
	"paybridge-transaction-service/internal/server/middleware"
	"paybridge-transaction-service/pkg/response"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Handler struct {
	service Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{service: svc}
}

func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.POST("/wallet", h.Create, middleware.ValidateInternalToken)
}

// CreateWalletHandler godoc
// @Summary Create wallet
// @Description Create a new wallet for a user
// @Tags Wallet
// @Accept json
// @Produce json
// @Param request body CreateWalletReq true "Wallet creation payload"
// @Success 200 {object} response.SwaggerSuccessResponse
// @Failure 400 {object} response.SwaggerErrorResponse
// @Failure 500 {object} response.SwaggerErrorResponse
// @Security InternalTokenAuth
// @Router /wallet [post]
func (h *Handler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	traceID := c.Request().Header.Get("X-Trace-Id")
	log.Println("Trace ID from other service: ", traceID)
	var req CreateWalletReq

	if err := c.Bind(&req); err != nil {
		logger.Error(ctx, "invalid request body", err)
		return c.JSON(
			http.StatusBadRequest,
			response.Error("invalid request body", http.StatusBadRequest),
		)
	}

	resp, err := h.service.CreateWallet(c.Request().Context(), req)
	if err != nil {
		logger.Error(ctx, "failed to create wallet", err)
		return c.JSON(
			http.StatusInternalServerError,
			response.Error("failed to create wallet", http.StatusInternalServerError),
		)
	}

	logger.Info(ctx, "wallet created successfully",
		zap.String("user_id", resp.UserID),
		zap.String("wallet_id", resp.ID),
	)

	return c.JSON(http.StatusOK,
		response.Success("wallet created", resp, http.StatusOK),
	)
}
