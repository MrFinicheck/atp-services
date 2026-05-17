package app

import (
	"sync"

	"atp-services/internal/ports"
	"atp-services/internal/services"
	"atp-services/internal/store"
)

// Container — composition root (DIP): собирает зависимости по интерфейсам.
type Container struct {
	mu      sync.Mutex
	dataDir string

	uow     ports.UnitOfWork
	Auth    ports.AuthService
	Tariffs ports.TariffCalculator
	Orders  ports.OrderService
	Waybill ports.WaybillService
	Report  ports.ReportService
}

func NewContainer(dataDir string) *Container {
	return &Container{dataDir: resolveDataDir(dataDir)}
}

func (c *Container) Init() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.uow != nil {
		return nil
	}

	st, err := store.Open(c.dataDir)
	if err != nil {
		return err
	}

	c.uow = st
	c.Auth = services.NewAuthService(st)
	c.Tariffs = services.NewTariffService(st)
	c.Orders = services.NewOrderService(st, c.Tariffs)
	c.Waybill = services.NewWaybillService(st)
	c.Report = services.NewReportService(st)

	_ = services.NewSeeder(st, c.Auth).Seed()
	return nil
}

func (c *Container) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.uow == nil {
		return nil
	}
	err := c.uow.Close()
	c.uow = nil
	c.Auth = nil
	c.Tariffs = nil
	c.Orders = nil
	c.Waybill = nil
	c.Report = nil
	return err
}

func (c *Container) DataDir() string { return c.dataDir }

func (c *Container) Store() ports.UnitOfWork {
	_ = c.Init()
	return c.uow
}
