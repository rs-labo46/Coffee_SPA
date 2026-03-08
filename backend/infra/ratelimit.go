package infra

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// Rule は rate limit 1件分
type Rule struct {
	Limit  int64
	Window time.Duration
}

// RateLimiter は Redis + Lua の rate limiter
type RateLimiter struct {
	rdb     *redis.Client
	login   Rule
	refresh Rule
	resend  Rule
	forgot  Rule
}

// Lua:
// 現在値を INCR
// 初回だけ EXPIRE 設定
// limit 超過なら deny
// retry-after 秒も返す
var allowScript = redis.NewScript(`
local current = redis.call("INCR", KEYS[1])
if current == 1 then
  redis.call("EXPIRE", KEYS[1], ARGV[2])
end

local ttl = redis.call("TTL", KEYS[1])

if current > tonumber(ARGV[1]) then
  return {0, ttl}
end

return {1, ttl}
`)

func NewRateLimiter(
	rdb *redis.Client,
	login Rule,
	refresh Rule,
	resend Rule,
	forgot Rule,
) *RateLimiter {
	return &RateLimiter{
		rdb:     rdb,
		login:   login,
		refresh: refresh,
		resend:  resend,
		forgot:  forgot,
	}
}

// login用rate limit
func (r *RateLimiter) AllowLogin(ip string) (bool, int, error) {
	key := "rl:login:ip:" + ip
	return r.allow(key, r.login)
}

// refresh用 rate limit
func (r *RateLimiter) AllowRefresh(ip string) (bool, int, error) {
	key := "rl:refresh:ip:" + ip
	return r.allow(key, r.refresh)
}

// verify再送用 rate limit
func (r *RateLimiter) AllowResend(ip string, emailHash string) (bool, int, error) {
	key := "rl:resend:ip:" + ip + ":email:" + emailHash
	return r.allow(key, r.resend)
}

// forgot password用rate limit
func (r *RateLimiter) AllowForgot(ip string, emailHash string) (bool, int, error) {
	key := "rl:forgot:ip:" + ip + ":email:" + emailHash
	return r.allow(key, r.forgot)
}

// 共通の実行
func (r *RateLimiter) allow(key string, rule Rule) (bool, int, error) {
	ctx := context.Background()

	out, err := allowScript.Run(
		ctx,
		r.rdb,
		[]string{key},
		rule.Limit,
		int64(rule.Window/time.Second),
	).Result()
	if err != nil {
		return false, 0, err
	}

	xs, ok := out.([]interface{})
	if !ok || len(xs) != 2 {
		return false, 0, redis.ErrClosed
	}

	allow, err := toI64(xs[0])
	if err != nil {
		return false, 0, err
	}

	retry, err := toI64(xs[1])
	if err != nil {
		return false, 0, err
	}

	return allow == 1, int(retry), nil
}

// Lua結果をint64に寄せる
func toI64(v interface{}) (int64, error) {
	switch x := v.(type) {
	case int64:
		return x, nil
	case string:
		return strconv.ParseInt(x, 10, 64)
	default:
		return 0, redis.ErrClosed
	}
}
