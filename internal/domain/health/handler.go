package health

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	Service Service
}

func NewHandler(Service Service) *Handler {
	return &Handler{Service: Service}
}

func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.GET("/health/db", h.CheckDB)
}

func (h *Handler) CheckDB(c echo.Context) error {
	if err := h.Service.CheckDB(c.Request().Context()); err != nil {
		return c.JSON(http.StatusServiceUnavailable, echo.Map{
			"status": "DOWN",
			"error":  err.Error(),
		})

	}
	return c.JSON(http.StatusOK, echo.Map{
		"status": "UP",
	})

}
