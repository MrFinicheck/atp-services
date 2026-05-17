package services

import (
	"time"

	"atp-services/internal/models"
	"atp-services/internal/ports"
)

type ReportService struct {
	orders   ports.OrderRepository
	vehicles ports.VehicleRepository
	waybills ports.WaybillRepository
	users    ports.UserRepository
}

func NewReportService(uow ports.UnitOfWork) *ReportService {
	return &ReportService{orders: uow, vehicles: uow, waybills: uow, users: uow}
}

func (rs *ReportService) Dashboard() (*models.DashboardStats, error) {
	orders, _ := rs.orders.ListOrders()
	vehicles, _ := rs.vehicles.ListVehicles()
	waybills, _ := rs.waybills.ListWaybills()

	today := time.Now().Format("2006-01-02")
	month := time.Now().Format("2006-01")
	stats := &models.DashboardStats{}

	for _, o := range orders {
		if o.ScheduledAt.Format("2006-01-02") == today {
			stats.OrdersToday++
		}
		if o.ScheduledAt.Format("2006-01") == month && o.Status == models.OrderCompleted {
			stats.RevenueMonth += o.Price
		}
	}
	for _, v := range vehicles {
		if v.Active {
			stats.ActiveVehicles++
		}
	}
	for _, w := range waybills {
		if !w.Closed {
			stats.OpenWaybills++
		}
		if w.OverPercent > fuelOverrunTolerance*100 {
			stats.FuelOverruns++
		}
	}
	return stats, nil
}

func (rs *ReportService) DriverRating() ([]map[string]any, error) {
	orders, err := rs.orders.ListOrders()
	if err != nil {
		return nil, err
	}
	counts := map[string]int{}
	completed := map[string]int{}
	for _, o := range orders {
		counts[o.DriverID]++
		if o.Status == models.OrderCompleted {
			completed[o.DriverID]++
		}
	}
	var result []map[string]any
	for id, total := range counts {
		u, err := rs.users.FindUserByID(id)
		name := id
		if err == nil {
			name = u.LastName + " " + u.FirstName
		}
		rate := 0.0
		if total > 0 {
			rate = float64(completed[id]) / float64(total) * 100
		}
		result = append(result, map[string]any{
			"driverId": id, "name": name, "total": total,
			"completed": completed[id], "completionRate": rate,
		})
	}
	return result, nil
}
