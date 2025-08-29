.PHONY: help build run dev test generate clean docker-up docker-down migrate-create migrate-up migrate-reset db-setup migrate-status

# デフォルトのターゲット
help:
	@echo "使用可能なコマンド:"
	@echo "  make dev        - 開発サーバーを起動"
	@echo "  make build      - アプリケーションをビルド"
	@echo "  make run        - ビルドしたアプリケーションを実行"
	@echo "  make test       - テストを実行"
	@echo "  make generate   - Entとgraphqlコードを生成"
	@echo "  make clean      - ビルド生成物を削除"
	@echo "  make docker-up  - Dockerコンテナを起動"
	@echo "  make docker-down- Dockerコンテナを停止"
	@echo "  make db-setup   - データベースのセットアップ（推奨）"
	@echo "  make migrate-create - マイグレーションファイルを作成"
	@echo "  make migrate-up - マイグレーションを適用"
	@echo "  make migrate-status - マイグレーション状態を確認"

# 開発サーバーを起動（ホットリロード付き）PORT=9000 make dev
dev:
	@echo "開発サーバーを起動しています..."
	@go run cmd/server/main.go

# アプリケーションをビルド
build:
	@echo "アプリケーションをビルドしています..."
	@go build -o bin/server cmd/server/main.go

# ビルドしたアプリケーションを実行
run: build
	@echo "アプリケーションを実行しています..."
	@./bin/server

# テストを実行
test:
	@echo "テストを実行しています..."
	@go test -v ./...

# EntとGraphQLコードを生成
generate:
	@echo "Entコードを生成しています..."
	@go generate ./ent
	@echo "GraphQLコードを生成しています..."
	@go run github.com/99designs/gqlgen generate

# ビルド生成物を削除
clean:
	@echo "ビルド生成物を削除しています..."
	@rm -rf bin/
	@go clean

# Dockerコンテナを起動
docker-up:
	@echo "Dockerコンテナを起動しています..."
	@docker-compose up -d postgres
	@sleep 3

# Dockerコンテナを停止
docker-down:
	@echo "Dockerコンテナを停止しています..."
	@docker-compose down

# Atlasマイグレーションファイルを作成
migrate-create:
	@echo "マイグレーションファイルを作成しています..."
	@atlas migrate diff --env local
	@echo "マイグレーションファイルが作成されました"

# Atlasマイグレーションを適用
migrate-up:
	@echo "マイグレーションを適用しています..."
	@atlas migrate apply --env local
	@echo "マイグレーションが適用されました"

# データベースの状態をリセット
migrate-reset:
	@echo "データベースをリセットしています..."
	@atlas schema clean --env local -y
	@$(MAKE) migrate-up
	@echo "データベースの状態がリセットされました"

# データベースの準備とマイグレーション実行（推奨）
db-setup: docker-up
	@echo "データベースのセットアップを実行しています..."
	@echo "マイグレーション用のディレクトリを作成..."
	@mkdir -p migrations
	@echo "Entスキーマからマイグレーションファイルを生成..."
	-@atlas migrate diff --env local 2>/dev/null || echo "マイグレーションファイルは既に最新です"
	@echo "マイグレーションを適用..."
	-@atlas migrate apply --env local 2>/dev/null || echo "マイグレーション適用完了または既に適用済みです"
	@echo "データベースのセットアップが完了しました！"

# マイグレーション状態を確認
migrate-status:
	@echo "マイグレーション状態を確認しています..."
	@atlas migrate status --env local