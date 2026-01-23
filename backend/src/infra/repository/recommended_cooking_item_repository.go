package repository

import (
	"time"

	"gorm.io/gorm"

	"context"
	"kondait-backend/domain/aggregation"
	"kondait-backend/domain/entity"
	domainrepo "kondait-backend/domain/repository"
)

type recommendedCookingItemRepository struct {
	db *gorm.DB
}

type recommendedCookingItemRow struct {
	Code           string    `gorm:"column:code"`
	Name           string    `gorm:"column:name"`
	CookCount      int64     `gorm:"column:cook_count"`
	LastCookedDate time.Time `gorm:"column:last_cooked_date"`
}

func NewRecommendedCookingItemRepository(db *gorm.DB) domainrepo.IRecommendedCookingItemRepository {
	return &recommendedCookingItemRepository{
		db: db,
	}
}

func (repo *recommendedCookingItemRepository) FetchByUserCode(ctx context.Context, usrCode string) (*aggregation.RecommendedCookingItemList, error) {
	var rows []recommendedCookingItemRow

	baseQuery := repo.db.WithContext(ctx).
		Table("cooking_items").
		Select("cooking_items.code AS code, cooking_items.name AS name, COALESCE(COUNT(cooking_histories.code), 0) AS cook_count, MAX(cooking_histories.cooked_at) AS last_cooked_date").
		Joins("LEFT JOIN cooking_histories ON cooking_histories.cooking_item_code = cooking_items.code").
		Where("cooking_items.owner_code = ?", usrCode).
		Group("cooking_items.code, cooking_items.name").
		Order("cook_count ASC").
		Order("last_cooked_date ASC NULLS FIRST").
		Limit(5)

	err := baseQuery.Find(&rows).Error
	if err != nil {
		return nil, err
	}

	items := make([]*entity.RecommendedCookingItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, &entity.RecommendedCookingItem{
			Code:           row.Code,
			Name:           row.Name,
			CookCount:      uint(row.CookCount),
			LastCookedDate: row.LastCookedDate,
		})
	}

	return aggregation.NewRecommendedCookingItemList(items)
}
