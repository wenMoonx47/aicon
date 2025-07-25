package system

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type SystemHandler struct{}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

// Health はヘルスチェックエンドポイントです
func (h *SystemHandler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "healthy",
		"message": "Application is running successfully",
	})
} 