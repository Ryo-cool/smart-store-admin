import { api } from './api';

export interface DeliveryQuery {
  page?: number;
  limit?: number;
  status?: string;
  search?: string;
}

export interface DeliveryResponse {
  deliveries: Delivery[];
  total: number;
}

export interface Delivery {
  id: string;
  deliveryType: string;
  address: string;
  estimatedDeliveryTime: string;
  actualDeliveryTime?: string;
  status: string;
  notes?: string;
  trackingInfo?: {
    currentLocation?: {
      latitude: number;
      longitude: number;
    };
    batteryLevel?: number;
    speed?: number;
  };
}

export interface DeliveryUpdateRequest {
  status?: string;
  notes?: string;
  actualDeliveryTime?: string;
  trackingInfo?: {
    currentLocation?: {
      latitude: number;
      longitude: number;
    };
    batteryLevel?: number;
    speed?: number;
  };
}

export interface DeliveryHistory {
  history: {
    status: string;
    timestamp: string;
    location?: {
      latitude: number;
      longitude: number;
    };
    note?: string;
  }[];
}

export const deliveriesApi = {
  getDeliveries: (query?: DeliveryQuery) =>
    api.get<DeliveryResponse>('/deliveries', { params: query }).then((res) => res.data),

  getDelivery: (id: string) =>
    api.get<Delivery>(`/deliveries/${id}`).then((res) => res.data),

  updateDelivery: (id: string, data: DeliveryUpdateRequest) =>
    api.patch<Delivery>(`/deliveries/${id}`, data).then((res) => res.data),

  updateDeliveryStatus: (id: string, status: string) =>
    api.patch<Delivery>(`/deliveries/${id}/status`, { status }).then((res) => res.data),

  getDeliveryHistory: (id: string) =>
    api.get<DeliveryHistory>(`/deliveries/${id}/history`).then((res) => res.data),
}; 