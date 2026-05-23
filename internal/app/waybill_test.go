package app

import (
	"testing"

	"atp-services/internal/models"
)

func TestCloseShiftDriver(t *testing.T) {
	core, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer core.Close()

	d, err := core.Login(models.LoginRequest{Login: "driver1", Password: "drv123"})
	if err != nil {
		t.Fatal(err)
	}
	vehicles, err := core.ListVehicles(d.Token)
	if err != nil {
		t.Fatal(err)
	}
	if len(vehicles) == 0 {
		t.Fatal("no vehicles in seed")
	}

	res, err := core.CloseShift(d.Token, models.CloseShiftRequest{
		EndOdometer:  1005,
		FuelEnd:      39.5,
		FuelRefilled: 0,
		Comment:      "",
	})
	if err != nil {
		t.Fatalf("CloseShift: %v", err)
	}
	if res.Message == "" {
		t.Fatal("empty message")
	}

	_, err = core.CloseShift(d.Token, models.CloseShiftRequest{
		EndOdometer:  1010,
		FuelEnd:      39.5,
		FuelRefilled: 0,
	})
	if err == nil {
		t.Fatal("expected error on second close")
	}
}
