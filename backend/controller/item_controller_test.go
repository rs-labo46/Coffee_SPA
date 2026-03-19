package controller

import (
	"bytes"
	"coffee-spa/entity"
	"coffee-spa/usecase"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

type mockItemUC struct {
	topFn    func(limit int) (usecase.TopItems, error)
	searchFn func(q usecase.ItemQ) ([]entity.Item, error)
	addFn    func(actor usecase.Actor, in usecase.AddItemIn) (entity.Item, error)

	addCalled bool
}

func (m *mockItemUC) Top(limit int) (usecase.TopItems, error) {
	if m.topFn != nil {
		return m.topFn(limit)
	}
	return usecase.TopItems{}, nil
}

func (m *mockItemUC) Search(q usecase.ItemQ) ([]entity.Item, error) {
	if m.searchFn != nil {
		return m.searchFn(q)
	}
	return nil, nil
}

func (m *mockItemUC) Add(actor usecase.Actor, in usecase.AddItemIn) (entity.Item, error) {
	m.addCalled = true

	if m.addFn != nil {
		return m.addFn(actor, in)
	}
	return entity.Item{}, nil
}

// GET /items/top の成功確認。
// queryがusecaseに渡り、4キーのJSONshapeが返ることを確認。
func TestItemCtlTop_OK(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/items/top?limit=0", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	want := usecase.TopItems{
		News:   []entity.Item{{ID: 1, Title: "n1", Kind: "news"}},
		Recipe: []entity.Item{{ID: 2, Title: "r1", Kind: "recipe"}},
		Deal:   []entity.Item{{ID: 3, Title: "d1", Kind: "deal"}},
		Shop:   []entity.Item{{ID: 4, Title: "s1", Kind: "shop"}},
	}

	ctl := NewItemCtl(&mockItemUC{
		topFn: func(limit int) (usecase.TopItems, error) {
			if limit != 0 {
				t.Fatalf("limit = %d, want 0", limit)
			}
			return want, nil
		},
	})

	if err := ctl.Top(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	var body TopItemsRes
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(body.News) != 1 || len(body.Recipe) != 1 || len(body.Deal) != 1 || len(body.Shop) != 1 {
		t.Fatalf("unexpected top item group lengths: %+v", body)
	}
	if body.News[0].ID != want.News[0].ID ||
		body.Recipe[0].ID != want.Recipe[0].ID ||
		body.Deal[0].ID != want.Deal[0].ID ||
		body.Shop[0].ID != want.Shop[0].ID {
		t.Fatalf("unexpected top items body: %+v", body)
	}
}

func TestItemCtlTop_InvalidLimit_ReturnsBadRequest(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/items/top?limit=bad", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	called := false
	ctl := NewItemCtl(&mockItemUC{
		topFn: func(limit int) (usecase.TopItems, error) {
			called = true
			return usecase.TopItems{}, nil
		},
	})

	if err := ctl.Top(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	if called {
		t.Fatalf("usecase.Top must not be called on invalid query")
	}

	var body ErrRes
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if body.Error != "invalid_request" {
		t.Fatalf("unexpected error body: %+v", body)
	}
}

func TestItemCtlList_OK(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/items?q=coffee&kind=news&limit=5&offset=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	wantTime := time.Date(2026, 3, 18, 12, 0, 0, 0, time.UTC)
	want := entity.Item{
		ID:        1,
		Title:     "news",
		Kind:      "news",
		CreatedAt: wantTime,
	}

	ctl := NewItemCtl(&mockItemUC{
		searchFn: func(q usecase.ItemQ) ([]entity.Item, error) {
			if q.Q != "coffee" || q.Kind != "news" || q.Limit != 5 || q.Offset != 10 {
				t.Fatalf("unexpected query: %+v", q)
			}
			return []entity.Item{want}, nil
		},
	})

	if err := ctl.List(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	var body ItemListRes
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(body.Items) != 1 {
		t.Fatalf("unexpected items length: %d", len(body.Items))
	}
	got := body.Items[0]
	if got.ID != want.ID || got.Title != want.Title || got.Kind != want.Kind || !got.CreatedAt.Equal(want.CreatedAt) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestItemCtlList_InvalidOffset_ReturnsBadRequest(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/items?offset=bad", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	called := false
	ctl := NewItemCtl(&mockItemUC{
		searchFn: func(q usecase.ItemQ) ([]entity.Item, error) {
			called = true
			return nil, nil
		},
	})

	if err := ctl.List(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	if called {
		t.Fatalf("usecase.Search must not be called on invalid query")
	}

	var body ErrRes
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if body.Error != "invalid_request" {
		t.Fatalf("unexpected error body: %+v", body)
	}
}

func TestItemCtlCreate_OK(t *testing.T) {
	t.Parallel()

	e := echo.New()
	payload := []byte(`{"title":"coffee news","kind":"news"}`)
	req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", int64(7))
	c.Set("role", "admin")

	ctl := NewItemCtl(&mockItemUC{
		addFn: func(actor usecase.Actor, in usecase.AddItemIn) (entity.Item, error) {
			if actor.UserID != 7 || actor.Role != "admin" {
				t.Fatalf("unexpected actor: %+v", actor)
			}
			if in.Title != "coffee news" || in.Kind != "news" {
				t.Fatalf("unexpected input: %+v", in)
			}
			return entity.Item{ID: 10, Title: in.Title, Kind: in.Kind}, nil
		},
	})

	if err := ctl.Create(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201", rec.Code)
	}

	var body ItemRes
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if body.Item.ID != 10 || body.Item.Title != "coffee news" || body.Item.Kind != "news" {
		t.Fatalf("unexpected body: %+v", body)
	}
}

func TestItemCtlCreate_InvalidJSON_ReturnsBadRequest(t *testing.T) {
	e := echo.New()

	mockUC := &mockItemUC{}
	ctl := NewItemCtl(mockUC)

	req := httptest.NewRequest(
		http.MethodPost,
		"/admin/items",
		strings.NewReader(`{"title":"coffee news",`),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", int64(1))
	c.Set("role", "admin")

	err := ctl.Create(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if mockUC.addCalled {
		t.Fatal("usecase.Add must not be called on invalid json")
	}
}
