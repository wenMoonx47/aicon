# 所持品管理API

高級品やコレクションアイテムを管理するREST APIサーバーです。

## 📋 API仕様

### ベースURL
```
http://localhost:8080
```

### エンドポイント一覧

| メソッド | パス | 説明 | ステータスコード |
|---------|------|------|-----------------|
| GET | `/health` | ヘルスチェック | 200 |
| GET | `/items` | 全アイテム取得 | 200 |
| POST | `/items` | アイテム登録 | 201, 400 |
| GET | `/items/{id}` | 特定アイテム取得 | 200, 404 |
| PATCH | `/items/{id}` | アイテム部分更新 | 200, 400, 404 |
| DELETE | `/items/{id}` | アイテム削除 | 204, 404 |
| GET | `/items/summary` | カテゴリー別集計 | 200 |

### データ形式

#### アイテム (Item)
```json
{
  "id": 1,
  "name": "ロレックス デイトナ",
  "category": "時計",
  "brand": "ROLEX",
  "purchase_price": 1500000,
  "purchase_date": "2023-01-15",
  "created_at": "2023-01-15T10:00:00Z",
  "updated_at": "2023-01-15T10:00:00Z"
}
```

#### 有効なカテゴリー
- `時計`
- `バッグ`
- `ジュエリー`
- `靴`
- `その他`

### バリデーションルール

| フィールド | 必須 | 制限 |
|-----------|------|------|
| name | ✓ | 100文字以内 |
| category | ✓ | 有効なカテゴリーのみ |
| brand | ✓ | 100文字以内 |
| purchase_price | ✓ | 0以上の整数 |
| purchase_date | ✓ | YYYY-MM-DD形式 |

### API使用例

#### 1. 全アイテム取得
```bash
curl -X GET http://localhost:8080/items
```

**レスポンス:**
```json
[
  {
    "id": 1,
    "name": "ロレックス デイトナ",
    "category": "時計",
    "brand": "ROLEX",
    "purchase_price": 1500000,
    "purchase_date": "2023-01-15",
    "created_at": "2023-01-15T10:00:00Z",
    "updated_at": "2023-01-15T10:00:00Z"
  }
]
```

#### 2. アイテム登録
```bash
curl -X POST http://localhost:8080/items \
  -H "Content-Type: application/json" \
  -d '{
    "name": "エルメス バーキン",
    "category": "バッグ",
    "brand": "HERMÈS",
    "purchase_price": 2000000,
    "purchase_date": "2023-02-20"
  }'
```

#### 3. 特定アイテム取得
```bash
curl -X GET http://localhost:8080/items/1
```

#### 4. アイテム部分更新
```bash
curl -X PATCH http://localhost:8080/items/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "更新されたロレックス デイトナ",
    "purchase_price": 1800000
  }'
```

**更新可能フィールド:**
- `name`: アイテム名
- `brand`: ブランド名  
- `purchase_price`: 購入価格

**不変フィールド（更新不可）:**
- `id`: アイテムID
- `category`: カテゴリー
- `purchase_date`: 購入日
- `created_at`: 作成日時
- `updated_at`: 更新日時（自動更新）

#### 5. アイテム削除
```bash
curl -X DELETE http://localhost:8080/items/1
```

#### 6. カテゴリー別集計
```bash
curl -X GET http://localhost:8080/items/summary
```

**レスポンス:**
```json
{
  "categories": {
    "時計": 2,
    "バッグ": 1,
    "ジュエリー": 3,
    "靴": 0,
    "その他": 1
  },
  "total": 7
}
```

### エラーレスポンス形式

```json
{
  "error": "validation failed",
  "details": [
    "name is required",
    "purchase_price must be 0 or greater"
  ]
}
```

## 🛠️ 技術スタック

- **言語**: Go 1.23
- **フレームワーク**: Echo v4
- **データベース**: MySQL 8.0
- **コンテナ**: Docker & Docker Compose

## 📁 プロジェクト構成

```
.
├── cmd/
│   └── main.go                 # エントリーポイント
├── internal/
│   ├── domain/
│   │   ├── entity/            # ドメインエンティティ
│   │   └── errors/            # ドメインエラー
│   ├── infrastructure/
│   │   ├── config/            # 設定管理
│   │   ├── database/          # データベース接続
│   │   └── server/            # HTTPサーバー
│   ├── interfaces/
│   │   ├── controller/        # HTTPハンドラー
│   │   └── database/          # リポジトリ
│   └── usecase/              # ビジネスロジック
├── sql/
│   └── init.sql              # データベース初期化
├── docker-compose.yml
├── Dockerfile
├── .env.example
└── README.md
```

## 🔧 開発環境

### 前提条件
- Docker
- Docker Compose

### ローカル開発
```bash
# 依存関係をインストール
go mod download

# ローカルでMySQLを起動（docker-compose経由）
docker-compose up -d mysql

# 環境変数を設定（ローカル開発用）
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=password
export DB_NAME=items_db

# アプリケーションを起動
go run cmd/main.go
```

### テストデータ

初期データとして以下のアイテムが登録されています：

1. ロレックス デイトナ (時計)
2. エルメス バーキン (バッグ)
3. ティファニー ネックレス (ジュエリー)
4. ルブタン パンプス (靴)
5. アップルウォッチ (その他)
