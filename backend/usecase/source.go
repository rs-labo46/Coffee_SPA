package usecase

import (
	"encoding/json"

	"coffee-spa/entity"
	"coffee-spa/repository"

	"gorm.io/datatypes"
)

type SourceVal interface {
	NewSource(in AddSourceIn) error
}

type SourceUC struct {
	source repository.SourceRepository
	audit  repository.AuditRepository
	val    SourceVal
}

type sourceCreateMeta struct {
	UserID   int64 `json:"user_id"`
	SourceID int64 `json:"source_id"`
}

func NewSourceUC(
	source repository.SourceRepository,
	audit repository.AuditRepository,
	val SourceVal,
) SourceUsecase {
	return &SourceUC{
		source: source,
		audit:  audit,
		val:    val,
	}
}

func (u *SourceUC) Add(actor Actor, in AddSourceIn) (entity.Source, error) {
	if err := u.val.NewSource(in); err != nil {
		return entity.Source{}, ErrInvalidRequest
	}

	src, err := u.source.Create(entity.Source{
		Name:    in.Name,
		SiteURL: in.SiteURL,
	})
	if err != nil {
		return entity.Source{}, mapRepoErr(err)
	}

	b, err := json.Marshal(sourceCreateMeta{
		UserID:   actor.UserID,
		SourceID: src.ID,
	})
	if err != nil {
		return entity.Source{}, ErrInternal
	}

	auditErr := u.audit.Create(entity.AuditLog{
		Type:     "admin.sources.create",
		UserID:   toI64Ptr(actor.UserID),
		IP:       actor.IP,
		UA:       actor.UA,
		MetaJSON: datatypes.JSON(b),
	})
	if auditErr != nil {
		return entity.Source{}, mapRepoErr(auditErr)
	}

	return src, nil
}

func (u *SourceUC) List() ([]entity.Source, error) {
	xs, err := u.source.List()
	if err != nil {
		return nil, mapRepoErr(err)
	}

	return xs, nil
}
