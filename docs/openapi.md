# OpenAPI — go-api

La documentación de **go-api** se genera automáticamente con [swaggo/swag](https://github.com/swaggo/swag).

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

## Swagger UI

Con el servicio corriendo: http://localhost:8080/docs/index.html

## node-api

Ver [openapi-node.md](openapi-node.md) — spec generada con swagger-jsdoc.
