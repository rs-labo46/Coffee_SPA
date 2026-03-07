package dbrepository

import (
	"errors"
	"time"

	"coffee-spa/entity"
	"coffee-spa/repository"

	"gorm.io/gorm"
)

type RtRepository struct {
	db *gorm.DB
}

func NewRtRepository(db *gorm.DB) repository.RtRepository {
	return &RtRepository{db: db}
}

func (r *RtRepository) Create(rt entity.RefreshToken) (entity.RefreshToken, error) {
	err := r.db.Create(&rt).Error
	if err != nil {
		if isDup(err) || isFK(err) {
			return entity.RefreshToken{}, repository.ErrConflict
		}
		return entity.RefreshToken{}, repository.ErrInternal
	}

	return rt, nil
}

func (r *RtRepository) GetByTokenHash(hash string) (entity.RefreshToken, error) {
	var rt entity.RefreshToken

	err := r.db.
		Where("token_hash = ?", hash).
		First(&rt).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.RefreshToken{}, repository.ErrNotFound
		}
		return entity.RefreshToken{}, repository.ErrInternal
	}

	return rt, nil
}

func (r *RtRepository) Revoke(id int64) error {
	now := time.Now()

	res := r.db.
		Model(&entity.RefreshToken{}).
		Where("id = ? AND revoked_at IS NULL", id).
		Update("revoked_at", now)

	if res.Error != nil {
		return repository.ErrInternal
	}
	if res.RowsAffected == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (r *RtRepository) MarkUsed(id int64) error {
	now := time.Now()

	res := r.db.
		Model(&entity.RefreshToken{}).
		Where("id = ? AND used_at IS NULL", id).
		Update("used_at", now)

	if res.Error != nil {
		return repository.ErrInternal
	}
	if res.RowsAffected == 0 {
		return repository.ErrConflict
	}

	return nil
}

func (r *RtRepository) SetReplacedBy(id int64, newID int64) error {
	res := r.db.
		Model(&entity.RefreshToken{}).
		Where("id = ?", id).
		Update("replaced_by_id", newID)

	if res.Error != nil {
		if isFK(res.Error) {
			return repository.ErrConflict
		}
		return repository.ErrInternal
	}
	if res.RowsAffected == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (r *RtRepository) RevokeByFamilyID(familyID string) error {
	now := time.Now()

	res := r.db.
		Model(&entity.RefreshToken{}).
		Where("family_id = ? AND revoked_at IS NULL", familyID).
		Update("revoked_at", now)

	if res.Error != nil {
		return repository.ErrInternal
	}

	return nil
}

func (r *RtRepository) RevokeAllByUser(userID int64) error {
	now := time.Now()

	res := r.db.
		Model(&entity.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now)

	if res.Error != nil {
		return repository.ErrInternal
	}

	return nil
}
