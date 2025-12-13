# Sprint1 API

Go言語で実装されたタスク管理APIサーバーです。

## 技術スタック

- Go 1.21+
- MySQL
- gorilla/mux (ルーティング)
- go-sql-driver/mysql (MySQLドライバ)
- rs/cors (CORS対応)

## API仕様

| エンドポイント | メソッド | 説明 |
|---------------|---------|------|
| `/api/tasks` | GET | タスク一覧取得 |
| `/api/tasks/{id}` | GET | タスク詳細取得 |
| `/api/tasks` | POST | タスク作成 |
| `/api/tasks/{id}` | PUT | タスク更新 |
| `/api/tasks/{id}` | DELETE | タスク削除 |

## データベース

### tasksテーブル

| カラム | 型 | 説明 |
|--------|-----|------|
| id | INT | 主キー (AUTO_INCREMENT) |
| task | VARCHAR(255) | タスク内容 |
| limit_date | DATE | 期限日 |
| status | VARCHAR(20) | ステータス (pending/completed) |
| created_at | DATETIME | 作成日時 |
| updated_at | DATETIME | 更新日時 |

## 環境変数

| 変数名 | デフォルト値 | 説明 |
|--------|-------------|------|
| DB_HOST | localhost | MySQLホスト |
| DB_PORT | 3306 | MySQLポート |
| DB_USER | root | MySQLユーザー |
| DB_PASSWORD | password | MySQLパスワード |
| DB_NAME | taskdb | データベース名 |

## ローカル開発

### 前提条件

- Go 1.21以上
- MySQL 8.0以上

### セットアップ

1. データベースを作成

```sql
CREATE DATABASE taskdb;
```

2. 依存関係をインストール

```bash
go mod tidy
```

3. 環境変数を設定して起動

```bash
export DB_HOST=localhost
export DB_PASSWORD=yourpassword
go run main.go
```

サーバーが `http://localhost:8080` で起動します。

### ビルド

```bash
go build -o task-api main.go
./task-api
```

## APIの使用例

### タスク作成

```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{"task": "Buy groceries", "limit_date": "2024-12-31"}'
```

### タスク一覧取得

```bash
curl http://localhost:8080/api/tasks
```

### タスク更新

```bash
curl -X PUT http://localhost:8080/api/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"task": "Buy groceries", "limit_date": "2024-12-31", "status": "completed"}'
```

### タスク削除

```bash
curl -X DELETE http://localhost:8080/api/tasks/1
```
