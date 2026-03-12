package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"coffee-spa/entity"

	"github.com/labstack/echo/v4"
)

// refresh token repositoryの代役。
// middlewareがtoken hashからuserIDを引けるかを差し替えるために使う。
type mockRtRepo struct {
	getByTokenHashFn func(hash string) (entity.RefreshToken, error)
}

func (m *mockRtRepo) Create(rt entity.RefreshToken) (entity.RefreshToken, error) {
	return entity.RefreshToken{}, nil
}
func (m *mockRtRepo) GetByTokenHash(hash string) (entity.RefreshToken, error) {
	return m.getByTokenHashFn(hash)
}
func (m *mockRtRepo) Revoke(id int64) error                     { return nil }
func (m *mockRtRepo) MarkUsed(id int64) error                   { return nil }
func (m *mockRtRepo) SetReplacedBy(id int64, newID int64) error { return nil }
func (m *mockRtRepo) RevokeByFamilyID(familyID string) error    { return nil }
func (m *mockRtRepo) RevokeAllByUser(userID int64) error        { return nil }

// limiterの代役。
// 指定ユーザーに対して許可か拒否かをテストごとに切り替える。
type mockRefreshLimiter struct {
	allowFn func(userID int64) (bool, int, error)
}

func (m *mockRefreshLimiter) AllowRefresh(userID int64) (bool, int, error) {
	return m.allowFn(userID)
}

// rate limitに引っかかったとき429とRetry-Afterが返ることを確認。
func TestRefreshRateLimit_TooManyRequests(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)

	//refresh token cookieを付けて、通常のrefreshリクエストを再現。
	req.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: "raw-token",
	})

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mw := RefreshRateLimit(
		&mockRtRepo{
			getByTokenHashFn: func(hash string) (entity.RefreshToken, error) {
				//token hashからuserID=7を引けた想定を返す。
				return entity.RefreshToken{
					UserID: 7,
				}, nil
			},
		},
		&mockRefreshLimiter{
			allowFn: func(userID int64) (bool, int, error) {
				//正しいuserIDでlimiterが呼ばれているか確認。
				if userID != 7 {
					t.Fatalf("userID = %d, want 7", userID)
				}

				//rate limit 拒否のケースを作る。
				return false, 9, nil
			},
		},
	)

	//middlewareで止まるので呼ばれない想定。
	h := mw(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	if err := h(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	//rate limit 拒否なので429を期待。
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want 429", rec.Code)
	}

	//再試行までの秒数が header に入っていることを確認。
	if got := rec.Header().Get("Retry-After"); got != "9" {
		t.Fatalf("Retry-After = %q, want 9", got)
	}
}
