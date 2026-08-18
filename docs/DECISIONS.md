# Decisiones técnicas

Registro de decisiones relevantes del proyecto.

---

## QR vs rotación de matriz

**Contexto:** El enunciado menciona rotación en la arquitectura general, pero especifica factorización QR en los requisitos funcionales.

**Decisión:** Implementar factorización QR.

**Motivo:** Es el requerimiento con mayor especificidad. Produce matrices Q y R con un contrato claro para la comunicación con node-api.

---

## Go como orquestador

**Contexto:** Dos APIs con responsabilidades distintas.

**Decisión:** go-api es la API pública. node-api es servicio interno.

**Motivo:** El cliente interactúa con un solo endpoint. Go coordina QR y estadísticas en una respuesta atómica.

---

## Gonum para factorización QR

**Contexto:** Implementar Householder desde cero vs usar librería.

**Decisión:** Usar `gonum.org/v1/gonum/mat`.

**Motivo:** Correctitud numérica probada. En producción se priorizan librerías auditadas sobre reimplementaciones.

---

## Sin base de datos

**Contexto:** El challenge no menciona persistencia.

**Decisión:** No usar BD ni capa de repository.

**Motivo:** Todo el procesamiento es en memoria. Agregar repository sería abstracción sin port real.

---

## Respuesta atómica ante fallo de Node

**Contexto:** go-api puede devolver QR aunque Node falle.

**Decisión:** Si node-api no responde, go-api devuelve 502 sin QR parcial.

**Motivo:** Respuesta consistente para el cliente. Evita estados intermedios difíciles de interpretar.

---

## Matriz diagonal

**Contexto:** Valores flotantes de QR pueden no ser exactamente cero.

**Decisión:** Matriz cuadrada con off-diagonal ≈ 0 usando epsilon `1e-10`.

**Motivo:** Tolerancia numérica estándar para comparaciones con floating point.

---

## Autenticación JWT

**Contexto:** El challenge lista JWT como funcionalidad opcional.

**Decisión:** JWT en middleware HTTP de ambas APIs. Mismo `JWT_SECRET` compartido. Go reenvía el token al llamar a Node.

**Motivo:** Protege endpoints de negocio sin acoplar auth al dominio. `/health` y `/auth/token` permanecen públicos en go-api.
