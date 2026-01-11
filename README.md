# kondait

## Testing

### Docker (recommended)
```
docker compose -f docker-compose.testing.yml up --build --abort-on-container-exit --exit-code-from backend-test
```

### Local (PostgreSQL running on the host)
```
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=kondait_test
export DB_SSLMODE=disable
export DB_MIGRATIONS_PATH=./backend/src/infra/migrations
export TESTDATA_PATH=./backend/src/infra/testdata

cd backend/src
go test ./...
```
