// Optional: テストで使用するグローバルなマッチャーを追加
import '@testing-library/jest-dom'

// Optional: テスト環境のグローバル設定
jest.setTimeout(10000) // タイムアウトを10秒に設定

// Optional: フェッチのモック
global.fetch = jest.fn()

// Optional: IntersectionObserverのモック
global.IntersectionObserver = class IntersectionObserver {
  constructor() {}
  observe() {
    return null
  }
  unobserve() {
    return null
  }
  disconnect() {
    return null
  }
} 