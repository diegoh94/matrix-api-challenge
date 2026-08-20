# OpenAPI — go-api

La documentación de **go-api** se genera automáticamente con [swaggo/swag](https://github.com/swaggo/swag).

## Endpoints documentados

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/health` | Health check |
| POST | `/auth/token` | Intercambio API key → JWT |
| POST | `/api/v1/matrix/qr` | Factorización QR + estadísticas |
| POST | `/api/v1/matrix/rotate` | Rotación 90/180/270° + estadísticas |

## Archivos generados

| Archivo | Descripción |
|---------|-------------|
| `go-api/docs/swagger.yaml` | Spec OpenAPI (fuente generada) |
| `go-api/docs/swagger.json` | Spec JSON |
| `go-api/docs/docs.go` | Embed para runtime |

## Regenerar tras cambiar handlers o DTOs

```bash
cd go-api
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal
```

Sin Go local (Docker):

```bash
cd go-api
docker run --rm -v "%cd%:/app" -w /app golang:1.24-alpine sh -c \
  "go install github.com/swaggo/swag/cmd/swag@latest && swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal"
```

## Swagger UI

Con el servicio corriendo: http://localhost:8080/docs/index.html

Producción: https://go-api-production-c0df.up.railway.app/docs/index.html

## node-api

Contrato estático en [openapi-node.yaml](openapi-node.yaml) — `POST /api/v1/statistics` con `{ matrices: [{ name, data }] }`.
