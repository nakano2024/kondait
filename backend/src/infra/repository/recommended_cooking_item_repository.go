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

func (repo *recommendedCookingItemRepository) FetchByUserCode(ctx context.Context, uCode string) (*aggregation.RecommendedCookingItemList, error) {
	var rows []recommendedCookingItemRow

	err := repo.db.WithContext(ctx).
		Table("cooking_items").
		Select("code, name, cook_count, last_cooked_date").
		Where("cooking_items.owner_code = ?", uCode).
		Order("cook_count ASC").
		Order("last_cooked_date ASC NULLS FIRST").
		Limit(5).
		Find(&rows).Error
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
