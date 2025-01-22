export type Role = 'admin' | 'staff' | 'viewer'

export interface User {
  id: string
  email: string | null
  name: string | null
  picture: string | null
  role: Role
}

export interface Session {
  user: User
  expires: string
}

export interface AuthResponse {
  token: string
  user: User
} 