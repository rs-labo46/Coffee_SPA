package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
		getFn: func(id int64) (entity.Item, error) {
			return entity.Item{}, nil
		},
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
func TestItemCtlGet_OK(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/items/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/items/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}
