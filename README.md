# Todo アプリ - GraphQL + Ent + Atlas

Go言語で構築されたTODO APIです。以下の技術スタックを使用しています：

- **Ent**: スキーマファーストなGo エンティティフレームワーク
- **GraphQL**: クエリ言語とランタイム (graph-gophers/graphql-go)
- **Atlas**: データベーススキーママイグレーションツール
- **PostgreSQL**: データベース
- **Chi**: HTTPルーター

## 機能

- ✅ Todo の作成（CREATE）
- ✅ Todo の取得（READ）
- ✅ Todo の更新（UPDATE）
- ✅ Todo の削除（DELETE）
- ✅ GraphQL Playground によるAPIテスト
- ✅ Atlas によるスキーママイグレーション
- ✅ Docker コンテナ対応
- ✅ REST API 互換性

## プロジェクト構造

```
├── cmd/server/          # メインアプリケーション
├── ent/                 # Ent エンティティとスキーマ
│   └── schema/
├── internal/graph/      # GraphQL実装
│   ├── schema.go        # GraphQLスキーマ定義
│   └── resolver.go      # GraphQL リゾルバー
├── migrations/          # Atlas マイグレーションファイル
├── atlas.hcl           # Atlas 設定ファイル
├── docker-compose.yml  # Docker Compose 設定
└── Makefile           # 開発用コマンド
```

## セットアップ

### 1. 依存関係のインストール

```bash
go mod download
```

### 2. データベースの起動

```bash
make docker-up
```

### 3. データベースのセットアップ

```bash
make db-setup
```

### 4. サーバーの起動

```bash
make dev
```

アプリケーションは `http://localhost:9000` で起動します。

## API エンドポイント

- **GraphQL Playground**: `http://localhost:9000/` 🎮
- **GraphQL API**: `http://localhost:9000/graphql` 🔗
- **Health Check**: `http://localhost:9000/health` ❤️

### GraphQL API

メインのAPIはGraphQLで提供されています。GraphQL Playgroundを使って直感的にAPIを探索・テストできます。

## GraphQL API使用例

### GraphQL Playgroundでの操作

ブラウザで `http://localhost:9000/` にアクセスすると、GraphQL Playgroundが開きます。
以下のクエリやミューテーションを試すことができます。

### クエリ例

**全てのTodoを取得:**
```graphql
query {
  todos {
    id
    title
    description
    completed
    createdAt
    updatedAt
  }
}
```

**特定のTodoを取得:**
```graphql
query {
  todo(id: "1") {
    id
    title
    description
    completed
    createdAt
    updatedAt
  }
}
```

### ミューテーション例

**新しいTodoを作成:**
```graphql
mutation {
  createTodo(input: {
    title: "買い物に行く"
    description: "牛乳とパンを買う"
  }) {
    id
    title
    description
    completed
    createdAt
    updatedAt
  }
}
```

**Todoを更新:**
```graphql
mutation {
  updateTodo(id: "1", input: {
    title: "更新されたタイトル"
    completed: true
  }) {
    id
    title
    description
    completed
    updatedAt
  }
}
```

**Todoを削除:**
```graphql
mutation {
  deleteTodo(id: "1")
}
```

### cURLでのGraphQL API使用

```bash
# GraphQLクエリを実行
curl -X POST http://localhost:9000/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "query { todos { id title completed } }"
  }'

# GraphQLミューテーションを実行
curl -X POST http://localhost:9000/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation($input: CreateTodoInput!) { createTodo(input: $input) { id title } }",
    "variables": {
      "input": {
        "title": "新しいタスク",
        "description": "GraphQLで作成"
      }
    }
  }'
```

### REST API エンドポイント（互換性のため）

GraphQLの他に、REST APIも利用できます：

- **GET /todos** - 全てのTodoを取得
- **POST /todos** - 新しいTodoを作成
- **GET /todos/{id}** - 特定のTodoを取得
- **PUT /todos/{id}** - Todoを更新
- **DELETE /todos/{id}** - Todoを削除

## 開発用コマンド

```bash
# 使用可能なコマンドを表示
make help

# 開発サーバーを起動
make dev

# アプリケーションをビルド
make build

# テストを実行
make test

# EntとGraphQLコードを再生成
make generate

# データベースのセットアップ
make db-setup

# データベースコンテナを起動
make docker-up

# データベースコンテナを停止
make docker-down
```

## 環境変数

| 変数名 | デフォルト値 | 説明 |
|--------|--------------|------|
| `PORT` | `9000` | サーバーのポート番号 |
| `DATABASE_URL` | `postgres://postgres:password@localhost:5432/todoapp?sslmode=disable` | データベース接続文字列 |

## Docker を使用した実行

```bash
# 全体をDockerで実行
docker-compose up --build

# データベースのみ起動
docker-compose up postgres
```

## Atlas マイグレーション

Atlas を使用してデータベーススキーマを管理できます：

```bash
# 新しいマイグレーションを作成
atlas migrate diff --env local

# マイグレーションを適用
atlas migrate apply --env local

# データベーススキーマの状態を確認
atlas migrate status --env local
```

## 技術仕様

- **GraphQLライブラリ**: `github.com/graph-gophers/graphql-go`
- **HTTPルーター**: Chi v5
- **ORM**: Ent
- **マイグレーション**: Atlas
- **データベース**: PostgreSQL
- **Go Version**: 1.21+

## ライセンス

MIT License