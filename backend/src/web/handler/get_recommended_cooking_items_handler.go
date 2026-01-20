package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"kondait-backend/application/usecase"
	"kondait-backend/web/dto"
	"kondait-backend/web/util"
)

type getRecommendedCookingItemsHandler struct {
	getRecCookingItmUsecase usecase.IGetRecommendedCookingItemsUsecase
}

func NewGetRecommendedCookingItemsHandler(getRecCookingItmUsecase usecase.IGetRecommendedCookingItemsUsecase) *getRecommendedCookingItemsHandler {
	return &getRecommendedCookingItemsHandler{
		getRecCookingItmUsecase: getRecCookingItmUsecase,
	}
}

type responseItem struct {
	Code           string    `json:"code"`
	Name           string    `json:"name"`
	CookCount      uint      `json:"cook_count"`
	LastCookedDate time.Time `json:"last_cooked_date,omitempty"`
}

type response struct {
	RecommendedCookingItems []responseItem `json:"cooking_items"`
}

func (handler *getRecommendedCookingItemsHandler) Handle(c echo.Context) error {
	principal, ok := c.Get(dto.PrincipalContextKeyName).(dto.Principal)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "principal is not set")
	}

	allowdScopes := []string{
		dto.ScopeCookingItem,
		dto.ScopeCookingItemRead,
		dto.ScopeCookingItemWrite,
		dto.ScopeCookingItemDelete,
	}
	if !util.HasAnyScope(allowdScopes, principal.Scopes) {
		return echo.NewHTTPError(http.StatusForbidden)
	}

	output, err := handler.getRecCookingItmUsecase.Exec(c.Request().Context(), usecase.ReccomendedCookingListFetchCondition{
		UserCode: principal.ActorCode,
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
