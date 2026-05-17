package services

import (
	"errors"
	"time"

	"atp-services/internal/models"
	"atp-services/internal/ports"
)

type OrderService struct {
	orders   ports.OrderRepository
	vehicles ports.VehicleRepository
	audit    ports.AuditRepository
	users    ports.UserRepository
	tariff   ports.TariffCalculator
}

func NewOrderService(uow ports.UnitOfWork, ts ports.TariffCalculator) *OrderService {
	return &OrderService{orders: uow, vehicles: uow, audit: uow, users: uow, tariff: ts}
}

func (os *OrderService) Create(req models.CreateOrderRequest, actorID string) (*models.Order, error) {
	price, err := os.tariff.CalculatePrice(req.TariffID, req.DistanceKm, req.IdleHours, req.Urgent)
	if err != nil {
		return nil, err
	}
	scheduled, _ := time.Parse(time.RFC3339, req.ScheduledAt)
	if scheduled.IsZero() {
		scheduled = time.Now()
	}
	o := &models.Order{
		ClientID:    req.ClientID,
		VehicleID:   req.VehicleID,
		DriverID:    req.DriverID,
		FromAddr:    req.FromAddr,
		ToAddr:      req.ToAddr,
		DistanceKm:  req.DistanceKm,
		IdleHours:   req.IdleHours,
		Urgent:      req.Urgent,
		TariffID:    req.TariffID,
		Price:       price,
		Status:      models.OrderAssigned,
		ScheduledAt: scheduled,
		CreatedAt:   time.Now(),
	}
	if err := os.orders.SaveOrder(o); err != nil {
		return nil, err
	}
	_ = os.audit.AddAudit(actorID, "create", "order", o.ID, req.FromAddr+" -> "+req.ToAddr)
	return o, nil
}

func (os *OrderService) ListForRole(user *models.User) ([]models.Order, error) {
	all, err := os.orders.ListOrders()
	if err != nil {
		return nil, err
	}
	switch user.Role {
	case models.RoleDriver:
		var filtered []models.Order
		today := time.Now().Format("2006-01-02")
		for _, o := range all {
			if o.DriverID == user.ID && o.ScheduledAt.Format("2006-01-02") == today {
				filtered = append(filtered, o)
			}
		}
		return filtered, nil
	default:
		return all, nil
	}
}

func (os *OrderService) UpdateStatus(orderID string, status models.OrderStatus, actorID string) (*models.Order, error) {
	o, err := os.orders.FindOrder(orderID)
	if err != nil {
		return nil, err
	}
	if user, _ := os.users.FindUserByID(actorID); user != nil && user.Role == models.RoleDriver && o.DriverID != actorID {
		return nil, errors.New("access denied")
	}
	now := time.Now()
	switch status {
	case models.OrderInProgress:
		o.StartedAt = &now
	case models.OrderCompleted:
		o.CompletedAt = &now
	}
	o.Status = status
	if err := os.orders.SaveOrder(o); err != nil {
		return nil, err
	}
	_ = os.audit.AddAudit(actorID, "status", "order", o.ID, string(status))
	return o, nil
}

func (os *OrderService) ScheduleToday() ([]models.VehicleScheduleItem, error) {
	vehicles, err := os.vehicles.ListVehicles()
	if err != nil {
		return nil, err
	}
	orders, err := os.orders.ListOrders()
	if err != nil {
		return nil, err
	}
	today := time.Now().Format("2006-01-02")
	var result []models.VehicleScheduleItem
	for _, v := range vehicles {
		if !v.Active {
			continue
		}
		item := models.VehicleScheduleItem{VehicleID: v.ID, Plate: v.Plate}
		for _, o := range orders {
			if o.VehicleID == v.ID && o.ScheduledAt.Format("2006-01-02") == today {
				item.Orders = append(item.Orders, o)
			}
		}
		result = append(result, item)
	}
	return result, nil
}

func (os *OrderService) PreviewPrice(tariffID string, distanceKm, idleHours float64, urgent bool) (float64, error) {
	return os.tariff.CalculatePrice(tariffID, distanceKm, idleHours, urgent)
}
