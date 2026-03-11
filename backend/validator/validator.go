package validator

import (
	"strings"

	"coffee-spa/policy"
	"coffee-spa/usecase"
)

// 認証のvalidator
type AuthVal interface {
	Signup(email string, pw string) error
	Login(email string, pw string) error
	NewPw(pw string) error
}

// itemのvalidator
type ItemVal interface {
	NewItem(input usecase.AddItemIn) error
	ListItem(q usecase.ItemQ) error
}

// sourceのvalidator
type SourceVal interface {
	NewSource(input usecase.AddSourceIn) error
}

type AuthValidator struct {
	email policy.EmailPol
	pw    policy.PwPol
}

type ItemValidator struct {
	kind policy.KindPol
	url  policy.URLPol
	page policy.PagePol
}

type SourceValidator struct {
	url policy.URLPol
}

func NewAuthValidator(
	email policy.EmailPol,
	pw policy.PwPol,
) *AuthValidator {
	return &AuthValidator{
		email: email,
		pw:    pw,
	}
}

func NewItemValidator(
	kind policy.KindPol,
	url policy.URLPol,
	page policy.PagePol,
) *ItemValidator {
	return &ItemValidator{
		kind: kind,
		url:  url,
		page: page,
	}
}

func NewSourceValidator(
	url policy.URLPol,
) *SourceValidator {
	return &SourceValidator{
		url: url,
	}
}

func (v *AuthValidator) Signup(email string, pw string) error {
	if err := v.email.Ok(email); err != nil {
		return err
	}
	if err := v.pw.Ok(pw); err != nil {
		return err
	}
	return nil
}

func (v *AuthValidator) Login(email string, pw string) error {
	return v.Signup(email, pw)
}

func (v *AuthValidator) NewPw(pw string) error {
	return v.pw.Ok(pw)
}

func (v *ItemValidator) NewItem(itemInput usecase.AddItemIn) error {
	title := strings.TrimSpace(itemInput.Title)
	if title == "" || len(title) > 120 {
		return usecase.ErrInvalidRequest
	}

	if itemInput.Summary != nil && len(*itemInput.Summary) > 500 {
		return usecase.ErrInvalidRequest
	}

	if itemInput.URL != nil {
		if err := v.url.Ok(*itemInput.URL); err != nil {
			return err
		}
	}

	if itemInput.ImageURL != nil {
		if err := v.url.Ok(*itemInput.ImageURL); err != nil {
			return err
		}
	}

	if err := v.kind.Ok(itemInput.Kind); err != nil {
		return err
	}

	if itemInput.SourceID <= 0 {
		return usecase.ErrInvalidRequest
	}

	return nil
}

func (v *ItemValidator) ListItem(q usecase.ItemQ) error {
	if err := v.page.Ok(q.Limit, q.Offset); err != nil {
		return err
	}

	if q.Kind == "" {
		return nil
	}
	return v.kind.Ok(q.Kind)
}

func (v *SourceValidator) NewSource(sourceInput usecase.AddSourceIn) error {
	name := strings.TrimSpace(sourceInput.Name)
	if name == "" || len(name) > 80 {
		return usecase.ErrInvalidRequest
	}

	if sourceInput.SiteURL == nil {
		return nil
	}

	return v.url.Ok(*sourceInput.SiteURL)
}
