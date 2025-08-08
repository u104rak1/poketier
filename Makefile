.PHONY: help build up down logs shell test clean

help: ## ヘルプを表示
	@grep -E '(^##|^[a-zA-Z_-]+:.*?##)' $(MAKEFILE_LIST) | \
		awk '/^##/ {print substr($$0, 4)} /^[a-zA-Z_-]+:/ {split($$0, a, ":.*?## "); printf "\033[36m%-30s\033[0m %s\n", a[1], a[2]}'

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

clean: ## 不要なDockerリソースを削除
	docker system prune -f
	docker volume prune -f
