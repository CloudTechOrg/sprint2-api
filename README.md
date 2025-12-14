# 名簿管理API (Meibo API)

Go言語で実装された名簿管理APIサーバーです。

## 技術スタック

- Go 1.21+
- MySQL 8.0
- gorilla/mux (ルーティング)
- go-sql-driver/mysql (MySQLドライバ)
- rs/cors (CORS対応)

## API仕様

| エンドポイント | メソッド | 説明 |
|---------------|---------|------|
| `/` | GET | ヘルスチェック |
| `/api/health` | GET | ヘルスチェック |
| `/api/persons` | GET | 名簿一覧取得 |
| `/api/persons/{id}` | GET | 名簿詳細取得 |
| `/api/persons` | POST | 名簿登録 |
| `/api/persons/{id}` | PUT | 名簿更新 |
| `/api/persons/{id}` | DELETE | 名簿削除 |

## ステータスコード一覧

### GET `/` , `/api/health` - ヘルスチェック

| ステータスコード | 説明 |
|-----------------|------|
| 200 OK | 正常動作中 |

**レスポンス例:**
```json
{
  "status": "ok",
  "service": "meibo-api"
}
```

### GET `/api/persons` - 名簿一覧取得

| ステータスコード | 説明 |
|-----------------|------|
| 200 OK | 一覧取得成功 |
| 500 Internal Server Error | データベースエラー |

**成功レスポンス例 (200):**
```json
[
  {
    "id": 1,
    "name": "山田太郎",
    "email": "yamada@example.com",
    "phone": "090-1234-5678"
  },
  {
    "id": 2,
    "name": "鈴木花子",
    "email": "suzuki@example.com",
    "phone": "090-8765-4321"
  }
]
```

### GET `/api/persons/{id}` - 名簿詳細取得

| ステータスコード | 説明 |
|-----------------|------|
| 200 OK | 取得成功 |
| 400 Bad Request | IDが無効（数値でない） |
| 404 Not Found | 指定されたIDの名簿が存在しない |
| 500 Internal Server Error | データベースエラー |

**成功レスポンス例 (200):**
```json
{
  "id": 1,
  "name": "山田太郎",
  "email": "yamada@example.com",
  "phone": "090-1234-5678"
}
```

**エラーレスポンス例 (400):**
```
Invalid ID
```

**エラーレスポンス例 (404):**
```
Person not found
```

### POST `/api/persons` - 名簿登録

| ステータスコード | 説明 |
|-----------------|------|
| 201 Created | 登録成功 |
| 400 Bad Request | リクエストボディが無効、または名前が未入力 |
| 500 Internal Server Error | データベースエラー |

**リクエスト例:**
```json
{
  "name": "田中一郎",
  "email": "tanaka@example.com",
  "phone": "090-1111-2222"
}
```

**成功レスポンス例 (201):**
```json
{
  "id": 3,
  "name": "田中一郎",
  "email": "tanaka@example.com",
  "phone": "090-1111-2222"
}
```

**エラーレスポンス例 (400):**
```
Invalid request body
```
または
```
Name is required
```

### PUT `/api/persons/{id}` - 名簿更新

| ステータスコード | 説明 |
|-----------------|------|
| 200 OK | 更新成功 |
| 400 Bad Request | IDが無効、リクエストボディが無効、または名前が未入力 |
| 404 Not Found | 指定されたIDの名簿が存在しない |
| 500 Internal Server Error | データベースエラー |

**リクエスト例:**
```json
{
  "name": "田中一郎（更新）",
  "email": "tanaka-new@example.com",
  "phone": "090-3333-4444"
}
```

**成功レスポンス例 (200):**
```json
{
  "id": 3,
  "name": "田中一郎（更新）",
  "email": "tanaka-new@example.com",
  "phone": "090-3333-4444"
}
```

**エラーレスポンス例 (404):**
```
Person not found
```

### DELETE `/api/persons/{id}` - 名簿削除

| ステータスコード | 説明 |
|-----------------|------|
| 204 No Content | 削除成功（レスポンスボディなし） |
| 400 Bad Request | IDが無効（数値でない） |
| 404 Not Found | 指定されたIDの名簿が存在しない |
| 500 Internal Server Error | データベースエラー |

## データベース

### personsテーブル

| カラム | 型 | 説明 |
|--------|-----|------|
| id | INT | 主キー (AUTO_INCREMENT) |
| name | VARCHAR(255) | 名前（必須） |
| email | VARCHAR(255) | メールアドレス |
| phone | VARCHAR(50) | 電話番号 |
| created_at | TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP | 更新日時 |

## 環境変数

| 変数名 | 必須 | 説明 |
|--------|-----|------|
| DB_HOST | ○ | MySQLホスト（RDSエンドポイント） |
| DB_USER | ○ | MySQLユーザー |
| DB_PASSWORD | ○ | MySQLパスワード |
| DB_NAME | ○ | データベース名（meibo） |

## ローカル開発

### 前提条件

- Go 1.21以上
- MySQL 8.0以上

### セットアップ

1. データベースを作成

```sql
CREATE DATABASE meibo;
```

2. 依存関係をインストール

```bash
go mod tidy
```

3. 環境変数を設定して起動

```bash
export DB_HOST=localhost
export DB_USER=root
export DB_PASSWORD=yourpassword
export DB_NAME=meibo
go run main.go
```

サーバーが `http://localhost:80` で起動します。

### ビルド

```bash
go build -o meibo-api main.go
./meibo-api
```

## APIの使用例

### 名簿登録

```bash
curl -X POST http://localhost/api/persons \
  -H "Content-Type: application/json" \
  -d '{"name": "山田太郎", "email": "yamada@example.com", "phone": "090-1234-5678"}'
```

### 名簿一覧取得

```bash
curl http://localhost/api/persons
```

### 名簿詳細取得

```bash
curl http://localhost/api/persons/1
```

### 名簿更新

```bash
curl -X PUT http://localhost/api/persons/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "山田太郎（更新）", "email": "yamada-new@example.com", "phone": "090-9999-0000"}'
```

### 名簿削除

```bash
curl -X DELETE http://localhost/api/persons/1
```
