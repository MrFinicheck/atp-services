import type { Order, VehicleScheduleItem } from '../types';

export function asArray<T>(v: unknown): T[] {
  return Array.isArray(v) ? v : [];
}

export function normalizeSchedule(raw: unknown): VehicleScheduleItem[] {
  return asArray<Record<string, unknown>>(raw).map((s) => ({
    vehicleId: String(s.vehicleId ?? s.VehicleID ?? ''),
    plate: String(s.plate ?? s.Plate ?? '—'),
    orders: normalizeOrders(s.orders ?? s.Orders),
  }));
}

export function normalizeOrders(raw: unknown): Order[] {
  return asArray<Record<string, unknown>>(raw).map((o) => ({
    id: String(o.id ?? o.ID ?? ''),
    clientId: String(o.clientId ?? o.ClientID ?? ''),
    vehicleId: String(o.vehicleId ?? o.VehicleID ?? ''),
    driverId: String(o.driverId ?? o.DriverID ?? ''),
    fromAddr: String(o.fromAddr ?? o.FromAddr ?? ''),
    toAddr: String(o.toAddr ?? o.ToAddr ?? ''),
    distanceKm: Number(o.distanceKm ?? o.DistanceKm ?? 0),
    idleHours: Number(o.idleHours ?? o.IdleHours ?? 0),
    urgent: Boolean(o.urgent ?? o.Urgent),
    tariffId: String(o.tariffId ?? o.TariffID ?? ''),
    price: Number(o.price ?? o.Price ?? 0),
    status: String(o.status ?? o.Status ?? 'new'),
    scheduledAt: String(o.scheduledAt ?? o.ScheduledAt ?? ''),
    createdAt: String(o.createdAt ?? o.CreatedAt ?? ''),
  }));
}
