package dbrepository

import (
	"errors"

	"coffee-spa/entity"
	"coffee-spa/repository"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(u entity.User) (entity.User, error) {
	err := r.db.Create(&u).Error
	if err != nil {
		if isDup(err) {
			return entity.User{}, repository.ErrConflict
		}
		return entity.User{}, repository.ErrInternal
	}

	return u, nil
}

func (r *UserRepository) GetByEmail(email string) (entity.User, error) {
	var u entity.User

	err := r.db.
		Where("email = ?", email).
		First(&u).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, repository.ErrNotFound
		}
		return entity.User{}, repository.ErrInternal
	}

	return u, nil
}

func (r *UserRepository) GetByID(id int64) (entity.User, error) {
	var u entity.User

	err := r.db.
		First(&u, id).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, repository.ErrNotFound
		}
		return entity.User{}, repository.ErrInternal
	}

	return u, nil
}

func (r *UserRepository) SetEmailVerified(userID int64) error {
	res := r.db.
		Model(&entity.User{}).
		Where("id = ?", userID).
		Update("email_verified", true)

	if res.Error != nil {
		return repository.ErrInternal
	}
	if res.RowsAffected == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (r *UserRepository) UpdatePassHash(userID int64, newHash string) error {
	res := r.db.
		Model(&entity.User{}).
		Where("id = ?", userID).
		Update("pass_hash", newHash)

	if res.Error != nil {
		return repository.ErrInternal
	}
	if res.RowsAffected == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (r *UserRepository) BumpTokenVer(userID int64) (int, error) {
	res := r.db.
		Model(&entity.User{}).
		Where("id = ?", userID).
		Update("token_ver", gorm.Expr("token_ver + 1"))

	if res.Error != nil {
		return 0, repository.ErrInternal
	}
	if res.RowsAffected == 0 {
		return 0, repository.ErrNotFound
	}

	var u entity.User
	err := r.db.
		Select("token_ver").
		First(&u, userID).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, repository.ErrNotFound
		}
		return 0, repository.ErrInternal
	}

	return u.TokenVer, nil
}
