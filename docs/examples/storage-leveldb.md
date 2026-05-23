# Хранение данных в LevelDB

## Расположение файлов

```
{ATP_DATA_DIR}/
├── .atp.lock          # блокировка процесса (flock)
└── data/              # каталог LevelDB
    ├── CURRENT
    ├── MANIFEST-*
    └── *.ldb
```

По умолчанию `{ATP_DATA_DIR}` = `./data` в каталоге запуска.

## Схема ключей

| Префикс ключа | Сущность | Пример ID |
|---------------|----------|-----------|
| `user:{id}` | Пользователь | UUID |
| `user:login:{login}` | Индекс → user ID | `user:login:dispatcher` |
| `session:{token}` | Сессия | hex token |
| `client:{id}` | Клиент | UUID |
| `vehicle:{id}` | ТС | UUID |
| `tariff:{id}` | Тариф | UUID |
| `order:{id}` | Заявка | UUID |
| `waybill:{id}` | Путевой лист | UUID |
| `audit:{id}` | Аудит | UUID |
| `meta:seeded` | Флаг начального seed | — |

Все значения — **JSON**, сериализация `encoding/json`.

## Пример: сохранение заявки

1. `POST /api/orders` → `Application.CreateOrder`
2. `OrderService.Create` рассчитывает `Price` через `TariffService`
3. `Store.SaveOrder` → `PUT order:{uuid}`

Фрагмент JSON в БД:

```json
{
  "id": "a1b2c3d4-...",
  "clientId": "...",
  "vehicleId": "...",
  "driverId": "...",
  "fromAddr": "Склад",
  "toAddr": "Магазин",
  "distanceKm": 15,
  "idleHours": 0,
  "urgent": false,
  "tariffId": "...",
  "price": 2850,
  "status": "assigned",
  "scheduledAt": "2026-05-17T10:00:00+03:00",
  "createdAt": "2026-05-17T09:15:00+03:00"
}
```

4. `Audit` — запись `audit:{uuid}` с `action: "create"`, `entityType: "order"`.

## Пример: сессия

```json
{
  "token": "64hex...",
  "userId": "user-uuid",
  "role": "dispatcher",
  "expiresAt": "2026-05-18T09:00:00Z"
}
```

При `FindSession` истёкшие сессии удаляются.

## Пример: путевой лист

После `POST /api/shift/close`:

```json
{
  "driverId": "...",
  "vehicleId": "...",
  "startOdometer": 1000,
  "endOdometer": 1045,
  "fuelStart": 40,
  "fuelEnd": 8,
  "fuelRefilled": 12,
  "actualConsumption": 44,
  "normConsumption": 41.2,
  "overPercent": 6.8,
  "comment": "простой с двигателем",
  "closed": true
}
```

## Seed (демо-данные)

При первом `Container.Init()` вызывается `SeedDemoData`:

- пользователи admin, dispatcher, driver1, driver2;
- клиенты, ТС, тарифы;
- `meta:seeded` = true.

Пароли — bcrypt в поле `passwordHash` у `User`.

## Резервное копирование

```bash
# остановить приложение
tar -czf atp-backup-$(date +%F).tar.gz -C /path/to data/
```

## Важно для frontend

Go `nil` slice `[]Order` сериализуется как JSON `null`. В `ScheduleToday` используйте `Orders: []models.Order{}`, иначе UI падает на `orders.length`.
