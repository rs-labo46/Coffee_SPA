package dbrepository

import (
	"coffee-spa/entity"
	"coffee-spa/repository"
	"errors"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepository{db}
}

// ユーザーの新規作成
func (r *UserRepository) Create(u entity.User) (entity.User, error) {
	//users テーブルへINSERTする
	err := r.db.Create(&u).Error
	if err != nil {
		if isDup(err) {
			return entity.User{}, repository.ErrConflict
		}
		return entity.User{}, repository.ErrInternal
	}
	return u, nil
}

// emailでユーザーを一件取得する
func (r *UserRepository) GetByEmail(email string) (entity.User, error) {
	var u entity.User

	//emailで完全一致で先頭の一件を取得
	err := r.db.Where("email = ?", email).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, repository.ErrNotFound
		}
		return entity.User{}, repository.ErrInternal
	}
	return u, nil
}

// idでユーザーを一件取得する
func (r *UserRepository) GetByID(id uint64) (entity.User, error) {
	var u entity.User
	//idで検索する
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

// email_verifiedをtrueに更新する
func (r *UserRepository) SetEmailVerified(userID uint64) error {
	//対象のユーザー1件のemail_verifiedをtrue にする
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

// パスワードハッシュを更新する
func (r *UserRepository) UpdatePassHash(userID uint64, newHash string) error {
	// 対象のユーザー1件のpass_hashを新しい値にする
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

// token_verを1増やして、その新しい値を返す
func (r *UserRepository) BumpTokenVer(userID uint64) (int, error) {
	//token_ver=token_ver+1をDB側で実行する
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

	//更新後のtoken_verを取り直す
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
