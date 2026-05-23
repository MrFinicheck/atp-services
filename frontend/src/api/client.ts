import { asArray, normalizeOrders, normalizeSchedule } from './normalize';
import type {
  Client,
  CloseShiftResult,
  ShiftStatus,
  OpenShiftResult,
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
  return !!(window as any)?.go?.main?.App;
}

function wailsError(err: unknown): string {
  if (typeof err === 'string') return err;
  if (err instanceof Error) return err.message;
  if (err && typeof err === 'object' && 'message' in err) {
    return String((err as { message: unknown }).message);
  }
  return 'Неизвестная ошибка';
}

function pickToken(res: Record<string, unknown>): string {
  const t = res.token ?? res.Token;
  return typeof t === 'string' ? t : '';
}

function pickUser(res: Record<string, unknown>): User {
  const u = (res.user ?? res.User) as User;
  return u;
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
  const body = await res.json().catch(() => ({}));
  if (!res.ok) {
    throw new Error((body as { message?: string }).message || res.statusText);
  }
  return body as T;
}

export const api = {
  async login(login: string, password: string) {
    const trimmedLogin = login.trim();
    const trimmedPass = password;

    if (isWails()) {
      const { LoginWithCredentials } = await import('../../wailsjs/go/main/App');
      try {
        const res = (await LoginWithCredentials(trimmedLogin, trimmedPass)) as unknown as Record<string, unknown>;
        const token = pickToken(res);
        if (!token) {
          throw new Error('Сервер не вернул токен сессии');
        }
        return { token, user: pickUser(res) };
      } catch (e) {
        // Fallback: старый метод со struct
        try {
          const { Login } = await import('../../wailsjs/go/main/App');
          const res = (await Login({ login: trimmedLogin, password: trimmedPass })) as unknown as Record<string, unknown>;
          const token = pickToken(res);
          if (!token) throw new Error('Сервер не вернул токен сессии');
          return { token, user: pickUser(res) };
        } catch {
          throw new Error(wailsError(e));
        }
      }
    }

    return httpCall<{ token: string; user: User }>('/api/login', {
      method: 'POST',
      body: JSON.stringify({ login: trimmedLogin, password: trimmedPass }),
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
    const raw = isWails()
      ? await wailsCall<unknown>('ListClients', token)
      : await httpCall<unknown>('/api/clients');
    return asArray<Client>(raw);
  },

  async saveClient(c: Client): Promise<Client> {
    const token = getToken();
    const payload = {
      name: String(c.name ?? '').trim(),
      phone: String(c.phone ?? '').trim(),
      debtLimit: Number(c.debtLimit ?? 0),
    };
    if (isWails()) {
      const { SaveClient } = await import('../../wailsjs/go/main/App');
      const { models } = await import('../../wailsjs/go/models');
      return SaveClient(token, new models.Client(payload));
    }
    return httpCall('/api/clients', { method: 'POST', body: JSON.stringify(payload) });
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

  async createUser(
    user: Pick<User, 'login' | 'role' | 'firstName' | 'lastName' | 'phone'> & { active?: boolean },
    password: string
  ): Promise<User> {
    const token = getToken();
    const payload = { ...user, active: user.active !== false };
    if (isWails()) return wailsCall('CreateUser', token, payload, password);
    return httpCall('/api/users', {
      method: 'POST',
      body: JSON.stringify({ user: payload, password }),
    });
  },

  async deleteUser(userId: string): Promise<void> {
    const token = getToken();
    if (isWails()) return wailsCall('DeleteUser', token, userId);
    await httpCall(`/api/users?id=${encodeURIComponent(userId)}`, { method: 'DELETE' });
  },

  async listOrders(): Promise<Order[]> {
    const token = getToken();
    const raw = isWails()
      ? await wailsCall<unknown>('ListOrders', token)
      : await httpCall<unknown>('/api/orders');
    return normalizeOrders(raw);
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
    let raw: unknown;
    if (isWails()) raw = await wailsCall('VehicleSchedule', token);
    else raw = await httpCall('/api/schedule');
    return normalizeSchedule(raw);
  },

  async shiftStatus(): Promise<ShiftStatus> {
    const token = getToken();
    if (isWails()) return wailsCall('ShiftStatus', token);
    return httpCall('/api/shift/status');
  },

  async listDriversAvailable(): Promise<User[]> {
    const token = getToken();
    const raw = isWails()
      ? await wailsCall<unknown>('ListDriversAvailable', token)
      : await httpCall<unknown>('/api/drivers/available');
    return asArray<User>(raw);
  },

  async openShift(req: Record<string, unknown>): Promise<OpenShiftResult> {
    const token = getToken();
    const payload = {
      vehicleId: String(req.vehicleId ?? ''),
      startOdometer: Number(req.startOdometer),
      fuelStart: Number(req.fuelStart),
    };
    if (isWails()) {
      const { OpenShift } = await import('../../wailsjs/go/main/App');
      const { models } = await import('../../wailsjs/go/models');
      return OpenShift(token, new models.OpenShiftRequest(payload));
    }
    return httpCall('/api/shift/open', { method: 'POST', body: JSON.stringify(payload) });
  },

  async closeShift(req: Record<string, unknown>): Promise<CloseShiftResult> {
    const token = getToken();
    const payload = {
      vehicleId: '',
      startOdometer: 0,
      endOdometer: Number(req.endOdometer),
      fuelStart: 0,
      fuelEnd: Number(req.fuelEnd),
      fuelRefilled: Number(req.fuelRefilled ?? 0),
      comment: String(req.comment ?? ''),
    };
    if (isWails()) {
      const { CloseShift } = await import('../../wailsjs/go/main/App');
      const { models } = await import('../../wailsjs/go/models');
      return CloseShift(token, new models.CloseShiftRequest(payload));
    }
    return httpCall('/api/shift/close', { method: 'POST', body: JSON.stringify(payload) });
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

export { wailsError };
