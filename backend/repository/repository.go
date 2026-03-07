package repository

import (
	"errors"

	"coffee-spa/entity"
)

var ErrNotFound = errors.New("not found")
var ErrConflict = errors.New("conflict")
var ErrInternal = errors.New("internal")

type ItemQ struct {
	Q      string
	Kind   string
	Limit  int
	Offset int
}

type TopItems struct {
	News   []entity.Item
	Recipe []entity.Item
	Deal   []entity.Item
	Shop   []entity.Item
}

type UserRepository interface {
	Create(u entity.User) (entity.User, error)
	GetByEmail(email string) (entity.User, error)
	GetByID(id uint64) (entity.User, error)
	SetEmailVerified(userID uint64) error
	UpdatePassHash(userID uint64, newHash string) error
	BumpTokenVer(userID uint64) (int, error)
}

type EvRepository interface {
	Create(ev entity.EmailVerify) error
	GetByTokenHash(hash string) (entity.EmailVerify, error)
	Use(id uint64) error
}

type PwRepository interface {
	Create(pw entity.PwReset) error
	GetByTokenHash(hash string) (entity.PwReset, error)
	Use(id uint64) error
}

type RtRepository interface {
	Create(rt entity.RefreshToken) (entity.RefreshToken, error)
	GetByTokenHash(hash string) (entity.RefreshToken, error)
	Revoke(id uint64) error
	MarkUsed(id uint64) error
	SetReplacedBy(id uint64, newID uint64) error
	RevokeByFamilyID(familyID string) error
	RevokeAllByUser(userID uint64) error
}

type SourceRepository interface {
	Create(s entity.Source) (entity.Source, error)
	GetByID(id uint64) (entity.Source, error)
	GetByName(name string) (entity.Source, error)
	List() ([]entity.Source, error)
}

type ItemRepository interface {
	Create(i entity.Item) (entity.Item, error)
	GetByID(id uint64) (entity.Item, error)
	List(q ItemQ) ([]entity.Item, error)
	Top(cap int) (TopItems, error)
}

type AuditRepository interface {
	Create(a entity.AuditLog) error
}
