package aggregation

import (
	"errors"

	"kondait-backend/domain/entity"
)

const maxCount = 5

type RecommendedCookingItemList struct {
	Items []*entity.RecommendedCookingItem
}

func NewRecommendedCookingItemList(items []*entity.RecommendedCookingItem) (*RecommendedCookingItemList, error) {
	if maxCount < len(items) {
		return nil, errors.New("recommended cooking items must be 5 or fewer")
	}

	return &RecommendedCookingItemList{
		Items: items,
	}, nil
}
