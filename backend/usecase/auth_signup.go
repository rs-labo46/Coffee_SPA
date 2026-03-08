package usecase

import (
	"coffee-spa/entity"
	"time"
)

// user作成 + verify token発行 + メール送信
func (u *AuthUC) Signup(in SignupIn) (entity.User, error) {
	if err := u.val.Signup(in.Email, in.Pw); err != nil {
		return entity.User{}, ErrInvalidRequest
	}

	email := normEmail(in.Email)

	passHash, err := u.ph.Hash(in.Pw)
	if err != nil {
		return entity.User{}, ErrInternal
	}

	user, err := u.user.Create(entity.User{
		Email:         email,
		PassHash:      passHash,
		Role:          string(entity.RoleUser),
		TokenVer:      1,
		EmailVerified: false,
	})
	if err != nil {
		return entity.User{}, mapRepoErr(err)
	}

	raw, err := u.tk.NewOpaque()
	if err != nil {
		return entity.User{}, ErrInternal
	}

	if err := u.ev.RevokeUnusedByUser(user.ID); err != nil {
		return entity.User{}, mapRepoErr(err)
	}

	err = u.ev.Create(entity.EmailVerify{
		UserID:    user.ID,
		TokenHash: sha256Hex(raw),
		ExpiresAt: time.Now().Add(5 * time.Minute),
	})
	if err != nil {
		return entity.User{}, mapRepoErr(err)
	}

	if err := u.mail.SendVerify(user.Email, raw); err != nil {
		_ = u.writeAudit(
			"auth.verify_email.mail_failed",
			toI64Ptr(user.ID),
			in.IP,
			in.UA,
			verifyMeta{UserID: user.ID},
		)
	}

	if err := u.writeAudit(
		"auth.signup",
		toI64Ptr(user.ID),
		in.IP,
		in.UA,
		signupMeta{Email: maskEmail(user.Email)},
	); err != nil {
		return entity.User{}, err
	}

	return user, nil
}
