FROM golang:1.23-alpine AS builder

ENV CGO_ENABLED=0 \
	GOTOOLCHAIN=local

WORKDIR /app

COPY backend/src/go.mod backend/src/go.sum ./
RUN go mod download

COPY backend/src ./
RUN go build -o /bin/server ./main.go

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /bin/server /app/server
COPY backend/src/infra/migrations /app/infra/migrations

ENV DB_MIGRATIONS_PATH=/app/infra/migrations

EXPOSE 8080

CMD ["/app/server"]
