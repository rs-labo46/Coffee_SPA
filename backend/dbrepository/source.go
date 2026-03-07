package dbrepository

import (
	"errors"

	"coffee-spa/entity"
	"coffee-spa/repository"

	"gorm.io/gorm"
)

type SourceRepository struct {
	db *gorm.DB
}

func NewSourceRepository(db *gorm.DB) repository.SourceRepository {
	return &SourceRepository{db: db}
}

func (r *SourceRepository) Create(s entity.Source) (entity.Source, error) {
	err := r.db.Create(&s).Error
	if err != nil {
		if isDup(err) {
			return entity.Source{}, repository.ErrConflict
		}
		return entity.Source{}, repository.ErrInternal
	}

	return s, nil
}

func (r *SourceRepository) GetByID(id int64) (entity.Source, error) {
	var s entity.Source

	err := r.db.
		First(&s, id).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Source{}, repository.ErrNotFound
		}
		return entity.Source{}, repository.ErrInternal
	}

	return s, nil
}

func (r *SourceRepository) GetByName(name string) (entity.Source, error) {
	var s entity.Source

	err := r.db.
		Where("name = ?", name).
		First(&s).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Source{}, repository.ErrNotFound
		}
		return entity.Source{}, repository.ErrInternal
	}

	return s, nil
}

func (r *SourceRepository) List() ([]entity.Source, error) {
	var xs []entity.Source

	err := r.db.
		Order("id ASC").
		Find(&xs).
		Error
	if err != nil {
		return nil, repository.ErrInternal
	}

	return xs, nil
}
