package entity

import (
	"time"

	"gorm.io/datatypes"
)

// user.roleの値
type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// item.kindの値
type ItemKind string

const (
	KindNews   ItemKind = "news"   // ニュース
	KindRecipe ItemKind = "recipe" // レシピ
	KindDeal   ItemKind = "deal"   // セール
	KindShop   ItemKind = "shop"   // 店舗ショップ
)

// Userはusersテーブル
type User struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Email         string    `gorm:"column:email;type:varchar;not null;uniqueIndex"`
	PassHash      string    `gorm:"column:pass_hash;type:varchar;not null"`
	Role          string    `gorm:"column:role;type:varchar;not null"`
	TokenVer      int       `gorm:"column:token_ver;not null;default:1"`
	EmailVerified bool      `gorm:"column:email_verified;not null;default:false"`
	CreatedAt     time.Time `gorm:"column:created_at;not null;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;not null;autoUpdateTime"`
}

// EmailVerifyはemail_verifiesテーブル
type EmailVerify struct {
	ID        int64      `gorm:"column:id;primaryKey;autoIncrement"`
	UserID    int64      `gorm:"column:user_id;not null;index"`
	TokenHash string     `gorm:"column:token_hash;type:varchar;not null;uniqueIndex"`
	ExpiresAt time.Time  `gorm:"column:expires_at;not null"`
	UsedAt    *time.Time `gorm:"column:used_at"`
}

// PwResetはpw_resetsテーブル
type PwReset struct {
	ID        int64      `gorm:"column:id;primaryKey;autoIncrement"`
	UserID    int64      `gorm:"column:user_id;not null;index"`
	TokenHash string     `gorm:"column:token_hash;type:varchar;not null;uniqueIndex"`
	ExpiresAt time.Time  `gorm:"column:expires_at;not null"`
	UsedAt    *time.Time `gorm:"column:used_at"`
	CreatedAt time.Time  `gorm:"column:created_at;not null;autoCreateTime"`
}

// RefreshTokenはrefresh_tokensテーブル
type RefreshToken struct {
	ID           int64      `gorm:"column:id;primaryKey;autoIncrement"`
	UserID       int64      `gorm:"column:user_id;not null;index"`
	FamilyID     string     `gorm:"column:family_id;type:varchar;not null;index"`
	TokenHash    string     `gorm:"column:token_hash;type:varchar;not null;uniqueIndex"`
	ExpiresAt    time.Time  `gorm:"column:expires_at;not null"`
	RevokedAt    *time.Time `gorm:"column:revoked_at"`
	UsedAt       *time.Time `gorm:"column:used_at"`
	ReplacedByID *int64     `gorm:"column:replaced_by_id"`
	CreatedAt    time.Time  `gorm:"column:created_at;not null;autoCreateTime"`
}

// Sourceはsourcesテーブル
type Source struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string    `gorm:"column:name;type:varchar;not null;uniqueIndex"`
	SiteURL   *string   `gorm:"column:site_url;type:varchar"`
	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime"`
}

// Itemはitemsテーブル
type Item struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Title       string    `gorm:"column:title;type:varchar;not null"`
	Summary     *string   `gorm:"column:summary;type:varchar"`
	URL         *string   `gorm:"column:url;type:varchar"`
	ImageURL    *string   `gorm:"column:image_url;type:varchar"`
	Kind        string    `gorm:"column:kind;type:varchar;not null;index"`
	SourceID    int64     `gorm:"column:source_id;not null;index"`
	PublishedAt time.Time `gorm:"column:published_at;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;autoCreateTime"`
}

// AuditLogはaudit_logsテーブル
type AuditLog struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement"`
	Type      string         `gorm:"column:type;type:varchar;not null;index"`
	UserID    *int64         `gorm:"column:user_id;index"`
	IP        string         `gorm:"column:ip;type:varchar;not null"`
	UA        string         `gorm:"column:ua;type:varchar;not null"`
	MetaJSON  datatypes.JSON `gorm:"column:meta_json;type:jsonb;not null;default:'{}'"`
	CreatedAt time.Time      `gorm:"column:created_at;not null;autoCreateTime;index"`
}
