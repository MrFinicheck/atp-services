# АТП — Система учёта услуг автотранспортного предприятия

Desktop-приложение (Wails) и веб-режим для учёта заявок, путевых листов и контроля расхода топлива.

## Стек

- **Backend:** Go, LevelDB
- **Desktop:** Wails v2
- **Frontend:** TypeScript, Vite, Bootstrap 5 (адаптивный UI для ПК и смартфона)

## Возможности

- Роли: администратор, диспетчер, водитель
- Автоматический расчёт стоимости заявки по тарифу
- График занятости автомобилей
- Мобильный интерфейс водителя: начало/конец рейса
- Блокировка закрытия смены при перерасходе топлива без комментария
- Отчёты и журнал аудита

## Демо-учётные записи

| Логин | Пароль | Роль |
|-------|--------|------|
| admin | admin123 | Администратор |
| dispatcher | disp123 | Диспетчер |
| driver1 | drv123 | Водитель |

## Запуск (desktop)

```bash
cd frontend && npm install && npm run build
cd .. && wails dev
```

Сборка:

```bash
wails build
```

## Запуск (веб)

```bash
cd frontend && npm run build
cd .. && go run ./cmd/web -addr :8080
```

Откройте http://localhost:8080 в браузере (на телефоне — по IP сервера в локальной сети).

Данные LevelDB по умолчанию хранятся в `./data` внутри каталога проекта (удобно для `wails dev`). Переопределение: переменная `ATP_DATA_DIR` или флаг `-data` для веб-сервера.

Если видите `resource temporarily unavailable` или «каталог данных занят»:

```bash
pkill -f "atp-services" 2>/dev/null
pkill -f "cmd/web" 2>/dev/null
rm -f data/.atp.lock data/data/LOCK
wails dev -tags=webkit2_41
```

## Запуск (Docker)

Веб-режим в контейнере: сборка frontend, Go API и LevelDB в volume `atp-data`.

**Требования:** Docker 24+, Docker Compose v2.

```bash
# из корня репозитория
docker compose -f docker/docker-compose.yml up --build -d
```

Откройте http://localhost:8080 (те же демо-логины).

Полезные команды:

```bash
# логи
docker compose -f docker/docker-compose.yml logs -f atp-web

# остановка
docker compose -f docker/docker-compose.yml down

# сброс данных (удалит демо-базу в volume)
docker compose -f docker/docker-compose.yml down -v
```

Сборка образа без Compose:

```bash
docker build -f docker/Dockerfile -t atp-web:latest .
docker run --rm -p 8080:8080 -v atp-data:/data atp-web:latest
```

Подробности: [docker/README.md](docker/README.md).

## Структура проекта

```
atp-services/
├── cmd/web/           # HTTP-сервер для веб-режима
├── docker/            # Dockerfile и docker-compose
├── docs/              # руководства пользователя и разработчика
├── internal/
│   ├── api/           # REST API
│   ├── app/           # Бизнес-логика
│   ├── models/
│   ├── services/
│   └── store/         # LevelDB
├── frontend/          # TypeScript UI
├── app.go             # Wails bindings
└── main.go
```

## Документация и тесты

| Материал | Путь |
|----------|------|
| Оглавление документации | [docs/README.md](docs/README.md) |
| Руководство пользователя (все роли) | [docs/user-guide-overview.md](docs/user-guide-overview.md) |
| Роли (детально) | [docs/user-guide/](docs/user-guide/) |
| Документация разработчика | [docs/developer/](docs/developer/) |
| Примеры API | [docs/examples/](docs/examples/) |

**Unit-тесты:**

```bash
go test ./...
```
