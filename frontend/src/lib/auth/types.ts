export type Role = 'admin' | 'staff' | 'viewer'

export interface User {
  id: string
  email: string
  name: string
  picture: string
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