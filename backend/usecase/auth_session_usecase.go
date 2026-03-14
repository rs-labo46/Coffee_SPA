package usecase

import (
	"coffee-spa/entity"
	"strings"
	"time"
)

// Login は access / refresh / csrf を発行する
func (u *AuthUC) Login(LoginInput LoginIn) (AuthOut, error) {
	if err := u.val.Login(LoginInput.Email, LoginInput.Pw); err != nil {
		return AuthOut{}, ErrInvalidRequest
	}

	email := normEmail(LoginInput.Email)
	emailHash := sha256Hex(email)

	ok, retry, err := u.rl.AllowLogin(emailHash)
	if err != nil {
		return AuthOut{}, ErrInternal
	}
	if !ok {
		return AuthOut{}, ErrRateLimited{RetryAfterSec: retry}
	}

	user, err := u.user.GetByEmail(email)
	if err != nil {
		_ = u.writeAudit(
			"auth.login.fail",
			nil,
			LoginInput.IP,
			LoginInput.UA,
			loginFailMeta{Reason: "unauthorized"},
		)
		return AuthOut{}, ErrUnauthorized
	}

	if !user.EmailVerified {
		_ = u.writeAudit(
			"auth.login.fail",
			int64Pointer(user.ID),
			LoginInput.IP,
			LoginInput.UA,
			loginFailMeta{Reason: "unauthorized"},
		)
		return AuthOut{}, ErrUnauthorized
	}

	if err := u.ph.Compare(user.PassHash, LoginInput.Pw); err != nil {
		_ = u.writeAudit(
			"auth.login.fail",
			int64Pointer(user.ID),
			LoginInput.IP,
			LoginInput.UA,
			loginFailMeta{Reason: "unauthorized"},
		)
		return AuthOut{}, ErrUnauthorized
	}

	familyID, err := u.tk.NewFamilyID()
	if err != nil {
		return AuthOut{}, ErrInternal
	}

	rawRefresh, err := u.tk.NewOpaque()
	if err != nil {
		return AuthOut{}, ErrInternal
	}

	_, err = u.rt.Create(entity.RefreshToken{
		UserID:    user.ID,
		FamilyID:  familyID,
		TokenHash: sha256Hex(rawRefresh),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		return AuthOut{}, mapRepoErr(err)
	}

	access, err := u.tk.NewAccess(user.ID, user.Role, user.TokenVer)
	if err != nil {
		return AuthOut{}, ErrInternal
	}

	csrf, err := u.tk.NewCSRF()
	if err != nil {
		return AuthOut{}, ErrInternal
	}

	if err := u.writeAudit(
		"auth.login.success",
		int64Pointer(user.ID),
		LoginInput.IP,
		LoginInput.UA,
		refreshOKMeta{
			UserID:   user.ID,
			FamilyID: familyID,
		},
	); err != nil {
		return AuthOut{}, err
	}

	return AuthOut{
		AccessToken:  access,
		RefreshToken: rawRefresh,
		CsrfToken:    csrf,
		User:         user,
	}, nil
}

// refresh token を回転させる
func (u *AuthUC) Refresh(RefreshInput RefreshIn) (AuthOut, error) {
	if strings.TrimSpace(RefreshInput.RefreshToken) == "" {
		_ = u.writeAudit("auth.refresh.fail", nil, RefreshInput.IP, RefreshInput.UA, nil)
		return AuthOut{}, ErrUnauthorized
	}

	rt, err := u.rt.GetByTokenHash(sha256Hex(RefreshInput.RefreshToken))
	if err != nil {
		_ = u.writeAudit("auth.refresh.fail", nil, RefreshInput.IP, RefreshInput.UA, nil)
		return AuthOut{}, ErrUnauthorized
	}

	if time.Now().After(rt.ExpiresAt) || rt.RevokedAt != nil {
		_ = u.writeAudit("auth.refresh.fail", int64Pointer(rt.UserID), RefreshInput.IP, RefreshInput.UA, nil)
		return AuthOut{}, ErrUnauthorized
	}

	if rt.UsedAt != nil {
		_ = u.rt.RevokeByFamilyID(rt.FamilyID)
		_, _ = u.user.BumpTokenVer(rt.UserID)

		_ = u.writeAudit(
			"auth.refresh.reuse_detected",
			int64Pointer(rt.UserID),
			RefreshInput.IP,
			RefreshInput.UA,
			refreshReuseMeta{
				UserID:   rt.UserID,
				FamilyID: rt.FamilyID,
				RtID:     rt.ID,
			},
		)

		return AuthOut{}, ErrUnauthorized
	}

	user, err := u.user.GetByID(rt.UserID)
	if err != nil {
		return AuthOut{}, mapRepoErr(err)
	}

	rawRefresh, err := u.tk.NewOpaque()
	if err != nil {
		return AuthOut{}, ErrInternal
	}

	newRT, err := u.rt.Create(entity.RefreshToken{
		UserID:    rt.UserID,
		FamilyID:  rt.FamilyID,
		TokenHash: sha256Hex(rawRefresh),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		return AuthOut{}, mapRepoErr(err)
	}

	if err := u.rt.MarkUsed(rt.ID); err != nil {
		_ = u.rt.RevokeByFamilyID(rt.FamilyID)
		_, _ = u.user.BumpTokenVer(rt.UserID)

		_ = u.writeAudit(
			"auth.refresh.reuse_detected",
			int64Pointer(rt.UserID),
			RefreshInput.IP,
			RefreshInput.UA,
			refreshReuseMeta{
				UserID:   rt.UserID,
				FamilyID: rt.FamilyID,
				RtID:     rt.ID,
			},
		)

		return AuthOut{}, ErrUnauthorized
	}

	if err := u.rt.Revoke(rt.ID); err != nil {
		return AuthOut{}, mapRepoErr(err)
	}

	if err := u.rt.SetReplacedBy(rt.ID, newRT.ID); err != nil {
		return AuthOut{}, mapRepoErr(err)
	}

	access, err := u.tk.NewAccess(user.ID, user.Role, user.TokenVer)
	if err != nil {
		return AuthOut{}, ErrInternal
	}

	csrf, err := u.tk.NewCSRF()
	if err != nil {
		return AuthOut{}, ErrInternal
	}

	if err := u.writeAudit(
		"auth.refresh.success",
		int64Pointer(user.ID),
		RefreshInput.IP,
		RefreshInput.UA,
		refreshOKMeta{
			UserID:   user.ID,
			FamilyID: rt.FamilyID,
		},
	); err != nil {
		return AuthOut{}, err
	}

	return AuthOut{
		AccessToken:  access,
		RefreshToken: rawRefresh,
		CsrfToken:    csrf,
		User:         user,
	}, nil
}

// Logoutはtoken_verを上げて、refresh familyを失効する
func (u *AuthUC) Logout(in LogoutIn) error {
	_, err := u.user.BumpTokenVer(in.UserID)
	if err != nil {
		return mapRepoErr(err)
	}

	if in.RefreshToken != "" {
		rt, err := u.rt.GetByTokenHash(sha256Hex(in.RefreshToken))
		if err == nil {
			if err := u.rt.RevokeByFamilyID(rt.FamilyID); err != nil {
				return mapRepoErr(err)
			}
		}
	}

	if err := u.writeAudit(
		"auth.logout",
		int64Pointer(in.UserID),
		in.IP,
		in.UA,
		logoutMeta{UserID: in.UserID},
	); err != nil {
		return err
	}

	return nil
}
