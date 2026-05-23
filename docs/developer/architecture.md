# Архитектура системы

## Назначение

Автоматизация учёта услуг малого АТП: заявки, тарификация, диспетчеризация, мобильный кабинет водителя, контроль ГСМ.

## Диаграмма слоёв

```
┌─────────────────────────────────────────────────────────┐
│  Frontend (TypeScript, Vite, TransitOS UI)              │
│  api/client.ts — Wails bindings ИЛИ REST                │
└──────────────────────────┬──────────────────────────────┘
                           │
         ┌─────────────────┴─────────────────┐
         ▼                                   ▼
┌─────────────────┐                 ┌─────────────────┐
│  app.go (Wails) │                 │  api.Server     │
│  thin adapter   │                 │  REST + SPA     │
└────────┬────────┘                 └────────┬────────┘
         │                                   │
         └─────────────────┬─────────────────┘
                           ▼
                 ┌───────────────────┐
                 │ app.Application   │  RBAC, делегирование
                 └─────────┬─────────┘
                           ▼
                 ┌───────────────────┐
                 │ app.Container     │  composition root
                 └─────────┬─────────┘
                           │
     ┌─────────────────────┼─────────────────────┐
     ▼                     ▼                     ▼
 AuthService          OrderService          WaybillService
 TariffService        ReportService         Seeder
     │                     │                     │
     └─────────────────────┼─────────────────────┘
                           ▼
                 ┌───────────────────┐
                 │ ports.UnitOfWork  │  интерфейсы
                 └─────────┬─────────┘
                           ▼
                 ┌───────────────────┐
                 │ store.Store       │  LevelDB + JSON
                 └───────────────────┘
```

## Ключевые компоненты

### `internal/app/Application`

Единая точка use-case для HTTP и Wails. Методы вида `ListOrders(token)`:

1. `requireUser(token)` — валидация сессии.
2. Проверка роли (RBAC).
3. Вызов сервиса.

### `internal/services`

| Сервис | Ответственность |
|--------|-----------------|
| `AuthService` | Login, bcrypt, сессии 24ч |
| `OrderService` | CRUD заявок, фильтр для водителя, `ScheduleToday` |
| `TariffService` | Расчёт цены, CRUD тарифов |
| `WaybillService` | Закрытие смены, перерасход 5% |
| `ReportService` | Dashboard, рейтинг водителей |

### `internal/store`

- Префиксные ключи LevelDB (`order:`, `user:`, …).
- UUID для новых сущностей.
- `flock` на `{dataDir}/.atp.lock` — один процесс на каталог данных.

### Frontend

- `renderApp` → login или shell (dock navigation).
- `normalizeSchedule` — защита от `null` в `orders` (Go nil slice → JSON null).

## Потоки данных

### Создание заявки

```
UI form → api.createOrder → Application.CreateOrder
  → OrderService.Create → TariffService.CalculatePrice
  → Store.SaveOrder → Audit
```

### Аутентификация

```
POST /api/login → AuthService.Login → SaveSession
→ Bearer token → FindSession → FindUser
```

## Деплой-формы

| Режим | Entry | UI |
|-------|-------|-----|
| Desktop | `main.go` + Wails | WebView2 / WebKit |
| Web | `cmd/web` | Static SPA + API |
| Docker | `atp-web` binary | То же, volume `/data` |

## SOLID (кратко)

Детали в `atp-services/docs/ARCHITECTURE.md` и [decisions.md](decisions.md).
