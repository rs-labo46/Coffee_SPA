package dbrepository

import (
	"coffee-spa/entity"
	"coffee-spa/repository"

	"gorm.io/gorm"
)

type AuditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) repository.AuditRepository {
	return &AuditRepository{db: db}
}

// 監査ログを新規作成する
func (r *AuditRepository) Create(a entity.AuditLog) error {
	err := r.db.Create(&a).Error
	if err != nil {
		if isFK(err) {
			return repository.ErrConflict
		}

		return repository.ErrInternal
	}

	return nil
}
