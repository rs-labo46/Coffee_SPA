package repository

import (
	"errors"

	"coffee-spa/entity"

	"gorm.io/gorm"
)

type sourceRepository struct {
	db *gorm.DB
}

func NewSourceRepository(db *gorm.DB) SourceRepository {
	return &sourceRepository{db}
}

func (r *sourceRepository) Create(s entity.Source) (entity.Source, error) {
	err := r.db.Create(&s).Error
	if err != nil {
		if isDup(err) {
			return entity.Source{}, ErrConflict
		}
		return entity.Source{}, ErrInternal
	}

	return s, nil
}

func (r *sourceRepository) GetByID(id int64) (entity.Source, error) {
	var s entity.Source

	err := r.db.
		First(&s, id).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Source{}, ErrNotFound
		}
		return entity.Source{}, ErrInternal
	}

	return s, nil
}

func (r *sourceRepository) GetByName(name string) (entity.Source, error) {
	var s entity.Source

	err := r.db.
		Where("name = ?", name).
		First(&s).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Source{}, ErrNotFound
		}
		return entity.Source{}, ErrInternal
	}

	return s, nil
}

func (r *sourceRepository) List() ([]entity.Source, error) {
	var xs []entity.Source

	err := r.db.
		Order("id ASC").
		Find(&xs).
		Error
	if err != nil {
		return nil, ErrInternal
	}

	return xs, nil
}
