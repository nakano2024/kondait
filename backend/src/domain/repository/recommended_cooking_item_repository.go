package repository

import (
	"context"
	"kondait-backend/domain/aggregation"
)

type IRecommendedCookingItemRepository interface {
	FetchByUserCode(ctx context.Context, uCode string) (*aggregation.RecommendedCookingItemList, error)
}
