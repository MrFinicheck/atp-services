package app

import (
	"sync"

	"atp-services/internal/models"
	"atp-services/internal/services"
	"atp-services/internal/store"
)

type Application struct {
	mu      sync.RWMutex
	once    sync.Once
	initErr error
	store   *store.Store
	auth    *services.AuthService
	orders  *services.OrderService
	tariffs *services.TariffService
	waybill *services.WaybillService
	report  *services.ReportService
	dataDir string
}

func New(dataDir string) (*Application, error) {
	a := &Application{dataDir: resolveDataDir(dataDir)}
	if err := a.init(); err != nil {
		return nil, err
	}
	return a, nil
}

// NewLazy defers database open until the first API call (used by Wails startup).
func NewLazy(dataDir string) *Application {
	return &Application{dataDir: resolveDataDir(dataDir)}
}

func (a *Application) init() error {
	a.once.Do(func() {
		st, err := store.Open(a.dataDir)
		if err != nil {
			a.initErr = err
			return
		}
		auth := services.NewAuthService(st)
		tariffs := services.NewTariffService(st)
		a.store = st
		a.auth = auth
		a.tariffs = tariffs
		a.orders = services.NewOrderService(st, tariffs)
		a.waybill = services.NewWaybillService(st)
		a.report = services.NewReportService(st)
		_ = services.SeedDemoData(st, auth)
	})
	return a.initErr
}

func (a *Application) ensure() error {
	return a.init()
}

// EnsureReady opens the database if it is not open yet (for Wails lazy startup).
func (a *Application) EnsureReady() error {
	return a.ensure()
}

func (a *Application) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.store == nil {
		return nil
	}
	err := a.store.Close()
	a.store = nil
	return err
}

func (a *Application) DataDir() string {
	return a.dataDir
}

func (a *Application) requireUser(token string) (*models.User, error) {
	if err := a.ensure(); err != nil {
		return nil, err
	}
	return a.auth.Validate(token)
}

func (a *Application) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	if err := a.ensure(); err != nil {
		return nil, err
	}
	return a.auth.Login(req)
}

func (a *Application) Logout(token string) error {
	if err := a.ensure(); err != nil {
		return err
	}
	return a.auth.Logout(token)
}

func (a *Application) Me(token string) (*models.User, error) {
	return a.requireUser(token)
}

func (a *Application) ListClients(token string) ([]models.Client, error) {
	if _, err := a.requireUser(token); err != nil {
		return nil, err
	}
	return a.store.ListClients()
}

func (a *Application) SaveClient(token string, c models.Client) (*models.Client, error) {
	u, err := a.requireUser(token)
	if err != nil {
		return nil, err
	}
	if u.Role != models.RoleAdmin && u.Role != models.RoleDispatcher {
		return nil, errAccessDenied()
	}
	if err := a.store.SaveClient(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (a *Application) ListVehicles(token string) ([]models.Vehicle, error) {
	if _, err := a.requireUser(token); err != nil {
		return nil, err
	}
	return a.store.ListVehicles()
}

func (a *Application) SaveVehicle(token string, v models.Vehicle) (*models.Vehicle, error) {
	u, err := a.requireUser(token)
	if err != nil {
		return nil, err
	}
	if u.Role != models.RoleAdmin {
		return nil, errAccessDenied()
	}
	if err := a.store.SaveVehicle(&v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (a *Application) ListTariffs(token string) ([]models.Tariff, error) {
	if _, err := a.requireUser(token); err != nil {
		return nil, err
	}
	return a.tariffs.List()
}

func (a *Application) SaveTariff(token string, t models.Tariff) (*models.Tariff, error) {
	u, err := a.requireUser(token)
	if err != nil {
		return nil, err
	}
	if u.Role != models.RoleAdmin {
		return nil, errAccessDenied()
	}
	if err := a.tariffs.Save(&t); err != nil {
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
	users, err := a.store.ListUsers()
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
	hash, err := a.auth.HashPassword(password)
	if err != nil {
		return nil, err
	}
	u.PasswordHash = hash
	u.Active = true
	if err := a.store.SaveUser(&u); err != nil {
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
	return a.orders.ListForRole(u)
}

func (a *Application) CreateOrder(token string, req models.CreateOrderRequest) (*models.Order, error) {
	u, err := a.requireUser(token)
	if err != nil {
		return nil, err
	}
	if u.Role != models.RoleAdmin && u.Role != models.RoleDispatcher {
		return nil, errAccessDenied()
	}
	return a.orders.Create(req, u.ID)
}

func (a *Application) UpdateOrderStatus(token, orderID, status string) (*models.Order, error) {
	u, err := a.requireUser(token)
	if err != nil {
		return nil, err
	}
	return a.orders.UpdateStatus(orderID, models.OrderStatus(status), u.ID)
}

func (a *Application) PreviewPrice(token, tariffID string, distanceKm, idleHours float64, urgent bool) (float64, error) {
	if _, err := a.requireUser(token); err != nil {
		return 0, err
	}
	return a.orders.PreviewPrice(tariffID, distanceKm, idleHours, urgent)
}

func (a *Application) VehicleSchedule(token string) ([]models.VehicleScheduleItem, error) {
	u, err := a.requireUser(token)
	if err != nil {
		return nil, err
	}
	if u.Role == models.RoleDriver {
		return nil, errAccessDenied()
	}
	return a.orders.ScheduleToday()
}

func (a *Application) CloseShift(token string, req models.CloseShiftRequest) (*models.CloseShiftResult, error) {
	u, err := a.requireUser(token)
	if err != nil {
		return nil, err
	}
	if u.Role != models.RoleDriver && u.Role != models.RoleAdmin {
		return nil, errAccessDenied()
	}
	driverID := u.ID
	if u.Role == models.RoleAdmin && req.VehicleID != "" {
		driverID = u.ID
	}
	return a.waybill.CloseShift(driverID, req)
}

func (a *Application) ListWaybills(token string) ([]models.Waybill, error) {
	u, err := a.requireUser(token)
	if err != nil {
		return nil, err
	}
	if u.Role != models.RoleAdmin && u.Role != models.RoleDispatcher {
		return nil, errAccessDenied()
	}
	return a.waybill.List()
}

func (a *Application) Dashboard(token string) (*models.DashboardStats, error) {
	u, err := a.requireUser(token)
	if err != nil {
		return nil, err
	}
	if u.Role == models.RoleDriver {
		return nil, errAccessDenied()
	}
	return a.report.Dashboard()
}

func (a *Application) DriverRating(token string) ([]map[string]any, error) {
	u, err := a.requireUser(token)
	if err != nil {
		return nil, err
	}
	if u.Role != models.RoleAdmin {
		return nil, errAccessDenied()
	}
	return a.report.DriverRating()
}

func (a *Application) ListAudit(token string) ([]models.AuditEntry, error) {
	u, err := a.requireUser(token)
	if err != nil {
		return nil, err
	}
	if u.Role != models.RoleAdmin {
		return nil, errAccessDenied()
	}
	return a.store.ListAudit(100)
}

func errAccessDenied() error {
	return &AppError{Code: "access_denied", Message: "Недостаточно прав"}
}

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string { return e.Message }
