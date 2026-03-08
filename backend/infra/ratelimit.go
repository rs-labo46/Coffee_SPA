package infra

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// Ruleはrate limit 1件分
type Rule struct {
	Limit  int64
	Window time.Duration
}

// RateLimiterはRedis + Lua のrate limiter
type RateLimiter struct {
	rdb        *redis.Client
	signupIP   Rule
	loginMail  Rule
	refreshUID Rule
	resendIP   Rule
	resendMail Rule
	forgotIP   Rule
	forgotMail Rule
}

// 秒単位 fixed window。
// return {allowed, retry_after_sec}
var allowScript = redis.NewScript(`
local current = redis.call("INCR", KEYS[1])
if current == 1 then
  redis.call("EXPIRE", KEYS[1], ARGV[2])
end

if current <= tonumber(ARGV[1]) then
  return {1, 0}
end

local ttl = redis.call("TTL", KEYS[1])
if ttl < 0 then
  ttl = tonumber(ARGV[2])
end

return {0, ttl}
`)

func NewRateLimiter(
	rdb *redis.Client,
	signupIP Rule,
	loginMail Rule,
	refreshUID Rule,
	resendIP Rule,
	resendMail Rule,
	forgotIP Rule,
	forgotMail Rule,
) *RateLimiter {
	return &RateLimiter{
		rdb:        rdb,
		signupIP:   signupIP,
		loginMail:  loginMail,
		refreshUID: refreshUID,
		resendIP:   resendIP,
		resendMail: resendMail,
		forgotIP:   forgotIP,
		forgotMail: forgotMail,
	}
}

func (r *RateLimiter) AllowSignup(ip string) (bool, int, error) {
	key := "rl:signup:ip:" + ip
	return r.allow(key, r.signupIP)
}

func (r *RateLimiter) AllowLogin(emailHash string) (bool, int, error) {
	key := "rl:login:mail:" + emailHash
	return r.allow(key, r.loginMail)
}

func (r *RateLimiter) AllowRefresh(userID int64) (bool, int, error) {
	key := fmt.Sprintf("rl:refresh:uid:%d", userID)
	return r.allow(key, r.refreshUID)
}

func (r *RateLimiter) AllowResendIP(ip string) (bool, int, error) {
	key := "rl:resend:ip:" + ip
	return r.allow(key, r.resendIP)
}

func (r *RateLimiter) AllowResendMail(emailHash string) (bool, int, error) {
	key := "rl:resend:mail:" + emailHash
	return r.allow(key, r.resendMail)
}

func (r *RateLimiter) AllowForgotIP(ip string) (bool, int, error) {
	key := "rl:forgot:ip:" + ip
	return r.allow(key, r.forgotIP)
}

func (r *RateLimiter) AllowForgotMail(emailHash string) (bool, int, error) {
	key := "rl:forgot:mail:" + emailHash
	return r.allow(key, r.forgotMail)
}

func (r *RateLimiter) allow(key string, rule Rule) (bool, int, error) {
	if r == nil || r.rdb == nil {
		return false, 0, errors.New("redis client is nil")
	}
	if rule.Limit <= 0 {
		return false, 0, errors.New("invalid rate limit")
	}

	sec := int64(rule.Window / time.Second)
	if sec <= 0 {
		sec = 1
	}

	out, err := allowScript.Run(
		context.Background(),
		r.rdb,
		[]string{key},
		rule.Limit,
		sec,
	).Result()
	if err != nil {
		return false, 0, err
	}

	xs, ok := out.([]interface{})
	if !ok || len(xs) != 2 {
		return false, 0, errors.New("invalid lua result")
	}

	allowed, err := toI64(xs[0])
	if err != nil {
		return false, 0, err
	}

	retry, err := toI64(xs[1])
	if err != nil {
		return false, 0, err
	}

	return allowed == 1, int(retry), nil
}

func toI64(v interface{}) (int64, error) {
	switch x := v.(type) {
	case int64:
		return x, nil
	case string:
		return strconv.ParseInt(x, 10, 64)
	default:
		return 0, errors.New("invalid lua value")
	}
}
