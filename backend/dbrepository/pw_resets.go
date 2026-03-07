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

// PwRepositoryを作る
func NewPwRepository(db *gorm.DB) repository.PwRepository {
	return &PwRepository{db: db}
}

// passwordresettokenを新規作成する
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

// token_hashで1件取得する
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

// 未使用レコードだけused_atを埋める
func (r *PwRepository) Use(id uint64) error {
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
