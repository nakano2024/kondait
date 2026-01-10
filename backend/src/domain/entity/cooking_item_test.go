package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewExistingCookingItem_Normal(t *testing.T) {
	testTable := []struct {
		name  string
		input struct {
			code           string
			itemName       string
			cookCount      uint
			lastCookedDate time.Time
		}
		expected *RecommendedCookingItem
	}{
		{
			name: "set fields",
			input: struct {
				code           string
				itemName       string
				cookCount      uint
				lastCookedDate time.Time
			}{
				code:           "A1",
				itemName:       "Rice",
				cookCount:      2,
				lastCookedDate: time.Date(2024, 2, 3, 4, 5, 6, 0, time.UTC),
			},
			expected: &RecommendedCookingItem{
				Code:           "A1",
				Name:           "Rice",
				CookCount:      2,
				LastCookedDate: time.Date(2024, 2, 3, 4, 5, 6, 0, time.UTC),
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			got := NewReccomendedCookingItem(
				tt.input.code,
				tt.input.itemName,
				tt.input.cookCount,
				tt.input.lastCookedDate,
			)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestRecommendedCookingItem_IsCooked_Normal(t *testing.T) {
	testTable := []struct {
		name string
		item *RecommendedCookingItem
	}{
		{
			name: "non-zero time",
			item: &RecommendedCookingItem{
				Code:           "A1",
				Name:           "Rice",
				CookCount:      1,
				LastCookedDate: time.Date(2024, 3, 4, 5, 6, 7, 0, time.UTC),
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.item.IsCooked()
			assert.True(t, got)
		})
	}
}

func TestRecommendedCookingItem_IsCooked_Abnormal(t *testing.T) {
	testTable := []struct {
		name string
		item *RecommendedCookingItem
	}{
		{
			name: "zero time",
			item: &RecommendedCookingItem{
				Code:           "A1",
				Name:           "Rice",
				CookCount:      0,
				LastCookedDate: time.Time{},
			},
		},
		{
			name: "zero time",
			item: &RecommendedCookingItem{
				Code:           "A1",
				Name:           "Rice",
				CookCount:      0,
				LastCookedDate: time.Time{},
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.item.IsCooked()
			assert.False(t, got)
		})
	}
}
