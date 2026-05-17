package models

import "time"

type Role string

const (
	RoleAdmin      Role = "admin"
	RoleDispatcher Role = "dispatcher"
	RoleDriver     Role = "driver"
)

type User struct {
	ID           string `json:"id"`
	Login        string `json:"login"`
	PasswordHash string `json:"-"`
	Role         Role   `json:"role"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Phone        string `json:"phone"`
	Active       bool   `json:"active"`
}

type Client struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Phone     string  `json:"phone"`
	DebtLimit float64 `json:"debtLimit"`
}

type Vehicle struct {
	ID       string  `json:"id"`
	Plate    string  `json:"plate"`
	Model    string  `json:"model"`
	FuelNorm float64 `json:"fuelNorm"`
	Active   bool    `json:"active"`
}

type Tariff struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	BaseFee         float64 `json:"baseFee"`
	PricePerKm      float64 `json:"pricePerKm"`
	PricePerIdleHr  float64 `json:"pricePerIdleHr"`
	UrgencyCoeff    float64 `json:"urgencyCoeff"`
	Active          bool    `json:"active"`
}

type OrderStatus string

const (
	OrderNew        OrderStatus = "new"
	OrderAssigned   OrderStatus = "assigned"
	OrderInProgress OrderStatus = "in_progress"
	OrderCompleted  OrderStatus = "completed"
	OrderCancelled  OrderStatus = "cancelled"
)

type Order struct {
	ID          string      `json:"id"`
	ClientID    string      `json:"clientId"`
	VehicleID   string      `json:"vehicleId"`
	DriverID    string      `json:"driverId"`
	FromAddr    string      `json:"fromAddr"`
	ToAddr      string      `json:"toAddr"`
	DistanceKm  float64     `json:"distanceKm"`
	IdleHours   float64     `json:"idleHours"`
	Urgent      bool        `json:"urgent"`
	TariffID    string      `json:"tariffId"`
	Price       float64     `json:"price"`
	Status      OrderStatus `json:"status"`
	ScheduledAt time.Time   `json:"scheduledAt"`
	CreatedAt   time.Time   `json:"createdAt"`
	StartedAt   *time.Time  `json:"startedAt,omitempty"`
	CompletedAt *time.Time  `json:"completedAt,omitempty"`
}

type Waybill struct {
	ID                string     `json:"id"`
	DriverID          string     `json:"driverId"`
	VehicleID         string     `json:"vehicleId"`
	Date              string     `json:"date"`
	StartOdometer     float64    `json:"startOdometer"`
	EndOdometer       float64    `json:"endOdometer"`
	FuelStart         float64    `json:"fuelStart"`
	FuelEnd           float64    `json:"fuelEnd"`
	FuelRefilled      float64    `json:"fuelRefilled"`
	ActualConsumption float64    `json:"actualConsumption"`
	NormConsumption   float64    `json:"normConsumption"`
	OverPercent       float64    `json:"overPercent"`
	Comment           string     `json:"comment"`
	Closed            bool       `json:"closed"`
	CreatedAt         time.Time  `json:"createdAt"`
	ClosedAt          *time.Time `json:"closedAt,omitempty"`
}

type AuditEntry struct {
	ID         string    `json:"id"`
	UserID     string    `json:"userId"`
	Action     string    `json:"action"`
	EntityType string    `json:"entityType"`
	EntityID   string    `json:"entityId"`
	Details    string    `json:"details"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Session struct {
	Token     string    `json:"token"`
	UserID    string    `json:"userId"`
	Role      Role      `json:"role"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type CreateOrderRequest struct {
	ClientID    string  `json:"clientId"`
	VehicleID   string  `json:"vehicleId"`
	DriverID    string  `json:"driverId"`
	FromAddr    string  `json:"fromAddr"`
	ToAddr      string  `json:"toAddr"`
	DistanceKm  float64 `json:"distanceKm"`
	IdleHours   float64 `json:"idleHours"`
	Urgent      bool    `json:"urgent"`
	TariffID    string  `json:"tariffId"`
	ScheduledAt string  `json:"scheduledAt"`
}

type CloseShiftRequest struct {
	VehicleID     string  `json:"vehicleId"`
	StartOdometer float64 `json:"startOdometer"`
	EndOdometer   float64 `json:"endOdometer"`
	FuelStart     float64 `json:"fuelStart"`
	FuelEnd       float64 `json:"fuelEnd"`
	FuelRefilled  float64 `json:"fuelRefilled"`
	Comment       string  `json:"comment"`
}

type CloseShiftResult struct {
	Waybill       Waybill `json:"waybill"`
	Blocked       bool    `json:"blocked"`
	RequireComment bool   `json:"requireComment"`
	Message       string  `json:"message"`
}

type DashboardStats struct {
	OrdersToday      int     `json:"ordersToday"`
	ActiveVehicles   int     `json:"activeVehicles"`
	OpenWaybills     int     `json:"openWaybills"`
	RevenueMonth     float64 `json:"revenueMonth"`
	FuelOverruns     int     `json:"fuelOverruns"`
}

type VehicleScheduleItem struct {
	VehicleID string  `json:"vehicleId"`
	Plate     string  `json:"plate"`
	Orders    []Order `json:"orders"`
}
