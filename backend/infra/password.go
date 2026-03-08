package infra

import "golang.org/x/crypto/bcrypt"

// bcryptを使うパスワードハッシュ実装
type BcryptHasher struct{}

// BcryptHasherを作る
func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

// 平文パスワードをbcrypt hashにする
func (h *BcryptHasher) Hash(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// hashと平文パスワードを照合する
func (h *BcryptHasher) Compare(hash string, pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
}
