# Matrix API

Go factoriza QR; Node calcula estadísticas. Frontend React para probar la API.

## Producción (Railway)

| | |
|---|---|
| API | https://go-api-production-c0df.up.railway.app |
| Swagger | https://go-api-production-c0df.up.railway.app/docs/index.html |

Tres servicios: `go-api` (público), `node-api` (interno), `frontend` (UI).

## Local

```bash
docker compose up --build
```

| | URL |
|---|-----|
| Frontend | http://localhost:3001 |
| go-api | http://localhost:8080 |
| Swagger | http://localhost:8080/docs/index.html |

## Frontend en Railway

Root Directory: `frontend` · Dockerfile · Generate Domain

```env
VITE_API_URL=https://go-api-production-c0df.up.railway.app
VITE_API_KEY=<misma API_KEY que go-api>
```

Variables embebidas al compilar (`npm run build`).

## Tests

```bash
cd go-api && go test ./...
cd node-api && npm test
```

## Docs

- [Arquitectura](docs/ARCHITECTURE.md)
- [Decisiones](docs/DECISIONS.md)
- [OpenAPI](docs/openapi.md)
