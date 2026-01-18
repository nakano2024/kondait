package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"kondait-backend/domain/aggregation"
	"kondait-backend/domain/entity"
	"kondait-backend/domain/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type MockRecommendedCookingItemRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRecommendedCookingItemRepositoryMockRecorder
}

type MockRecommendedCookingItemRepositoryMockRecorder struct {
	mock *MockRecommendedCookingItemRepository
}

func NewMockRecommendedCookingItemRepository(ctrl *gomock.Controller) *MockRecommendedCookingItemRepository {
	mock := &MockRecommendedCookingItemRepository{ctrl: ctrl}
	mock.recorder = &MockRecommendedCookingItemRepositoryMockRecorder{mock}
	return mock
}

func (m *MockRecommendedCookingItemRepository) EXPECT() *MockRecommendedCookingItemRepositoryMockRecorder {
	return m.recorder
}

func (m *MockRecommendedCookingItemRepository) FetchByUserCode(ctx context.Context, uCode string) (*aggregation.RecommendedCookingItemList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchByUserCode", ctx, uCode)
	ret0 := ret[0].(*aggregation.RecommendedCookingItemList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockRecommendedCookingItemRepositoryMockRecorder) FetchByUserCode(ctx interface{}, uCode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchByUserCode", reflect.TypeOf((*MockRecommendedCookingItemRepository)(nil).FetchByUserCode), ctx, uCode)
}

func TestGetRecommendedCookingItemsUsecase_Exec_Success(t *testing.T) {
	testTable := []struct {
		name      string
		ctx       context.Context
		fetchCond ReccomendedCookingListFetchCondition
		expected  ReccomendedCookingListItemOutput
		setupMock func(t *testing.T, ctrl *gomock.Controller, ctx context.Context) repository.IRecommendedCookingItemRepository
	}{
		{
			name:      "要素が空の場合、空の一覧を返すこと",
			ctx:       context.WithValue(context.Background(), "ctx-key-1", "ctx-1"),
			fetchCond: ReccomendedCookingListFetchCondition{UserCode: "user-1"},
			expected:  ReccomendedCookingListItemOutput{List: []ReccomendedCookingOutputItem{}},
			setupMock: func(t *testing.T, ctrl *gomock.Controller, ctx context.Context) repository.IRecommendedCookingItemRepository {
				repo := NewMockRecommendedCookingItemRepository(ctrl)
				list, err := aggregation.NewRecommendedCookingItemList([]*entity.RecommendedCookingItem{})
				require.NoError(t, err)
				repo.EXPECT().FetchByUserCode(ctx, "user-1").Return(list, error(nil))
				return repo
			},
		},
		{
			name:      "5件の場合、5件の一覧を返すこと",
			ctx:       context.WithValue(context.Background(), "ctx-key-2", "ctx-2"),
			fetchCond: ReccomendedCookingListFetchCondition{UserCode: "user-2"},
			expected: ReccomendedCookingListItemOutput{List: []ReccomendedCookingOutputItem{
				{Code: "A1", Name: "Rice", CookCount: 1, LastCookedDate: time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)},
				{Code: "B2", Name: "Curry", CookCount: 2, LastCookedDate: time.Date(2024, 1, 3, 3, 4, 5, 0, time.UTC)},
				{Code: "C3", Name: "Soup", CookCount: 3, LastCookedDate: time.Date(2024, 1, 4, 3, 4, 5, 0, time.UTC)},
				{Code: "D4", Name: "Pasta", CookCount: 4, LastCookedDate: time.Date(2024, 1, 5, 3, 4, 5, 0, time.UTC)},
				{Code: "E5", Name: "Salad", CookCount: 5, LastCookedDate: time.Date(2024, 1, 6, 3, 4, 5, 0, time.UTC)},
			}},
			setupMock: func(t *testing.T, ctrl *gomock.Controller, ctx context.Context) repository.IRecommendedCookingItemRepository {
				repo := NewMockRecommendedCookingItemRepository(ctrl)
				list, err := aggregation.NewRecommendedCookingItemList([]*entity.RecommendedCookingItem{
					{Code: "A1", Name: "Rice", CookCount: 1, LastCookedDate: time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)},
					{Code: "B2", Name: "Curry", CookCount: 2, LastCookedDate: time.Date(2024, 1, 3, 3, 4, 5, 0, time.UTC)},
					{Code: "C3", Name: "Soup", CookCount: 3, LastCookedDate: time.Date(2024, 1, 4, 3, 4, 5, 0, time.UTC)},
					{Code: "D4", Name: "Pasta", CookCount: 4, LastCookedDate: time.Date(2024, 1, 5, 3, 4, 5, 0, time.UTC)},
					{Code: "E5", Name: "Salad", CookCount: 5, LastCookedDate: time.Date(2024, 1, 6, 3, 4, 5, 0, time.UTC)},
				})
				require.NoError(t, err)
				repo.EXPECT().FetchByUserCode(ctx, "user-2").Return(list, error(nil))
				return repo
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := tt.setupMock(t, ctrl, tt.ctx)
			usecase := NewGetRecommendedCookingItemsUsecase(repo)

			got, err := usecase.Exec(tt.ctx, tt.fetchCond)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestGetRecommendedCookingItemsUsecase_Exec_Failure(t *testing.T) {
	testTable := []struct {
		name      string
		ctx       context.Context
		fetchCond ReccomendedCookingListFetchCondition
		setupMock func(t *testing.T, ctrl *gomock.Controller, ctx context.Context) repository.IRecommendedCookingItemRepository
	}{
		{
			name:      "リポジトリエラーの場合、エラーが返ること",
			ctx:       context.WithValue(context.Background(), "ctx-key-3", "ctx-3"),
			fetchCond: ReccomendedCookingListFetchCondition{UserCode: "user-1"},
			setupMock: func(t *testing.T, ctrl *gomock.Controller, ctx context.Context) repository.IRecommendedCookingItemRepository {
				repo := NewMockRecommendedCookingItemRepository(ctrl)
				repo.EXPECT().FetchByUserCode(ctx, "user-1").
					Return((*aggregation.RecommendedCookingItemList)(nil), errors.New("repository error"))
				return repo
			},
		},
		{
			name:      "6件の場合、エラーが返ること",
			ctx:       context.WithValue(context.Background(), "ctx-key-4", "ctx-4"),
			fetchCond: ReccomendedCookingListFetchCondition{UserCode: "user-2"},
			setupMock: func(t *testing.T, ctrl *gomock.Controller, ctx context.Context) repository.IRecommendedCookingItemRepository {
				repo := NewMockRecommendedCookingItemRepository(ctrl)
				list, err := aggregation.NewRecommendedCookingItemList([]*entity.RecommendedCookingItem{
					{Code: "A1", Name: "Rice", CookCount: 1, LastCookedDate: time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)},
					{Code: "B2", Name: "Curry", CookCount: 2, LastCookedDate: time.Date(2024, 1, 3, 3, 4, 5, 0, time.UTC)},
					{Code: "C3", Name: "Soup", CookCount: 3, LastCookedDate: time.Date(2024, 1, 4, 3, 4, 5, 0, time.UTC)},
					{Code: "D4", Name: "Pasta", CookCount: 4, LastCookedDate: time.Date(2024, 1, 5, 3, 4, 5, 0, time.UTC)},
					{Code: "E5", Name: "Salad", CookCount: 5, LastCookedDate: time.Date(2024, 1, 6, 3, 4, 5, 0, time.UTC)},
					{Code: "F6", Name: "Stew", CookCount: 6, LastCookedDate: time.Date(2024, 1, 7, 3, 4, 5, 0, time.UTC)},
				})
				require.Error(t, err)
				repo.EXPECT().FetchByUserCode(ctx, "user-2").Return(list, err)
				return repo
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := tt.setupMock(t, ctrl, tt.ctx)
			usecase := NewGetRecommendedCookingItemsUsecase(repo)

			got, err := usecase.Exec(tt.ctx, tt.fetchCond)
			require.Error(t, err)
			assert.Equal(t, ReccomendedCookingListItemOutput{}, got)
			assert.Nil(t, got.List)
		})
	}
}
