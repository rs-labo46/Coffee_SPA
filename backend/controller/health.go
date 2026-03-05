package controller

import (
	"coffee-spa/usecase"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type HealthCtl struct {
	uc usecase.HealthUC // Usecase
}

func NewHealthCtl(uc usecase.HealthUC) HealthCtl {
	return HealthCtl{uc: uc} // DI
}

func (h HealthCtl) Get(c echo.Context) error {
	ctx := c.Request().Context()
	_ = ctx
	t := time.Now() // 現在時刻
	_ = t
	err := h.uc.Check(c.Request().Context()) // DBが動いているか
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"status": "ng",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}
