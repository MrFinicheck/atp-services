export type Role = 'admin' | 'dispatcher' | 'driver';

export interface User {
  id: string;
  login: string;
  role: Role;
  firstName: string;
  lastName: string;
  phone: string;
  active: boolean;
}

export interface Client {
  id?: string;
  name: string;
  phone: string;
  debtLimit: number;
}

export interface Vehicle {
  id?: string;
  plate: string;
  model: string;
  fuelNorm: number;
  active: boolean;
}

export interface Tariff {
  id?: string;
  name: string;
  baseFee: number;
  pricePerKm: number;
  pricePerIdleHr: number;
  urgencyCoeff: number;
  active: boolean;
}

export interface Order {
  id: string;
  clientId: string;
  vehicleId: string;
  driverId: string;
  fromAddr: string;
  toAddr: string;
  distanceKm: number;
  idleHours: number;
  urgent: boolean;
  tariffId: string;
  price: number;
  status: string;
  scheduledAt: string;
  createdAt: string;
}

export interface Waybill {
  id: string;
  driverId: string;
  vehicleId: string;
  date: string;
  actualConsumption: number;
  normConsumption: number;
  overPercent: number;
  closed: boolean;
}

export interface DashboardStats {
  ordersToday: number;
  activeVehicles: number;
  openWaybills: number;
  revenueMonth: number;
  fuelOverruns: number;
}

export interface VehicleScheduleItem {
  vehicleId: string;
  plate: string;
  orders: Order[];
}

export interface CloseShiftResult {
  waybill: Waybill;
  blocked: boolean;
  requireComment: boolean;
  message: string;
}
