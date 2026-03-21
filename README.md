# コーヒーSPA

コーヒー情報を 1 つの SPA で閲覧・管理するアプリです。  
フロントエンドは React + TypeScript + Vite、バックエンドは Go + Echo + GORM で構成しています。

公開ユーザー向けの記事閲覧に加えて、認証、メール確認、パスワード再設定、管理者による Item / Source 作成までを実装しています。

---

## 概要

| 項目           | 内容                                                    |
| -------------- | ------------------------------------------------------- |
| アプリ名       | コーヒーSPA                                             |
| フロントエンド | React / TypeScript / Vite / React Router / Tailwind CSS |
| バックエンド   | Go / Echo / GORM / PostgreSQL / Redis / JWT             |
| 主な用途       | コーヒー関連記事の表示、認証、管理画面、記事登録        |
| 開発方針       | シンプルで読みやすい構成、責務を分けた実務寄りの実装    |

---

## 主な機能

| 区分         | 機能                                                              |
| ------------ | ----------------------------------------------------------------- |
| 認証         | サインアップ / メール確認 / ログイン / リフレッシュ / ログアウト  |
| パスワード   | パスワード再設定メール送信 / 再設定                               |
| 公開画面     | トップ表示 / 記事一覧 / 記事詳細 / Source一覧                     |
| 管理画面     | Admin による Item 作成 / Source 作成                              |
| UI           | モーダルでの記事プレビュー / 詳細ページで全文表示                 |
| セキュリティ | JWT / Refresh Token / CSRF / CORS / Security Headers / Rate Limit |

---

## 技術スタック

### Backend

| 技術       | 用途                         |
| ---------- | ---------------------------- |
| Go         | API 実装                     |
| Echo       | ルーティング / HTTP ハンドラ |
| GORM       | DB アクセス                  |
| PostgreSQL | 永続データ保存               |
| Redis      | Rate Limit / 補助用途        |
| JWT        | Access Token                 |
| Cookie     | Refresh Token / CSRF Token   |

### Frontend

| 技術         | 用途                |
| ------------ | ------------------- |
| React        | UI                  |
| TypeScript   | 型安全な実装        |
| Vite         | 開発サーバ / ビルド |
| React Router | 画面遷移            |
| Tailwind CSS | スタイリング        |

---

## ディレクトリ構成

| パス                 | 内容                                     |
| -------------------- | ---------------------------------------- |
| `backend/`           | API サーバ、ユースケース、リポジトリ、DB |
| `frontend/`          | SPA フロントエンド                       |
| `docs/`              | 仕様書・補助資料                         |
| `docker-compose.yml` | 開発用コンテナ定義                       |
| `.env`               | Docker / backend 用の環境変数            |
| `frontend/.env`      | frontend ローカル起動用の環境変数        |

---

## 環境変数

### 1. ルートの `.env`

これは **Docker Composeとbackend 用** です。

| 変数名                | 例                      | 用途                     |
| --------------------- | ----------------------- | ------------------------ |
| `PORT`                | `8080`                  | API待受ポート            |
| `POSTGRES_USER`       | `myuser`                | DBユーザー               |
| `POSTGRES_PASSWORD`   | `mypassword`            | DBパスワード             |
| `POSTGRES_DB`         | `mydb`                  | DB名                     |
| `POSTGRES_PORT`       | `5433`                  | ホスト側DBポート         |
| `POSTGRES_HOST`       | `localhost`             | ローカル実行時のDB接続先 |
| `JWT_SECRET`          | `***`                   | JWT署名鍵                |
| `GO_ENV`              | `dev`                   | 実行環境                 |
| `API_DOMAIN`          | `localhost`             | Cookie / ドメイン設定用  |
| `FE_URL`              | `http://localhost:3000` | フロントURL              |
| `SEED_ADMIN_EMAIL`    | `admin@test.com`        | seed用管理者メール       |
| `SEED_ADMIN_PASSWORD` | `AdminPass123!`         | seed用管理者パスワード   |

### 2. `frontend/.env`

これは **frontendをローカルで `npm run dev` するとき専用** です。

| 変数名              | 例                      | 用途         |
| ------------------- | ----------------------- | ------------ |
| `VITE_API_BASE_URL` | `http://localhost:8080` | API の接続先 |

---

## `.env` が2つある理由

| ファイル        | 役割                        | 読み手                    |
| --------------- | --------------------------- | ------------------------- |
| `.env`          | backend / Docker Compose 用 | Go アプリ、docker compose |
| `frontend/.env` | Vite 用                     | React アプリ              |

---

## 起動方法

### Docker でまとめて起動

```bash
docker compose up --build
```

起動後のURL:

| サービス    | URL                     |
| ----------- | ----------------------- |
| Frontend    | `http://localhost:3000` |
| Backend API | `http://localhost:8080` |
| PostgreSQL  | `localhost:5433`        |
| Redis       | `localhost:6379`        |

---

## ローカル実行

### frontend

```bash
cd frontend
npm install
npm run dev
```

### backend

```bash
cd backend
go run main.go
```

---

## 動作確認コマンド

### Frontend確認

| コマンド        | 内容           |
| --------------- | -------------- |
| `npx tsc -b`    | 型チェック     |
| `npm run lint`  | ESLint         |
| `npm run build` | 本番ビルド確認 |

```bash
cd frontend
npx tsc -b
npm run lint
npm run build
```

### Backend テスト

| コマンド        | 内容                     |
| --------------- | ------------------------ |
| `go test ./...` | 単体テスト含む全体テスト |

```bash
cd backend
go test ./...
```

### Backend E2E テスト

E2E 実行前にPostgreSQL / Redis / APIが起動している必要があります。

| コマンド                                                         | 内容      |
| ---------------------------------------------------------------- | --------- |
| `BASE_URL=http://127.0.0.1:8080 go test ./tests/e2e -v -count=1` | E2Eテスト |

```bash
cd backend
BASE_URL=http://127.0.0.1:8080 go test ./tests/e2e -v -count=1
```

---

## 確認済みステータス

| 区分        | 確認内容                                                         |
| ----------- | ---------------------------------------------------------------- |
| Frontend    | `npx tsc -b && npm run lint && npm run build`                    |
| Backend     | `go test ./...`                                                  |
| Backend E2E | `BASE_URL=http://127.0.0.1:8080 go test ./tests/e2e -v -count=1` |

---

## 実装上のポイント

### 記事表示

| 項目         | 内容                             |
| ------------ | -------------------------------- |
| トップページ | カテゴリ切り替え型で表示         |
| 一覧UI       | カード形式で表示                 |
| プレビュー   | 記事クリックでモーダル表示       |
| 詳細導線     | 「もっと見る」で詳細ページへ遷移 |
| 詳細ページ   | 本文を全文表示                   |

### 管理画面

| 項目     | 内容                                        |
| -------- | ------------------------------------------- |
| 権限制御 | Admin ログイン時のみItem / Sourceを作成可能 |
| 認証維持 | Access Token切れ時はrefreshを挟んで再試行   |

### データ

| 項目       | 内容                             |
| ---------- | -------------------------------- |
| 初期データ | seed 記事データを投入            |
| 管理者     | `GO_ENV=dev`時にseed管理者を作成 |

---

## 開発用の初期アカウント

`.env`の値を使ってseedされます。

| 項目     | 値                    |
| -------- | --------------------- |
| Email    | `SEED_ADMIN_EMAIL`    |
| Password | `SEED_ADMIN_PASSWORD` |

---
