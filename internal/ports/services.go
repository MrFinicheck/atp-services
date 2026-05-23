package ports

import "atp-services/internal/models"

// AuthService — аутентификация (DIP: зависит от репозиториев, не от LevelDB).
type AuthService interface {
	Login(req models.LoginRequest) (*models.LoginResponse, error)
	Logout(token string) error
	Validate(token string) (*models.User, error)
	HashPassword(password string) (string, error)
}

// TariffCalculator — расчёт стоимости (SRP + Strategy).
type TariffCalculator interface {
	CalculatePrice(tariffID string, distanceKm, idleHours float64, urgent bool) (float64, error)
	List() ([]models.Tariff, error)
	Save(t *models.Tariff) error
}

// OrderService — заявки.
type OrderService interface {
	Create(req models.CreateOrderRequest, actorID string) (*models.Order, error)
	ListForRole(user *models.User) ([]models.Order, error)
	UpdateStatus(orderID string, status models.OrderStatus, actorID string) (*models.Order, error)
	ScheduleToday() ([]models.VehicleScheduleItem, error)
	PreviewPrice(tariffID string, distanceKm, idleHours float64, urgent bool) (float64, error)
}

// WaybillService — закрытие смены и контроль топлива.
type WaybillService interface {
	OpenShift(driverID string, req models.OpenShiftRequest) (*models.OpenShiftResult, error)
	CloseShift(driverID string, req models.CloseShiftRequest) (*models.CloseShiftResult, error)
	ShiftStatus(driverID string) (*models.ShiftStatus, error)
	DriverShiftOpen(driverID, date string) (bool, error)
	DriverShiftClosed(driverID, date string) (bool, error)
	List() ([]models.Waybill, error)
}

// ReportService — отчёты.
type ReportService interface {
	Dashboard() (*models.DashboardStats, error)
	DriverRating() ([]map[string]any, error)
}

// Seeder — начальные данные.
type Seeder interface {
	Seed() error
}
