# コーヒーSPA

コーヒー系トピックを扱うSPAです。  
Go + Echo + GORM のバックエンドと、React + TypeScript + Vite のフロントエンドで構成しています。

このリポジトリでは、以下を実装しています。

- 認証（サインアップ / メール確認 / ログイン / リフレッシュ / ログアウト / パスワード再設定）
- 記事一覧 / トップ表示 / 記事詳細
- Source 一覧 / Source 作成
- Admin による Item / Source 作成
- モーダルでのプレビュー表示
- 記事詳細ページでの全文表示

---

## 技術スタック

### Backend

- Go
- Echo
- GORM
- PostgreSQL
- Redis
- JWT
- CSRF（refresh / logout）

### Frontend

- React
- TypeScript
- Vite
- React Router
- Tailwind CSS

---

## ディレクトリ構成

```txt
backend/   API サーバ・ドメインロジック・DB
frontend/  SPA フロントエンド
docs/      仕様書・補助資料
```

---

## 起動方法

### Docker 起動

```bash
docker compose up --build
```

---

## 動作確認コマンド

### Frontend 品質確認

```bash
cd frontend
npx tsc -b
npm run lint
npm run build
```

### Backend テスト

```bash
cd backend
go test ./...
```

### Backend E2E テスト

E2E実行前にPostgreSQL / Redis / APIが起動している必要があります。

```bash
cd backend
BASE_URL=http://127.0.0.1:8080 go test ./tests/e2e -v -count=1
```

---

## 確認済みステータス

現時点で以下を通過済みです。

- Frontend: `npx tsc -b && npm run lint && npm run build`
- Backend: `go test ./...`
- Backend E2E: `BASE_URL=http://127.0.0.1:8080 go test ./tests/e2e -v -count=1`

---

## 実装上のポイント

### 記事表示

- トップページはカテゴリ切り替え型
- 記事クリックでモーダル表示
- 「もっと見る」で詳細ページへ遷移
- 詳細ページでは本文を全文表示

### 管理画面

- Admin ログイン時のみ Item / Source を作成可能
- Access Token 切れに対しては refresh を挟んで再試行する構成

### データ

- seed で記事データを投入

---
