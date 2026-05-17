package services

import (
	"atp-services/internal/ports"
)

type Seeder struct {
	store ports.UnitOfWork
	auth  ports.AuthService
}

func NewSeeder(store ports.UnitOfWork, auth ports.AuthService) *Seeder {
	return &Seeder{store: store, auth: auth}
}

func (s *Seeder) Seed() error {
	return SeedDemoData(s.store, s.auth)
}
