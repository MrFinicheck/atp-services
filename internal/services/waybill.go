package services

import (
	"errors"
	"strings"
	"time"

	"atp-services/internal/models"
	"atp-services/internal/store"
)

type WaybillService struct {
	store *store.Store
}

func NewWaybillService(s *store.Store) *WaybillService {
	return &WaybillService{store: s}
}

func (ws *WaybillService) CloseShift(driverID string, req models.CloseShiftRequest) (*models.CloseShiftResult, error) {
	if req.EndOdometer <= req.StartOdometer {
		return nil, errors.New("end odometer must be greater than start")
	}
	v, err := ws.store.FindVehicle(req.VehicleID)
	if err != nil {
		return nil, err
	}
	mileage := req.EndOdometer - req.StartOdometer
	actual := req.FuelStart + req.FuelRefilled - req.FuelEnd
	if actual < 0 {
		actual = 0
	}
	norm := mileage / 100 * v.FuelNorm
	overPercent := 0.0
	if norm > 0 {
		overPercent = ((actual - norm) / norm) * 100
	}
	date := time.Now().Format("2006-01-02")
	wb, err := ws.store.FindOpenWaybill(driverID, date)
	if err != nil {
		wb = &models.Waybill{
			DriverID:  driverID,
			VehicleID: req.VehicleID,
			Date:      date,
		}
	}
	wb.VehicleID = req.VehicleID
	wb.StartOdometer = req.StartOdometer
	wb.EndOdometer = req.EndOdometer
	wb.FuelStart = req.FuelStart
	wb.FuelEnd = req.FuelEnd
	wb.FuelRefilled = req.FuelRefilled
	wb.ActualConsumption = actual
	wb.NormConsumption = norm
	wb.OverPercent = overPercent
	wb.Comment = strings.TrimSpace(req.Comment)

	blocked := overPercent > fuelOverrunTolerance*100
	requireComment := blocked && wb.Comment == ""

	result := &models.CloseShiftResult{
		Waybill:        *wb,
		Blocked:        requireComment,
		RequireComment: requireComment,
	}
	if requireComment {
		result.Message = "Перерасход топлива. Укажите комментарий для закрытия смены."
		if err := ws.store.SaveWaybill(wb); err != nil {
			return nil, err
		}
		return result, nil
	}
	now := time.Now()
	wb.Closed = true
	wb.ClosedAt = &now
	if err := ws.store.SaveWaybill(wb); err != nil {
		return nil, err
	}
	_ = ws.store.AddAudit(driverID, "close_shift", "waybill", wb.ID, result.Message)
	result.Message = "Смена успешно закрыта"
	result.Blocked = false
	return result, nil
}

func (ws *WaybillService) List() ([]models.Waybill, error) {
	return ws.store.ListWaybills()
}
