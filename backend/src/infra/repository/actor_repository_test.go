package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"kondait-backend/domain/entity"
	"kondait-backend/infra/model"
)

func TestActorRepository_FetchBySub(t *testing.T) {
	testTable := []struct {
		name     string
		ctx      context.Context
		sub      string
		seedData func(t *testing.T, db *gorm.DB)
		expected *entity.Actor
	}{
		{
			name: "ユーザーが存在する場合、Actorが取得できること",
			ctx:  context.WithValue(context.Background(), "ctx-key-1", "ctx-1"),
			sub:  "sub-1",
			seedData: func(t *testing.T, db *gorm.DB) {
				require.NoError(t, db.Exec("INSERT INTO users (code, sub) VALUES (?, ?)", "11111111-1111-1111-1111-111111111111", "sub-1").Error)
			},
			expected: &entity.Actor{
				Code: "11111111-1111-1111-1111-111111111111",
				Sub:  "sub-1",
			},
		},
		{
			name: "ユーザーが存在しない場合、nilが返ること",
			ctx:  context.WithValue(context.Background(), "ctx-key-2", "ctx-2"),
			sub:  "sub-2",
			seedData: func(t *testing.T, db *gorm.DB) {
				require.NoError(t, db.Exec("INSERT INTO users (code, sub) VALUES (?, ?)", "22222222-2222-2222-2222-222222222222", "sub-3").Error)
			},
			expected: (*entity.Actor)(nil),
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

			tt.seedData(t, tx)

			repo := NewActorRepository(tx)
			got, err := repo.FetchBySub(tt.ctx, tt.sub)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestActorRepository_Save(t *testing.T) {
	testTable := []struct {
		name          string
		ctx           context.Context
		actor         *entity.Actor
		seedData      func(t *testing.T, db *gorm.DB)
		expectedUser  model.User
		expectedCount int64
	}{
		{
			name:  "新規の場合、レコードが作成されること",
			ctx:   context.WithValue(context.Background(), "ctx-key-3", "ctx-3"),
			actor: &entity.Actor{Code: "33333333-3333-3333-3333-333333333333", Sub: "sub-1"},
			seedData: func(t *testing.T, db *gorm.DB) {
				require.NoError(t, db.Exec("DELETE FROM users").Error)
			},
			expectedUser:  model.User{Code: "33333333-3333-3333-3333-333333333333", Sub: "sub-1"},
			expectedCount: 1,
		},
		{
			name:  "codeが存在する場合、subが更新されること",
			ctx:   context.WithValue(context.Background(), "ctx-key-4", "ctx-4"),
			actor: &entity.Actor{Code: "44444444-4444-4444-4444-444444444444", Sub: "sub-2"},
			seedData: func(t *testing.T, db *gorm.DB) {
				require.NoError(t, db.Exec("INSERT INTO users (code, sub) VALUES (?, ?)", "44444444-4444-4444-4444-444444444444", "sub-old").Error)
			},
			expectedUser:  model.User{Code: "44444444-4444-4444-4444-444444444444", Sub: "sub-2"},
			expectedCount: 1,
		},
		{
			name:  "subが存在する場合、codeが更新されること",
			ctx:   context.WithValue(context.Background(), "ctx-key-5", "ctx-5"),
			actor: &entity.Actor{Code: "66666666-6666-6666-6666-666666666666", Sub: "sub-3"},
			seedData: func(t *testing.T, db *gorm.DB) {
				require.NoError(t, db.Exec("INSERT INTO users (code, sub) VALUES (?, ?)", "55555555-5555-5555-5555-555555555555", "sub-3").Error)
			},
			expectedUser:  model.User{Code: "66666666-6666-6666-6666-666666666666", Sub: "sub-3"},
			expectedCount: 1,
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

			tt.seedData(t, tx)

			repo := NewActorRepository(tx)
			err := repo.Save(tt.ctx, tt.actor)
			require.NoError(t, err)

			var gotUser model.User
			err = tx.Where("sub = ?", tt.expectedUser.Sub).First(&gotUser).Error
			require.NoError(t, err)
			assert.Equal(t, tt.expectedUser, gotUser)

			var count int64
			require.NoError(t, tx.Table("users").Count(&count).Error)
			assert.Equal(t, tt.expectedCount, count)
		})
	}
}

func TestActorRepository_Save_ActorNil(t *testing.T) {
	testTable := []struct {
		name string
		ctx  context.Context
	}{
		{
			name: "Actorがnilの場合、エラーが返ること",
			ctx:  context.WithValue(context.Background(), "ctx-key-6", "ctx-6"),
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

			repo := NewActorRepository(tx)
			err := repo.Save(tt.ctx, nil)
			require.Error(t, err)
		})
	}
}
