/// <reference types="vite" />

interface ImportMetaEnv {
  VITE_API_BASE_URL: string
  MODE: string
  DEV: boolean
  PROD: boolean
  BASE_URL: string
  NEXT_PUBLIC_API_URL?: string
}

interface ImportMeta {
  env: ImportMetaEnv
}

declare module 'node' {
  interface ProcessEnv extends ImportMetaEnv {}
}

export {}
