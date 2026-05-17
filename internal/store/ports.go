package store

import "atp-services/internal/ports"

// Compile-time check: Store implements UnitOfWork (Liskov).
var _ ports.UnitOfWork = (*Store)(nil)
