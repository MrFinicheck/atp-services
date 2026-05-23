# Документация проекта «АТП — TransitOS»

Полный комплект материалов к дипломному проекту **atp-services** (учёт перевозок малого автотранспортного предприятия).

Исходный код приложения — корень этого репозитория (`/`).

## Содержание

| Раздел | Файл | Для кого |
|--------|------|----------|
| **Руководство пользователя (все роли)** | [user-guide-overview.md](user-guide-overview.md) | Администратор, диспетчер, водитель |
| Администратор (детально) | [user-guide/administrator.md](user-guide/administrator.md) | Админ |
| Диспетчер | [user-guide/dispatcher.md](user-guide/dispatcher.md) | Диспетчер |
| Водитель | [user-guide/driver.md](user-guide/driver.md) | Водитель |
| Разработчик — обзор | [developer/README.md](developer/README.md) | Программист |
| Архитектура | [developer/architecture.md](developer/architecture.md) | Программист, защита |
| Обоснование решений (ADR) | [developer/decisions.md](developer/decisions.md) | Программист, защита |
| API | [developer/api-reference.md](developer/api-reference.md) | Интегратор |
| Сборка и тесты | [developer/setup-and-testing.md](developer/setup-and-testing.md) | DevOps |
| Docker | [../docker/README.md](../docker/README.md) | DevOps |
| Примеры HTTP | [examples/api.http](examples/api.http) | Разработчик |
| Хранение LevelDB | [examples/storage-leveldb.md](examples/storage-leveldb.md) | Разработчик |

## Быстрый старт

### Desktop (Wails)

```bash
cd frontend && npm install && npm run build
cd .. && wails dev
```

### Web

```bash
cd frontend && npm run build
go run ./cmd/web -addr :8080
```

### Docker

```bash
docker compose -f docker/docker-compose.yml up --build -d
```

Откройте http://localhost:8080

### Тесты

```bash
go test ./...
```

## Демо-учётные записи

| Логин | Пароль | Роль |
|-------|--------|------|
| admin | admin123 | Администратор |
| dispatcher | disp123 | Диспетчер |
| driver1 | drv123 | Водитель |
