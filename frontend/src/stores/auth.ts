import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { authApi, type AuthResponse } from '@/lib/api/auth';

interface AuthState {
  token: string | null;
  user: AuthResponse['user'] | null;
  isAuthenticated: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  setUser: (user: AuthResponse['user']) => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      token: null,
      user: null,
      isAuthenticated: false,

      login: async (email: string, password: string) => {
        const response = await authApi.login({ email, password });
        set({
          token: response.token,
          user: response.user,
          isAuthenticated: true,
        });
      },

      logout: async () => {
        await authApi.logout();
        set({
          token: null,
          user: null,
          isAuthenticated: false,
        });
      },

      setUser: (user) => {
        set({
          user,
          isAuthenticated: true,
        });
      },
    }),
    {
      name: 'auth-storage',
    }
  )
); 