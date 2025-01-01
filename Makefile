.PHONY: all frontend dev down

# デフォルトターゲット：すべてのサービスを起動
all: dev

# フロントエンドサーバーの起動
frontend:
	cd frontend && npm run dev

# Docker Composeでの開発環境起動
dev:
	docker compose up -d

# Docker Compose環境の停止
down:
	docker compose down

# ヘルプコマンド
help:
	@echo "使用可能なコマンド:"
	@echo "  make all      - Docker Composeですべてのサービスを起動"
	@echo "  make frontend - フロントエンドサーバーをローカルで起動"
	@echo "  make dev      - Docker Composeで開発環境を起動"
	@echo "  make down     - Docker Compose環境を停止" 