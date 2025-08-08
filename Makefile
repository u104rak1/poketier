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

mod-tidy: ## Go modulesを更新
	docker-compose exec poketier-backend go mod tidy

test: ## テストを実行
	docker-compose exec poketier-backend go test ./...

clean: ## 不要なDockerリソースを削除
	docker system prune -f
	docker volume prune -f
