# Matrix API Challenge

Solución al challenge de Interseguro: el cliente envía una matriz, go-api la procesa (factorización QR o rotación), node-api calcula estadísticas, y todo vuelve en una sola respuesta.

## Operaciones

| Operación | Endpoint | Descripción |
|-----------|----------|-------------|
| **QR** | `POST /api/v1/matrix/qr` | Descomposición A = Q×R con Gonum (Householder) |
| **Rotación** | `POST /api/v1/matrix/rotate` | 90°, 180° o 270° horario por reordenamiento O(m×n) |

Las decisiones de diseño están en [DECISIONS.md](docs/DECISIONS.md).

## Demo en vivo

| | |
|---|---|
| UI | https://frontend-production-69af.up.railway.app |
| API | https://go-api-production-c0df.up.railway.app |
| Swagger | https://go-api-production-c0df.up.railway.app/docs/index.html |

Railway corre tres servicios: `go-api` y `frontend` expuestos, `node-api` accesible solo por red privada.

El frontend ofrece dos acciones sobre la misma matriz de entrada: **Calcular QR** y **Rotar** (selector 90° / 180° / 270°). Comparte autenticación JWT con cache en memoria.

## Cómo está armado

```text
Cliente → go-api (Go + Fiber) ──HTTP──► node-api (Node + Express)
              │                              │
              │  QR (Gonum) o rotación       │  max, min, avg, sum, diagonal
              └──────── { input, result, statistics } ◄┘
```

Go es la API pública y orquesta el flujo. El cliente nunca llama a node-api directo.

Cada servicio tiene capas `domain` → `application` → `infrastructure`. La factorización usa [Gonum](https://gonum.org/v1/gonum) (`mat.QR`, Householder). La rotación vive en `domain/matrix_rotation.go` (reordenamiento por índices, sin trigonometría). Sin base de datos: todo el procesamiento es en memoria.

Si node-api no responde, go-api devuelve **502** y no entrega resultado parcial (ni QR ni matriz rotada sin estadísticas).

### Complejidad

| Parte | Complejidad | Notas |
|-------|-------------|-------|
| QR (Householder) | O(m·n²) | Cuello de botella en matrices grandes |
| Rotación | O(m·n) | Un pase por elemento |
| Estadísticas (Node) | O(m·n) | Sobre los valores enviados (Q+R o matriz rotada) |

## JWT

Implementado en ambas APIs con el mismo `JWT_SECRET`. Go reenvía el token a Node en cada llamada interna.

```text
POST /auth/token             →  { "apiKey": "..." }  →  JWT (HS256)
POST /api/v1/matrix/qr       →  Authorization: Bearer <token>
POST /api/v1/matrix/rotate   →  Authorization: Bearer <token>
```

Ejemplo rotación (default `degrees: 90` si se omite):

```json
{
  "matrix": [[1, 2], [3, 4], [5, 6]],
  "degrees": 180
}
```

El frontend obtiene el token al compilar (`VITE_API_KEY`) y lo **cachea en memoria** usando `expiresIn`, reutilizándolo en operaciones siguientes y renovándolo si recibe 401. En producción usaría `POST /auth/session` en go-api para no embeber la API key en el bundle.

## Correr en local

```bash
docker compose up --build
```

| | URL |
|---|-----|
| Frontend | http://localhost:3001 |
| go-api | http://localhost:8080 |
| Swagger | http://localhost:8080/docs/index.html |

Tras cambios en el frontend hay que **reconstruir la imagen** (`--build`); un simple restart no actualiza el bundle de Vite.

## Tests

```bash
cd go-api && go test ./...
cd node-api && npm test
```

## Documentación

- [Arquitectura](docs/ARCHITECTURE.md) — capas, flujos QR y rotación, Docker
- [Decisiones](docs/DECISIONS.md) — QR vs rotación, algoritmo, respuesta atómica, etc.
- [OpenAPI](docs/openapi.md) — contrato y regeneración de Swagger
