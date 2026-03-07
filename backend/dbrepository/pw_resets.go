package dbrepository

import (
	"errors"
	"time"

	"coffee-spa/entity"
	"coffee-spa/repository"

	"gorm.io/gorm"
)

type PwRepository struct {
	db *gorm.DB
}

func NewPwRepository(db *gorm.DB) repository.PwRepository {
	return &PwRepository{db: db}
}

func (r *PwRepository) Create(pw entity.PwReset) error {
	err := r.db.Create(&pw).Error
	if err != nil {
		if isDup(err) || isFK(err) {
			return repository.ErrConflict
		}
		return repository.ErrInternal
	}

	return nil
}

func (r *PwRepository) GetByTokenHash(hash string) (entity.PwReset, error) {
	var pw entity.PwReset

	err := r.db.
		Where("token_hash = ?", hash).
		First(&pw).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.PwReset{}, repository.ErrNotFound
		}
		return entity.PwReset{}, repository.ErrInternal
	}

	return pw, nil
}

func (r *PwRepository) Use(id int64) error {
	now := time.Now()

	res := r.db.
		Model(&entity.PwReset{}).
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

func (r *PwRepository) RevokeUnusedByUser(userID int64) error {
	now := time.Now()

	res := r.db.
		Model(&entity.PwReset{}).
		Where("user_id = ? AND used_at IS NULL", userID).
		Update("used_at", now)

	if res.Error != nil {
		return repository.ErrInternal
	}

	return nil
}
