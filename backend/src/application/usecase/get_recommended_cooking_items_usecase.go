package usecase

import (
	"context"
	"time"

	"kondait-backend/domain/entity"
	"kondait-backend/domain/repository"
)

type ReccomendedCookingOutputItem struct {
	Code           string
	Name           string
	CookCount      uint
	LastCookedDate time.Time
}

type ReccomendedCookingListItemOutput struct {
	List []ReccomendedCookingOutputItem
}

type ReccomendedCookingListFetchCondition struct {
	UserCode string
}

type IGetRecommendedCookingItemsUsecase interface {
	Exec(ctx context.Context, fCond ReccomendedCookingListFetchCondition) (ReccomendedCookingListItemOutput, error)
}

type getRecommendedCookingItemsUsecase struct {
	repo repository.IRecommendedCookingItemRepository
}

func NewGetRecommendedCookingItemsUsecase(
	repository repository.IRecommendedCookingItemRepository,
) IGetRecommendedCookingItemsUsecase {
	return &getRecommendedCookingItemsUsecase{
		repo: repository,
	}
}

func (usecase *getRecommendedCookingItemsUsecase) Exec(ctx context.Context, fCond ReccomendedCookingListFetchCondition) (ReccomendedCookingListItemOutput, error) {
	list, err := usecase.repo.FetchByUserCode(ctx, fCond.UserCode)
	if err != nil {
		return ReccomendedCookingListItemOutput{}, err
	}

	return convertRecommendedCookingListItems(list.Items), nil
}

func convertRecommendedCookingListItems(items []*entity.RecommendedCookingItem) ReccomendedCookingListItemOutput {
	outputItems := make([]ReccomendedCookingOutputItem, 0, len(items))
	for _, item := range items {
		outputItems = append(outputItems, ReccomendedCookingOutputItem{
			Code:           item.Code,
			Name:           item.Name,
			CookCount:      item.CookCount,
			LastCookedDate: item.LastCookedDate,
		})
	}

	return ReccomendedCookingListItemOutput{
		List: outputItems,
	}
}
