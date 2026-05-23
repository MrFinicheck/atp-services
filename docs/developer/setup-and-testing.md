# Сборка, запуск и тестирование

## Локальная разработка

### 1. Клонирование и зависимости

```bash
cd atp-services
cd frontend && npm install && cd ..
go mod download
```

### 2. Frontend

```bash
cd frontend
npm run dev      # Vite :5173 (API нужно проксировать или использовать wails dev)
npm run build    # артефакт в frontend/dist
```

### 3. Desktop

```bash
wails dev
# Linux при проблемах WebKit:
wails dev -tags=webkit2_41
```

### 4. Web

```bash
go run ./cmd/web -addr :8080 -static frontend/dist
```

### 5. Переменные окружения

```bash
export ATP_DATA_DIR=/var/lib/atp/data
go run ./cmd/web
```

## Unit-тесты

Расположение: `atp-services/internal/**/**/*_test.go`

| Пакет | Что проверяет |
|-------|----------------|
| `services` | Login, tariff price, schedule nil-slice, driver filter |
| `store` | Client/order/session persistence |
| `api` | Health, login/me, 401 |
| `app` | RBAC driver/dashboard, dispatcher create order |

```bash
cd atp-services
go test ./...
go test -v ./internal/services/...
go test -coverprofile=coverage.out ./internal/...
go tool cover -html=coverage.out
```

Тесты используют `t.TempDir()` — отдельная LevelDB на каждый тест.

## Docker

См. [../docker/README.md](../docker/README.md).

## CI (рекомендация)

```yaml
- run: cd frontend && npm ci && npm run build
- run: go test ./...
- run: go build -o /dev/null ./cmd/web
```

## Отладка LevelDB

```bash
ls -la data/data/
# Остановить приложение перед копированием каталога
```

Не редактируйте LevelDB вручную при запущенном процессе.
