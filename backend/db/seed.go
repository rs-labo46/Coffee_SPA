package db

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"coffee-spa/entity"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedDev(d DB, adminEmail string, adminPassword string) error {
	//まずadminを作る。
	if err := SeedDevAdmin(d, adminEmail, adminPassword); err != nil {
		return err
	}

	//sourceを作る。
	sourceMap, err := SeedDevSources(d)
	if err != nil {
		return err
	}

	//itemを作る。
	if err := SeedDevItems(d, sourceMap); err != nil {
		return err
	}

	return nil
}

// adminユーザーを作成する。
// すでに同じemailがあれば何もしない。
func SeedDevAdmin(d DB, email string, password string) error {
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)

	//空なら何もしない。
	if email == "" || password == "" {
		return nil
	}

	//既存確認。
	var exists entity.User
	err := d.G.
		Where("email = ?", email).
		First(&exists).
		Error

	if err == nil {
		//既に存在するから何もしない。
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("find seed admin: %w", err)
	}

	//パスワードをハッシュ化する。
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash seed admin password: %w", err)
	}

	u := entity.User{
		Email:         email,
		PassHash:      string(hash),
		Role:          string(entity.RoleAdmin),
		TokenVer:      1,
		EmailVerified: true,
	}

	if err := d.G.Create(&u).Error; err != nil {
		return fmt.Errorf("create seed admin: %w", err)
	}

	return nil
}

func SeedDevSources(d DB) (map[string]entity.Source, error) {
	defs := []entity.Source{
		{
			Name:    "Coffee Daily",
			SiteURL: strPtr("https://example.com/coffee-daily"),
		},
		{
			Name:    "Roastery Journal",
			SiteURL: strPtr("https://example.com/roastery-journal"),
		},
		{
			Name:    "Bean Market",
			SiteURL: strPtr("https://example.com/bean-market"),
		},
		{
			Name:    "Cafe Guide",
			SiteURL: strPtr("https://example.com/cafe-guide"),
		},
	}

	//結果を詰める。
	res := make(map[string]entity.Source, len(defs))

	for _, def := range defs {
		var src entity.Source

		err := d.G.
			Where("name = ?", def.Name).
			First(&src).
			Error

		if err == nil {
			//既存なら採用して次へ。
			res[src.Name] = src
			continue
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find seed source %s: %w", def.Name, err)
		}

		//なければ作成する。
		if err := d.G.Create(&def).Error; err != nil {
			return nil, fmt.Errorf("create seed source %s: %w", def.Name, err)
		}

		res[def.Name] = def
	}

	return res, nil
}

func SeedDevItems(d DB, sourceMap map[string]entity.Source) error {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)

	defs := []seedItemDef{
		//news
		{
			SourceName: "Coffee Daily",
			Title:      "今週のおすすめ豆 3選",
			Summary:    "春向けの軽やかな味を楽しめる豆を3つ紹介します。",
			URL:        "https://example.com/coffee-daily/recommend-beans",
			ImageURL:   "https://images.unsplash.com/photo-1495474472287-4d71bcdd2085",
			Kind:       string(entity.KindNews),
			PublishedAt: time.Date(
				2026, 3, 16, 9, 0, 0, 0, jst,
			),
		},
		{
			SourceName: "Coffee Daily",
			Title:      "深煎り豆の選び方",
			Summary:    "苦味だけでなく香りと後味で選ぶための簡単な基準を解説します。",
			URL:        "https://example.com/coffee-daily/dark-roast-guide",
			ImageURL:   "https://images.unsplash.com/photo-1461023058943-07fcbe16d735",
			Kind:       string(entity.KindNews),
			PublishedAt: time.Date(
				2026, 3, 15, 7, 45, 0, 0, jst,
			),
		},
		{
			SourceName: "Coffee Daily",
			Title:      "朝に合う浅煎りブレンド特集",
			Summary:    "軽い口当たりで飲みやすい朝向けブレンドをまとめました。",
			URL:        "https://example.com/coffee-daily/light-blend-morning",
			ImageURL:   "https://images.unsplash.com/photo-1517701604599-bb29b565090c",
			Kind:       string(entity.KindNews),
			PublishedAt: time.Date(
				2026, 3, 14, 8, 20, 0, 0, jst,
			),
		},
		{
			SourceName: "Coffee Daily",
			Title:      "初心者向けコーヒー器具まとめ",
			Summary:    "最初に揃えるなら何が必要かをわかりやすく整理しました。",
			URL:        "https://example.com/coffee-daily/beginner-tools",
			ImageURL:   "https://images.unsplash.com/photo-1447933601403-0c6688de566e",
			Kind:       string(entity.KindNews),
			PublishedAt: time.Date(
				2026, 3, 13, 10, 10, 0, 0, jst,
			),
		},

		//recipe
		{
			SourceName: "Roastery Journal",
			Title:      "ハンドドリップ基本レシピ",
			Summary:    "お湯の温度と蒸らし時間だけで味が安定する基本レシピです。",
			URL:        "https://example.com/roastery-journal/basic-drip-recipe",
			ImageURL:   "https://images.unsplash.com/photo-1447933601403-0c6688de566e",
			Kind:       string(entity.KindRecipe),
			PublishedAt: time.Date(
				2026, 3, 12, 8, 30, 0, 0, jst,
			),
		},
		{
			SourceName: "Roastery Journal",
			Title:      "ミルクブリューの作り方",
			Summary:    "牛乳に一晩浸して甘みを引き出す簡単レシピです。",
			URL:        "https://example.com/roastery-journal/milk-brew",
			ImageURL:   "https://images.unsplash.com/photo-1511920170033-f8396924c348",
			Kind:       string(entity.KindRecipe),
			PublishedAt: time.Date(
				2026, 3, 11, 18, 15, 0, 0, jst,
			),
		},
		{
			SourceName: "Roastery Journal",
			Title:      "アイスコーヒー急冷レシピ",
			Summary:    "氷で一気に冷やして香りを残す作り方です。",
			URL:        "https://example.com/roastery-journal/iced-flash-brew",
			ImageURL:   "https://images.unsplash.com/photo-1498804103079-a6351b050096",
			Kind:       string(entity.KindRecipe),
			PublishedAt: time.Date(
				2026, 3, 10, 13, 0, 0, 0, jst,
			),
		},
		{
			SourceName: "Roastery Journal",
			Title:      "フレンチプレス入門",
			Summary:    "道具が少なくてもコクを出しやすい抽出方法を紹介します。",
			URL:        "https://example.com/roastery-journal/french-press-guide",
			ImageURL:   "https://images.unsplash.com/photo-1509042239860-f550ce710b93",
			Kind:       string(entity.KindRecipe),
			PublishedAt: time.Date(
				2026, 3, 9, 9, 40, 0, 0, jst,
			),
		},

		//deal
		{
			SourceName: "Bean Market",
			Title:      "週末セール開催中",
			Summary:    "定番ブレンド豆が期間限定で10%オフになります。",
			URL:        "https://example.com/bean-market/weekend-sale",
			ImageURL:   "https://images.unsplash.com/photo-1509042239860-f550ce710b93",
			Kind:       string(entity.KindDeal),
			PublishedAt: time.Date(
				2026, 3, 8, 10, 0, 0, 0, jst,
			),
		},
		{
			SourceName: "Bean Market",
			Title:      "新生活応援セット割",
			Summary:    "ドリッパーと豆のセットをまとめ買いしやすくしました。",
			URL:        "https://example.com/bean-market/new-life-set",
			ImageURL:   "https://images.unsplash.com/photo-1459755486867-b55449bb39ff",
			Kind:       string(entity.KindDeal),
			PublishedAt: time.Date(
				2026, 3, 7, 12, 15, 0, 0, jst,
			),
		},
		{
			SourceName: "Bean Market",
			Title:      "送料無料キャンペーン",
			Summary:    "一定額以上の購入で送料が無料になる期間限定施策です。",
			URL:        "https://example.com/bean-market/free-shipping",
			ImageURL:   "https://images.unsplash.com/photo-1497636577773-f1231844b336",
			Kind:       string(entity.KindDeal),
			PublishedAt: time.Date(
				2026, 3, 6, 16, 30, 0, 0, jst,
			),
		},
		{
			SourceName: "Bean Market",
			Title:      "会員限定クーポン配布",
			Summary:    "ログイン会員向けに今月使える割引クーポンを配布中です。",
			URL:        "https://example.com/bean-market/member-coupon",
			ImageURL:   "https://images.unsplash.com/photo-1517048676732-d65bc937f952",
			Kind:       string(entity.KindDeal),
			PublishedAt: time.Date(
				2026, 3, 5, 11, 5, 0, 0, jst,
			),
		},

		//shop
		{
			SourceName: "Cafe Guide",
			Title:      "渋谷の新店舗オープン",
			Summary:    "駅近で立ち寄りやすい新しいコーヒースタンドを紹介します。",
			URL:        "https://example.com/cafe-guide/shibuya-new-shop",
			ImageURL:   "https://images.unsplash.com/photo-1501339847302-ac426a4a7cbb",
			Kind:       string(entity.KindShop),
			PublishedAt: time.Date(
				2026, 3, 4, 11, 0, 0, 0, jst,
			),
		},
		{
			SourceName: "Cafe Guide",
			Title:      "下北沢の隠れ家カフェ",
			Summary:    "静かに過ごしたい日に向いている落ち着いた店舗です。",
			URL:        "https://example.com/cafe-guide/shimokitazawa-cafe",
			ImageURL:   "https://images.unsplash.com/photo-1453614512568-c4024d13c247",
			Kind:       string(entity.KindShop),
			PublishedAt: time.Date(
				2026, 3, 3, 14, 45, 0, 0, jst,
			),
		},
		{
			SourceName: "Cafe Guide",
			Title:      "新宿で朝早く開く一杯",
			Summary:    "出勤前にも寄りやすい朝営業のコーヒースポットを紹介します。",
			URL:        "https://example.com/cafe-guide/shinjuku-morning",
			ImageURL:   "https://images.unsplash.com/photo-1509042239860-f550ce710b93",
			Kind:       string(entity.KindShop),
			PublishedAt: time.Date(
				2026, 3, 2, 7, 30, 0, 0, jst,
			),
		},
		{
			SourceName: "Cafe Guide",
			Title:      "浅草のレトロ喫茶まとめ",
			Summary:    "観光ついでに立ち寄れる老舗寄りの喫茶店を整理しました。",
			URL:        "https://example.com/cafe-guide/asakusa-retro",
			ImageURL:   "https://images.unsplash.com/photo-1442512595331-e89e73853f31",
			Kind:       string(entity.KindShop),
			PublishedAt: time.Date(
				2026, 3, 1, 15, 20, 0, 0, jst,
			),
		},
	}

	for _, def := range defs {
		src, ok := sourceMap[def.SourceName]
		if !ok {
			return fmt.Errorf("seed item source not found: %s", def.SourceName)
		}

		var exists entity.Item
		err := d.G.
			Where("title = ? AND kind = ? AND source_id = ?", def.Title, def.Kind, src.ID).
			First(&exists).
			Error

		if err == nil {
			continue
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("find seed item %s: %w", def.Title, err)
		}

		item := entity.Item{
			Title:       def.Title,
			Summary:     strPtr(def.Summary),
			URL:         strPtr(def.URL),
			ImageURL:    strPtr(def.ImageURL),
			Kind:        def.Kind,
			SourceID:    src.ID,
			PublishedAt: def.PublishedAt,
		}

		if err := d.G.Create(&item).Error; err != nil {
			return fmt.Errorf("create seed item %s: %w", def.Title, err)
		}
	}

	return nil
}

type seedItemDef struct {
	SourceName  string
	Title       string
	Summary     string
	URL         string
	ImageURL    string
	Kind        string
	PublishedAt time.Time
}

func strPtr(s string) *string {
	v := strings.TrimSpace(s)
	if v == "" {
		return nil
	}
	return &v
}
