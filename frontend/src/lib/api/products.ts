import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

export interface Product {
  id: string;
  name: string;
  sku: string;
  price: number;
  stock: number;
  status: '販売中' | '在庫切れ' | '入荷待ち' | '在庫少';
  description?: string;
  category?: string;
  weight?: string;
  dimensions?: string;
  images?: string[];
  createdAt: string;
  updatedAt: string;
}

export interface ProductsResponse {
  products: Product[];
  total: number;
  page: number;
  perPage: number;
}

export interface ProductsQuery {
  page?: number;
  perPage?: number;
  search?: string;
  category?: string;
  status?: string;
  minPrice?: number;
  maxPrice?: number;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
}

export const productsApi = {
  getProducts: async (query: ProductsQuery = {}): Promise<ProductsResponse> => {
    const response = await axios.get(`${API_BASE_URL}/products`, {
      params: query,
    });
    return response.data;
  },

  getProduct: async (id: string): Promise<Product> => {
    const response = await axios.get(`${API_BASE_URL}/products/${id}`);
    return response.data;
  },

  createProduct: async (product: Omit<Product, 'id' | 'createdAt' | 'updatedAt'>): Promise<Product> => {
    const response = await axios.post(`${API_BASE_URL}/products`, product);
    return response.data;
  },

  updateProduct: async (id: string, product: Partial<Product>): Promise<Product> => {
    const response = await axios.patch(`${API_BASE_URL}/products/${id}`, product);
    return response.data;
  },

  deleteProduct: async (id: string): Promise<void> => {
    await axios.delete(`${API_BASE_URL}/products/${id}`);
  },
}; 