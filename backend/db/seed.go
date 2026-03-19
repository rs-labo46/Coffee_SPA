package db

import (
	"fmt"
	"time"

	"coffee-spa/entity"

	"gorm.io/gorm"
)

func SeedDev(db *gorm.DB) error {
	if err := seedSources(db); err != nil {
		return err
	}

	if err := seedItems(db); err != nil {
		return err
	}

	return nil
}

func seedSources(db *gorm.DB) error {
	var count int64
	if err := db.Model(&entity.Source{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	sources := []entity.Source{
		{
			Name:    "Coffee Daily",
			SiteURL: strPtrIfNotEmpty("https://example.com/coffee-daily"),
		},
		{
			Name:    "Roast Journal",
			SiteURL: strPtrIfNotEmpty("https://example.com/roast-journal"),
		},
		{
			Name:    "Home Brew Note",
			SiteURL: strPtrIfNotEmpty("https://example.com/home-brew-note"),
		},
		{
			Name:    "Cafe Guide",
			SiteURL: strPtrIfNotEmpty("https://example.com/cafe-guide"),
		},
	}

	for _, src := range sources {
		if err := db.Create(&src).Error; err != nil {
			return err
		}
	}

	return nil
}

func seedItems(db *gorm.DB) error {
	var count int64
	if err := db.Model(&entity.Item{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	var sources []entity.Source
	if err := db.Order("id asc").Find(&sources).Error; err != nil {
		return err
	}

	if len(sources) == 0 {
		return fmt.Errorf("seed source not found")
	}

	now := time.Now()

	items := make([]entity.Item, 0, 40)

	items = append(items, buildNewsItems(now, sources[0].ID)...)
	items = append(items, buildRecipeItems(now, sources[1%len(sources)].ID)...)
	items = append(items, buildDealItems(now, sources[2%len(sources)].ID)...)
	items = append(items, buildShopItems(now, sources[3%len(sources)].ID)...)

	for _, item := range items {
		if err := db.Create(&item).Error; err != nil {
			return err
		}
	}

	return nil
}

func min4(a int, b int, c int, d int) int {
	m := a
	if b < m {
		m = b
	}
	if c < m {
		m = c
	}
	if d < m {
		m = d
	}
	return m
}

func buildNewsItems(now time.Time, sourceID int64) []entity.Item {
	titles := []string{
		"スペシャルティコーヒー市場で浅煎り需要が再び拡大",
		"都市型ロースタリーがサブスク会員向け焙煎便を開始",
		"ペーパーフィルター価格の見直しで家庭抽出のコスト感に変化",
		"エチオピア新豆の入荷が始まりフローラル系の注目が上昇",
		"カフェ運営者の間で小型焙煎機の導入相談が増加",
		"コーヒーイベントで抽出器具の比較展示が話題に",
		"リユースカップ運用を進める店舗が都心部で増えている",
		"ミルの粒度安定性を重視した家庭用モデルが人気",
		"豆価格の変動を受けて、定番ブレンドの構成比を調整する店舗が増加",
	}

	summaries := []string{
		"業界トレンドの確認用ダミーデータです。短めの概要文を入れています。",
		"会員制モデルと定期配送の組み合わせが、小規模ロースターでも試され始めています。",
		"",
		"花のような香りや柑橘系の明るさを前面に出した構成が目立っています。",
		"導入コストを抑えながら自家焙煎へ移行したい事業者向けの話題です。",
		"",
		"実店舗での体験価値と環境配慮を両立する取り組みとして注目されています。",
		"刃の違い、回転数、清掃性など、比較軸が一般消費者にも広がっています。",
		"",
	}

	images := []string{
		"https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?auto=format&fit=crop&w=1200&q=80",
		"https://images.unsplash.com/photo-1447933601403-0c6688de566e?auto=format&fit=crop&w=1200&q=80",
		"",
		"https://images.unsplash.com/photo-1517701604599-bb29b565090c?auto=format&fit=crop&w=1200&q=80",
		"",
		"https://images.unsplash.com/photo-1459755486867-b55449bb39ff?auto=format&fit=crop&w=1200&q=80",
		"",
		"https://images.unsplash.com/photo-1461988320302-91bde64fc8e4?auto=format&fit=crop&w=1200&q=80",
		"",
		"https://images.unsplash.com/photo-1494314671902-399b18174975?auto=format&fit=crop&w=1200&q=80",
	}

	urls := []string{
		"https://example.com/news/1",
		"https://example.com/news/2",
		"",
		"https://example.com/news/4",
		"",
		"https://example.com/news/6",
		"",
		"https://example.com/news/8",
		"",
		"https://example.com/news/10",
	}

	n := min4(len(titles), len(summaries), len(images), len(urls))

	items := make([]entity.Item, 0, n)

	for i := 0; i < n; i++ {
		items = append(items, entity.Item{
			Title:       titles[i],
			Summary:     strPtrIfNotEmpty(summaries[i]),
			URL:         strPtrIfNotEmpty(urls[i]),
			ImageURL:    strPtrIfNotEmpty(images[i]),
			Kind:        string(entity.KindNews),
			SourceID:    sourceID,
			PublishedAt: now.Add(time.Duration(-(i + 1)) * 6 * time.Hour),
		})
	}

	return items
}

func buildRecipeItems(now time.Time, sourceID int64) []entity.Item {
	titles := []string{
		"ハンドドリップの基本比率を見直して甘さを出すレシピ",
		"アイスコーヒー向けに濃度を上げた抽出手順",
		"フレンチプレスで雑味を抑える湯温の考え方",
		"朝の一杯を早く淹れるための時短ドリップ構成",
		"中煎り豆でバランスを崩しにくい家庭向けレシピ",
		"エアロプレスで酸味を丸くする短時間抽出",
		"少量抽出でも味を薄くしにくい一人分レシピ",
		"来客時に安定して淹れやすい二杯取りの基準",
		"牛乳に合わせやすい深煎り向けの濃い抽出レシピ",
	}

	summaries := []string{
		"粉量、湯量、抽出時間の基本を見直したダミーレシピです。",
		"",
		"湯温を少し下げるだけでも口当たりが穏やかになります。",
		"",
		"失敗しにくいレシピを置いて、一覧の見え方を確認します。",
		"短時間でも薄くなりすぎないように攪拌を調整する想定です。",
		"",
		"抽出量が増えた時に味がぶれやすい人向けです。",
		"",
	}

	images := []string{
		"https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?auto=format&fit=crop&w=1200&q=80",
		"",
		"https://images.unsplash.com/photo-1517701604599-bb29b565090c?auto=format&fit=crop&w=1200&q=80",
		"",
		"https://images.unsplash.com/photo-1442512595331-e89e73853f31?auto=format&fit=crop&w=1200&q=80",
		"",
		"https://images.unsplash.com/photo-1509042239860-f550ce710b93?auto=format&fit=crop&w=1200&q=80",
		"",
		"https://images.unsplash.com/photo-1461023058943-07fcbe16d735?auto=format&fit=crop&w=1200&q=80",
		"",
	}

	urls := []string{
		"",
		"https://example.com/recipe/2",
		"",
		"https://example.com/recipe/4",
		"",
		"https://example.com/recipe/6",
		"",
		"https://example.com/recipe/8",
		"",
		"https://example.com/recipe/10",
	}

	n := min4(len(titles), len(summaries), len(images), len(urls))

	items := make([]entity.Item, 0, n)

	for i := 0; i < n; i++ {
		items = append(items, entity.Item{
			Title:       titles[i],
			Summary:     strPtrIfNotEmpty(summaries[i]),
			URL:         strPtrIfNotEmpty(urls[i]),
			ImageURL:    strPtrIfNotEmpty(images[i]),
			Kind:        string(entity.KindNews),
			SourceID:    sourceID,
			PublishedAt: now.Add(time.Duration(-(i + 1)) * 6 * time.Hour),
		})
	}

	return items
}

func buildDealItems(now time.Time, sourceID int64) []entity.Item {
	titles := []string{
		"週末限定でドリッパーが10%オフ",
		"初回購入向けの送料無料キャンペーン",
		"深煎りセットのまとめ買い値引き",
		"春の新生活向けコーヒー器具セール",
		"ミルとケトルの同時購入で割引適用",
		"定期便スタート記念のクーポン配布",
		"アイスコーヒー器具の季節セール",
		"店舗受け取り限定の豆セット特価",
		"レビュー投稿で次回使えるクーポン配布",
	}

	summaries := []string{
		"価格表示の見え方確認用ダミーデータ。",
		"",
		"まとめ買い導線がある時の一覧密度を確かめます。",
		"",
		"複数商品を組み合わせた訴求の見え方確認用です。",
		"",
		"季節キャンペーンの短い説明文です。",
		"",
		"",
	}

	images := []string{
		"",
		"https://images.unsplash.com/photo-1512568400610-62da28bc8a13?auto=format&fit=crop&w=1200&q=80",
		"",
		"https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?auto=format&fit=crop&w=1200&q=80",
		"",
		"https://images.unsplash.com/photo-1447933601403-0c6688de566e?auto=format&fit=crop&w=1200&q=80",
		"",
		"https://images.unsplash.com/photo-1453614512568-c4024d13c247?auto=format&fit=crop&w=1200&q=80",
		"",
		"https://images.unsplash.com/photo-1461988320302-91bde64fc8e4?auto=format&fit=crop&w=1200&q=80",
	}

	urls := []string{
		"https://example.com/deal/1",
		"",
		"https://example.com/deal/3",
		"",
		"https://example.com/deal/5",
		"",
		"https://example.com/deal/7",
		"",
		"https://example.com/deal/9",
		"",
	}

	n := min4(len(titles), len(summaries), len(images), len(urls))

	items := make([]entity.Item, 0, n)

	for i := 0; i < n; i++ {
		items = append(items, entity.Item{
			Title:       titles[i],
			Summary:     strPtrIfNotEmpty(summaries[i]),
			URL:         strPtrIfNotEmpty(urls[i]),
			ImageURL:    strPtrIfNotEmpty(images[i]),
			Kind:        string(entity.KindNews),
			SourceID:    sourceID,
			PublishedAt: now.Add(time.Duration(-(i + 1)) * 6 * time.Hour),
		})
	}

	return items
}

func buildShopItems(now time.Time, sourceID int64) []entity.Item {
	titles := []string{
		"駅前に小型ロースタリー併設店がオープン",
		"朝営業に強いカフェの新店舗情報",
		"自家製スイーツと相性が良い人気店",
		"深夜まで営業する作業向けカフェ",
		"豆の量り売りに対応した地域密着店",
		"静かな空間でハンドドリップを味わえる店",
		"テイクアウト需要に強いスタンド型ショップ",
		"焙煎体験イベントを行う店舗の紹介",
		"地方ロースターの豆を週替わりで出す店",
	}

	summaries := []string{
		"新店カードの見え方確認用です。",
		"朝利用しやすい店舗情報を想定した短い説明です。",
		"",
		"作業利用、席数、電源有無などが気になる人向けの想定です。",
		"",
		"静かな店のニーズ確認用です。",
		"",
		"イベント性のある店舗情報が混ざった時の見え方確認です。",
		"",
	}

	images := []string{
		"https://images.unsplash.com/photo-1442512595331-e89e73853f31?auto=format&fit=crop&w=1200&q=80",
		"",
		"https://images.unsplash.com/photo-1453614512568-c4024d13c247?auto=format&fit=crop&w=1200&q=80",
		"",
		"https://images.unsplash.com/photo-1494314671902-399b18174975?auto=format&fit=crop&w=1200&q=80",
		"",
		"https://images.unsplash.com/photo-1509042239860-f550ce710b93?auto=format&fit=crop&w=1200&q=80",
		"",
		"https://images.unsplash.com/photo-1517701604599-bb29b565090c?auto=format&fit=crop&w=1200&q=80",
		"",
	}

	urls := []string{
		"",
		"https://example.com/shop/2",
		"",
		"https://example.com/shop/4",
		"",
		"https://example.com/shop/6",
		"",
		"https://example.com/shop/8",
		"",
		"https://example.com/shop/10",
	}

	n := min4(len(titles), len(summaries), len(images), len(urls))

	items := make([]entity.Item, 0, n)

	for i := 0; i < n; i++ {
		items = append(items, entity.Item{
			Title:       titles[i],
			Summary:     strPtrIfNotEmpty(summaries[i]),
			URL:         strPtrIfNotEmpty(urls[i]),
			ImageURL:    strPtrIfNotEmpty(images[i]),
			Kind:        string(entity.KindNews),
			SourceID:    sourceID,
			PublishedAt: now.Add(time.Duration(-(i + 1)) * 6 * time.Hour),
		})
	}

	return items
}

func strPtrIfNotEmpty(v string) *string {
	if v == "" {
		return nil
	}

	return &v
}
