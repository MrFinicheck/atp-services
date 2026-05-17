import type {
  Client,
  CloseShiftResult,
  DashboardStats,
  Order,
  Tariff,
  User,
  Vehicle,
  VehicleScheduleItem,
  Waybill,
} from '../types';

const TOKEN_KEY = 'atp_session_token';

export function getToken(): string {
  return localStorage.getItem(TOKEN_KEY) || '';
}

export function setToken(token: string) {
  localStorage.setItem(TOKEN_KEY, token);
}

export function clearToken() {
  localStorage.removeItem(TOKEN_KEY);
}

function isWails(): boolean {
  return !!(window as any).go?.main?.App;
}

async function wailsCall<T>(method: string, ...args: unknown[]): Promise<T> {
  const app = (window as any).go.main.App;
  return app[method](...args) as T;
}

async function httpCall<T>(path: string, init?: RequestInit): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(init?.headers as Record<string, string>),
  };
  const token = getToken();
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await fetch(path, { ...init, headers });
  if (!res.ok) {
    const err = await res.json().catch(() => ({ message: res.statusText }));
    throw new Error(err.message || 'Request failed');
  }
  return res.json();
}

export const api = {
  async login(login: string, password: string) {
    if (isWails()) {
      const { Login } = await import('../../wailsjs/go/main/App');
      const wailsModels = await import('../../wailsjs/go/models');
      const res = await Login(new wailsModels.models.LoginRequest({ login, password }));
      return { token: res.token, user: res.user as User };
    }
    return httpCall<{ token: string; user: User }>('/api/login', {
      method: 'POST',
      body: JSON.stringify({ login, password }),
    });
  },

  async logout() {
    const token = getToken();
    if (isWails()) return wailsCall<void>('Logout', token);
    return httpCall('/api/logout', { method: 'POST' });
  },

  async me(): Promise<User> {
    const token = getToken();
    if (isWails()) return wailsCall('Me', token);
    return httpCall('/api/me');
  },

  async listClients(): Promise<Client[]> {
    const token = getToken();
    if (isWails()) return wailsCall('ListClients', token);
    return httpCall('/api/clients');
  },

  async listVehicles(): Promise<Vehicle[]> {
    const token = getToken();
    if (isWails()) return wailsCall('ListVehicles', token);
    return httpCall('/api/vehicles');
  },

  async listTariffs(): Promise<Tariff[]> {
    const token = getToken();
    if (isWails()) return wailsCall('ListTariffs', token);
    return httpCall('/api/tariffs');
  },

  async listUsers(): Promise<User[]> {
    const token = getToken();
    if (isWails()) return wailsCall('ListUsers', token);
    return httpCall('/api/users');
  },

  async listOrders(): Promise<Order[]> {
    const token = getToken();
    if (isWails()) return wailsCall('ListOrders', token);
    return httpCall('/api/orders');
  },

  async createOrder(req: Record<string, unknown>): Promise<Order> {
    const token = getToken();
    if (isWails()) return wailsCall('CreateOrder', token, req);
    return httpCall('/api/orders', { method: 'POST', body: JSON.stringify(req) });
  },

  async updateOrderStatus(orderId: string, status: string): Promise<Order> {
    const token = getToken();
    if (isWails()) return wailsCall('UpdateOrderStatus', token, orderId, status);
    return httpCall('/api/orders/status', {
      method: 'POST',
      body: JSON.stringify({ orderId, status }),
    });
  },

  async previewPrice(tariffId: string, distanceKm: number, idleHours: number, urgent: boolean): Promise<number> {
    const token = getToken();
    if (isWails()) return wailsCall('PreviewPrice', token, tariffId, distanceKm, idleHours, urgent);
    const r = await httpCall<{ price: number }>('/api/orders/preview-price', {
      method: 'POST',
      body: JSON.stringify({ tariffId, distanceKm, idleHours, urgent }),
    });
    return r.price;
  },

  async vehicleSchedule(): Promise<VehicleScheduleItem[]> {
    const token = getToken();
    if (isWails()) return wailsCall('VehicleSchedule', token);
    return httpCall('/api/schedule');
  },

  async closeShift(req: Record<string, unknown>): Promise<CloseShiftResult> {
    const token = getToken();
    if (isWails()) return wailsCall('CloseShift', token, req);
    return httpCall('/api/shift/close', { method: 'POST', body: JSON.stringify(req) });
  },

  async dashboard(): Promise<DashboardStats> {
    const token = getToken();
    if (isWails()) return wailsCall('Dashboard', token);
    return httpCall('/api/dashboard');
  },

  async driverRating(): Promise<Record<string, unknown>[]> {
    const token = getToken();
    if (isWails()) return wailsCall('DriverRating', token);
    return httpCall('/api/reports/drivers');
  },

  async saveTariff(t: Tariff): Promise<Tariff> {
    const token = getToken();
    if (isWails()) return wailsCall('SaveTariff', token, t);
    return httpCall('/api/tariffs', { method: 'POST', body: JSON.stringify(t) });
  },

  async saveVehicle(v: Vehicle): Promise<Vehicle> {
    const token = getToken();
    if (isWails()) return wailsCall('SaveVehicle', token, v);
    return httpCall('/api/vehicles', { method: 'POST', body: JSON.stringify(v) });
  },

  async listWaybills(): Promise<Waybill[]> {
    const token = getToken();
    if (isWails()) return wailsCall('ListWaybills', token);
    return httpCall('/api/waybills');
  },
};
