package main

import (
	"context"

	"atp-services/internal/app"
	"atp-services/internal/models"
)

type App struct {
	ctx context.Context
	core *app.Application
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	core, err := app.New("")
	if err != nil {
		panic(err)
	}
	a.core = core
}

func (a *App) shutdown(ctx context.Context) {
	if a.core != nil {
		_ = a.core.Close()
	}
}

func (a *App) GetDataDir() string {
	if a.core == nil {
		return ""
	}
	return a.core.DataDir()
}

func (a *App) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	return a.core.Login(req)
}

func (a *App) Logout(token string) error {
	return a.core.Logout(token)
}

func (a *App) Me(token string) (*models.User, error) {
	return a.core.Me(token)
}

func (a *App) ListClients(token string) ([]models.Client, error) {
	return a.core.ListClients(token)
}

func (a *App) SaveClient(token string, c models.Client) (*models.Client, error) {
	return a.core.SaveClient(token, c)
}

func (a *App) ListVehicles(token string) ([]models.Vehicle, error) {
	return a.core.ListVehicles(token)
}

func (a *App) SaveVehicle(token string, v models.Vehicle) (*models.Vehicle, error) {
	return a.core.SaveVehicle(token, v)
}

func (a *App) ListTariffs(token string) ([]models.Tariff, error) {
	return a.core.ListTariffs(token)
}

func (a *App) SaveTariff(token string, t models.Tariff) (*models.Tariff, error) {
	return a.core.SaveTariff(token, t)
}

func (a *App) ListUsers(token string) ([]models.User, error) {
	return a.core.ListUsers(token)
}

func (a *App) CreateUser(token string, u models.User, password string) (*models.User, error) {
	return a.core.CreateUser(token, u, password)
}

func (a *App) ListOrders(token string) ([]models.Order, error) {
	return a.core.ListOrders(token)
}

func (a *App) CreateOrder(token string, req models.CreateOrderRequest) (*models.Order, error) {
	return a.core.CreateOrder(token, req)
}

func (a *App) UpdateOrderStatus(token, orderID, status string) (*models.Order, error) {
	return a.core.UpdateOrderStatus(token, orderID, status)
}

func (a *App) PreviewPrice(token, tariffID string, distanceKm, idleHours float64, urgent bool) (float64, error) {
	return a.core.PreviewPrice(token, tariffID, distanceKm, idleHours, urgent)
}

func (a *App) VehicleSchedule(token string) ([]models.VehicleScheduleItem, error) {
	return a.core.VehicleSchedule(token)
}

func (a *App) CloseShift(token string, req models.CloseShiftRequest) (*models.CloseShiftResult, error) {
	return a.core.CloseShift(token, req)
}

func (a *App) ListWaybills(token string) ([]models.Waybill, error) {
	return a.core.ListWaybills(token)
}

func (a *App) Dashboard(token string) (*models.DashboardStats, error) {
	return a.core.Dashboard(token)
}

func (a *App) DriverRating(token string) ([]map[string]any, error) {
	return a.core.DriverRating(token)
}

func (a *App) ListAudit(token string) ([]models.AuditEntry, error) {
	return a.core.ListAudit(token)
}
