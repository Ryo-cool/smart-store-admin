import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface AuthResponse {
  token: string;
  user: {
    id: string;
    email: string;
    name: string;
    role: string;
  };
}

export const authApi = {
  login: async (credentials: LoginCredentials): Promise<AuthResponse> => {
    const response = await axios.post(`${API_BASE_URL}/auth/login`, credentials);
    return response.data;
  },

  logout: async (): Promise<void> => {
    await axios.post(`${API_BASE_URL}/auth/logout`);
  },

  getCurrentUser: async (): Promise<AuthResponse['user']> => {
    const response = await axios.get(`${API_BASE_URL}/auth/me`);
    return response.data;
  },
}; 