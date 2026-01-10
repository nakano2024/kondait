package usecase

import (
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
	list []ReccomendedCookingOutputItem
}

type ReccomendedCookingListFetchCondition struct {
	userCode string
}

type GetRecommendedCookingItemsUsecase interface {
	Exec(fCond ReccomendedCookingListFetchCondition) (ReccomendedCookingListItemOutput, error)
}

type getRecommendedCookingItemsUsecase struct {
	repo repository.IRecommendedCookingItemRepository
}

func NewGetRecommendedCookingItemsUsecase(
	repository repository.IRecommendedCookingItemRepository,
) GetRecommendedCookingItemsUsecase {
	return &getRecommendedCookingItemsUsecase{
		repo: repository,
	}
}

func (usecase *getRecommendedCookingItemsUsecase) Exec(fCond ReccomendedCookingListFetchCondition) (ReccomendedCookingListItemOutput, error) {
	list, err := usecase.repo.FetchByUserCode(fCond.userCode)
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
		list: outputItems,
	}
}
