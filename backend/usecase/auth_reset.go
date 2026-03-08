package usecase

import (
	"coffee-spa/entity"
	"coffee-spa/repository"
	"errors"
	"strings"
	"time"
)

// reset token を発行してメール送信する
func (u *AuthUC) ForgotPw(in ForgotPwIn) error {
	email := normEmail(in.Email)
	if err := checkEmailOnly(email); err != nil {
		return err
	}

	emailHash := sha256Hex(email)

	ok, _, err := u.rl.AllowForgot(in.IP, emailHash)
	if err != nil {
		return ErrInternal
	}
	if !ok {
		_ = u.writeAudit(
			"auth.password.forgot.rate_limited",
			nil,
			in.IP,
			in.UA,
			forgotMeta{Email: maskEmail(email)},
		)
		return nil
	}

	user, err := u.user.GetByEmail(email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			_ = u.writeAudit(
				"auth.password.forgot",
				nil,
				in.IP,
				in.UA,
				forgotMeta{Email: maskEmail(email)},
			)
			return nil
		}
		return ErrInternal
	}

	raw, err := u.tk.NewOpaque()
	if err != nil {
		return ErrInternal
	}

	if err := u.pw.RevokeUnusedByUser(user.ID); err != nil {
		return mapRepoErr(err)
	}

	err = u.pw.Create(entity.PwReset{
		UserID:    user.ID,
		TokenHash: sha256Hex(raw),
		ExpiresAt: time.Now().Add(5 * time.Minute),
	})
	if err != nil {
		return mapRepoErr(err)
	}

	//送信失敗でも成功扱い。
	if err := u.mail.SendReset(user.Email, raw); err != nil {
		_ = u.writeAudit(
			"auth.password.forgot.mail_failed",
			toI64Ptr(user.ID),
			in.IP,
			in.UA,
			verifyMeta{UserID: user.ID},
		)
	}

	if err := u.writeAudit(
		"auth.password.forgot",
		toI64Ptr(user.ID),
		in.IP,
		in.UA,
		forgotMeta{Email: maskEmail(user.Email)},
	); err != nil {
		return err
	}

	return nil
}

// password更新 + token消費 + 全refresh失効
func (u *AuthUC) ResetPw(in ResetPwIn) error {
	if strings.TrimSpace(in.Token) == "" {
		return ErrInvalidRequest
	}

	if err := u.val.NewPw(in.NewPw); err != nil {
		return ErrInvalidRequest
	}

	pw, err := u.pw.GetByTokenHash(sha256Hex(in.Token))
	if err != nil {
		return ErrUnauthorized
	}

	if pw.UsedAt != nil || time.Now().After(pw.ExpiresAt) {
		return ErrUnauthorized
	}

	newHash, err := u.ph.Hash(in.NewPw)
	if err != nil {
		return ErrInternal
	}

	if err := u.user.UpdatePassHash(pw.UserID, newHash); err != nil {
		return mapRepoErr(err)
	}

	if err := u.pw.Use(pw.ID); err != nil {
		return mapRepoErr(err)
	}

	if _, err := u.user.BumpTokenVer(pw.UserID); err != nil {
		return mapRepoErr(err)
	}

	if err := u.rt.RevokeAllByUser(pw.UserID); err != nil {
		return mapRepoErr(err)
	}

	if err := u.writeAudit(
		"auth.password.reset",
		toI64Ptr(pw.UserID),
		in.IP,
		in.UA,
		resetMeta{UserID: pw.UserID},
	); err != nil {
		return err
	}

	return nil
}

// Meはuserを返す
func (u *AuthUC) Me(userID int64) (entity.User, error) {
	user, err := u.user.GetByID(userID)
	if err != nil {
		return entity.User{}, mapRepoErr(err)
	}

	return user, nil
}
