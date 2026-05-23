package services

import (
	"path/filepath"
	"testing"

	"atp-services/internal/models"
	"atp-services/internal/store"
)

func TestCalculatePrice(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "tariff")
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	auth := NewAuthService(st)
	if err := NewSeeder(st, auth).Seed(); err != nil {
		t.Fatal(err)
	}

	tariffs, err := st.ListTariffs()
	if err != nil || len(tariffs) == 0 {
		t.Fatal("no tariffs after seed")
	}
	tariffID := tariffs[0].ID
	ts := NewTariffService(st)

	price, err := ts.CalculatePrice(tariffID, 10, 0, false)
	if err != nil {
		t.Fatal(err)
	}
	if price <= 0 {
		t.Fatalf("expected positive price, got %v", price)
	}

	urgent, err := ts.CalculatePrice(tariffID, 10, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if urgent <= price {
		t.Fatalf("urgent price %v should exceed normal %v", urgent, price)
	}
}

func TestCalculatePriceInactiveTariff(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	ts := NewTariffService(st)
	inactive := &models.Tariff{
		Name: "test", BaseFee: 100, PricePerKm: 10, PricePerIdleHr: 50,
		UrgencyCoeff: 1.5, Active: false,
	}
	if err := st.SaveTariff(inactive); err != nil {
		t.Fatal(err)
	}
	_, err = ts.CalculatePrice(inactive.ID, 5, 0, false)
	if err == nil {
		t.Fatal("expected error for inactive tariff")
	}
}
