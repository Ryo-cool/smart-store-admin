/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL: string
  readonly MODE: 'development' | 'production'
  readonly DEV: boolean
  readonly PROD: boolean
  readonly BASE_URL: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}

declare module '@types/vite/client' {
  interface ImportMetaEnv {
    readonly VITE_API_BASE_URL: string
    readonly MODE: 'development' | 'production'
    readonly DEV: boolean
    readonly PROD: boolean
    readonly BASE_URL: string
  }
}

export {} 