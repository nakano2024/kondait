# kondait

## Backend

### Testing

#### Docker (recommended)
```
docker compose up --build --abort-on-container-exit --exit-code-from backend-test
```

#### Local (PostgreSQL running on the host)
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

## Frontend

### Setup
```
cd frontend
npm install
```

### Development
```
cd frontend
npm run dev
```

### Production build
```
cd frontend
npm run build
npm run preview
```

### Docker (dev)
```
docker compose -f frontend/docker-compose.yml up --build
```
