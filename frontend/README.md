# Shadcn Admin Dashboard

Admin Dashboard UI crafted with Shadcn and Vite. Built with responsiveness and accessibility in mind.

![alt text](public/images/shadcn-admin.png)

I've been creating dashboard UIs at work and for my personal projects. I always wanted to make a reusable collection of dashboard UI for future projects; and here it is now. While I've created a few custom components, some of the code is directly adapted from ShadcnUI examples.

> This is not a starter project (template) though. I'll probably make one in the future.

## Features

- Light/dark mode
- Responsive
- Accessible
- With built-in Sidebar component
- Global Search Command
- 10+ pages
- Extra custom components

## Tech Stack

**UI:** [ShadcnUI](https://ui.shadcn.com) (TailwindCSS + RadixUI)

**Build Tool:** [Vite](https://vitejs.dev/)

**Routing:** [TanStack Router](https://tanstack.com/router/latest)

**Type Checking:** [TypeScript](https://www.typescriptlang.org/)

**Linting/Formatting:** [Eslint](https://eslint.org/) & [Prettier](https://prettier.io/)

**Icons:** [Tabler Icons](https://tabler.io/icons)

## Run Locally

Clone the project

```bash
  git clone https://github.com/satnaing/shadcn-admin.git
```

Go to the project directory

```bash
  cd shadcn-admin
```

Install dependencies

```bash
  pnpm install
```

Start the server

```bash
  pnpm run dev
```

## Author

Crafted with 🤍 by [@satnaing](https://github.com/satnaing)

## License

Licensed under the [MIT License](https://choosealicense.com/licenses/mit/)

## Google OAuth設定手順

1. [Google Cloud Console](https://console.cloud.google.com/)にアクセス
2. 新しいプロジェクトを作成または既存のプロジェクトを選択
3. 左側のメニューから「APIとサービス」→「認証情報」を選択
4. 「認証情報を作成」→「OAuth クライアントID」を選択
5. アプリケーションの種類で「ウェブアプリケーション」を選択
6. 以下の情報を入力：
   - 名前：「NEXT MART 2030 Admin」
   - 承認済みのJavaScript生成元：`http://localhost:3000`
   - 承認済みのリダイレクトURI：`http://localhost:3000/api/auth/callback/google`
7. 「作成」をクリックし、表示されるクライアントIDとクライアントシークレットを保存
8. `.env.local`ファイルに以下の情報を設定：
   ```
   GOOGLE_CLIENT_ID=取得したクライアントID
   GOOGLE_CLIENT_SECRET=取得したクライアントシークレット
   ALLOWED_EMAIL_DOMAINS=許可するメールドメイン
   ```

## 環境変数の設定

開発環境で必要な環境変数：

- `NEXTAUTH_URL`: アプリケーションのベースURL（開発環境では`http://localhost:3000`）
- `NEXTAUTH_SECRET`: NextAuth.jsの暗号化キー（`openssl rand -base64 32`で生成可能）
- `GOOGLE_CLIENT_ID`: Google OAuth クライアントID
- `GOOGLE_CLIENT_SECRET`: Google OAuth クライアントシークレット
- `ALLOWED_EMAIL_DOMAINS`: ログインを許可するメールドメイン（カンマ区切りで複数指定可能）

## セキュリティに関する注意

- `.env.local`ファイルはGitにコミットしないでください
- 本番環境では適切なドメインとリダイレクトURIを設定してください
- `NEXTAUTH_SECRET`は本番環境ごとに異なる値を使用してください
