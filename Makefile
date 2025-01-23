.PHONY: all frontend dev down lint format type-check test clean install help build

# デフォルトターゲット：すべてのサービスを起動
all: dev

# フロントエンドの依存関係インストール
install:
	cd frontend && npm install

# フロントエンドサーバーの起動
frontend:
	cd frontend && npm run dev

# Docker Composeでの開発環境起動
dev:
	docker compose up -d

# Docker Compose環境の停止
down:
	docker compose down

# フロントエンドのlintチェック
lint:
	cd frontend && npm run lint

# フロントエンドのコード整形
format:
	cd frontend && npm run format

# フロントエンドの型チェック
type-check:
	cd frontend && npm run type-check

# フロントエンドのテスト実行
test:
	cd frontend && npm run test

# バックエンドのテスト実行
test-backend:
	cd backend && go test ./... -v

# フロントエンドのビルド
build:
	cd frontend && npm run build

# 生成されたファイルの削除
clean:
	cd frontend && rm -rf .next out
	cd backend && go clean

# すべての品質チェックを実行
check: lint type-check test

# ヘルプコマンド
help:
	@echo "使用可能なコマンド:"
	@echo "  make all         - Docker Composeですべてのサービスを起動"
	@echo "  make install     - フロントエンドの依存関係をインストール"
	@echo "  make frontend    - フロントエンドサーバーをローカルで起動"
	@echo "  make dev        - Docker Composeで開発環境を起動"
	@echo "  make down       - Docker Compose環境を停止"
	@echo "  make lint       - フロントエンドのlintチェックを実行"
	@echo "  make format     - フロントエンドのコード整形を実行"
	@echo "  make type-check - フロントエンドの型チェックを実行"
	@echo "  make test       - フロントエンドのテストを実行"
	@echo "  make test-backend - バックエンドのテストを実行"
	@echo "  make build      - フロントエンドのビルドを実行"
	@echo "  make clean      - 生成されたファイルを削除"
	@echo "  make check      - すべての品質チェックを実行" 