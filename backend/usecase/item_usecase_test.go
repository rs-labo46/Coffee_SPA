package usecase

import (
	"errors"
	"testing"
	"time"

	"coffee-spa/entity"
	"coffee-spa/repository"
)

// item repositoryの代役。
// DB を使わずに、usecaseが期待どおりrepositoryを呼ぶかだけを確認。
type mockItemRepo struct {
	topFn    func(cap int) (repository.TopItems, error)
	createFn func(i entity.Item) (entity.Item, error)
	listFn   func(q repository.ItemQ) ([]entity.Item, error)
}

func (m *mockItemRepo) Create(i entity.Item) (entity.Item, error) {
	return m.createFn(i)
}

func (m *mockItemRepo) GetByID(id int64) (entity.Item, error) {
	return entity.Item{}, nil
}

func (m *mockItemRepo) List(q repository.ItemQ) ([]entity.Item, error) {
	return m.listFn(q)
}

func (m *mockItemRepo) Top(cap int) (repository.TopItems, error) {
	return m.topFn(cap)
}

// source repositoryの代役。
// Add時にSourceIDの存在確認
type mockSourceRepo struct {
	getByIDFn func(id int64) (entity.Source, error)
}

func (m *mockSourceRepo) Create(s entity.Source) (entity.Source, error) {
	return entity.Source{}, nil
}

func (m *mockSourceRepo) GetByID(id int64) (entity.Source, error) {
	return m.getByIDFn(id)
}

func (m *mockSourceRepo) GetByName(name string) (entity.Source, error) {
	return entity.Source{}, nil
}

func (m *mockSourceRepo) List() ([]entity.Source, error) {
	return nil, nil
}

// audit repositoryの代役。
// 成功時にauditlogを残しているかを確認。
type mockAuditRepo struct {
	createFn func(a entity.AuditLog) error
}

func (m *mockAuditRepo) Create(a entity.AuditLog) error {
	return m.createFn(a)
}

// mockItemValはvalidatorの代役。
// usecaseがvalidatorを先に通しているか。
type mockItemVal struct {
	newItemFn  func(input AddItemIn) error
	listItemFn func(q ItemQ) error
}

func (m *mockItemVal) NewItem(input AddItemIn) error {
	return m.newItemFn(input)
}

func (m *mockItemVal) ListItem(q ItemQ) error {
	return m.listItemFn(q)
}

// Top(cap=0) のときに 4キー固定の空配列が返ることを確認する。
func TestItemUCTop_ZeroCap_ReturnsEmptyGroups(t *testing.T) {
	t.Parallel()

	uc := &ItemUC{
		item: &mockItemRepo{
			topFn: func(cap int) (repository.TopItems, error) {
				//usecaseがrepositoryにcap=0 をそのまま渡しているか確認する。
				if cap != 0 {
					t.Fatalf("cap = %d, want 0", cap)
				}

				//repositoryが4キー固定の空配列を返した想定。
				return repository.TopItems{
					News:   []entity.Item{},
					Recipe: []entity.Item{},
					Deal:   []entity.Item{},
					Shop:   []entity.Item{},
				}, nil
			},
		},
		source: &mockSourceRepo{},
		audit: &mockAuditRepo{
			createFn: func(a entity.AuditLog) error { return nil },
		},
		val: &mockItemVal{
			newItemFn: func(input AddItemIn) error { return nil },
			listItemFn: func(q ItemQ) error {
				return nil
			},
		},
	}

	got, err := uc.Top(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	//APIでは[]とnull が違うため、ここを確認する。
	if got.News == nil || got.Recipe == nil || got.Deal == nil || got.Shop == nil {
		t.Fatalf("expected all groups to be non-nil slices")
	}

	//4グループとも 0件であることを確認する。
	if len(got.News) != 0 || len(got.Recipe) != 0 || len(got.Deal) != 0 || len(got.Shop) != 0 {
		t.Fatalf("expected all groups empty")
	}
}

// repositoryエラーがusecaseエラーに適切に変換されることを確認する。
func TestItemUCTop_RepoError_Mapped(t *testing.T) {
	t.Parallel()

	uc := &ItemUC{
		item: &mockItemRepo{
			topFn: func(cap int) (repository.TopItems, error) {
				//repositoryで内部エラーが起きた想定。
				return repository.TopItems{}, repository.ErrInternal
			},
		},
		source: &mockSourceRepo{},
		audit: &mockAuditRepo{
			createFn: func(a entity.AuditLog) error { return nil },
		},
		val: &mockItemVal{
			newItemFn: func(input AddItemIn) error { return nil },
			listItemFn: func(q ItemQ) error {
				return nil
			},
		},
	}

	_, err := uc.Top(3)
	if !errors.Is(err, ErrInternal) {
		t.Fatalf("err = %v, want ErrInternal", err)
	}
}

// Add 成功時に item 作成と audit 作成が呼ばれることを確認する。
func TestItemUCAdd_OK(t *testing.T) {
	t.Parallel()

	var created bool
	var audited bool

	uc := &ItemUC{
		item: &mockItemRepo{
			createFn: func(i entity.Item) (entity.Item, error) {
				created = true
				//DB保存後にIDが付いた想定。
				i.ID = 10
				return i, nil
			},
			topFn: func(cap int) (repository.TopItems, error) {
				return repository.TopItems{}, nil
			},
			listFn: func(q repository.ItemQ) ([]entity.Item, error) {
				return nil, nil
			},
		},
		source: &mockSourceRepo{
			getByIDFn: func(id int64) (entity.Source, error) {
				//SourceID が存在している前提を返す。
				return entity.Source{ID: id, Name: "test"}, nil
			},
		},
		audit: &mockAuditRepo{
			createFn: func(a entity.AuditLog) error {
				audited = true
				return nil
			},
		},
		val: &mockItemVal{
			newItemFn: func(input AddItemIn) error { return nil },
			listItemFn: func(q ItemQ) error {
				return nil
			},
		},
	}

	item, err := uc.Add(Actor{
		UserID: 1,
		IP:     "127.0.0.1",
		UA:     "test",
	}, AddItemIn{
		Title:       "coffee",
		Kind:        "news",
		SourceID:    1,
		PublishedAt: time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	//item repositoryのCreateが実行されたことを確認。
	if !created {
		t.Fatalf("item create was not called")
	}

	//audit repositoryのCreateが実行されたことを確認。
	if !audited {
		t.Fatalf("audit create was not called")
	}

	//保存後のIDが返ってくることを確認。
	if item.ID != 10 {
		t.Fatalf("item.ID = %d, want 10", item.ID)
	}
}
