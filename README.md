# Matrix API Challenge

Solución al challenge de Interseguro: el cliente envía una matriz, go-api la procesa (factorización QR o rotación), node-api calcula estadísticas, y todo vuelve en una sola respuesta.

Operaciones disponibles:
- **QR** — descomposición A = Q×R con Gonum
- **Rotación** — 90°, 180° o 270° horario por reordenamiento O(m×n)

Las decisiones de diseño están en [DECISIONS.md](docs/DECISIONS.md).

## Demo en vivo

| | |
|---|---|
| UI | https://frontend-production-69af.up.railway.app |
| API | https://go-api-production-c0df.up.railway.app |
| Swagger | https://go-api-production-c0df.up.railway.app/docs/index.html |

Railway corre tres servicios: `go-api` y `frontend` expuestos, `node-api` accesible solo por red privada.

## Cómo está armado

```text
Cliente → go-api (Go + Fiber) ──HTTP──► node-api (Node + Express)
              │                              │
              │  QR (Gonum) o rotación       │  max, min, avg, sum, diagonal
              └──────── { input, result, statistics } ◄┘
```

Go es la API pública y orquesta el flujo. El cliente nunca llama a node-api directo.

Cada servicio tiene capas `domain` → `application` → `infrastructure`. La factorización usa [Gonum](https://gonum.org/v1/gonum) (`mat.QR`, Householder) — preferí una librería numérica probada antes que reimplementar el algoritmo. Sin base de datos: el procesamiento es en memoria.

Si node-api no responde, go-api devuelve **502** y no entrega un QR incompleto.

El cuello de botella matemático es la factorización QR (O(m·n²) con Householder); las estadísticas en Node son O(m·n) y no pesan en matrices pequeñas. Para el tamaño típico del challenge no hace falta optimizar, aunque habría margen: un solo pase para max/min/sum, menos copias al convertir entre `[][]float64` y `mat.Dense`, o serialización binaria si Q y R fueran muy grandes.

## JWT

Implementado en ambas APIs con el mismo `JWT_SECRET`. Go reenvía el token a Node en cada llamada interna.

```text
POST /auth/token          →  { "apiKey": "..." }  →  JWT (HS256)
POST /api/v1/matrix/qr    →  Authorization: Bearer <token>
POST /api/v1/matrix/rotate →  Authorization: Bearer <token>  (degrees: 90|180|270, default 90)
```

El frontend obtiene el token automáticamente al compilar (`VITE_API_KEY`) y lo **cachea en memoria** usando `expiresIn` del backend, reutilizándolo en operaciones siguientes y renovándolo si recibe 401. En producción usaría `POST /auth/session` en go-api para no embeber la API key en el bundle.

## Correr en local

```bash
docker compose up --build
```

| | URL |
|---|-----|
| Frontend | http://localhost:3001 |
| go-api | http://localhost:8080 |
| Swagger | http://localhost:8080/docs/index.html |

## Tests

```bash
cd go-api && go test ./...
cd node-api && npm test
```

## Documentación

- [Arquitectura](docs/ARCHITECTURE.md) — capas, flujo de petición, Docker
- [Decisiones](docs/DECISIONS.md) — por qué QR, Gonum, respuesta atómica, etc.
- [OpenAPI](docs/openapi.md) — contrato y regeneración de Swagger
