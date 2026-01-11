package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"kondait-backend/application/usecase"
)

type getRecommendedCookingItemsHandler struct {
	getRecCookingItmUsecase usecase.GetRecommendedCookingItemsUsecase
}

func NewGetRecommendedCookingItemsHandler(getRecCookingItmUsecase usecase.GetRecommendedCookingItemsUsecase) *getRecommendedCookingItemsHandler {
	return &getRecommendedCookingItemsHandler{
		getRecCookingItmUsecase: getRecCookingItmUsecase,
	}
}

type responseItem struct {
	Code           string    `json:"code"`
	Name           string    `json:"name"`
	CookCount      uint      `json:"cookCount"`
	LastCookedDate time.Time `json:"last_cooked_date,omitempty"`
}

type response struct {
	RecommendedCookingItems []responseItem `json:"recommended_cooking_items"`
}

func (handler *getRecommendedCookingItemsHandler) Handle(c echo.Context) error {
	authUserId, ok := c.Get("auth_user_id").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	output, err := handler.getRecCookingItmUsecase.Exec(usecase.ReccomendedCookingListFetchCondition{
		UserCode: authUserId,
	})

	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, parseOutputToJson(output))
}

func parseOutputToJson(output usecase.ReccomendedCookingListItemOutput) response {
	resItems := make([]responseItem, 0, len(output.List))
	for _, item := range output.List {
		resItems = append(resItems, responseItem{
			Code:           item.Code,
			Name:           item.Name,
			CookCount:      item.CookCount,
			LastCookedDate: item.LastCookedDate,
		})
	}

	return response{
		RecommendedCookingItems: resItems,
	}
}
