package main

import (
	"coffee-spa/config"
	"coffee-spa/controller"
	"coffee-spa/db"
	"coffee-spa/router"
	"coffee-spa/usecase"
	"log"

	"github.com/labstack/echo/v4"
)

func main() {
	// 設定を読む
	c, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// DB接続を作る
	d, err := db.Open(c)
	if err != nil {
		log.Fatal(err)
	}

	//テーブル作成とindexの作成
	if err := db.Migrate(d); err != nil {
		log.Fatal(err)
	}

	// Echo本体
	e := echo.New()

	// Usecase
	uc := usecase.NewHealthUC(d.S)

	// Controller
	ctl := controller.NewHealthCtl(uc)

	// Router
	router.New(e, ctl)

	// 起動
	log.Fatal(e.Start(":" + c.Port))
}
