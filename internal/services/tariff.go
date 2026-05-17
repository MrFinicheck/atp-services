package services

import (
	"errors"

	"atp-services/internal/models"
	"atp-services/internal/ports"
)

const fuelOverrunTolerance = 0.05

type TariffService struct {
	tariffs ports.TariffRepository
}

func NewTariffService(uow ports.TariffRepository) *TariffService {
	return &TariffService{tariffs: uow}
}

func (ts *TariffService) CalculatePrice(tariffID string, distanceKm, idleHours float64, urgent bool) (float64, error) {
	t, err := ts.tariffs.FindTariff(tariffID)
	if err != nil {
		return 0, err
	}
	if !t.Active {
		return 0, errors.New("tariff is inactive")
	}
	coeff := 1.0
	if urgent {
		coeff = t.UrgencyCoeff
		if coeff < 1 {
			coeff = 1.5
		}
	}
	price := t.BaseFee + distanceKm*t.PricePerKm + idleHours*t.PricePerIdleHr
	return price * coeff, nil
}

func (ts *TariffService) List() ([]models.Tariff, error) {
	return ts.tariffs.ListTariffs()
}

func (ts *TariffService) Save(t *models.Tariff) error {
	if t.UrgencyCoeff <= 0 {
		t.UrgencyCoeff = 1.5
	}
	return ts.tariffs.SaveTariff(t)
}
