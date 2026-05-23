package services

import (
	"testing"
	"time"

	"atp-services/internal/models"
	"atp-services/internal/store"
)

func TestScheduleTodayEmptyOrdersSlice(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	auth := NewAuthService(st)
	if err := NewSeeder(st, auth).Seed(); err != nil {
		t.Fatal(err)
	}

	ts := NewTariffService(st)
	os := NewOrderService(st, ts)

	schedule, err := os.ScheduleToday()
	if err != nil {
		t.Fatal(err)
	}
	if len(schedule) == 0 {
		t.Fatal("expected at least one vehicle in schedule")
	}
	for _, item := range schedule {
		if item.Orders == nil {
			t.Fatalf("vehicle %s: Orders must not be nil (JSON null breaks UI)", item.Plate)
		}
	}
}

func TestListForRoleDriverFiltersToday(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	auth := NewAuthService(st)
	if err := NewSeeder(st, auth).Seed(); err != nil {
		t.Fatal(err)
	}

	driver, err := st.FindUserByLogin("driver1")
	if err != nil {
		t.Fatal(err)
	}
	vehicles, _ := st.ListVehicles()
	clients, _ := st.ListClients()
	tariffs, _ := st.ListTariffs()
	if len(vehicles) == 0 || len(clients) == 0 || len(tariffs) == 0 {
		t.Fatal("seed incomplete")
	}

	ts := NewTariffService(st)
	os := NewOrderService(st, ts)

	today := time.Now()
	_, err = os.Create(models.CreateOrderRequest{
		ClientID: clients[0].ID, VehicleID: vehicles[0].ID, DriverID: driver.ID,
		FromAddr: "A", ToAddr: "B", DistanceKm: 5, TariffID: tariffs[0].ID,
		ScheduledAt: today.Format(time.RFC3339),
	}, driver.ID)
	if err != nil {
		t.Fatal(err)
	}

	list, err := os.ListForRole(driver)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) < 1 {
		t.Fatalf("driver should see at least 1 order today, got %d", len(list))
	}
	var found bool
	for _, o := range list {
		if o.FromAddr == "A" && o.ToAddr == "B" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("created order not in driver's list")
	}
}
