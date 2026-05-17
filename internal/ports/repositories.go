package ports

import "atp-services/internal/models"

// UserRepository — учётные записи (ISP).
type UserRepository interface {
	SaveUser(u *models.User) error
	FindUserByLogin(login string) (*models.User, error)
	FindUserByID(id string) (*models.User, error)
	ListUsers() ([]models.User, error)
}

// SessionRepository — сессии авторизации (ISP).
type SessionRepository interface {
	SaveSession(sess *models.Session) error
	FindSession(token string) (*models.Session, error)
	DeleteSession(token string) error
}

// ClientRepository — справочник клиентов.
type ClientRepository interface {
	SaveClient(c *models.Client) error
	ListClients() ([]models.Client, error)
}

// VehicleRepository — автопарк.
type VehicleRepository interface {
	SaveVehicle(v *models.Vehicle) error
	FindVehicle(id string) (*models.Vehicle, error)
	ListVehicles() ([]models.Vehicle, error)
}

// TariffRepository — тарифы.
type TariffRepository interface {
	SaveTariff(t *models.Tariff) error
	FindTariff(id string) (*models.Tariff, error)
	ListTariffs() ([]models.Tariff, error)
}

// OrderRepository — заявки.
type OrderRepository interface {
	SaveOrder(o *models.Order) error
	FindOrder(id string) (*models.Order, error)
	ListOrders() ([]models.Order, error)
}

// WaybillRepository — путевые листы.
type WaybillRepository interface {
	SaveWaybill(w *models.Waybill) error
	ListWaybills() ([]models.Waybill, error)
	FindOpenWaybill(driverID, date string) (*models.Waybill, error)
}

// AuditRepository — журнал аудита.
type AuditRepository interface {
	AddAudit(userID, action, entityType, entityID, details string) error
	ListAudit(limit int) ([]models.AuditEntry, error)
}

// MetaRepository — служебные флаги БД.
type MetaRepository interface {
	IsSeeded() (bool, error)
	MarkSeeded() error
}

// UnitOfWork объединяет репозитории для инициализации (facade).
type UnitOfWork interface {
	UserRepository
	SessionRepository
	ClientRepository
	VehicleRepository
	TariffRepository
	OrderRepository
	WaybillRepository
	AuditRepository
	MetaRepository
	Close() error
}
