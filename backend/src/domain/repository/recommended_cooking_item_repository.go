package repository

import "kondait-backend/domain/aggregation"

type IRecommendedCookingItemRepository interface {
	FetchByUserCode(uCode string) (*aggregation.RecommendedCookingItemList, error)
}
