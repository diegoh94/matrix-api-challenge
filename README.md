# Matrix API Challenge

Solución para el coding challenge de Interseguro: dos APIs REST que reciben una matriz, calculan su factorización QR y devuelven estadísticas sobre las matrices resultantes.

## Stack

| Servicio | Tecnología | Rol |
|----------|------------|-----|
| `go-api` | Go 1.24 + Fiber | API pública, factorización QR, orquestación |
| `node-api` | Node.js 24 + Express | Servicio interno de estadísticas |
| Infra | Docker + Docker Compose | Ejecución local y despliegue |

## Requisitos

- Docker Desktop
- Node.js >= 24 (solo desarrollo local sin Docker)
- Go 1.24 (solo desarrollo local sin Docker)

## Inicio rápido

```bash
docker compose up --build
```

Verificar:

```bash
curl http://localhost:8080/health
```

Obtener token JWT:

```bash
curl -X POST http://localhost:8080/auth/token \
  -H "Content-Type: application/json" \
  -d '{"apiKey": "change-me-api-key"}'
```

Procesar una matriz (usar el token obtenido):

```bash
curl -X POST http://localhost:8080/api/v1/matrix/qr \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"matrix": [[1, 2], [3, 4], [5, 6]]}'
```

Detener:

```bash
docker compose down
```

## API pública (go-api)

### `GET /health`

```json
{ "status": "ok", "service": "go-api" }
```

### `POST /auth/token`

Emite un JWT válido para consumir la API.

**Request**

```json
{ "apiKey": "change-me-api-key" }
```

**Response 200**

```json
{ "token": "...", "expiresIn": 86400 }
```

### `POST /api/v1/matrix/qr`

Requiere header `Authorization: Bearer <token>` cuando JWT está habilitado.

**Request**

```json
{
  "matrix": [
    [1, 2],
    [3, 4],
    [5, 6]
  ]
}
```

**Response 200**

```json
{
  "input": { "rows": 3, "cols": 2 },
  "qr": {
    "Q": [[...], [...], [...]],
    "R": [[...], [...], [...]]
  },
  "statistics": {
    "max": 0.897085,
    "min": -7.437357,
    "average": -0.881237,
    "sum": -13.218558,
    "hasDiagonalMatrix": false,
    "diagonalMatrices": []
  }
}
```

**Errores**

| Código | Caso |
|--------|------|
| 400 | JSON inválido o matriz inválida |
| 422 | Matriz no factorizable |
| 502 | Servicio de estadísticas no disponible |
| 500 | Error interno |

## API interna (node-api)

Expuesta solo en la red Docker. Go la consume vía HTTP.

### `POST /api/v1/statistics`

Recibe las matrices Q y R, devuelve max, min, average, sum y detección de matriz diagonal.

## Variables de entorno

Copiar `.env.example` como referencia.

| Variable | Servicio | Default | Descripción |
|----------|----------|---------|-------------|
| `GO_API_PORT` | go-api | `8080` | Puerto HTTP |
| `NODE_API_URL` | go-api | `http://localhost:3000` | URL de node-api |
| `NODE_API_TIMEOUT_MS` | go-api | `5000` | Timeout hacia node-api |
| `NODE_API_PORT` | node-api | `3000` | Puerto HTTP |
| `JWT_SECRET` | ambos | — | Secreto compartido para firmar/validar JWT |
| `JWT_EXPIRATION_HOURS` | go-api | `24` | Duración del token |
| `API_KEY` | go-api | — | Clave para obtener token en `/auth/token` |

En Docker Compose, `NODE_API_URL` se configura como `http://node-api:3000`.

## Estructura del proyecto

```text
matrix-api-challenge/
├── go-api/                 # Go + Fiber
│   ├── cmd/server/         # entry point
│   └── internal/
│       ├── domain/         # entidades, ports
│       ├── application/    # use cases
│       └── infrastructure/ # HTTP, Gonum, gateway Node
├── node-api/               # Node + Express
│   └── src/
│       ├── domain/
│       ├── application/
│       └── infrastructure/
├── docker-compose.yml
└── docs/
    ├── ARCHITECTURE.md
    └── DECISIONS.md
```

## Desarrollo local (sin Docker)

**node-api**

```bash
cd node-api
npm install
npm run dev
```

**go-api**

```bash
cd go-api
go run ./cmd/server
```

Asegurar `NODE_API_URL=http://localhost:3000`.

## Documentación

- [Arquitectura](docs/ARCHITECTURE.md)
- [Decisiones técnicas](docs/DECISIONS.md)
- [OpenAPI — generación](docs/openapi.md)
- **Swagger UI (go-api):** http://localhost:8080/docs/index.html

### Regenerar OpenAPI go-api

```bash
cd go-api
swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal
```

## Próximos pasos

- Despliegue en cloud (Railway / Render)
- Frontend (opcional)
