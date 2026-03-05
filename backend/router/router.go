package router

import (
	"coffee-spa/controller"

	"github.com/labstack/echo/v4"
)

func New(e *echo.Echo, h controller.HealthCtl) {
	e.GET("/health", h.Get) // ヘルスチェック
}
