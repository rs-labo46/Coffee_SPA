package main

import (
	"coffee-spa/config"
	"coffee-spa/controller"
	"coffee-spa/db"
	"coffee-spa/dbrepository"
	"coffee-spa/infra"
	"coffee-spa/policy"
	"coffee-spa/router"
	"coffee-spa/usecase"
	"coffee-spa/validator"
	"log"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func main() {
	//設定を読む
	c, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	//DB接続を作る
	d, err := db.Open(c)
	if err != nil {
		log.Fatal(err)
	}

	//テーブル作成とindexの作成
	if err := db.Migrate(d); err != nil {
		log.Fatal(err)
	}

	//Echo本体
	e := echo.New()

	//Redis接続
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	//Repository
	userRepo := dbrepository.NewUserRepository(d.G)
	evRepo := dbrepository.NewEvRepository(d.G)
	pwRepo := dbrepository.NewPwRepository(d.G)
	rtRepo := dbrepository.NewRtRepository(d.G)
	sourceRepo := dbrepository.NewSourceRepository(d.G)
	itemRepo := dbrepository.NewItemRepository(d.G)
	auditRepo := dbrepository.NewAuditRepository(d.G)

	//Policy
	pwPol := policy.NewPwPol()
	emailPol := policy.NewEmailPol()
	kindPol := policy.NewKindPol()
	urlPol := policy.NewURLPol()
	pagePol := policy.NewPagePol()

	//Validator
	authVal := validator.NewAuthValidator(emailPol, pwPol)
	itemVal := validator.NewItemValidator(kindPol, urlPol, pagePol)
	sourceVal := validator.NewSourceValidator(urlPol)

	//Infra concrete
	ph := infra.NewBcryptHasher()
	tk := infra.NewJWTMaker(c.JWTSecret)
	mail := infra.NewLogMailer(c.FEURL)
	rl := infra.NewRateLimiter(
		rdb,
		infra.Rule{
			Limit:  5,
			Window: 60 * time.Second,
		},
		infra.Rule{
			Limit:  10,
			Window: 60 * time.Second,
		},
		infra.Rule{
			Limit:  3,
			Window: 300 * time.Second,
		},
		infra.Rule{
			Limit:  3,
			Window: 300 * time.Second,
		},
	)

	//Usecase
	healthUC := usecase.NewHealthUC(d.S)
	authUC := usecase.NewAuthUC(
		userRepo,
		evRepo,
		pwRepo,
		rtRepo,
		auditRepo,
		authVal,
		ph,
		tk,
		mail,
		rl,
	)
	itemUC := usecase.NewItemUC(
		itemRepo,
		sourceRepo,
		auditRepo,
		itemVal,
	)
	sourceUC := usecase.NewSourceUC(
		sourceRepo,
		auditRepo,
		sourceVal,
	)

	//Controller
	healthCtl := controller.NewHealthCtl(healthUC)

	//controller / router が未実装の間だけ未使用回避
	_ = authUC
	_ = itemUC
	_ = sourceUC

	//Router
	router.New(e, healthCtl)

	//起動
	log.Fatal(e.Start(":" + c.Port))
}
