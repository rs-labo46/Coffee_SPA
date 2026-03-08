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
	"context"
	"log"
	"os"
	"strings"
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

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr(),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal(err)
	}

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
			Limit:  1,
			Window: 5 * time.Second,
		},
		infra.Rule{
			Limit:  1,
			Window: 5 * time.Second,
		},
		infra.Rule{
			Limit:  1,
			Window: 2 * time.Second,
		},
		infra.Rule{
			Limit:  1,
			Window: 2 * time.Second,
		},
		infra.Rule{
			Limit:  1,
			Window: 4 * time.Second,
		},
		infra.Rule{
			Limit:  1,
			Window: 2 * time.Second,
		},
		infra.Rule{
			Limit:  1,
			Window: 4 * time.Second,
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
		c.JWTSecret,
		userRepo,
		c.FEURL,
	)

	log.Fatal(e.Start(":" + c.Port))
}

func redisAddr() string {
	addr := strings.TrimSpace(os.Getenv("REDIS_ADDR"))
	if addr != "" {
		return addr
	}

	host := strings.TrimSpace(os.Getenv("REDIS_HOST"))
	if host == "" {
		host = "localhost"
	}

	port := strings.TrimSpace(os.Getenv("REDIS_PORT"))
	if port == "" {
		port = "6379"
	}

	return host + ":" + port
}
