.PHONY: help build up down logs shell test clean

help: ## ヘルプを表示
	@grep -E '(^##|^[a-zA-Z_-]+:.*?##)' $(MAKEFILE_LIST) | \
		awk '/^##/ {print substr($$0, 4)} /^[a-zA-Z_-]+:/ {split($$0, a, ":.*?## "); printf "\033[36m%-16s\033[0m %s\n", a[1], a[2]}'

build: ## Docker イメージをビルド
	docker-compose build

up: ## Docker コンテナを起動
	docker-compose up -d

down: ## Docker コンテナを停止・削除
	docker-compose down

logs: ## Docker バックエンドコンテナのログを表示
	docker-compose logs -f poketier-backend

shell: ## バックエンドコンテナにシェルでアクセス
	docker-compose exec poketier-backend sh

# アプリケーション制御コマンド
run: ## DelveデバッガーとAPIサーバーを同時起動
	docker-compose exec -d poketier-backend dlv debug ./cmd/poketier --headless --listen=:2345 --api-version=2 --accept-multiclient --continue
	sleep 2

stop: ## DelveとAPIサーバーを完全停止（コンテナは永続化維持）
	docker-compose exec poketier-backend pkill -f "go run" || true
	docker-compose exec poketier-backend pkill -f "__debug_bin" || true
	docker-compose exec poketier-backend pkill -f "poketier" || true
	docker-compose exec poketier-backend pkill -f "dlv" || true
	sleep 2

rerun: ## make stop → make run
	make stop && make run

health: ## バックエンドコンテナのヘルスチェック
	curl -f http://localhost:28080/health

mod-tidy: ## Go modulesを更新
	docker-compose exec poketier-backend go mod tidy

test: ## テストを実行
	docker-compose exec poketier-backend go test ./...

mockgen: ## モックを生成（例: make mockgen PATH=./pkg/vo/id/id.go）
	@if [ -z "$(PATH)" ]; then \
		echo "Error: PATH is required. Usage: make mockgen PATH=./path/to/file.go"; \
		exit 1; \
	fi; \
	if [ ! -f "backend/$(PATH)" ]; then \
		echo "Error: File backend/$(PATH) does not exist"; \
		exit 1; \
	fi; \
	PACKAGE_NAME=$$(basename $$(dirname $(PATH))); \
	OUTPUT_FILE=$$(dirname $(PATH))/$${PACKAGE_NAME}_mock_test.go; \
	docker-compose exec poketier-backend mockgen -source=$(PATH) -destination=$${OUTPUT_FILE} -package=$${PACKAGE_NAME}_test && \
	echo "Mock generated: $${OUTPUT_FILE}"

lint: ## golangci-lintでコードチェックを実行
	docker-compose exec poketier-backend golangci-lint run --config .golangci.json ./...

lint-fix: ## golangci-lintで自動修正可能な問題を修正
	docker-compose exec poketier-backend golangci-lint run --config .golangci.json --fix ./...

# データベース関連コマンド
migrate-up: ## マイグレーションを適用（UP）
	docker-compose exec poketier-backend migrate -path ./sqlc/migrations -database "postgresql://local_user:password@postgres:5432/POCGO_LOCAL_DB?sslmode=disable" up

migrate-down: ## マイグレーションを1つ戻す（DOWN）
	docker-compose exec poketier-backend migrate -path ./sqlc/migrations -database "postgresql://local_user:password@postgres:5432/POCGO_LOCAL_DB?sslmode=disable" down 1

migrate-force: ## マイグレーションバージョンを強制設定（例: make migrate-force VERSION=1）
	docker-compose exec poketier-backend migrate -path ./sqlc/migrations -database "postgresql://local_user:password@postgres:5432/POCGO_LOCAL_DB?sslmode=disable" force $(VERSION)

migrate-version: ## 現在のマイグレーションバージョンを表示
	docker-compose exec poketier-backend migrate -path ./sqlc/migrations -database "postgresql://local_user:password@postgres:5432/POCGO_LOCAL_DB?sslmode=disable" version

sqlc-generate: ## SQLCでGoコードを生成
	docker-compose exec poketier-backend sqlc generate -f ./sqlc/sqlc.json

sqlc-vet: ## SQLCで設定とクエリをチェック
	docker-compose exec poketier-backend sqlc vet -f ./sqlc/sqlc.json

# 開発用ショートカットコマンド
db-reset: ## データベースを初期化（DOWN→UP→SQLCコード生成）
	make migrate-down || true
	make migrate-up
	make sqlc-generate

db-setup: ## 初回データベースセットアップ（マイグレーション適用→SQLCコード生成）
	make migrate-up
	make sqlc-generate

clean: ## 不要なDockerリソースを削除
	docker system prune -f
	docker volume prune -f
