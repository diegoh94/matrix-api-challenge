# Arquitectura

## Visión general

Dos microservicios HTTP desacoplados. Go expone la API pública y orquesta el flujo. Node procesa estadísticas como servicio interno.

```text
Cliente
   │
   ▼
go-api (Fiber) ──HTTP──► node-api (Express)
   │                           │
   │  QR (Gonum)               │  max, min, avg, sum, diagonal
   │                           │
   └──── respuesta unificada ◄─┘
         { input, qr, statistics }
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
| `domain` | Matrix, QRDecomposition, validación, errores |
| `domain/ports` | `QRFactorizer`, `StatisticsGateway` |
| `application` | `FactorizeMatrixUseCase` — orquesta QR + stats |
| `infrastructure/qr` | `GonumQRFactorizer` |
| `infrastructure/statistics` | `NodeStatisticsGateway` |
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

## Flujo de una petición

1. Cliente envía `POST /api/v1/matrix/qr` con una matriz.
2. Go valida la matriz (`domain.NewMatrix`).
3. Go calcula QR con Gonum → matrices Q y R.
4. Go envía Q y R a Node vía `POST /api/v1/statistics`.
5. Node calcula estadísticas sobre todos los valores de Q y R.
6. Go devuelve QR + statistics al cliente.

Si Node falla, Go responde **502** sin devolver QR parcial.

## Comunicación HTTP

| Origen | Destino | Endpoint |
|--------|---------|----------|
| Cliente | go-api | `POST /api/v1/matrix/qr` |
| go-api | node-api | `POST /api/v1/statistics` |

Timeout configurable: `NODE_API_TIMEOUT_MS` (default 5s).

## Docker

```text
docker-compose
├── node-api   (red interna, puerto 3000)
└── go-api     (puerto 8080 expuesto al host)
```

go-api espera a que node-api esté healthy antes de iniciar.

## Decisiones de diseño

Ver [DECISIONS.md](DECISIONS.md) para el detalle de cada decisión.
