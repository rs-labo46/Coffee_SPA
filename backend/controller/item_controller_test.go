package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"coffee-spa/entity"
	"coffee-spa/usecase"

	"github.com/labstack/echo/v4"
)

type mockItemUC struct {
	getFn    func(id int64) (entity.Item, error)
	topFn    func(limit int) (usecase.TopItems, error)
	searchFn func(q usecase.ItemQ) ([]entity.Item, error)
	addFn    func(actor usecase.Actor, in usecase.AddItemIn) (entity.Item, error)
}

func (m *mockItemUC) Get(id int64) (entity.Item, error) {
	if m.getFn == nil {
		return entity.Item{}, nil
	}
	return m.getFn(id)
}

func (m *mockItemUC) Add(actor usecase.Actor, in usecase.AddItemIn) (entity.Item, error) {
	return m.addFn(actor, in)
}

func (m *mockItemUC) Search(q usecase.ItemQ) ([]entity.Item, error) {
	return m.searchFn(q)
}

func (m *mockItemUC) Top(limit int) (usecase.TopItems, error) {
	return m.topFn(limit)
}

// GET /items/top の成功確認。
// JSONの4キーが返ることを確認。
func TestItemCtlTop_OK(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/items/top?limit=0", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	ctl := NewItemCtl(&mockItemUC{
		topFn: func(limit int) (usecase.TopItems, error) {
			if limit != 0 {
				t.Fatalf("limit = %d, want 0", limit)
			}
			return usecase.TopItems{
				News:   []entity.Item{},
				Recipe: []entity.Item{},
				Deal:   []entity.Item{},
				Shop:   []entity.Item{},
			}, nil
		},
		searchFn: func(q usecase.ItemQ) ([]entity.Item, error) {
			return nil, nil
		},
		addFn: func(actor usecase.Actor, in usecase.AddItemIn) (entity.Item, error) {
			return entity.Item{}, nil
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
	if body.News == nil || body.Recipe == nil || body.Deal == nil || body.Shop == nil {
		t.Fatalf("expected all top item groups to exist")
	}
}

func TestItemCtlList_OK(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/items?q=coffee&kind=news&limit=5&offset=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	wantTime := time.Date(2026, 3, 18, 12, 0, 0, 0, time.UTC)
	ctl := NewItemCtl(&mockItemUC{
		topFn: func(limit int) (usecase.TopItems, error) {
			return usecase.TopItems{}, nil
		},
		searchFn: func(q usecase.ItemQ) ([]entity.Item, error) {
			if q.Q != "coffee" || q.Kind != "news" || q.Limit != 5 || q.Offset != 10 {
				t.Fatalf("unexpected query: %+v", q)
			}
			return []entity.Item{{
				ID:        1,
				Title:     "news",
				Kind:      "news",
				CreatedAt: wantTime,
			}}, nil
		},
		addFn: func(actor usecase.Actor, in usecase.AddItemIn) (entity.Item, error) {
			return entity.Item{}, nil
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
	if body.Items[0].ID != 1 || body.Items[0].Title != "news" || body.Items[0].Kind != "news" {
		t.Fatalf("unexpected body: %+v", body)
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
		topFn: func(limit int) (usecase.TopItems, error) {
			return usecase.TopItems{}, nil
		},
		searchFn: func(q usecase.ItemQ) ([]entity.Item, error) {
			return nil, nil
		},
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
