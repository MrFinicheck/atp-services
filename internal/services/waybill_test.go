package services

import (
	"testing"
	"time"

	"atp-services/internal/models"
	"atp-services/internal/store"
)

func TestOpenShiftAndClose(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	auth := NewAuthService(st)
	if err := NewSeeder(st, auth).Seed(); err != nil {
		t.Fatal(err)
	}

	driver1, _ := st.FindUserByLogin("driver1")
	driver2, _ := st.FindUserByLogin("driver2")
	vehicles, _ := st.ListVehicles()
	ws := NewWaybillService(st)

	// driver1 shift opened in seed
	_, err = ws.OpenShift(driver1.ID, models.OpenShiftRequest{
		VehicleID: vehicles[0].ID, StartOdometer: 100, FuelStart: 50,
	})
	if err == nil {
		t.Fatal("expected error: shift already open for driver1")
	}

	_, err = ws.OpenShift(driver2.ID, models.OpenShiftRequest{
		VehicleID: vehicles[0].ID, StartOdometer: 200, FuelStart: 60,
	})
	if err != nil {
		t.Fatalf("open driver2: %v", err)
	}

	_, err = ws.CloseShift(driver1.ID, models.CloseShiftRequest{
		EndOdometer: 1005, FuelEnd: 39.5, FuelRefilled: 0,
	})
	if err != nil {
		t.Fatalf("close driver1: %v", err)
	}
}

func TestCloseShiftTwiceRejected(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	auth := NewAuthService(st)
	if err := NewSeeder(st, auth).Seed(); err != nil {
		t.Fatal(err)
	}

	driver, _ := st.FindUserByLogin("driver1")
	ws := NewWaybillService(st)

	closeReq := models.CloseShiftRequest{EndOdometer: 1005, FuelEnd: 39.5, FuelRefilled: 0}
	if _, err := ws.CloseShift(driver.ID, closeReq); err != nil {
		t.Fatalf("first close: %v", err)
	}
	if _, err := ws.CloseShift(driver.ID, closeReq); err == nil {
		t.Fatal("expected error on second close")
	}
}

func TestCloseShiftWithoutOpenRejected(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	auth := NewAuthService(st)
	_ = NewSeeder(st, auth).Seed()
	driver2, _ := st.FindUserByLogin("driver2")
	ws := NewWaybillService(st)

	_, err = ws.CloseShift(driver2.ID, models.CloseShiftRequest{
		EndOdometer: 205, FuelEnd: 49.5,
	})
	if err == nil {
		t.Fatal("expected error closing without open")
	}
}

func TestCreateOrderRejectedWhenDriverShiftClosed(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	auth := NewAuthService(st)
	if err := NewSeeder(st, auth).Seed(); err != nil {
		t.Fatal(err)
	}

	driver, _ := st.FindUserByLogin("driver1")
	vehicles, _ := st.ListVehicles()
	clients, _ := st.ListClients()
	tariffs, _ := st.ListTariffs()

	ws := NewWaybillService(st)
	_, err = ws.CloseShift(driver.ID, models.CloseShiftRequest{
		EndOdometer: 1005, FuelEnd: 39.5,
	})
	if err != nil {
		t.Fatal(err)
	}

	os := NewOrderService(st, NewTariffService(st))
	_, err = os.Create(models.CreateOrderRequest{
		ClientID: clients[0].ID, VehicleID: vehicles[0].ID, DriverID: driver.ID,
		FromAddr: "A", ToAddr: "B", DistanceKm: 5, TariffID: tariffs[0].ID,
		ScheduledAt: time.Now().Format(time.RFC3339),
	}, "dispatcher")
	if err == nil {
		t.Fatal("expected error assigning order to driver with closed shift")
	}
}

func TestCreateOrderRejectedWhenDriverShiftNotOpen(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	auth := NewAuthService(st)
	_ = NewSeeder(st, auth).Seed()

	driver2, _ := st.FindUserByLogin("driver2")
	vehicles, _ := st.ListVehicles()
	clients, _ := st.ListClients()
	tariffs, _ := st.ListTariffs()

	os := NewOrderService(st, NewTariffService(st))
	_, err = os.Create(models.CreateOrderRequest{
		ClientID: clients[0].ID, VehicleID: vehicles[0].ID, DriverID: driver2.ID,
		FromAddr: "A", ToAddr: "B", DistanceKm: 5, TariffID: tariffs[0].ID,
		ScheduledAt: time.Now().Format(time.RFC3339),
	}, "dispatcher")
	if err == nil {
		t.Fatal("expected error assigning order without open shift")
	}
}

func TestCreateOrderAllowedWhenShiftOpen(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	auth := NewAuthService(st)
	_ = NewSeeder(st, auth).Seed()

	driver2, _ := st.FindUserByLogin("driver2")
	vehicles, _ := st.ListVehicles()
	clients, _ := st.ListClients()
	tariffs, _ := st.ListTariffs()

	ws := NewWaybillService(st)
	_, err = ws.OpenShift(driver2.ID, models.OpenShiftRequest{
		VehicleID: vehicles[0].ID, StartOdometer: 500, FuelStart: 30,
	})
	if err != nil {
		t.Fatal(err)
	}

	os := NewOrderService(st, NewTariffService(st))
	_, err = os.Create(models.CreateOrderRequest{
		ClientID: clients[0].ID, VehicleID: vehicles[0].ID, DriverID: driver2.ID,
		FromAddr: "X", ToAddr: "Y", DistanceKm: 5, TariffID: tariffs[0].ID,
		ScheduledAt: time.Now().Format(time.RFC3339),
	}, "dispatcher")
	if err != nil {
		t.Fatalf("create with open shift: %v", err)
	}
}
