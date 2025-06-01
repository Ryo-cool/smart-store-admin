# CLAUDE.md

このファイルは、このリポジトリでコードを扱う際のClaude Code (claude.ai/code) へのガイダンスを提供します。

## プロジェクト概要

Smart Store Adminは、在庫管理、売上、配送、環境負荷追跡を含むスマートストア運営を管理するフルスタックアプリケーションです。

**アーキテクチャ:**
- **バックエンド**: Go（Echoフレームワーク）、MongoDB、JWT認証
- **フロントエンド**: Next.js（TypeScript）、TanStack Router、TanStack Query、Zustandによる状態管理
- **データベース**: MongoDB
- **インフラ**: 開発環境用Docker Compose

## 開発コマンド

### Makeコマンド（推奨）
共通タスクには以下のMakeコマンドを使用してください：
```bash
make dev          # Docker Composeですべてのサービスを起動
make down         # Docker Composeサービスを停止
make install      # フロントエンドの依存関係をインストール
make frontend     # フロントエンド開発サーバーをローカルで実行
make lint         # フロントエンドのlintチェック
make format       # Prettierでフロントエンドコードをフォーマット
make type-check   # TypeScriptの型チェック
make test         # フロントエンドテストを実行
make test-backend # Goバックエンドテストを実行
make build        # フロントエンドの本番ビルド
make check        # すべての品質チェックを実行（lint + type-check + test）
make clean        # 生成されたファイルを削除
```

### フロントエンド（Next.js）
```bash
cd frontend
npm run dev              # 開発サーバー
npm run build            # 本番ビルド
npm run lint             # ESLintチェック
npm run lint:fix         # ESLintの問題を修正
npm run format           # Prettierでフォーマット
npm run type-check       # TypeScriptチェック
npm run test             # Jestテスト
npm run test:watch       # Jestをウォッチモードで実行
npm run test:coverage    # テストカバレッジレポート
```

### バックエンド（Go）
```bash
cd backend
go run main.go           # 開発サーバーを起動
go test ./... -v         # 詳細出力ですべてのテストを実行
go mod tidy              # 依存関係をクリーンアップ
```

## コードアーキテクチャ

### バックエンド構造
- **リポジトリパターン**: `repository/`にインターフェースを持つデータアクセス層
- **サービス層**: `service/`のビジネスロジックと対応するテスト
- **ハンドラー層**: APIエンドポイント用の`handler/`のHTTPハンドラー
- **モデル**: MongoDB BSOMタグを使用した`models/`のデータベースモデル
- **ミドルウェア**: 認証・認可ミドルウェア
- **ルーター**: Echoフレームワークグループを使用したAPIルーティング設定

### フロントエンド構造
- **TanStack Router**: `routes/`でのファイルベースルーティングとコード生成
- **状態管理**: 認証とアプリ状態用の`stores/`のZustandストア
- **API層**: インターセプター付きの`lib/api/`のAxiosベースAPIクライアント
- **コンポーネント**: `components/ui/`のshadcn/uiを使用した再利用可能なUIコンポーネント
- **機能**: 分離されたコンポーネントを持つ`features/`の機能ベース組織
- **フック**: `hooks/`のカスタムReactフック

### 主要パターン
- **エラーハンドリング**: トースト通知による集約エラーハンドリング
- **認証**: 自動トークンハンドリング付きJWTベース認証
- **型安全性**: Zodバリデーション付き厳密なTypeScript設定
- **テスト**: フロントエンドはJest、バックエンドはGoのテストパッケージ（モック付き）

### API統合
- フロントエンドはサーバー状態管理にTanStack Queryを使用
- Axiosインターセプターが認証トークンとエラーレスポンスを処理
- `VITE_API_BASE_URL`環境変数で設定可能なベースAPI URL
- バックエンドは`/api`でロールベースアクセス制御付きREST APIを提供

### 環境設定
- バックエンドはポート8080、MongoDBはポート27017で動作
- Docker Composeが完全な開発環境をセットアップ
- フロントエンドは設定可能なAPI Base URL経由でバックエンドに接続