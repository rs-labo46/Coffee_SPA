package usecase

import (
	"encoding/json"
	"errors"
	"time"

	"coffee-spa/entity"
	"coffee-spa/repository"

	"gorm.io/datatypes"
)

type ItemVal interface {
	NewItem(in AddItemIn) error
	ListItem(q ItemQ) error
}

type ItemUC struct {
	item   repository.ItemRepository
	source repository.SourceRepository
	audit  repository.AuditRepository
	val    ItemVal
}

type itemCreateMeta struct {
	UserID int64 `json:"user_id"`
	ItemID int64 `json:"item_id"`
}

func NewItemUC(
	item repository.ItemRepository,
	source repository.SourceRepository,
	audit repository.AuditRepository,
	val ItemVal,
) ItemUsecase {
	return &ItemUC{
		item:   item,
		source: source,
		audit:  audit,
		val:    val,
	}
}

func (u *ItemUC) Add(actor Actor, in AddItemIn) (entity.Item, error) {
	//認可はmiddlewareで実施。
	if err := u.val.NewItem(in); err != nil {
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

	err = u.audit.Create(entity.AuditLog{
		Type:     "admin.items.create",
		UserID:   toI64Ptr(actor.UserID),
		IP:       actor.IP,
		UA:       actor.UA,
		MetaJSON: datatypes.JSON(b),
	})
	if err != nil {
		return entity.Item{}, mapRepoErr(err)
	}

	return item, nil
}

func (u *ItemUC) Search(q ItemQ) ([]entity.Item, error) {
	if err := u.val.ListItem(q); err != nil {
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
