FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:3.21 AS runtime

WORKDIR /app

RUN apk add --no-cache ca-certificates wget \
	&& addgroup -S goapi \
	&& adduser -S goapi -G goapi

COPY --from=builder /app/server ./server

ENV GO_API_PORT=8080
ENV NODE_API_URL=http://node-api:3000
ENV NODE_API_TIMEOUT_MS=5000

USER goapi

EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
	CMD wget -qO- http://127.0.0.1:8080/health || exit 1

CMD ["./server"]
