package router

import (
	"coffee-spa/controller"

	"github.com/labstack/echo/v4"
)

func New(
	e *echo.Echo,
	healthCtl controller.HealthCtl,
	authCtl controller.AuthCtl,
	itemCtl controller.ItemCtl,
	srcCtl controller.SrcCtl,
) {
	// Health
	e.GET("/health", healthCtl.Get)
	pub := e.Group("")

	pub.POST("/auth/signup", authCtl.Signup)
	pub.POST("/auth/verify-email", authCtl.VerifyEmail)
	pub.POST("/auth/resend-verify", authCtl.ResendVerify)
	pub.POST("/auth/login", authCtl.Login)
	pub.POST("/auth/password/forgot", authCtl.ForgotPw)
	pub.POST("/auth/password/reset", authCtl.ResetPw)

	pub.GET("/items/top", itemCtl.Top)
	pub.GET("/items", itemCtl.List)
	pub.GET("/sources", srcCtl.List)

	csrf := e.Group("")
	csrf.POST("/auth/refresh", authCtl.Refresh)

	priv := e.Group("")
	priv.POST("/auth/logout", authCtl.Logout)
	priv.GET("/me", authCtl.Me)

	admin := e.Group("")
	admin.POST("/items", itemCtl.Create)
	admin.POST("/sources", srcCtl.Create)
}
