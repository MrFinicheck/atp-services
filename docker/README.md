# Docker — веб-режим TransitOS

## Требования

- Docker 24+
- Docker Compose v2

## Сборка и запуск

Из корня репозитория:

```bash
docker compose -f docker/docker-compose.yml up --build -d
```

Или из этой папки:

```bash
cd docker
docker compose up --build -d
```

Откройте http://localhost:8080

Логи:

```bash
docker compose -f docker/docker-compose.yml logs -f atp-web
```

Остановка:

```bash
docker compose -f docker/docker-compose.yml down
```

Данные LevelDB сохраняются в volume `atp-data`.

## Переменные

| Переменная | Значение в контейнере |
|------------|----------------------|
| `ATP_DATA_DIR` | `/data` |
| `TZ` | `Europe/Moscow` |

Порт хоста: `8080` (измените в `docker-compose.yml` при конфликте).

## Сборка без Compose

```bash
docker build -f docker/Dockerfile -t atp-web:latest .
docker run --rm -p 8080:8080 -v atp-data:/data atp-web:latest
```

## Healthcheck

`GET /api/health` → `{"status":"ok"}`

## Устранение неполадок

| Проблема | Решение |
|----------|---------|
| Пустой UI | `docker compose build --no-cache` |
| Сброс демо-данных | `docker compose down -v` (удалит volume) |
