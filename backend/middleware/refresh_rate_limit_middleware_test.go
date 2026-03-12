package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"coffee-spa/entity"

	"github.com/labstack/echo/v4"
)

type mockRtRepository struct {
	getByTokenHashFn func(hash string) (entity.RefreshToken, error)
}

func (m *mockRtRepository) Create(rt entity.RefreshToken) (entity.RefreshToken, error) {
	return entity.RefreshToken{}, nil
}
func (m *mockRtRepository) GetByTokenHash(hash string) (entity.RefreshToken, error) {
	return m.getByTokenHashFn(hash)
}
func (m *mockRtRepository) Revoke(id int64) error                     { return nil }
func (m *mockRtRepository) MarkUsed(id int64) error                   { return nil }
func (m *mockRtRepository) SetReplacedBy(id int64, newID int64) error { return nil }
func (m *mockRtRepository) RevokeByFamilyID(familyID string) error    { return nil }
func (m *mockRtRepository) RevokeAllByUser(userID int64) error        { return nil }

type mockRefreshLimiter struct {
	allowFn func(userID int64) (bool, int, error)
}

func (m *mockRefreshLimiter) AllowRefresh(userID int64) (bool, int, error) {
	return m.allowFn(userID)
}

// rate limit に引っかかったとき 429 と Retry-After が返ることを確認。
func TestRefreshRateLimit_TooManyRequests(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: "raw-token",
	})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mw := RefreshRateLimit(
		&mockRtRepository{
			getByTokenHashFn: func(hash string) (entity.RefreshToken, error) {
				return entity.RefreshToken{
					UserID: 7,
				}, nil
			},
		},
		&mockRefreshLimiter{
			allowFn: func(userID int64) (bool, int, error) {
				if userID != 7 {
					t.Fatalf("userID = %d, want 7", userID)
				}
				return false, 9, nil
			},
		},
	)

	h := mw(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	if err := h(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want 429", rec.Code)
	}
	if got := rec.Header().Get("Retry-After"); got != "9" {
		t.Fatalf("Retry-After = %q, want 9", got)
	}
}
