package usecase

import (
	"encoding/json"
	"errors"
	"time"

	"coffee-spa/entity"
	"coffee-spa/repository"

	"gorm.io/datatypes"
)

type ItemUC struct {
	item   repository.ItemRepository
	source repository.SourceRepository
	audit  repository.AuditRepository
}

type itemCreateMeta struct {
	UserID int64 `json:"user_id"`
	ItemID int64 `json:"item_id"`
}

func NewItemUC(
	item repository.ItemRepository,
	source repository.SourceRepository,
	audit repository.AuditRepository,
) ItemUsecase {
	return &ItemUC{
		item:   item,
		source: source,
		audit:  audit,
	}
}

func (u *ItemUC) Add(actor Actor, in AddItemIn) (entity.Item, error) {
	if actor.Role != string(entity.RoleAdmin) {
		return entity.Item{}, ErrForbidden
	}

	if in.Title == "" || in.Kind == "" || in.SourceID == 0 || in.PublishedAt == "" {
		return entity.Item{}, ErrInvalidRequest
	}

	if !isItemKind(in.Kind) {
		return entity.Item{}, ErrInvalidRequest
	}

	publishedAt, err := time.Parse(time.RFC3339, in.PublishedAt)
	if err != nil {
		return entity.Item{}, ErrInvalidRequest
	}

	_, err = u.source.GetByID(in.SourceID)
	if err != nil {
		return entity.Item{}, mapRepoErr(err)
	}

	item, err := u.item.Create(entity.Item{
		Title:       in.Title,
		Summary:     in.Summary,
		URL:         in.URL,
		ImageURL:    in.ImageURL,
		Kind:        in.Kind,
		SourceID:    in.SourceID,
		PublishedAt: publishedAt,
	})
	if err != nil {
		return entity.Item{}, mapRepoErr(err)
	}

	b, err := json.Marshal(itemCreateMeta{
		UserID: actor.UserID,
		ItemID: item.ID,
	})
	if err != nil {
		return entity.Item{}, ErrInternal
	}

	auditErr := u.audit.Create(entity.AuditLog{
		Type:     "admin.items.create",
		UserID:   toI64Ptr(actor.UserID),
		IP:       actor.IP,
		UA:       actor.UA,
		MetaJSON: datatypes.JSON(b),
	})
	if auditErr != nil {
		return entity.Item{}, mapRepoErr(auditErr)
	}

	return item, nil
}

func (u *ItemUC) Search(q ItemQ) ([]entity.Item, error) {
	if q.Kind != "" && !isItemKind(q.Kind) {
		return nil, ErrInvalidRequest
	}

	if q.Limit < 0 || q.Offset < 0 {
		return nil, ErrInvalidRequest
	}

	items, err := u.item.List(repository.ItemQ{
		Q:      q.Q,
		Kind:   q.Kind,
		Limit:  q.Limit,
		Offset: q.Offset,
	})
	if err != nil {
		return nil, mapRepoErr(err)
	}

	return items, nil
}

func (u *ItemUC) Top(limit int) (TopItems, error) {
	if limit < 0 {
		return TopItems{}, ErrInvalidRequest
	}

	out, err := u.item.Top(limit)
	if err != nil {
		return TopItems{}, mapRepoErr(err)
	}

	return TopItems{
		News:   out.News,
		Recipe: out.Recipe,
		Deal:   out.Deal,
		Shop:   out.Shop,
	}, nil
}

func isItemKind(kind string) bool {
	switch kind {
	case string(entity.KindNews),
		string(entity.KindRecipe),
		string(entity.KindDeal),
		string(entity.KindShop):
		return true
	default:
		return false
	}
}

func mapRepoErr(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return ErrNotFound
	case errors.Is(err, repository.ErrConflict):
		return ErrConflict
	case errors.Is(err, repository.ErrInternal):
		return ErrInternal
	default:
		return ErrInternal
	}
}

func toI64Ptr(id int64) *int64 {
	if id == 0 {
		return nil
	}
	return &id
}
