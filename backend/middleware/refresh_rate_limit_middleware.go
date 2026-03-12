package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"

	"coffee-spa/repository"

	"github.com/labstack/echo/v4"
)

type RefreshLimiter interface {
	AllowRefresh(userID int64) (bool, int, error)
}

type refreshErrRes struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func RefreshRateLimit(
	rtRepo repository.RtRepository,
	rl RefreshLimiter,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("refresh_token")
			if err != nil || cookie == nil {
				return c.JSON(http.StatusUnauthorized, refreshErrRes{
					Error:   "unauthorized",
					Message: "refresh token is missing",
				})
			}

			raw := strings.TrimSpace(cookie.Value)
			if raw == "" {
				return c.JSON(http.StatusUnauthorized, refreshErrRes{
					Error:   "unauthorized",
					Message: "refresh token is missing",
				})
			}

			rt, err := rtRepo.GetByTokenHash(sha256Hex(raw))
			if err != nil {
				return c.JSON(http.StatusUnauthorized, refreshErrRes{
					Error:   "unauthorized",
					Message: "refresh token is invalid",
				})
			}

			ok, retryAfter, err := rl.AllowRefresh(rt.UserID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, refreshErrRes{
					Error:   "internal",
					Message: "internal server error",
				})
			}

			if !ok {
				c.Response().Header().Set("Retry-After", strconv.Itoa(retryAfter))
				return c.JSON(http.StatusTooManyRequests, refreshErrRes{
					Error:   "rate_limited",
					Message: "too many refresh requests",
				})
			}

			return next(c)
		}
	}
}

func sha256Hex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}
