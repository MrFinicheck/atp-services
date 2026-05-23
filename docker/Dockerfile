# syntax=docker/dockerfile:1
# Контекст сборки: корень репозитория atp-services (см. docker-compose.yml)

FROM node:20-alpine AS frontend
WORKDIR /src/frontend
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm ci 2>/dev/null || npm install
COPY frontend/ ./
RUN npm run build

FROM golang:1.23-alpine AS backend
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /src/frontend/dist ./frontend/dist
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /atp-web ./cmd/web

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata wget
ENV TZ=Europe/Moscow
WORKDIR /app
COPY --from=backend /atp-web .
COPY --from=frontend /src/frontend/dist ./frontend/dist
ENV ATP_DATA_DIR=/data
VOLUME /data
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget -qO- http://127.0.0.1:8080/api/health || exit 1
ENTRYPOINT ["/app/atp-web"]
CMD ["-addr", ":8080", "-static", "frontend/dist", "-data", "/data"]
