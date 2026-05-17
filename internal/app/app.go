package app

import (
	"atp-services/internal/models"
)

// Application — фасад use-case слоя для HTTP и Wails (SRP: только делегирование).
type Application struct {
	ctr *Container
}

func New(dataDir string) (*Application, error) {
	ctr := NewContainer(dataDir)
	if err := ctr.Init(); err != nil {
		return nil, err
	}
	return &Application{ctr: ctr}, nil
}

func NewLazy(dataDir string) *Application {
	return &Application{ctr: NewContainer(dataDir)}
}

func (a *Application) EnsureReady() error {
	return a.ctr.Init()
}

func (a *Application) Close() error {
	return a.ctr.Close()
}

func (a *Application) DataDir() string {
	return a.ctr.DataDir()
}

func (a *Application) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	if err := a.ctr.Init(); err != nil {
		return nil, err
	}
	return a.ctr.Auth.Login(req)
}

// LoginWithCredentials — надёжный вход для Wails (два string вместо struct).
func (a *Application) LoginWithCredentials(login, password string) (*models.LoginResponse, error) {
	return a.Login(models.LoginRequest{Login: login, Password: password})
}

func (a *Application) Logout(token string) error {
	if err := a.ctr.Init(); err != nil {
		return err
	}
	return a.ctr.Auth.Logout(token)
}

func (a *Application) Me(token string) (*models.User, error) {
	if err := a.ctr.Init(); err != nil {
		return nil, err
	}
	return a.ctr.Auth.Validate(token)
}

func (a *Application) ListClients(token string) ([]models.Client, error) {
	if _, err := a.requireUser(token); err != nil {
		return nil, err
	}
	return a.ctr.Store().ListClients()
}

func (a *Application) SaveClient(token string, c models.Client) (*models.Client, error) {
	u, err := a.requireUser(token)
	if err != nil {
		return nil, err
	}
	if u.Role != models.RoleAdmin && u.Role != models.RoleDispatcher {
		return nil, errAccessDenied()
	}
	if err := a.ctr.Store().SaveClient(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (a *Application) ListVehicles(token string) ([]models.Vehicle, error) {
	if _, err := a.requireUser(token); err != nil {
		return nil, err
	}
	return a.ctr.Store().ListVehicles()
}

func (a *Application) SaveVehicle(token string, v models.Vehicle) (*models.Vehicle, error) {
	u, err := a.requireUser(token)
	if err != nil {
		return nil, err
	}
	if u.Role != models.RoleAdmin {
		return nil, errAccessDenied()
	}
	if err := a.ctr.Store().SaveVehicle(&v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (a *Application) ListTariffs(token string) ([]models.Tariff, error) {
	if _, err := a.requireUser(token); err != nil {
		return nil, err
	}
	return a.ctr.Tariffs.List()
}

func (a *Application) SaveTariff(token string, t models.Tariff) (*models.Tariff, error) {
	u, err := a.requireUser(token)
	if err != nil {
		return nil, err
	}
	if u.Role != models.RoleAdmin {
		return nil, errAccessDenied()
	}
	if err := a.ctr.Tariffs.Save(&t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (a *Application) ListUsers(token string) ([]models.User, error) {
	u, err := a.requireUser(token)
	if err != nil {
		return nil, err
	}
	if u.Role != models.RoleAdmin {
		return nil, errAccessDenied()
	}
	users, err := a.ctr.Store().ListUsers()
	if err != nil {
		return nil, err
	}
	for i := range users {
		users[i].PasswordHash = ""
	}
	return users, nil
}

func (a *Application) CreateUser(token string, u models.User, password string) (*models.User, error) {
	actor, err := a.requireUser(token)
	if err != nil {
		return nil, err
	}
	if actor.Role != models.RoleAdmin {
		return nil, errAccessDenied()
	}
	hash, err := a.ctr.Auth.HashPassword(password)
	if err != nil {
		return nil, err
	}
	u.PasswordHash = hash
	u.Active = true
	if err := a.ctr.Store().SaveUser(&u); err != nil {
		return nil, err
	}
	u.PasswordHash = ""
	return &u, nil
}

func (a *Application) ListOrders(token string) ([]models.Order, error) {
	u, err := a.requireUser(token)
	if err != nil {
		return nil, err
	}
	return a.ctr.Orders.ListForRole(u)
}

func (a *Application) CreateOrder(token string, req models.CreateOrderRequest) (*models.Order, error) {
	u, err := a.requireUser(token)
	if err != nil {
		return nil, err
	}
	if u.Role != models.RoleAdmin && u.Role != models.RoleDispatcher {
		return nil, errAccessDenied()
	}
	return a.ctr.Orders.Create(req, u.ID)
}

func (a *Application) UpdateOrderStatus(token, orderID, status string) (*models.Order, error) {
	u, err := a.requireUser(token)
	if err != nil {
		return nil, err
	}
	return a.ctr.Orders.UpdateStatus(orderID, models.OrderStatus(status), u.ID)
}

func (a *Application) PreviewPrice(token, tariffID string, distanceKm, idleHours float64, urgent bool) (float64, error) {
	if _, err := a.requireUser(token); err != nil {
		return 0, err
	}
	return a.ctr.Orders.PreviewPrice(tariffID, distanceKm, idleHours, urgent)
}

func (a *Application) VehicleSchedule(token string) ([]models.VehicleScheduleItem, error) {
	u, err := a.requireUser(token)
	if err != nil {
		return nil, err
	}
	if u.Role == models.RoleDriver {
		return nil, errAccessDenied()
	}
	return a.ctr.Orders.ScheduleToday()
}

func (a *Application) CloseShift(token string, req models.CloseShiftRequest) (*models.CloseShiftResult, error) {
	u, err := a.requireUser(token)
	if err != nil {
		return nil, err
	}
	if u.Role != models.RoleDriver && u.Role != models.RoleAdmin {
		return nil, errAccessDenied()
	}
	return a.ctr.Waybill.CloseShift(u.ID, req)
}

func (a *Application) ListWaybills(token string) ([]models.Waybill, error) {
	u, err := a.requireUser(token)
	if err != nil {
		return nil, err
	}
	if u.Role != models.RoleAdmin && u.Role != models.RoleDispatcher {
		return nil, errAccessDenied()
	}
	return a.ctr.Waybill.List()
}

func (a *Application) Dashboard(token string) (*models.DashboardStats, error) {
	u, err := a.requireUser(token)
	if err != nil {
		return nil, err
	}
	if u.Role == models.RoleDriver {
		return nil, errAccessDenied()
	}
	return a.ctr.Report.Dashboard()
}

func (a *Application) DriverRating(token string) ([]map[string]any, error) {
	u, err := a.requireUser(token)
	if err != nil {
		return nil, err
	}
	if u.Role != models.RoleAdmin {
		return nil, errAccessDenied()
	}
	return a.ctr.Report.DriverRating()
}

func (a *Application) ListAudit(token string) ([]models.AuditEntry, error) {
	u, err := a.requireUser(token)
	if err != nil {
		return nil, err
	}
	if u.Role != models.RoleAdmin {
		return nil, errAccessDenied()
	}
	return a.ctr.Store().ListAudit(100)
}

func (a *Application) requireUser(token string) (*models.User, error) {
	if err := a.ctr.Init(); err != nil {
		return nil, err
	}
	return a.ctr.Auth.Validate(token)
}

func errAccessDenied() error {
	return &AppError{Code: "access_denied", Message: "Недостаточно прав"}
}

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string { return e.Message }
