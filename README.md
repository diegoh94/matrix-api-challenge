# Matrix API Challenge

Solución al challenge de Interseguro: el cliente envía una matriz, go-api la factoriza en Q y R, node-api calcula estadísticas sobre ambas, y todo vuelve en una sola respuesta.

Cuando el enunciado dice **QR**, hablamos de factorización matricial (A = Q×R), no de códigos QR ni rotación. Las decisiones de diseño están en [DECISIONS.md](docs/DECISIONS.md).

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
              │  QR con Gonum                │  max, min, avg, sum, diagonal
              └──────── { input, qr, statistics } ◄┘
```

Go es la API pública y orquesta el flujo. El cliente nunca llama a node-api directo.

Cada servicio tiene capas `domain` → `application` → `infrastructure`. La factorización usa [Gonum](https://gonum.org/v1/gonum) (`mat.QR`, Householder) — preferí una librería numérica probada antes que reimplementar el algoritmo. Sin base de datos: el procesamiento es en memoria.

Si node-api no responde, go-api devuelve **502** y no entrega un QR incompleto.

El cuello de botella matemático es la factorización QR (O(m·n²) con Householder); las estadísticas en Node son O(m·n) y no pesan en matrices pequeñas. Para el tamaño típico del challenge no hace falta optimizar, aunque habría margen: un solo pase para max/min/sum, menos copias al convertir entre `[][]float64` y `mat.Dense`, o serialización binaria si Q y R fueran muy grandes.

## JWT

Implementado en ambas APIs con el mismo `JWT_SECRET`. Go reenvía el token a Node en cada llamada interna.

```text
POST /auth/token   →  { "apiKey": "..." }  →  JWT (HS256)
POST /api/v1/matrix/qr  →  Authorization: Bearer <token>
```

El frontend obtiene el token automáticamente al compilar (`VITE_API_KEY`). Sirve para la demo; en un entorno real la API key no iría en el bundle del navegador — quedaría solo en el servidor.

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
