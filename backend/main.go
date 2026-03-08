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
	c, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	d, err := db.Open(c)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Migrate(d); err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	userRepo := dbrepository.NewUserRepository(d.G)
	evRepo := dbrepository.NewEvRepository(d.G)
	pwRepo := dbrepository.NewPwRepository(d.G)
	rtRepo := dbrepository.NewRtRepository(d.G)
	sourceRepo := dbrepository.NewSourceRepository(d.G)
	itemRepo := dbrepository.NewItemRepository(d.G)
	auditRepo := dbrepository.NewAuditRepository(d.G)

	pwPol := policy.NewPwPol()
	emailPol := policy.NewEmailPol()
	kindPol := policy.NewKindPol()
	urlPol := policy.NewURLPol()
	pagePol := policy.NewPagePol()

	authVal := validator.NewAuthValidator(emailPol, pwPol)
	itemVal := validator.NewItemValidator(kindPol, urlPol, pagePol)
	sourceVal := validator.NewSourceValidator(urlPol)

	ph := infra.NewBcryptHasher()
	tk := infra.NewJWTMaker(c.JWTSecret)
	mail := infra.NewLogMailer(c.FEURL)
	rl := infra.NewRateLimiter(
		rdb,
		infra.Rule{
			Limit:  5,
			Window: 1 * time.Second,
		},
		infra.Rule{
			Limit:  5,
			Window: 1 * time.Second,
		},
		infra.Rule{
			Limit:  2,
			Window: 1 * time.Second,
		},
		infra.Rule{
			Limit:  2,
			Window: 1 * time.Second,
		},
		infra.Rule{
			Limit:  4,
			Window: 1 * time.Second,
		},
		infra.Rule{
			Limit:  2,
			Window: 1 * time.Second,
		},
		infra.Rule{
			Limit:  4,
			Window: 1 * time.Second,
		},
	)

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

	healthCtl := controller.NewHealthCtl(healthUC)
	authCtl := controller.NewAuthCtl(authUC)
	itemCtl := controller.NewItemCtl(itemUC)
	srcCtl := controller.NewSrcCtl(sourceUC)

	router.New(
		e,
		healthCtl,
		authCtl,
		itemCtl,
		srcCtl,
	)

	log.Fatal(e.Start(":" + c.Port))
}
