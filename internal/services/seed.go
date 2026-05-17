package services

import (
	"time"

	"atp-services/internal/models"
	"atp-services/internal/store"
)

func SeedDemoData(s *store.Store, auth *AuthService) error {
	ok, err := s.IsSeeded()
	if err != nil || ok {
		return err
	}

	users := []struct {
		login, pass, first, last, phone string
		role                            models.Role
	}{
		{"admin", "admin123", "Иван", "Петров", "+79001112233", models.RoleAdmin},
		{"dispatcher", "disp123", "Мария", "Сидорова", "+79002223344", models.RoleDispatcher},
		{"driver1", "drv123", "Алексей", "Козлов", "+79003334455", models.RoleDriver},
		{"driver2", "drv123", "Дмитрий", "Новиков", "+79004445566", models.RoleDriver},
	}
	var driverIDs []string
	for _, u := range users {
		hash, _ := auth.HashPassword(u.pass)
		user := &models.User{
			Login: u.login, PasswordHash: hash, Role: u.role,
			FirstName: u.first, LastName: u.last, Phone: u.phone, Active: true,
		}
		if err := s.SaveUser(user); err != nil {
			return err
		}
		if u.role == models.RoleDriver {
			driverIDs = append(driverIDs, user.ID)
		}
	}

	vehicles := []*models.Vehicle{
		{Plate: "А123ВС77", Model: "ГАЗель Next", FuelNorm: 12.5, Active: true},
		{Plate: "В456КМ77", Model: "Ford Transit", FuelNorm: 10.2, Active: true},
		{Plate: "С789НР77", Model: "Hyundai County", FuelNorm: 18.0, Active: true},
	}
	var vehicleIDs []string
	for _, v := range vehicles {
		if err := s.SaveVehicle(v); err != nil {
			return err
		}
		vehicleIDs = append(vehicleIDs, v.ID)
	}

	tariffs := []*models.Tariff{
		{Name: "Городской груз", BaseFee: 500, PricePerKm: 35, PricePerIdleHr: 400, UrgencyCoeff: 1.5, Active: true},
		{Name: "Пассажирский", BaseFee: 800, PricePerKm: 45, PricePerIdleHr: 500, UrgencyCoeff: 1.3, Active: true},
	}
	var tariffIDs []string
	for _, t := range tariffs {
		if err := s.SaveTariff(t); err != nil {
			return err
		}
		tariffIDs = append(tariffIDs, t.ID)
	}

	clients := []*models.Client{
		{Name: "ООО СтройМастер", Phone: "+74951234567", DebtLimit: 100000},
		{Name: "ИП Смирнов", Phone: "+74957654321", DebtLimit: 50000},
	}
	var clientIDs []string
	for _, c := range clients {
		if err := s.SaveClient(c); err != nil {
			return err
		}
		clientIDs = append(clientIDs, c.ID)
	}

	ts := NewTariffService(s)
	os := NewOrderService(s, ts)
	if len(clientIDs) > 0 && len(vehicleIDs) > 0 && len(driverIDs) > 0 {
		_, _ = os.Create(models.CreateOrderRequest{
			ClientID: clientIDs[0], VehicleID: vehicleIDs[0], DriverID: driverIDs[0],
			FromAddr: "ул. Ленина, 10", ToAddr: "пр. Мира, 25",
			DistanceKm: 15, IdleHours: 0.5, Urgent: false, TariffID: tariffIDs[0],
			ScheduledAt: time.Now().Format(time.RFC3339),
		}, "system")
	}

	return s.MarkSeeded()
}
