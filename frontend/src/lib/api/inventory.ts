import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

export interface InventoryHistory {
  id: string;
  productId: string;
  productName: string;
  type: '入庫' | '出庫' | '在庫調整';
  quantity: number;
  reason: string;
  note?: string;
  createdAt: string;
  createdBy: string;
}

export interface InventoryHistoryResponse {
  histories: InventoryHistory[];
  total: number;
  page: number;
  perPage: number;
}

export interface InventoryHistoryQuery {
  page?: number;
  perPage?: number;
  productId?: string;
  type?: string;
  startDate?: string;
  endDate?: string;
}

export interface InventoryUpdateRequest {
  productId: string;
  type: '入庫' | '出庫' | '在庫調整';
  quantity: number;
  reason: string;
  note?: string;
}

export const inventoryApi = {
  getHistories: async (query: InventoryHistoryQuery = {}): Promise<InventoryHistoryResponse> => {
    const response = await axios.get(`${API_BASE_URL}/inventory/histories`, {
      params: query,
    });
    return response.data;
  },

  getProductHistory: async (productId: string, query: Omit<InventoryHistoryQuery, 'productId'> = {}): Promise<InventoryHistoryResponse> => {
    const response = await axios.get(`${API_BASE_URL}/inventory/products/${productId}/histories`, {
      params: query,
    });
    return response.data;
  },

  updateInventory: async (data: InventoryUpdateRequest): Promise<void> => {
    await axios.post(`${API_BASE_URL}/inventory/update`, data);
  },
}; 