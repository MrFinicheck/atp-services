# Архитектура (SOLID)

## Слои

```
frontend/          — UI (адаптивный веб)
app.go             — адаптер Wails (тонкий фасад)
internal/app/      — use-case фасад + Container (composition root)
internal/services/ — бизнес-логика
internal/ports/    — интерфейсы (DIP, ISP)
internal/store/    — LevelDB (реализация репозиториев)
internal/models/   — доменные структуры
```

## Принципы SOLID

| Принцип | Реализация |
|---------|------------|
| **S** — единственная ответственность | `AuthService`, `OrderService`, `WaybillService`, `TariffService`, `ReportService` — отдельные типы |
| **O** — открытость/закрытость | Новые тарифы/отчёты через новые типы, без изменения `Container` |
| **L** — подстановка Лисков | `*store.Store` реализует `ports.UnitOfWork` (проверка compile-time) |
| **I** — разделение интерфейсов | `UserRepository`, `OrderRepository`, … вместо одного «бога»-интерфейса |
| **D** — инверсия зависимостей | Сервисы зависят от `ports.*`, не от LevelDB |

## Composition root

`internal/app/container.go` собирает зависимости и вызывает `Seeder` при инициализации.
