# Документация разработчика

## Репозиторий

```
atp-services/
├── main.go, app.go          # Wails entry + bindings
├── cmd/web/main.go          # HTTP server
├── internal/
│   ├── api/                 # REST handlers
│   ├── app/                 # Application facade + RBAC
│   ├── models/
│   ├── ports/               # interfaces (DIP)
│   ├── services/            # domain logic
│   └── store/               # LevelDB
└── frontend/                # TypeScript + Vite
```

## Требования

- Go 1.23+
- Node.js 18+ (для frontend)
- Wails v2 CLI (для desktop)
- Docker 24+ (опционально)

## Сборка

```bash
# Frontend
cd frontend && npm ci && npm run build

# Web binary
go build -o bin/atp-web ./cmd/web

# Desktop
wails build
```

## Тесты

```bash
go test ./...
go test -cover ./internal/...
```

Покрытие: auth, tariff, order schedule, store CRUD, HTTP login, RBAC.

См. [setup-and-testing.md](setup-and-testing.md)

## Конфигурация

| Переменная / флаг | Описание |
|-------------------|----------|
| `ATP_DATA_DIR` | Корень данных LevelDB |
| `-data` | То же для `cmd/web` |
| `-addr` | Адрес HTTP (default `:8080`) |
| `-static` | Путь к `frontend/dist` |

## Документы

- [architecture.md](architecture.md)
- [decisions.md](decisions.md)
- [api-reference.md](api-reference.md)
