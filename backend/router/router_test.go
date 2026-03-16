package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"coffee-spa/controller"
	"coffee-spa/entity"
	"coffee-spa/repository"
	"coffee-spa/usecase"

	"github.com/labstack/echo/v4"
)

// mock controller
// health controllerモック。
type mockHealthUC struct{}

func (m *mockHealthUC) Check() error {
	return nil
}

// auth controllerのモック。
type mockAuthUCForRouter struct{}

func (m *mockAuthUCForRouter) Signup(in usecase.SignupIn) (entity.User, error) {
	return entity.User{}, nil
}
func (m *mockAuthUCForRouter) VerifyEmail(in usecase.VerifyEmailIn) error {
	return nil
}
func (m *mockAuthUCForRouter) ResendVerify(in usecase.ResendVerifyIn) error {
	return nil
}
func (m *mockAuthUCForRouter) Login(in usecase.LoginIn) (usecase.AuthOut, error) {
	return usecase.AuthOut{}, nil
}
func (m *mockAuthUCForRouter) Refresh(in usecase.RefreshIn) (usecase.AuthOut, error) {
	return usecase.AuthOut{}, nil
}
func (m *mockAuthUCForRouter) Logout(in usecase.LogoutIn) error {
	return nil
}
func (m *mockAuthUCForRouter) ForgotPw(in usecase.ForgotPwIn) error {
	return nil
}
func (m *mockAuthUCForRouter) ResetPw(in usecase.ResetPwIn) error {
	return nil
}
func (m *mockAuthUCForRouter) Me(userID int64) (entity.User, error) {
	return entity.User{
		ID:            userID,
		Email:         "user@example.com",
		Role:          "user",
		EmailVerified: true,
	}, nil
}

// item controllerモック。
type mockItemUCForRouter struct{}

func (m *mockItemUCForRouter) Add(actor usecase.Actor, in usecase.AddItemIn) (entity.Item, error) {
	return entity.Item{}, nil
}
func (m *mockItemUCForRouter) Search(q usecase.ItemQ) ([]entity.Item, error) {
	return []entity.Item{}, nil
}
func (m *mockItemUCForRouter) Top(limit int) (usecase.TopItems, error) {
	return usecase.TopItems{
		News:   []entity.Item{},
		Recipe: []entity.Item{},
		Deal:   []entity.Item{},
		Shop:   []entity.Item{},
	}, nil
}

// source controllerモック。
type mockSourceUCForRouter struct{}

func (m *mockSourceUCForRouter) Add(actor usecase.Actor, in usecase.AddSourceIn) (entity.Source, error) {
	return entity.Source{}, nil
}
func (m *mockSourceUCForRouter) List() ([]entity.Source, error) {
	return []entity.Source{}, nil
}

//---- mock repository 群 ----

// TokenVersion middlewareモック。
type mockUserRepoForRouter struct{}

func (m *mockUserRepoForRouter) Create(u entity.User) (entity.User, error) {
	return entity.User{}, nil
}
func (m *mockUserRepoForRouter) GetByID(id int64) (entity.User, error) {
	return entity.User{
		ID:       id,
		Role:     "user",
		TokenVer: 1,
	}, nil
}
func (m *mockUserRepoForRouter) GetByEmail(email string) (entity.User, error) {
	return entity.User{}, repository.ErrNotFound
}
func (m *mockUserRepoForRouter) SetEmailVerified(userID int64) error {
	return nil
}
func (m *mockUserRepoForRouter) UpdateEmailVerified(userID int64, ok bool) error {
	return nil
}
func (m *mockUserRepoForRouter) UpdatePassHash(userID int64, passHash string) error {
	return nil
}
func (m *mockUserRepoForRouter) BumpTokenVer(userID int64) (int, error) {
	return 2, nil
}

type mockRtRepoForRouter struct{}

func (m *mockRtRepoForRouter) Create(rt entity.RefreshToken) (entity.RefreshToken, error) {
	return entity.RefreshToken{}, nil
}
func (m *mockRtRepoForRouter) GetByTokenHash(hash string) (entity.RefreshToken, error) {
	return entity.RefreshToken{
		ID:       1,
		UserID:   1,
		FamilyID: "fam-1",
	}, nil
}
func (m *mockRtRepoForRouter) Revoke(id int64) error {
	return nil
}
func (m *mockRtRepoForRouter) MarkUsed(id int64) error {
	return nil
}
func (m *mockRtRepoForRouter) SetReplacedBy(id int64, newID int64) error {
	return nil
}
func (m *mockRtRepoForRouter) RevokeByFamilyID(familyID string) error {
	return nil
}
func (m *mockRtRepoForRouter) RevokeAllByUser(userID int64) error {
	return nil
}

// router.Newを通したEchoを作るヘルパ。
func newTestEcho() *echo.Echo {
	e := echo.New()

	healthCtl := controller.NewHealthCtl(&mockHealthUC{})
	authCtl := controller.NewAuthCtl(&mockAuthUCForRouter{})
	itemCtl := controller.NewItemCtl(&mockItemUCForRouter{})
	srcCtl := controller.NewSrcCtl(&mockSourceUCForRouter{})

	//RefreshRateLimit用のlimiter
	rlStore := repository.NewRateLimitStore(nil)
	rl := usecase.NewRateLimitUC(
		rlStore,
		usecase.RateRule{},
		usecase.RateRule{},
		usecase.RateRule{},
		usecase.RateRule{},
		usecase.RateRule{},
		usecase.RateRule{},
		usecase.RateRule{},
	)

	New(
		e,
		healthCtl,
		authCtl,
		itemCtl,
		srcCtl,
		"test-secret",
		&mockUserRepoForRouter{},
		&mockRtRepoForRouter{},
		rl,
		"http://localhost:3000",
	)

	return e
}

//public endpoint である GET /items/top が存在ことを確認。

func TestRouter_PublicItemsTop_Exists(t *testing.T) {
	t.Parallel()

	e := newTestEcho()

	req := httptest.NewRequest(http.MethodGet, "/items/top?limit=0", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Fatalf("route not found")
	}
}

// public endpointであるGET /sourcesが存在ことを確認。
func TestRouter_PublicSources_Exists(t *testing.T) {
	t.Parallel()

	e := newTestEcho()

	req := httptest.NewRequest(http.MethodGet, "/sources", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Fatalf("route not found")
	}
}

// POST /auth/refresh が少なくともpublicな素通しではないことを確認。
func TestRouter_RefreshRoute_ProtectedByCSRFOrAuthLayer(t *testing.T) {
	t.Parallel()

	e := newTestEcho()

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	//成功してしまうのは危険。
	if rec.Code == http.StatusOK {
		t.Fatalf("refresh route should not pass without protection")
	}
	if rec.Code == http.StatusNotFound {
		t.Fatalf("refresh route not found")
	}
}

// POST /auth/logoutがJWT保護下にあることを確認。
func TestRouter_LogoutRoute_Protected(t *testing.T) {
	t.Parallel()

	e := newTestEcho()

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusOK {
		t.Fatalf("logout route should require auth")
	}
	if rec.Code == http.StatusNotFound {
		t.Fatalf("logout route not found")
	}
}

// POST /itemsがadmin側の保護下にあることを確認。
func TestRouter_AdminItemCreate_Protected(t *testing.T) {
	t.Parallel()

	e := newTestEcho()

	req := httptest.NewRequest(http.MethodPost, "/items", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusOK {
		t.Fatalf("admin create route should require auth/admin")
	}
	if rec.Code == http.StatusNotFound {
		t.Fatalf("admin create route not found")
	}
}
