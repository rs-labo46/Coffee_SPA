# Coffee SPA

Coffee SPA は、**コーヒー関連の情報を集約して閲覧できるシングルページアプリケーション**。  
フロントエンドは **React + Vite + TypeScript**、バックエンドは **Go + Echo + GORM**、データストアは **PostgreSQL + Redis** で構成しています。

このリポジトリでは、単なる画面表示だけではなく、**認証、メール確認、パスワード再設定、公開一覧、管理者による source / item 登録**までを一通り確認できます。

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
| 管理機能     | 管理者が source を登録                                     |
| 管理機能     | 管理者が item を登録                                       |
| 運用寄り機能 | Redis を使ったレート制限                                   |
| 品質確認     | 単体テストと E2E テストを実装                              |

---

## 技術スタック

### フロントエンド

| 項目         | 内容                         |
| ------------ | ---------------------------- |
| UI           | React 19                     |
| ビルド       | Vite 8                       |
| 言語         | TypeScript 5                 |
| ルーティング | React Router 7               |
| スタイリング | Tailwind CSS 4               |
| API 通信     | `fetch` ベースの独自ラッパー |
| 状態管理     | React Context                |

### バックエンド

| 項目               | 内容                            |
| ------------------ | ------------------------------- |
| 言語               | Go 1.25.0                       |
| Web フレームワーク | Echo                            |
| ORM                | GORM                            |
| DB                 | PostgreSQL 16                   |
| 制御               | Redis 7                         |
| 認証               | JWT Bearer + Refresh Token      |
| パスワード         | bcrypt                          |
| API仕様            | OpenAPI 3 (`docs/openapi.yaml`) |

### 実行環境

| 項目     | 内容                    |
| -------- | ----------------------- |
| コンテナ | Docker / Docker Compose |
| API      | Go アプリをコンテナ起動 |
| Frontend | Node 22 Alpine          |
| DB       | PostgreSQL 16           |
| Redis    | Redis 7                 |

---

## 画面構成

| パス               | 役割                              |
| ------------------ | --------------------------------- |
| `/`                | トップページ                      |
| `/items`           | 一覧ページ                        |
| `/signup`          | 新規登録                          |
| `/verify-email`    | メール確認                        |
| `/resend-verify`   | 確認メール再送                    |
| `/login`           | ログイン                          |
| `/forgot-password` | パスワード再設定メール送信        |
| `/reset-password`  | 新しいパスワード設定              |
| `/me`              | ログインユーザー情報表示          |
| `/admin`           | 管理者向け source / item 登録画面 |

`/me` はログイン必須、`/admin` は admin 権限必須。

---

## API 構成

| 区分   | エンドポイント               | 内容                   |
| ------ | ---------------------------- | ---------------------- |
| Health | `GET /health`                | ヘルスチェック         |
| Auth   | `POST /auth/signup`          | ユーザー登録           |
| Auth   | `POST /auth/verify-email`    | メール確認             |
| Auth   | `POST /auth/resend-verify`   | 確認メール再送         |
| Auth   | `POST /auth/login`           | ログイン               |
| Auth   | `POST /auth/password/forgot` | 再設定メール送信       |
| Auth   | `POST /auth/password/reset`  | パスワード更新         |
| Auth   | `POST /auth/refresh`         | アクセストークン再発行 |
| Auth   | `POST /auth/logout`          | ログアウト             |
| Auth   | `GET /me`                    | 現在ユーザー情報       |
| Public | `GET /items/top`             | 種別ごとのトップ表示   |
| Public | `GET /items`                 | アイテム一覧 / 検索    |
| Public | `GET /sources`               | source 一覧            |
| Admin  | `POST /items`                | item 作成              |
| Admin  | `POST /sources`              | source 作成            |

詳細は `docs/openapi.yaml` を参照してください。

---

## 認証の考え方

このプロジェクトは、SPA の認証で最低限必要な流れをまとめて確認できる構成。

| 項目                 | 方式                                        |
| -------------------- | ------------------------------------------- |
| アクセストークン     | Bearer JWT                                  |
| リフレッシュトークン | Cookie 保持                                 |
| CSRF 対策            | `csrf_token` Cookie + `X-CSRF-Token` ヘッダ |
| セッション失効       | `token_ver` による無効化                    |
| logout               | JWT 認証が必要                              |
| refresh              | Cookie + CSRF の両方が必要                  |

### 重要な挙動

- `POST /auth/refresh` は **refresh_token Cookie** と **CSRF トークン** の両方が必要。
- `POST /auth/logout` は JWT 認証が必要。
- `GET /me` は認証必須。
- メール送信は開発時には実送信ではなく、**バックエンドログにリンクを出力する方式**。

そのため、サインアップ後やパスワード再設定時は、**バックエンドログに出た verify / reset リンクを開く**必要がある。

---

## アーキテクチャ

バックエンドは、役割ごとに分けた構成。

| ディレクトリ         | 役割                            |
| -------------------- | ------------------------------- |
| `backend/entity`     | ドメイン構造体                  |
| `backend/db`         | DB 接続、マイグレーション、seed |
| `backend/repository` | 永続化処理                      |
| `backend/usecase`    | 業務ロジック                    |
| `backend/validator`  | 入力検証                        |
| `backend/policy`     | バリデーションルール            |
| `backend/controller` | HTTP 入出力                     |
| `backend/middleware` | JWT、CSRF、CORS など            |
| `backend/router`     | ルーティング定義                |
| `backend/tests/e2e`  | E2E テスト                      |
| `frontend/src/pages` | 画面                            |
| `frontend/src/auth`  | 認証状態管理                    |
| `frontend/src/lib`   | API クライアントや共通関数      |

依存の流れは概ね以下。

```text
router -> controller -> usecase -> repository -> db
                           -> validator / policy
```

---

## データモデル

主なテーブルは以下。

| テーブル         | 役割                     |
| ---------------- | ------------------------ |
| `users`          | ユーザー本体             |
| `email_verifies` | メール確認トークン       |
| `pw_resets`      | パスワード再設定トークン |
| `refresh_tokens` | リフレッシュトークン管理 |
| `sources`        | 出典情報                 |
| `items`          | コーヒー関連アイテム     |
| `audit_logs`     | 監査ログ                 |

補足:

- `items.kind` は `news / recipe / deal / shop` の制約がある。
- `users.role` は `user / admin` の制約がある。
- `items.title` と `items.summary` は検索を考慮したインデックスを使っています。

---

## ディレクトリ構成

```text
.
├── .env
├── docker-compose.yml
├── docs/
│   ├── openapi.yaml
│   ├── SPA_ER.pdf
│   └── アーキテクチャ.png
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
    ├── package.json
    └── README.md
```

---

## セットアップ手順

### 1. 前提

| 項目                    | 用途                                 |
| ----------------------- | ------------------------------------ |
| Docker / Docker Compose | 全体起動                             |
| Git                     | ソース管理                           |
| Go 1.25 系              | バックエンドをローカル起動する場合   |
| Node.js / npm           | フロントエンドをローカル起動する場合 |

まずは Docker 起動が最短。

### 2. 環境変数

ルートの `.env` では少なくとも以下を使います。

```env
PORT=8080
POSTGRES_USER=myuser
POSTGRES_PASSWORD=mypassword
POSTGRES_DB=mydb
POSTGRES_PORT=5433
POSTGRES_HOST=localhost
REDIS_HOST=localhost
REDIS_PORT=6379
JWT_SECRET=your-secret
GO_ENV=dev
API_DOMAIN=localhost
FE_URL=http://localhost:3000
SEED_ADMIN_EMAIL=admin@test.com
SEED_ADMIN_PASSWORD=AdminPass123!
```

フロントエンド側は `frontend/.env` で以下。

```env
VITE_API_BASE_URL=http://localhost:8080
```

### 3. Docker で起動

リポジトリルートで実行。

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

- `.env` の読み込み
- PostgreSQL 接続
- GORM によるマイグレーション
- `GO_ENV=dev` の場合は管理者 seed
- Redis 接続確認
- Echo サーバ起動

開発用の管理者アカウントは `.env` の以下。

```env
SEED_ADMIN_EMAIL=admin@test.com
SEED_ADMIN_PASSWORD=AdminPass123!
```

---

## ローカル起動手順

### バックエンド

```bash
cd backend
go run main.go
```

### フロントエンド

```bash
cd frontend
npm install
npm run dev
```

この場合でも、PostgreSQL と Redis は先に起動しておく必要がある。

---

## 開発フローの確認方法

### サインアップから利用開始まで

1. `/signup` で登録する
2. バックエンドログの verify リンクを開く
3. `/login` でログインする
4. `/me` で状態を確認する

### パスワード再設定

1. `/forgot-password` でメールアドレス送信
2. バックエンドログの reset リンクを開く
3. `/reset-password` で新しいパスワードを設定
4. 新しいパスワードでログインする

### 管理者として登録確認

1. 管理者アカウントでログイン
2. `/admin` を開く
3. 先に Source を作成
4. その Source を使って Item を作成
5. `/` または `/items` で公開表示を確認する

---

## テスト

### 現状の考え方

このリポジトリには **unit に近いテスト** と **E2E テスト** の両方がある。  
ただし、「全レイヤーを完全網羅している」状態ではない。  
実務目線では、**MVPとしては十分前進しているが、今後も auth / repository / middleware は継続補強したい段階**。

### 実行コマンド

```bash
cd backend
go test ./...
```

E2E は API 起動後に別ターミナルで実行。

```bash
cd backend
BASE_URL=http://127.0.0.1:8080 go test ./tests/e2e -v -count=1
```

`localhost` ではなく `127.0.0.1` を使うのは、環境によって IPv6 優先の問題を避けるため。

### テスト対象の例

| 種別               | 主な対象                                                      |
| ------------------ | ------------------------------------------------------------- |
| controller         | auth, item, source のHTTP入出力                               |
| middleware         | CSRF, JWT など                                                |
| policy / validator | バリデーションルール                                          |
| usecase            | auth, item, source                                            |
| router             | 主要ルート                                                    |
| E2E                | signup, verify, login, refresh, logout, reset, items, sources |

### 注意点

- DB と Redis が起動していないと E2E は通らない。
- Docker のディスク容量不足で PostgreSQL が起動失敗することがある。
- レート制限に引っかかると 429 が返ることがある。

---

## 注意点

| 詰みやすい点                 | 原因                           | 回避策                                              |
| ---------------------------- | ------------------------------ | --------------------------------------------------- |
| サインアップ後に先へ進めない | メールは実送信ではなくログ出力 | バックエンドログの verify リンクを開く              |
| reset が進まない             | reset リンクもログ出力方式     | バックエンドログの reset リンクを開く               |
| refresh が失敗する           | CSRF ヘッダ or Cookie 不足     | `credentials: include` と `X-CSRF-Token` を確認     |
| `/admin` に入れない          | admin 権限がない               | seed 管理者でログインする                           |
| API が起動しない             | `.env` 不足、DB / Redis 未起動 | ルート `.env` と `docker compose up --build` を確認 |
| DB 接続できない              | ホストとコンテナ内ポートを混同 | ホストは `5433`、apiコンテナからは `5432`           |
| `No space left on device`    | Docker のディスク容量不足      | 不要 volume / image を掃除する                      |
| 429 が返る                   | Redis レート制限               | 少し待って再試行する                                |

---

## ドキュメント

| ファイル                  | 内容     |
| ------------------------- | -------- |
| `docs/openapi.yaml`       | API 仕様 |
| `docs/SPA_ER.pdf`         | ER 図    |
| `docs/アーキテクチャ.png` | 構成図   |

---

## 今後の拡張ポイント

- item の更新 / 削除
- 管理画面の一覧 / 編集 / 削除
- OpenAPI からの型生成
- CI での自動テスト
- メール送信基盤の実装
- repository / auth の追加テスト

---
