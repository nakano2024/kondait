package aggregation

import (
	"testing"
	"time"

	"kondait-backend/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRecommendedCookingItemList_Success(t *testing.T) {
	testTable := []struct {
		name  string
		items []*entity.RecommendedCookingItem
	}{
		{
			name: "要素が空の場合、件数0で返ること",
			items: []*entity.RecommendedCookingItem{
				// empty
			},
		},
		{
			name: "5件の場合、件数5で返ること",
			items: []*entity.RecommendedCookingItem{
				{Code: "A1", Name: "Rice", CookCount: 1, LastCookedDate: time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)},
				{Code: "B2", Name: "Curry", CookCount: 2, LastCookedDate: time.Date(2024, 1, 3, 3, 4, 5, 0, time.UTC)},
				{Code: "C3", Name: "Soup", CookCount: 3, LastCookedDate: time.Date(2024, 1, 4, 3, 4, 5, 0, time.UTC)},
				{Code: "D4", Name: "Pasta", CookCount: 4, LastCookedDate: time.Date(2024, 1, 5, 3, 4, 5, 0, time.UTC)},
				{Code: "E5", Name: "Salad", CookCount: 5, LastCookedDate: time.Date(2024, 1, 6, 3, 4, 5, 0, time.UTC)},
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRecommendedCookingItemList(tt.items)
			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Len(t, got.Items, len(tt.items))
		})
	}
}

func TestNewRecommendedCookingItemList_Failure(t *testing.T) {
	testTable := []struct {
		name  string
		items []*entity.RecommendedCookingItem
	}{
		{
			name: "6件の場合、エラーが返ること",
			items: []*entity.RecommendedCookingItem{
				{Code: "A1", Name: "Rice", CookCount: 1, LastCookedDate: time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)},
				{Code: "B2", Name: "Curry", CookCount: 2, LastCookedDate: time.Date(2024, 1, 3, 3, 4, 5, 0, time.UTC)},
				{Code: "C3", Name: "Soup", CookCount: 3, LastCookedDate: time.Date(2024, 1, 4, 3, 4, 5, 0, time.UTC)},
				{Code: "D4", Name: "Pasta", CookCount: 4, LastCookedDate: time.Date(2024, 1, 5, 3, 4, 5, 0, time.UTC)},
				{Code: "E5", Name: "Salad", CookCount: 5, LastCookedDate: time.Date(2024, 1, 6, 3, 4, 5, 0, time.UTC)},
				{Code: "F6", Name: "Stew", CookCount: 6, LastCookedDate: time.Date(2024, 1, 7, 3, 4, 5, 0, time.UTC)},
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRecommendedCookingItemList(tt.items)
			require.Error(t, err)
			assert.Nil(t, got)
		})
	}
}
