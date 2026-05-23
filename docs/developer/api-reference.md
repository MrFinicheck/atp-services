# Справочник REST API

Базовый URL: `http://localhost:8080` (или хост Docker).

Аутентификация (кроме `/api/login`, `/api/health`):

```http
Authorization: Bearer <token>
```

Альтернативы: `X-Session-Token`, query `?token=`.

## Общие коды ответов

| Код | Значение |
|-----|----------|
| 200 | OK |
| 400 | Неверное тело запроса |
| 401 | Нет/неверный токен |
| 403 | Недостаточно прав |
| 405 | Метод не поддерживается |

## Endpoints

### `GET /api/health`

Проверка живости. Без авторизации.

```json
{ "status": "ok" }
```

### `POST /api/login`

```json
{ "login": "dispatcher", "password": "disp123" }
```

Ответ: `LoginResponse` — `token`, `user` (без passwordHash).

### `GET|POST /api/logout`

Инвалидирует сессию (если токен передан).

### `GET /api/me`

Текущий пользователь.

### `GET|POST /api/clients`

- GET — список клиентов.
- POST — создать/обновить (`id` опционален, генерируется).

**POST body:** `Client` — `name`, `phone`, `debtLimit`.

### `GET|POST /api/vehicles`

**POST:** admin only. `Vehicle` — `plate`, `model`, `fuelNorm`, `active`.

### `GET|POST /api/tariffs`

**POST:** admin only. `Tariff` — `name`, `baseFee`, `pricePerKm`, `pricePerIdleHr`, `urgencyCoeff`, `active`.

### `GET|POST|DELETE /api/users`

- GET, POST, DELETE — **admin only**.
- POST body: `{ "user": User, "password": "string" }`
- DELETE query: `?id={userId}` — удаление учётной записи (нельзя удалить себя или последнего активного администратора)

### `GET|POST /api/orders`

- GET — список (для driver — фильтр сегодня + свой driverId).
- POST — создать заявку (`CreateOrderRequest`).

**CreateOrderRequest:**

```json
{
  "clientId": "uuid",
  "vehicleId": "uuid",
  "driverId": "uuid",
  "fromAddr": "string",
  "toAddr": "string",
  "distanceKm": 10,
  "idleHours": 0,
  "urgent": false,
  "tariffId": "uuid",
  "scheduledAt": "2026-05-17T12:00:00Z"
}
```

### `POST /api/orders/status`

```json
{ "orderId": "uuid", "status": "in_progress" }
```

Статусы: `assigned`, `in_progress`, `completed`, `cancelled`.

### `GET|POST /api/orders/preview-price`

Query или POST: `tariffId`, `distanceKm`, `idleHours`, `urgent`.

### `GET /api/schedule`

Массив `VehicleScheduleItem`:

```json
[
  {
    "vehicleId": "uuid",
    "plate": "А123ВС",
    "orders": [ { "id": "...", "fromAddr": "...", "status": "assigned", ... } ]
  }
]
```

`orders` всегда массив (может быть пустым).

### `GET /api/waybills`

Список путевых листов (admin, dispatcher).

### `POST /api/shift/close`

**Driver** (и admin). Body: `CloseShiftRequest`.

```json
{
  "vehicleId": "uuid",
  "startOdometer": 1000,
  "endOdometer": 1040,
  "fuelStart": 45,
  "fuelEnd": 10,
  "fuelRefilled": 5,
  "comment": "пробки"
}
```

Ответ: `CloseShiftResult` — `requireComment`, `message`, `waybill`.

### `GET /api/dashboard`

`DashboardStats` — `ordersToday`, `activeVehicles`, `revenueMonth`, `fuelOverruns`.

### `GET /api/reports/drivers`

Рейтинг водителей (admin).

### `GET /api/audit`

Журнал аудита (admin).

## Wails bindings

Те же операции доступны как методы `main.App` — см. `frontend/wailsjs/go/main/App.d.ts`.

Генерация после изменения `app.go`:

```bash
wails generate module
```
