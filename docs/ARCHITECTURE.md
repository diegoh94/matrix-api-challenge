# Arquitectura

## Visión general

Tres servicios HTTP: frontend (demo), go-api (pública, orquestación) y node-api (interna, estadísticas).

```text
Cliente / Frontend
        │
        ▼
go-api (Fiber) ──HTTP──► node-api (Express)
   │                           │
   │  QR (Gonum)               │  max, min, avg, sum, diagonal
   │  o rotación (domain)      │
   └──── respuesta unificada ◄─┘
         { input, qr | rotated, statistics }
```

El cliente **no** llama a node-api directamente.

## Capas (Clean Architecture)

Ambas APIs siguen la misma estructura:

```text
infrastructure/http  →  adaptador de entrada (Fiber / Express)
application/         →  use cases
domain/              →  entidades, reglas, ports
infrastructure/*     →  adaptadores de salida (Gonum, HTTP client)
```

### go-api

| Capa | Responsabilidad |
|------|-----------------|
| `domain` | Matrix, QRDecomposition, `RotateMatrix()`, validación, errores |
| `domain/ports` | `QRFactorizer`, `StatisticsGateway` |
| `application` | `FactorizeMatrixUseCase`, `RotateMatrixUseCase` |
| `infrastructure/qr` | `GonumQRFactorizer` |
| `infrastructure/statistics` | `NodeStatisticsGateway` (payload genérico `NamedMatrix[]`) |
| `infrastructure/http` | Handlers Fiber, DTOs, mapeo de errores |
| `cmd/server` | Composition root — wiring de dependencias |

### node-api

| Capa | Responsabilidad |
|------|-----------------|
| `domain` | Matrix, MatrixStatistics, errores |
| `application` | `ComputeMatrixStatisticsUseCase` |
| `infrastructure/services` | StatisticsCalculator, DiagonalMatrixChecker |
| `infrastructure/http` | Routes, controller, DTOs Express |
| `main.js` | Composition root |

node-api **no conoce** QR ni rotación: recibe `{ matrices: [{ name, data }] }` y calcula estadísticas sobre todos los valores.

## Flujo QR

1. Cliente envía `POST /api/v1/matrix/qr` con una matriz.
2. Go valida la matriz (`domain.NewMatrix`).
3. Go calcula QR con Gonum → matrices Q y R.
4. Go envía Q y R a Node vía `POST /api/v1/statistics`.
5. Node calcula estadísticas sobre todos los valores de Q y R.
6. Go devuelve `{ input, qr, statistics }`.

## Flujo rotación

1. Cliente envía `POST /api/v1/matrix/rotate` con `{ matrix, degrees }` (`degrees`: 90, 180 o 270; default 90).
2. Go valida la matriz.
3. Go rota en dominio (`domain.RotateMatrix`) — reordenamiento O(m×n).
4. Go envía la matriz rotada a Node como `{ name: "rotated", data: [...] }`.
5. Go devuelve `{ input, degrees, rotated, statistics }`.

## Comunicación HTTP

| Origen | Destino | Endpoint |
|--------|---------|----------|
| Cliente | go-api | `POST /api/v1/matrix/qr` |
| Cliente | go-api | `POST /api/v1/matrix/rotate` |
| go-api | node-api | `POST /api/v1/statistics` |

Timeout configurable: `NODE_API_TIMEOUT_MS` (default 5s).

Si Node falla en cualquiera de los flujos, Go responde **502** sin devolver resultado parcial.

## Frontend

React + Vite. Dos acciones sobre la misma matriz:

- **Calcular QR** → `POST /api/v1/matrix/qr`
- **Rotar** → selector 90° / 180° / 270° + `POST /api/v1/matrix/rotate`

JWT cacheado en memoria (`expiresIn`, retry en 401). Ver `frontend/src/api.js` y `App.jsx`.

## Docker

```text
docker-compose
├── node-api   (red interna, puerto 3000)
├── go-api     (puerto 8080 expuesto al host)
└── frontend   (puerto 3001 expuesto; build Vite con VITE_API_URL / VITE_API_KEY)
```

go-api espera a que node-api esté healthy antes de iniciar. El frontend espera a go-api.

## Decisiones de diseño

Ver [DECISIONS.md](DECISIONS.md) para el detalle de cada decisión.

## Documentación API

| Recurso | Ubicación |
|---------|-----------|
| Swagger UI (go-api) | http://localhost:8080/docs/index.html |
| OpenAPI generado (go-api) | `go-api/docs/swagger.yaml` |
| Guía regeneración | [docs/openapi.md](openapi.md) |
