package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Cfg struct {
	Port   string
	PgUser string
	PgPass string
	PgDB   string
	PgHost string
	PgPort int
	GoEnv  string
}

func Load() (Cfg, error) {
	_ = godotenv.Load()

	c := Cfg{}
	c.Port = getenv("PORT", "8080")                 // APIのポート
	c.PgUser = mustGet("POSTGRES_USER")             // DBユーザ
	c.PgPass = mustGet("POSTGRES_PASSWORD")         // DBパスワード
	c.PgDB = mustGet("POSTGRES_DB")                 // DB名
	c.PgHost = getenv("POSTGRES_HOST", "localhost") // DBホスト
	c.GoEnv = getenv("GO_ENV", "dev")               // 実行環境

	p := getenv("POSTGRES_PORT", "5433")
	n, err := strconv.Atoi(p)
	if err != nil {
		n = 5433
	}
	c.PgPort = n

	return c, nil
}

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}

func mustGet(k string) string {
	v := os.Getenv(k)
	if v == "" {
		panic("missing env: " + k)
	}
	return v
}
