import { Role } from './types'

const roleHierarchy: Record<Role, number> = {
  admin: 3,
  staff: 2,
  viewer: 1,
}

export function hasRequiredRole(userRole: Role, requiredRole: Role): boolean {
  return roleHierarchy[userRole] >= roleHierarchy[requiredRole]
}

export const ROLE_BASED_ROUTES = {
  '/dashboard': 'viewer',
  '/products': 'staff',
  '/sales': 'staff',
  '/deliveries': 'staff',
  '/settings': 'admin',
} as const

export type ProtectedRoute = keyof typeof ROLE_BASED_ROUTES 