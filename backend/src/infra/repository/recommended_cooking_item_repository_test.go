package repository

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"kondait-backend/domain/aggregation"
	"kondait-backend/domain/entity"
)

const (
	testUserCode = "11111111-1111-1111-1111-111111111111"
)

func TestRecommendedCookingItemRepository_FetchByUserCode_Normal(t *testing.T) {
	testTable := []struct {
		name     string
		ctx      context.Context
		userCode string
		expected *aggregation.RecommendedCookingItemList
	}{
		{
			name:     "order by cook_count, then last_cooked_date (nulls first), take top 5",
			ctx:      context.WithValue(context.Background(), "ctx-key-1", "ctx-1"),
			userCode: testUserCode,
			expected: &aggregation.RecommendedCookingItemList{
				Items: []*entity.RecommendedCookingItem{
					{
						Code:           "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1",
						Name:           "Rice",
						CookCount:      0,
						LastCookedDate: time.Time{},
					},
					{
						Code:           "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb2",
						Name:           "Pasta",
						CookCount:      0,
						LastCookedDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						Code:           "cccccccc-cccc-cccc-cccc-ccccccccccc3",
						Name:           "Salad",
						CookCount:      1,
						LastCookedDate: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
					},
					{
						Code:           "dddddddd-dddd-dddd-dddd-ddddddddddd4",
						Name:           "Soup",
						CookCount:      1,
						LastCookedDate: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
					},
					{
						Code:           "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeee5",
						Name:           "Bread",
						CookCount:      2,
						LastCookedDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			migrateDB(t)
			db := openTestDB(t)

			tx := db.Begin()
			t.Cleanup(func() {
				_ = tx.Rollback().Error
			})

			loadTestData(t, tx)

			repo := NewRecommendedCookingItemRepository(tx)
			got, err := repo.FetchByUserCode(tt.ctx, tt.userCode)
			require.NoError(t, err)
			require.NotNil(t, got)

			assert.Len(t, got.Items, 5)

			assert.Equal(t, tt.expected, got)
		})
	}
}

func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sslMode := map[bool]string{true: os.Getenv("DB_SSLMODE"), false: "disable"}[os.Getenv("DB_SSLMODE") != ""]

	require.NotEmpty(t, host, "DB_HOST is not set")
	require.NotEmpty(t, port, "DB_PORT is not set")
	require.NotEmpty(t, user, "DB_USER is not set")
	require.NotEmpty(t, password, "DB_PASSWORD is not set")
	require.NotEmpty(t, dbName, "DB_NAME is not set")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		user,
		password,
		dbName,
		sslMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	return db
}

func migrateDB(t *testing.T) {
	t.Helper()

	migrationsPath := os.Getenv("DB_MIGRATIONS_PATH")
	require.NotEmpty(t, migrationsPath, "DB_MIGRATIONS_PATH is not set")

	absMigrationsPath, err := filepath.Abs(migrationsPath)
	require.NoError(t, err)

	sourceURL := "file://" + absMigrationsPath

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sslMode := map[bool]string{true: os.Getenv("DB_SSLMODE"), false: "disable"}[os.Getenv("DB_SSLMODE") != ""]

	require.NotEmpty(t, host, "DB_HOST is not set")
	require.NotEmpty(t, port, "DB_PORT is not set")
	require.NotEmpty(t, user, "DB_USER is not set")
	require.NotEmpty(t, password, "DB_PASSWORD is not set")
	require.NotEmpty(t, dbName, "DB_NAME is not set")

	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   dbName,
	}
	q := u.Query()
	q.Set("sslmode", sslMode)
	u.RawQuery = q.Encode()
	dbURL := u.String()

	migrator, err := migrate.New(sourceURL, dbURL)
	require.NoError(t, err)

	err = migrator.Up()
	require.True(t, err == nil || err == migrate.ErrNoChange, "migrate.Up: %v", err)
}

func loadTestData(t *testing.T, db *gorm.DB) {
	t.Helper()

	testdataPath := os.Getenv("TESTDATA_PATH")
	require.NotEmpty(t, testdataPath, "TESTDATA_PATH is not set")

	testdataFile := filepath.Join(testdataPath, "recommended_cooking_item_repository.sql")
	content, err := os.ReadFile(testdataFile)
	require.NoError(t, err)

	require.NoError(t, db.Exec(string(content)).Error)
}
