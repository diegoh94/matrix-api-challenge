# Matrix API

Go factoriza QR; Node calcula estadísticas. Frontend React para probar la API.

Factorización **QR matricial** (A = Q×R), no rotación ni código QR — ver [DECISIONS.md](docs/DECISIONS.md). Se usa [Gonum](https://gonum.org/v1/gonum) (`mat.QR`, Householder): correctitud numérica probada sin reimplementar el algoritmo.

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

## Autenticación (JWT)

Flujo mínimo implementado:

1. `POST /auth/token` con `{ "apiKey": "..." }` — go-api valida contra `API_KEY` y devuelve un JWT firmado (HS256, claims `sub`, `iat`, `exp`).
2. `POST /api/v1/matrix/qr` con `Authorization: Bearer <token>` — go-api valida el JWT y lo reenvía a node-api (mismo `JWT_SECRET`).

`/health` y `/auth/token` son públicos; el endpoint de negocio exige JWT.

### API key en el frontend

El frontend obtiene el token automáticamente usando `VITE_API_KEY` embebida al compilar. **Esto cumple el requisito JWT del challenge**, pero **no es una buena práctica en producción**: cualquiera puede ver la clave en DevTools o en el bundle JS.

| Cliente | Uso recomendado |
|---------|-----------------|
| curl, Swagger, integraciones | `POST /auth/token` con `apiKey` en el body |
| SPA (navegador) | La key no debería ir en el cliente; lo correcto sería un endpoint server-side (p. ej. `POST /auth/session`) donde go-api use `API_KEY` solo en el servidor |

Para este demo se prioriza simplicidad. En un entorno real, `API_KEY` y `JWT_SECRET` permanecen solo en go-api.

## Tests

```bash
cd go-api && go test ./...
cd node-api && npm test
```

## Docs

- [Arquitectura](docs/ARCHITECTURE.md)
- [Decisiones](docs/DECISIONS.md)
- [OpenAPI](docs/openapi.md)
