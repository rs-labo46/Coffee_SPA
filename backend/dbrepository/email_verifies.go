package dbrepository

import (
	"errors"
	"time"

	"coffee-spa/entity"
	"coffee-spa/repository"

	"gorm.io/gorm"
)

type EvRepository struct {
	db *gorm.DB
}

func NewEvRepository(db *gorm.DB) repository.EvRepository {
	return &EvRepository{db: db}
}

func (r *EvRepository) Create(ev entity.EmailVerify) error {
	err := r.db.Create(&ev).Error
	if err != nil {
		if isDup(err) || isFK(err) {
			return repository.ErrConflict
		}
		return repository.ErrInternal
	}

	return nil
}

func (r *EvRepository) GetByTokenHash(hash string) (entity.EmailVerify, error) {
	var ev entity.EmailVerify

	err := r.db.
		Where("token_hash = ?", hash).
		First(&ev).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.EmailVerify{}, repository.ErrNotFound
		}
		return entity.EmailVerify{}, repository.ErrInternal
	}

	return ev, nil
}

func (r *EvRepository) Use(id int64) error {
	now := time.Now()

	res := r.db.
		Model(&entity.EmailVerify{}).
		Where("id = ? AND used_at IS NULL", id).
		Update("used_at", now)

	if res.Error != nil {
		return repository.ErrInternal
	}
	if res.RowsAffected == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (r *EvRepository) RevokeUnusedByUser(userID int64) error {
	now := time.Now()

	res := r.db.
		Model(&entity.EmailVerify{}).
		Where("user_id = ? AND used_at IS NULL", userID).
		Update("used_at", now)

	if res.Error != nil {
		return repository.ErrInternal
	}

	return nil
}
