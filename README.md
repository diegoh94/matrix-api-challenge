# Matrix API Challenge

Solución para el coding challenge de Interseguro: dos APIs REST que reciben una matriz, calculan su factorización QR y devuelven estadísticas sobre las matrices resultantes.

## Despliegue en producción (Railway)

| | |
|---|---|
| **API pública** | https://go-api-production-c0df.up.railway.app |
| **Health** | https://go-api-production-c0df.up.railway.app/health |
| **Swagger** | https://go-api-production-c0df.up.railway.app/docs/index.html |

Arquitectura en cloud: dos servicios Railway (`go-api` expuesto, `node-api` en red privada). El cliente solo interactúa con go-api.

```bash
# Token (usar la API_KEY configurada en Railway)
curl -X POST https://go-api-production-c0df.up.railway.app/auth/token \
  -H "Content-Type: application/json" \
  -d '{"apiKey": "<API_KEY>"}'

# Factorización QR + estadísticas
curl -X POST https://go-api-production-c0df.up.railway.app/api/v1/matrix/qr \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"matrix": [[1, 2], [3, 4], [5, 6]]}'
```

## Cobertura del desafío

| Requisito | Estado |
|-----------|--------|
| API Go + Fiber — matriz → factorización QR | ✅ |
| API Node + Express — estadísticas sobre Q y R | ✅ |
| max, min, promedio, suma, matriz diagonal | ✅ |
| Comunicación HTTP entre APIs | ✅ |
| Docker + Docker Compose | ✅ |
| Despliegue en cloud | ✅ Railway |
| Documentación | ✅ README, arquitectura, decisiones, OpenAPI |
| JWT (opcional) | ✅ |
| Tests (opcional) | ✅ dominio, use cases, HTTP e integración Go↔Node |
| Frontend (opcional) | — |

**Decisión documentada:** el enunciado menciona rotación en la arquitectura general, pero especifica QR en los requisitos funcionales. Ver [docs/DECISIONS.md](docs/DECISIONS.md).

**Adicional (no pedido explícitamente):** Clean Architecture, Swagger UI, respuesta atómica (502 si Node falla, sin QR parcial), OpenAPI generado con swaggo.

## Stack

| Servicio | Tecnología | Rol |
|----------|------------|-----|
| `go-api` | Go 1.24 + Fiber | API pública, factorización QR, orquestación |
| `node-api` | Node.js 24 + Express | Servicio interno de estadísticas |
| Infra | Docker + Docker Compose | Ejecución local |
| Cloud | Railway | Producción (2 servicios) |

## Inicio rápido (local)

```bash
docker compose up --build
```

Swagger local: http://localhost:8080/docs/index.html

```bash
docker compose down
```

## Tests

Pruebas sobre lógica de negocio e integración (no smoke tests vacíos).

**go-api** — validación de matrices, QR (Q×R = original), orquestación use case, JWT, gateway HTTP con reenvío de token.

**node-api** — use case Q+R, detección diagonal con tolerancia numérica, endpoint HTTP + JWT.

```bash
# go-api
cd go-api && go test ./...

# node-api (Node >= 24)
cd node-api && npm test
```

Con Docker (sin Go/Node local):

```bash
docker run --rm -v "%cd%/go-api:/app" -w /app golang:1.24-alpine go test ./...
docker run --rm -v "%cd%/node-api:/app" -w /app node:24-alpine npm test
```

## API

Contratos, schemas y ejemplos en **Swagger**:

| Entorno | Swagger |
|---------|---------|
| Producción | https://go-api-production-c0df.up.railway.app/docs/index.html |
| Local | http://localhost:8080/docs/index.html |

Endpoints públicos (go-api): `GET /health`, `POST /auth/token`, `POST /api/v1/matrix/qr`.

node-api es interno (`POST /api/v1/statistics`); no expuesto al cliente.

## Variables de entorno

Ver `.env.example`. En Railway: `JWT_SECRET` compartido en ambos servicios; `NODE_API_URL` en go-api apunta a node-api por red privada.

| Variable | Servicio | Descripción |
|----------|----------|-------------|
| `API_KEY` | go-api | Clave para `/auth/token` |
| `JWT_SECRET` | ambos | Firma/validación JWT |
| `NODE_API_URL` | go-api | URL interna de node-api |
| `JWT_EXPIRATION_HOURS` | go-api | Duración del token |

## Documentación

- [Arquitectura](docs/ARCHITECTURE.md)
- [Decisiones técnicas](docs/DECISIONS.md)
- [OpenAPI — generación](docs/openapi.md)

Regenerar spec go-api: ver [docs/openapi.md](docs/openapi.md).
