package repository

import (
	"errors"

	"coffee-spa/entity"

	"gorm.io/gorm"
)

type itemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) ItemRepository {
	return &itemRepository{db}
}

func (r *itemRepository) Create(i entity.Item) (entity.Item, error) {
	err := r.db.Create(&i).Error
	if err != nil {
		if isDup(err) || isFK(err) {
			return entity.Item{}, ErrConflict
		}
		return entity.Item{}, ErrInternal
	}

	return i, nil
}

func (r *itemRepository) GetByID(id int64) (entity.Item, error) {
	var i entity.Item

	err := r.db.
		First(&i, id).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Item{}, ErrNotFound
		}
		return entity.Item{}, ErrInternal
	}

	return i, nil
}

func (r *itemRepository) List(q ItemQ) ([]entity.Item, error) {
	var xs []entity.Item

	tx := r.db.Model(&entity.Item{})

	if q.Kind != "" {
		tx = tx.Where("kind = ?", q.Kind)
	}

	if q.Q != "" {
		like := "%" + q.Q + "%"
		tx = tx.Where(
			"title ILIKE ? OR COALESCE(summary, '') ILIKE ?",
			like,
			like,
		)
	}

	lim := q.Limit
	if lim <= 0 {
		lim = 20
	}
	if lim > 50 {
		lim = 50
	}

	off := q.Offset
	if off < 0 {
		off = 0
	}

	err := tx.
		Order("published_at DESC").
		Order("created_at DESC").
		Limit(lim).
		Offset(off).
		Find(&xs).
		Error
	if err != nil {
		return nil, ErrInternal
	}

	return xs, nil
}

func (r *itemRepository) Top(cap int) (TopItems, error) {
	if cap < 0 {
		cap = 0
	}
	if cap > 50 {
		cap = 50
	}

	topItems := TopItems{
		News:   []entity.Item{},
		Recipe: []entity.Item{},
		Deal:   []entity.Item{},
		Shop:   []entity.Item{},
	}

	// cap=0 は4キー固定で空配列を返す。
	if cap == 0 {
		return topItems, nil
	}

	err := r.db.
		Where("kind = ?", "news").
		Order("published_at DESC").
		Order("created_at DESC").
		Limit(cap).
		Find(&topItems.News).
		Error
	if err != nil {
		return TopItems{}, ErrInternal
	}

	err = r.db.
		Where("kind = ?", "recipe").
		Order("published_at DESC").
		Order("created_at DESC").
		Limit(cap).
		Find(&topItems.Recipe).
		Error
	if err != nil {
		return TopItems{}, ErrInternal
	}

	err = r.db.
		Where("kind = ?", "deal").
		Order("published_at DESC").
		Order("created_at DESC").
		Limit(cap).
		Find(&topItems.Deal).
		Error
	if err != nil {
		return TopItems{}, ErrInternal
	}

	err = r.db.
		Where("kind = ?", "shop").
		Order("published_at DESC").
		Order("created_at DESC").
		Limit(cap).
		Find(&topItems.Shop).
		Error
	if err != nil {
		return TopItems{}, ErrInternal
	}

	return topItems, nil
}
