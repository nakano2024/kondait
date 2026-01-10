package entity

import "time"

type RecommendedCookingItem struct {
	Code           string
	Name           string
	CookCount      uint
	LastCookedDate time.Time
}

func NewReccomendedCookingItem(code string, name string, cookCnt uint, lastCookedTime time.Time) *RecommendedCookingItem {
	return &RecommendedCookingItem{
		Code:           code,
		Name:           name,
		CookCount:      cookCnt,
		LastCookedDate: lastCookedTime,
	}
}

func (cookingItm *RecommendedCookingItem) IsCooked() bool {
	return !cookingItm.LastCookedDate.IsZero()
}
