package usecase

import (
	"coffee-spa/entity"
	"coffee-spa/repository"
	"errors"
	"strings"
	"time"
)

// 認証系の入力検証
type AuthVal interface {
	Signup(email string, pw string) error
	Login(email string, pw string) error
	NewPw(pw string) error
}

// パスワードハッシュ
type PwHash interface {
	Hash(pw string) (string, error)
	Compare(hash string, pw string) error
}

// JWT/CSRF/ランダムtoken生成
type Tok interface {
	NewAccess(userID int64, role string, tokenVer int) (string, error)
	NewCSRF() (string, error)
	NewOpaque() (string, error)
	NewFamilyID() (string, error)
}

// メール送信
type Mailer interface {
	SendVerify(email string, token string) error
	SendReset(email string, token string) error
}

// レート制御
// signup=ip、login=email_hash、refresh=user_id を使う。
type RateLim interface {
	AllowSignup(ip string) (bool, int, error)
	AllowLogin(emailHash string) (bool, int, error)
	AllowRefresh(userID int64) (bool, int, error)
	AllowResendIP(ip string) (bool, int, error)
	AllowResendMail(emailHash string) (bool, int, error)
	AllowForgotIP(ip string) (bool, int, error)
	AllowForgotMail(emailHash string) (bool, int, error)
}

// 認証のusecase
type AuthUC struct {
	user  repository.UserRepository
	ev    repository.EvRepository
	pw    repository.PwRepository
	rt    repository.RtRepository
	audit repository.AuditRepository
	val   AuthVal
	ph    PwHash
	tk    Tok
	mail  Mailer
	rl    RateLim
}

type loginFailMeta struct {
	Reason string `json:"reason"`
}

type signupMeta struct {
	Email string `json:"email"`
}

type verifyMeta struct {
	UserID int64 `json:"user_id"`
}

type resendMeta struct {
	Email string `json:"email"`
}

type refreshOKMeta struct {
	UserID   int64  `json:"user_id"`
	FamilyID string `json:"family_id"`
}

type refreshReuseMeta struct {
	UserID   int64  `json:"user_id"`
	FamilyID string `json:"family_id"`
	RtID     int64  `json:"rt_id"`
}

type logoutMeta struct {
	UserID int64 `json:"user_id"`
}

type forgotMeta struct {
	Email string `json:"email"`
}

type resetMeta struct {
	UserID int64 `json:"user_id"`
}

// AuthUCを作る
func NewAuthUC(
	user repository.UserRepository,
	ev repository.EvRepository,
	pw repository.PwRepository,
	rt repository.RtRepository,
	audit repository.AuditRepository,
	val AuthVal,
	ph PwHash,
	tk Tok,
	mail Mailer,
	rl RateLim,
) AuthUsecase {
	return &AuthUC{
		user:  user,
		ev:    ev,
		pw:    pw,
		rt:    rt,
		audit: audit,
		val:   val,
		ph:    ph,
		tk:    tk,
		mail:  mail,
		rl:    rl,
	}
}

// verify tokenでemail_verifiedをtrueにする
func (u *AuthUC) VerifyEmail(in VerifyEmailIn) error {
	token := strings.TrimSpace(in.Token)
	if token == "" {
		return ErrInvalidRequest
	}

	ev, err := u.ev.GetByTokenHash(sha256Hex(token))
	if err != nil {
		return ErrUnauthorized
	}

	if ev.UsedAt != nil || time.Now().After(ev.ExpiresAt) {
		return ErrUnauthorized
	}

	if err := u.user.SetEmailVerified(ev.UserID); err != nil {
		return mapRepoErr(err)
	}

	if err := u.ev.Use(ev.ID); err != nil {
		return mapRepoErr(err)
	}

	if err := u.writeAudit(
		"auth.verify_email",
		toI64Ptr(ev.UserID),
		in.IP,
		in.UA,
		verifyMeta{UserID: ev.UserID},
	); err != nil {
		return err
	}

	return nil
}

// 未認証ユーザー向けにverify tokenを再送する
func (u *AuthUC) ResendVerify(in ResendVerifyIn) error {
	email := normEmail(in.Email)
	if err := checkEmailOnly(email); err != nil {
		return err
	}

	emailHash := sha256Hex(email)

	okIP, _, err := u.rl.AllowResendIP(in.IP)
	if err != nil {
		return ErrInternal
	}

	okMail, _, err := u.rl.AllowResendMail(emailHash)
	if err != nil {
		return ErrInternal
	}

	if !okIP || !okMail {
		_ = u.writeAudit(
			"auth.email.resend.rate_limited",
			nil,
			in.IP,
			in.UA,
			resendMeta{Email: maskEmail(email)},
		)
		return nil
	}

	user, err := u.user.GetByEmail(email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil
		}
		return ErrInternal
	}

	if user.EmailVerified {
		return nil
	}

	raw, err := u.tk.NewOpaque()
	if err != nil {
		return ErrInternal
	}

	if err := u.ev.RevokeUnusedByUser(user.ID); err != nil {
		return mapRepoErr(err)
	}

	err = u.ev.Create(entity.EmailVerify{
		UserID:    user.ID,
		TokenHash: sha256Hex(raw),
		ExpiresAt: time.Now().Add(5 * time.Minute),
	})
	if err != nil {
		return mapRepoErr(err)
	}

	if err := u.mail.SendVerify(user.Email, raw); err != nil {
		_ = u.writeAudit(
			"auth.email.resend.mail_failed",
			toI64Ptr(user.ID),
			in.IP,
			in.UA,
			verifyMeta{UserID: user.ID},
		)
	}

	if err := u.writeAudit(
		"auth.email.resend",
		toI64Ptr(user.ID),
		in.IP,
		in.UA,
		resendMeta{Email: maskEmail(user.Email)},
	); err != nil {
		return err
	}

	return nil
}

func checkEmailOnly(email string) error {
	if email == "" {
		return ErrInvalidRequest
	}
	if len(email) > 254 {
		return ErrInvalidRequest
	}
	if !strings.Contains(email, "@") {
		return ErrInvalidRequest
	}
	return nil
}
