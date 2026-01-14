package loan

import (
	"net/http"
	"paybridge-transaction-service/internal/loan/dto"
	"paybridge-transaction-service/internal/server/middleware"
	"paybridge-transaction-service/pkg/response"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

type Handler struct {
	service Service
	log     *zap.Logger
}

func NewHandler(svc Service, log *zap.Logger) *Handler {
	return &Handler{service: svc, log: log}
}

func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.POST("/loan-application", h.Create, middleware.ValidateInternalToken)
}

func (h *Handler) Create(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.LoanAppCreateReq

	if err := c.Bind(&req); err != nil {
		log.Error(ctx, "invalid request body", err)
		return c.JSON(
			http.StatusBadRequest,
			response.Error("invalid request body", http.StatusBadRequest),
		)
	}

	resp, err := h.service.Create(c.Request().Context(), req)
	if err != nil {
		log.Error(ctx, "failed to create loan application", err)
		return c.JSON(
			http.StatusInternalServerError,
			response.Error("failed to create wallet", http.StatusInternalServerError),
		)
	}

	log.Info(ctx, "loan application created successfully",
		zap.String("ID", resp.ID),
		zap.String("Status", resp.Status),
	)

	return c.JSON(http.StatusOK,
		response.Success("loan application created", resp, http.StatusOK),
	)
}
