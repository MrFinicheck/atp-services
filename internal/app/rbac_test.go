package app

import (
	"testing"
	"time"

	"atp-services/internal/models"
)

func TestDriverCannotAccessDashboard(t *testing.T) {
	core, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer core.Close()

	resp, err := core.Login(models.LoginRequest{Login: "driver1", Password: "drv123"})
	if err != nil {
		t.Fatal(err)
	}

	_, err = core.Dashboard(resp.Token)
	if err == nil {
		t.Fatal("driver should not access dashboard")
	}
}

func TestDispatcherCanCreateOrder(t *testing.T) {
	core, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer core.Close()

	resp, err := core.Login(models.LoginRequest{Login: "dispatcher", Password: "disp123"})
	if err != nil {
		t.Fatal(err)
	}

	clients, _ := core.ListClients(resp.Token)
	vehicles, _ := core.ListVehicles(resp.Token)
	tariffs, _ := core.ListTariffs(resp.Token)
	adminResp, err := core.Login(models.LoginRequest{Login: "admin", Password: "admin123"})
	if err != nil {
		t.Fatal(err)
	}
	users, _ := core.ListUsers(adminResp.Token)
	var driverID string
	for _, u := range users {
		if u.Login == "driver1" {
			driverID = u.ID
			break
		}
	}
	if driverID == "" || len(clients) == 0 || len(vehicles) == 0 || len(tariffs) == 0 {
		t.Fatal("seed data missing")
	}

	_, err = core.CreateOrder(resp.Token, models.CreateOrderRequest{
		ClientID: clients[0].ID, VehicleID: vehicles[0].ID, DriverID: driverID,
		FromAddr: "Склад", ToAddr: "Клиент", DistanceKm: 8,
		TariffID: tariffs[0].ID, ScheduledAt: time.Now().Format(time.RFC3339),
	})
	if err != nil {
		t.Fatalf("dispatcher create order: %v", err)
	}
}
