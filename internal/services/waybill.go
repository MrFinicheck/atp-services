package services

import (
	"errors"
	"strings"
	"time"

	"atp-services/internal/models"
	"atp-services/internal/ports"
	"atp-services/internal/store"
)

type WaybillService struct {
	waybills ports.WaybillRepository
	vehicles ports.VehicleRepository
	audit    ports.AuditRepository
}

func NewWaybillService(uow ports.UnitOfWork) *WaybillService {
	return &WaybillService{waybills: uow, vehicles: uow, audit: uow}
}

func (ws *WaybillService) OpenShift(driverID string, req models.OpenShiftRequest) (*models.OpenShiftResult, error) {
	if req.VehicleID == "" {
		return nil, errors.New("выберите автомобиль")
	}
	if req.StartOdometer < 0 || req.FuelStart < 0 {
		return nil, errors.New("некорректные показания на начало смены")
	}
	if _, err := ws.vehicles.FindVehicle(req.VehicleID); err != nil {
		return nil, err
	}
	date := time.Now().Format("2006-01-02")
	existing, err := ws.waybills.FindWaybillByDriverAndDate(driverID, date)
	if err == nil {
		if existing.Closed {
			return nil, errors.New("смена уже закрыта за сегодня")
		}
		return nil, errors.New("смена уже открыта")
	}
	if !errors.Is(err, store.ErrNotFound) {
		return nil, err
	}
	wb := &models.Waybill{
		DriverID:      driverID,
		VehicleID:     req.VehicleID,
		Date:          date,
		StartOdometer: req.StartOdometer,
		FuelStart:     req.FuelStart,
		Closed:        false,
	}
	if err := ws.waybills.SaveWaybill(wb); err != nil {
		return nil, err
	}
	_ = ws.audit.AddAudit(driverID, "open_shift", "waybill", wb.ID, wb.VehicleID)
	return &models.OpenShiftResult{
		Waybill: *wb,
		Message: "Смена открыта",
	}, nil
}

func (ws *WaybillService) CloseShift(driverID string, req models.CloseShiftRequest) (*models.CloseShiftResult, error) {
	date := time.Now().Format("2006-01-02")
	existing, err := ws.waybills.FindWaybillByDriverAndDate(driverID, date)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, errors.New("смена не открыта — сначала откройте смену")
		}
		return nil, err
	}
	if existing.Closed {
		return nil, errors.New("смена уже закрыта за сегодня")
	}
	wb := existing
	startOdometer := wb.StartOdometer
	fuelStart := wb.FuelStart
	if req.EndOdometer <= startOdometer {
		return nil, errors.New("конечный пробег должен быть больше начального")
	}
	v, err := ws.vehicles.FindVehicle(wb.VehicleID)
	if err != nil {
		return nil, err
	}
	mileage := req.EndOdometer - startOdometer
	actual := fuelStart + req.FuelRefilled - req.FuelEnd
	if actual < 0 {
		actual = 0
	}
	norm := mileage / 100 * v.FuelNorm
	overPercent := 0.0
	if norm > 0 {
		overPercent = ((actual - norm) / norm) * 100
	}
	wb.EndOdometer = req.EndOdometer
	wb.EndOdometer = req.EndOdometer
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
		if err := ws.waybills.SaveWaybill(wb); err != nil {
			return nil, err
		}
		return result, nil
	}
	now := time.Now()
	wb.Closed = true
	wb.ClosedAt = &now
	if err := ws.waybills.SaveWaybill(wb); err != nil {
		return nil, err
	}
	_ = ws.audit.AddAudit(driverID, "close_shift", "waybill", wb.ID, result.Message)
	result.Message = "Смена успешно закрыта"
	result.Blocked = false
	return result, nil
}

func (ws *WaybillService) List() ([]models.Waybill, error) {
	return ws.waybills.ListWaybills()
}

func (ws *WaybillService) ShiftStatus(driverID string) (*models.ShiftStatus, error) {
	date := time.Now().Format("2006-01-02")
	status := &models.ShiftStatus{Date: date}
	wb, err := ws.waybills.FindWaybillByDriverAndDate(driverID, date)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return status, nil
		}
		return nil, err
	}
	status.Closed = wb.Closed
	status.Opened = !wb.Closed
	if status.Opened {
		status.VehicleID = wb.VehicleID
		status.StartOdometer = wb.StartOdometer
		status.FuelStart = wb.FuelStart
	}
	return status, nil
}

func (ws *WaybillService) DriverShiftOpen(driverID, date string) (bool, error) {
	wb, err := ws.waybills.FindWaybillByDriverAndDate(driverID, date)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return !wb.Closed, nil
}

func (ws *WaybillService) DriverShiftClosed(driverID, date string) (bool, error) {
	wb, err := ws.waybills.FindWaybillByDriverAndDate(driverID, date)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return wb.Closed, nil
}
