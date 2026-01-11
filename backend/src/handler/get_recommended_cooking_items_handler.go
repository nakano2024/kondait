package handler

import (
	"net/http"
	"slices"
	"time"

	"github.com/labstack/echo/v4"

	"kondait-backend/application/usecase"
	"kondait-backend/dto/auth"
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
	principal, ok := c.Get("principal").(auth.Principal)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	if !slices.Contains(principal.Scopes, auth.ScopeCookingItemsRead) {
		return echo.NewHTTPError(http.StatusForbidden)
	}

	output, err := handler.getRecCookingItmUsecase.Exec(usecase.ReccomendedCookingListFetchCondition{
		UserCode: principal.UserCode,
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
