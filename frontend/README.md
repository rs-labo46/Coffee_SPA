# Coffee SPA

Coffee SPA は、**コーヒーに関する情報をシングルページに集約するアプリ**です。  
フロントエンドは **React + Vite + TypeScript**、バックエンドは **Go + Echo + GORM**、データストアは **PostgreSQL + Redis** で構成しています。
このリポジトリでは、**認証・メール確認・パスワード再設定・トークン更新・公開一覧・管理画面からの登録**まで、一連の基本フローを動かせるようにしています。

---

## このプロジェクトでできること

| 区分         | できること                                                 |
| ------------ | ---------------------------------------------------------- |
| 公開機能     | トップページで `news / recipe / deal / shop` の4分類を表示 |
| 公開機能     | 一覧ページでキーワード検索・種別絞り込み                   |
| 認証         | サインアップ、メール確認、ログイン、ログアウト             |
| 認証         | リフレッシュトークンによるアクセストークン再発行           |
| 認証         | 確認メール再送、パスワード再設定メール送信、パスワード更新 |
| 会員機能     | `/me` で自分のユーザー情報を確認                           |
| 管理機能     | 管理者がソース（出典）を登録                               |
| 管理機能     | 管理者が記事アイテムを登録                                 |
| 運用寄り機能 | 監査ログ保存、Redis によるレート制限                       |
| 品質担保     | E2E テストと一部ユニットテストを実装                       |

---

## 技術スタック

### フロントエンド

| 項目         | 内容                           |
| ------------ | ------------------------------ |
| UI           | React 19                       |
| ビルド       | Vite 8                         |
| 言語         | TypeScript                     |
| ルーティング | React Router                   |
| スタイリング | Tailwind CSS 4                 |
| API 通信     | `fetch` ベースのラッパー       |
| 認証状態管理 | React Context (`AuthProvider`) |

### バックエンド

| 項目               | 内容                              |
| ------------------ | --------------------------------- |
| 言語               | Go 1.25                           |
| Web フレームワーク | Echo                              |
| ORM                | GORM                              |
| DB                 | PostgreSQL 16                     |
| キャッシュ / 制御  | Redis 7                           |
| 認証               | JWT Bearer + Refresh Token        |
| パスワード         | bcrypt                            |
| API 仕様           | OpenAPI 3.0 (`docs/openapi.yaml`) |

### インフラ / 実行環境

| 項目           | 内容                    |
| -------------- | ----------------------- |
| コンテナ       | Docker / Docker Compose |
| API コンテナ   | distroless 実行イメージ |
| Front コンテナ | Node 22 Alpine          |
| DB コンテナ    | PostgreSQL 16           |
| Redis コンテナ | Redis 7                 |

---

## 画面構成

実装済みの主な画面は以下です。

| パス               | 役割                                       |
| ------------------ | ------------------------------------------ |
| `/`                | トップページ。種別ごとの最新アイテムを表示 |
| `/items`           | 一覧ページ。                               |
| `/signup`          | 新規登録                                   |
| `/verify-email`    | メール確認                                 |
| `/resend-verify`   | 確認メール再送                             |
| `/login`           | ログイン                                   |
| `/forgot-password` | パスワード再設定メール送信                 |
| `/reset-password`  | 新しいパスワード設定                       |
| `/me`              | ログインユーザー情報の確認                 |
| `/admin`           | 管理者向けの source / item 登録画面        |

`/me` はログイン必須、`/admin` は **admin 権限必須** 。

---

## API 構成

バックエンドの主要エンドポイントは以下です。

| 区分   | エンドポイント               | 内容                   |
| ------ | ---------------------------- | ---------------------- |
| Health | `GET /health`                | ヘルスチェック         |
| Auth   | `POST /auth/signup`          | ユーザー登録           |
| Auth   | `POST /auth/verify-email`    | メール確認             |
| Auth   | `POST /auth/resend-verify`   | 確認メール再送         |
| Auth   | `POST /auth/login`           | ログイン               |
| Auth   | `POST /auth/refresh`         | アクセストークン再発行 |
| Auth   | `POST /auth/logout`          | ログアウト             |
| Auth   | `GET /me`                    | 現在ユーザー情報       |
| Auth   | `POST /auth/password/forgot` | 再設定メール送信       |
| Auth   | `POST /auth/password/reset`  | パスワード更新         |
| Portal | `GET /items/top`             | 種別ごとのトップ表示   |
| Portal | `GET /items`                 | アイテム一覧 / 検索    |
| Portal | `GET /sources`               | ソース一覧             |
| Admin  | `POST /items`                | アイテム作成           |
| Admin  | `POST /sources`              | ソース作成             |

詳細は `docs/openapi.yaml` を参照。

---

## 認証の考え方

| 項目                 | 方式                                        |
| -------------------- | ------------------------------------------- |
| アクセストークン     | Bearer JWT をフロントで保持                 |
| リフレッシュトークン | Cookie で保持                               |
| CSRF 対策            | `csrf_token` Cookie + `X-CSRF-Token` ヘッダ |
| セッション失効       | `token_ver` による無効化                    |
| パスワード変更時     | 既存 refresh token を失効                   |
| リフレッシュ再利用   | 不正利用として 401 を返す                   |

### 重要な挙動

- `POST /auth/refresh` は **refresh_token cookie** と **CSRF トークン** の両方が必要。
- `POST /auth/logout` は JWT 保護。
- `GET /me` は `Cache-Control: no-store` を返す設計。
- メール送信は開発用として、**実メール送信の代わりにバックエンドログへリンクを出力**している。

そのため、サインアップ後やパスワード再設定時は、**バックエンドのログに出たリンクをブラウザで開く**流れになる。

---

## アーキテクチャ

| ディレクトリ         | 役割                                    |
| -------------------- | --------------------------------------- |
| `backend/entity`     | ドメインに近い構造体定義                |
| `backend/db`         | DB 接続、マイグレーション、開発用 seed  |
| `backend/repository` | 永続化処理                              |
| `backend/usecase`    | 業務ロジック                            |
| `backend/validator`  | 入力検証                                |
| `backend/policy`     | メール・パスワード・URL などのルール    |
| `backend/controller` | HTTP リクエスト / レスポンス処理        |
| `backend/middleware` | JWT、CSRF、CORS、セキュリティヘッダなど |
| `backend/router`     | ルーティング定義                        |
| `docs`               | OpenAPI、ER 図などのドキュメント        |
| `frontend/src/pages` | 各画面                                  |
| `frontend/src/auth`  | 認証状態管理                            |
| `frontend/src/lib`   | API クライアントや共通関数              |

## 依存の流れは

```text
router -> controller -> usecase -> repository -> db
                           -> validator / policy
```

## データモデル

| テーブル         | 役割                     |
| ---------------- | ------------------------ |
| `users`          | ユーザー本体             |
| `email_verifies` | メール確認トークン       |
| `pw_resets`      | パスワード再設定トークン |
| `refresh_tokens` | リフレッシュトークン管理 |
| `sources`        | 記事ソース / 出典        |
| `items`          | コーヒー関連アイテム     |
| `audit_logs`     | 監査ログ                 |

補足:

- `items.kind` は `news / recipe / deal / shop` の制約。
- `users.role` は `user / admin` の制約。
- `items.title` と `items.summary` には PostgreSQL の `pg_trgm` を使った検索用インデックスを作成。

---

## ディレクトリ構成

```text
.
├── docker-compose.yml
├── .env
├── docs/
│   ├── openapi.yaml
│   └── SPA_ER.pdf
├── backend/
│   ├── main.go
│   ├── config/
│   ├── controller/
│   ├── db/
│   ├── entity/
│   ├── middleware/
│   ├── policy/
│   ├── repository/
│   ├── router/
│   ├── tests/
│   │   └── e2e/
│   ├── usecase/
│   └── validator/
└── frontend/
    ├── src/
    │   ├── auth/
    │   ├── lib/
    │   └── pages/
    └── package.json
```

---

## セットアップ手順

### 1. 前提

必要なものは以下です。

| 項目                    | 用途                                 |
| ----------------------- | ------------------------------------ |
| Docker / Docker Compose | 全体起動                             |
| Git                     | ソース管理                           |
| Go                      | バックエンドをローカル起動する場合   |
| npm                     | フロントエンドをローカル起動する場合 |

まずは Docker 起動が最短です。

### 2. 環境変数を確認

ルートの `.env` では以下を使っています。

```env
PORT=8080
POSTGRES_USER=myuser
POSTGRES_PASSWORD=mypassword
POSTGRES_DB=mydb
POSTGRES_PORT=5433
POSTGRES_HOST=localhost
JWT_SECRET=your-secret
GO_ENV=dev
API_DOMAIN=localhost
FE_URL=http://localhost:3000
SEED_ADMIN_EMAIL=admin@test.com
SEED_ADMIN_PASSWORD=AdminPass123!
```

フロントエンド側は `frontend/.env` で以下を使う。

```env
VITE_API_BASE_URL=http://localhost:8080
```

### 3. Docker で起動

リポジトリルートで実行します。

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

### 4. 初回起動時に行われること

バックエンド起動時に以下が実行。

- `.env` の読み込み
- PostgreSQL 接続
- GORM によるマイグレーション
- `GO_ENV=dev` の場合、管理者アカウント seed
- Redis 接続確認
- Echo サーバ起動

開発用の管理者アカウントは `.env` の以下です。

```env
SEED_ADMIN_EMAIL=admin@test.com
SEED_ADMIN_PASSWORD=AdminPass123!
```

---

## 開発フローの確認方法

### サインアップから利用開始まで

1. `/signup` で登録する
2. バックエンドログに出力された verify リンクを開く
3. `/login` でログインする
4. `/me` で自分の状態を確認する

### パスワード再設定

1. `/forgot-password` でメールアドレス送信
2. バックエンドログに出力された reset リンクを開く
3. `/reset-password` で新しいパスワードを設定
4. 新しいパスワードでログインする

## 開発環境での確認メール / パスワード再設定メール

このプロジェクトでは、**開発環境では実メール送信を行いません**。  
その代わり、確認メール用リンク・パスワード再設定リンクは **backend(API) のログ** に出力されます。

### ログの確認手順

- docker compose logs -f api

### 確認メールのログ例

````text
coffee-spa-api  | 2026/03/20 04:25:41 [MAIL][VERIFY] to=testuser@test.com link=http://localhost:3000/verify-email?token=xxxxxxxx
### 管理者としてアイテム登録

1. 管理者アカウントでログイン
2. `/admin` を開く
3. 先に Source を作成
4. その Source を使って Item を作成
5. `/` または `/items` で公開表示を確認する

---

## テスト

### 実装しているテスト

| 種別              | 主な対象                                                                   |
| ----------------- | -------------------------------------------------------------------------- |
| E2E               | サインアップ、メール確認、ログイン、refresh、logout、password reset、items |
| Unit / Layer 単位 | controller、middleware、policy、usecase、validator、router                 |

実際に確認できたテストファイル例:

- `backend/tests/e2e/auth_test.go`
- `backend/tests/e2e/items_test.go`
- `backend/tests/e2e/refresh_logout_test.go`
- `backend/tests/e2e/reset_test.go`
- `backend/controller/auth_controller_test.go`
- `backend/usecase/item_usecase_test.go`
- `backend/validator/validator_test.go`

### 全テスト

```bash
cd backend
go test ./...
````

### E2E テスト

API を起動した状態で、別ターミナルから実行。

```bash
cd backend
BASE_URL=http://127.0.0.1:8080 go test ./tests/e2e
```

---

## ドキュメント

| ファイル                  | 内容             |
| ------------------------- | ---------------- |
| `docs/openapi.yaml`       | API仕様          |
| `docs/SPA_ER.pdf`         | ER図             |
| `docs/アーキテクチャ.png` | アーキテクチャ図 |

---

## 注意

| 詰みやすい点                 | 原因                           | 回避策                                              |
| ---------------------------- | ------------------------------ | --------------------------------------------------- |
| サインアップ後に先へ進めない | メールは実送信ではなくログ出力 | バックエンドログの verify リンクを開く              |
| refresh が失敗する           | CSRF ヘッダ or Cookie 不足     | `credentials: include` と `X-CSRF-Token` を確認     |
| `/admin` に入れない          | admin 権限がない               | seed 管理者でログインする                           |
| API が起動しない             | `.env` 不足、DB/Redis 未起動   | ルート `.env` と `docker compose up --build` を確認 |
| DB 接続できない              | ホスト / ポートの勘違い        | ローカルは `5433`、コンテナ内は `5432`              |
| 429 が返る                   | レート制限                     | 少し待つか、開発中の再試行間隔を空ける              |

---

## 今後の拡張ポイント

- item の更新・削除
- 管理画面の一覧 / 編集 / 削除
- 画像アップロード
- ページネーションの UI 強化
- 本物のメール送信基盤への差し替え
- RBAC の詳細化

---
