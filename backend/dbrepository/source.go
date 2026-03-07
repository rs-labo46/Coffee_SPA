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

// SourceRepositoryを作る
func NewSourceRepository(db *gorm.DB) repository.SourceRepository {
	return &SourceRepository{db: db}
}

// sourceを新規作成する
func (r *SourceRepository) Create(s entity.Source) (entity.Source, error) {
	//sourceテーブルへINSERTする
	err := r.db.Create(&s).Error
	if err != nil {
		if isDup(err) {
			return entity.Source{}, repository.ErrConflict
		}

		return entity.Source{}, repository.ErrInternal
	}

	return s, nil
}

// idでsourceを1件取得する
func (r *SourceRepository) GetByID(id uint64) (entity.Source, error) {
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

// nameでsourceを1件取得する
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

// source一覧を返す
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
