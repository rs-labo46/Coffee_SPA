package dbrepository

import (
	"errors"

	"coffee-spa/entity"
	"coffee-spa/repository"

	"gorm.io/gorm"
)

type ItemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) repository.ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) Create(i entity.Item) (entity.Item, error) {
	err := r.db.Create(&i).Error
	if err != nil {
		if isDup(err) || isFK(err) {
			return entity.Item{}, repository.ErrConflict
		}
		return entity.Item{}, repository.ErrInternal
	}

	return i, nil
}

func (r *ItemRepository) GetByID(id int64) (entity.Item, error) {
	var i entity.Item

	err := r.db.
		First(&i, id).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Item{}, repository.ErrNotFound
		}
		return entity.Item{}, repository.ErrInternal
	}

	return i, nil
}

func (r *ItemRepository) List(q repository.ItemQ) ([]entity.Item, error) {
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
		return nil, repository.ErrInternal
	}

	return xs, nil
}

func (r *ItemRepository) Top(cap int) (repository.TopItems, error) {
	if cap <= 0 {
		cap = 5
	}
	if cap > 50 {
		cap = 50
	}

	var out repository.TopItems

	err := r.db.
		Where("kind = ?", "news").
		Order("published_at DESC").
		Order("created_at DESC").
		Limit(cap).
		Find(&out.News).
		Error
	if err != nil {
		return repository.TopItems{}, repository.ErrInternal
	}

	err = r.db.
		Where("kind = ?", "recipe").
		Order("published_at DESC").
		Order("created_at DESC").
		Limit(cap).
		Find(&out.Recipe).
		Error
	if err != nil {
		return repository.TopItems{}, repository.ErrInternal
	}

	err = r.db.
		Where("kind = ?", "deal").
		Order("published_at DESC").
		Order("created_at DESC").
		Limit(cap).
		Find(&out.Deal).
		Error
	if err != nil {
		return repository.TopItems{}, repository.ErrInternal
	}

	err = r.db.
		Where("kind = ?", "shop").
		Order("published_at DESC").
		Order("created_at DESC").
		Limit(cap).
		Find(&out.Shop).
		Error
	if err != nil {
		return repository.TopItems{}, repository.ErrInternal
	}

	return out, nil
}
