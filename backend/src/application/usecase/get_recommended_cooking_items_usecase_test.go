package usecase

import (
	"errors"
	"testing"
	"time"

	"kondait-backend/domain/aggregation"
	"kondait-backend/domain/entity"
	"kondait-backend/domain/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type recommendedCookingItemRepositoryStub struct {
	t            *testing.T
	items        []*entity.RecommendedCookingItem
	err          error
	expectedCode string
}

func (stub *recommendedCookingItemRepositoryStub) FetchByUserCode(uCode string) (*aggregation.RecommendedCookingItemList, error) {
	assert.Equal(stub.t, stub.expectedCode, uCode)
	if stub.err != nil {
		return nil, stub.err
	}

	return aggregation.NewRecommendedCookingItemList(stub.items)
}

func TestGetRecommendedCookingItemsUsecase_Exec_Success(t *testing.T) {
	testTable := []struct {
		name         string
		items        []*entity.RecommendedCookingItem
		fetchCond    ReccomendedCookingListFetchCondition
		expected     ReccomendedCookingListItemOutput
		expectedCode string
	}{
		{
			name:  "empty items",
			items: []*entity.RecommendedCookingItem{},
			fetchCond: ReccomendedCookingListFetchCondition{
				UserCode: "user-1",
			},
			expected:     ReccomendedCookingListItemOutput{List: []ReccomendedCookingOutputItem{}},
			expectedCode: "user-1",
		},
		{
			name: "five items",
			items: []*entity.RecommendedCookingItem{
				{Code: "A1", Name: "Rice", CookCount: 1, LastCookedDate: time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)},
				{Code: "B2", Name: "Curry", CookCount: 2, LastCookedDate: time.Date(2024, 1, 3, 3, 4, 5, 0, time.UTC)},
				{Code: "C3", Name: "Soup", CookCount: 3, LastCookedDate: time.Date(2024, 1, 4, 3, 4, 5, 0, time.UTC)},
				{Code: "D4", Name: "Pasta", CookCount: 4, LastCookedDate: time.Date(2024, 1, 5, 3, 4, 5, 0, time.UTC)},
				{Code: "E5", Name: "Salad", CookCount: 5, LastCookedDate: time.Date(2024, 1, 6, 3, 4, 5, 0, time.UTC)},
			},
			fetchCond: ReccomendedCookingListFetchCondition{
				UserCode: "user-2",
			},
			expected: ReccomendedCookingListItemOutput{List: []ReccomendedCookingOutputItem{
				{Code: "A1", Name: "Rice", CookCount: 1, LastCookedDate: time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)},
				{Code: "B2", Name: "Curry", CookCount: 2, LastCookedDate: time.Date(2024, 1, 3, 3, 4, 5, 0, time.UTC)},
				{Code: "C3", Name: "Soup", CookCount: 3, LastCookedDate: time.Date(2024, 1, 4, 3, 4, 5, 0, time.UTC)},
				{Code: "D4", Name: "Pasta", CookCount: 4, LastCookedDate: time.Date(2024, 1, 5, 3, 4, 5, 0, time.UTC)},
				{Code: "E5", Name: "Salad", CookCount: 5, LastCookedDate: time.Date(2024, 1, 6, 3, 4, 5, 0, time.UTC)},
			}},
			expectedCode: "user-2",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			repo := &recommendedCookingItemRepositoryStub{
				t:            t,
				items:        tt.items,
				expectedCode: tt.expectedCode,
			}
			usecase := NewGetRecommendedCookingItemsUsecase(repo)

			got, err := usecase.Exec(tt.fetchCond)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestGetRecommendedCookingItemsUsecase_Exec_Failure(t *testing.T) {
	testTable := []struct {
		name      string
		repo      repository.IRecommendedCookingItemRepository
		fetchCond ReccomendedCookingListFetchCondition
	}{
		{
			name: "repository error",
			repo: &recommendedCookingItemRepositoryStub{
				t:            t,
				err:          errors.New("repository error"),
				expectedCode: "user-1",
			},
			fetchCond: ReccomendedCookingListFetchCondition{
				UserCode: "user-1",
			},
		},
		{
			name: "too many items",
			repo: &recommendedCookingItemRepositoryStub{
				t: t,
				items: []*entity.RecommendedCookingItem{
					{Code: "A1", Name: "Rice", CookCount: 1, LastCookedDate: time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)},
					{Code: "B2", Name: "Curry", CookCount: 2, LastCookedDate: time.Date(2024, 1, 3, 3, 4, 5, 0, time.UTC)},
					{Code: "C3", Name: "Soup", CookCount: 3, LastCookedDate: time.Date(2024, 1, 4, 3, 4, 5, 0, time.UTC)},
					{Code: "D4", Name: "Pasta", CookCount: 4, LastCookedDate: time.Date(2024, 1, 5, 3, 4, 5, 0, time.UTC)},
					{Code: "E5", Name: "Salad", CookCount: 5, LastCookedDate: time.Date(2024, 1, 6, 3, 4, 5, 0, time.UTC)},
					{Code: "F6", Name: "Stew", CookCount: 6, LastCookedDate: time.Date(2024, 1, 7, 3, 4, 5, 0, time.UTC)},
				},
				expectedCode: "user-2",
			},
			fetchCond: ReccomendedCookingListFetchCondition{
				UserCode: "user-2",
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			usecase := NewGetRecommendedCookingItemsUsecase(tt.repo)

			got, err := usecase.Exec(tt.fetchCond)
			require.Error(t, err)
			assert.Equal(t, ReccomendedCookingListItemOutput{}, got)
			assert.Nil(t, got.List)
		})
	}
}
