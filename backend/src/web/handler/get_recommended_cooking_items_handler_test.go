package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"kondait-backend/application/usecase"
	"kondait-backend/web/dto"
)

type MockGetRecommendedCookingItemsUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockGetRecommendedCookingItemsUsecaseMockRecorder
}

type MockGetRecommendedCookingItemsUsecaseMockRecorder struct {
	mock *MockGetRecommendedCookingItemsUsecase
}

func NewMockGetRecommendedCookingItemsUsecase(ctrl *gomock.Controller) *MockGetRecommendedCookingItemsUsecase {
	mock := &MockGetRecommendedCookingItemsUsecase{ctrl: ctrl}
	mock.recorder = &MockGetRecommendedCookingItemsUsecaseMockRecorder{mock}
	return mock
}

func (m *MockGetRecommendedCookingItemsUsecase) EXPECT() *MockGetRecommendedCookingItemsUsecaseMockRecorder {
	return m.recorder
}

func (m *MockGetRecommendedCookingItemsUsecase) Exec(ctx context.Context, input usecase.ReccomendedCookingListFetchCondition) (usecase.ReccomendedCookingListItemOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exec", ctx, input)
	ret0 := ret[0].(usecase.ReccomendedCookingListItemOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockGetRecommendedCookingItemsUsecaseMockRecorder) Exec(ctx interface{}, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockGetRecommendedCookingItemsUsecase)(nil).Exec), ctx, input)
}

func TestGetRecommendedCookingItemsHandler_Handle_Success(t *testing.T) {
	testTable := []struct {
		name           string
		ctx            context.Context
		principal      dto.Principal
		expectedStatus int
		expectedBody   string
		setupMock      func(t *testing.T, ctrl *gomock.Controller, ctx context.Context) usecase.IGetRecommendedCookingItemsUsecase
	}{
		{
			name: "一覧が空の場合、空のレスポンスを返すこと",
			ctx:  context.WithValue(context.Background(), "ctx-key-1", "ctx-1"),
			principal: dto.Principal{
				ActorCode: "actor-1",
				Scopes:    []string{dto.ScopeCookingItemsRead},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"recommended_cooking_items":[]}`,
			setupMock: func(t *testing.T, ctrl *gomock.Controller, ctx context.Context) usecase.IGetRecommendedCookingItemsUsecase {
				usecaseMock := NewMockGetRecommendedCookingItemsUsecase(ctrl)
				usecaseMock.EXPECT().
					Exec(ctx, usecase.ReccomendedCookingListFetchCondition{UserCode: "actor-1"}).
					Return(usecase.ReccomendedCookingListItemOutput{List: []usecase.ReccomendedCookingOutputItem{}}, error(nil))
				return usecaseMock
			},
		},
		{
			name: "一覧に要素がある場合、返却されること",
			ctx:  context.WithValue(context.Background(), "ctx-key-2", "ctx-2"),
			principal: dto.Principal{
				ActorCode: "actor-2",
				Scopes:    []string{dto.ScopeCookingItemsRead},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"recommended_cooking_items":[{"code":"A1","name":"Rice","cookCount":1,"last_cooked_date":"2024-01-02T03:04:05Z"}]}`,
			setupMock: func(t *testing.T, ctrl *gomock.Controller, ctx context.Context) usecase.IGetRecommendedCookingItemsUsecase {
				usecaseMock := NewMockGetRecommendedCookingItemsUsecase(ctrl)
				usecaseMock.EXPECT().
					Exec(ctx, usecase.ReccomendedCookingListFetchCondition{UserCode: "actor-2"}).
					Return(usecase.ReccomendedCookingListItemOutput{
						List: []usecase.ReccomendedCookingOutputItem{
							{
								Code:           "A1",
								Name:           "Rice",
								CookCount:      1,
								LastCookedDate: time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC),
							},
						},
					}, error(nil))
				return usecaseMock
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			usecaseMock := tt.setupMock(t, ctrl, tt.ctx)
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/private/cooking-items/recommends", nil).WithContext(tt.ctx)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("principal", tt.principal)
			handler := NewGetRecommendedCookingItemsHandler(usecaseMock)

			err := handler.Handle(c)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
		})
	}
}

func TestGetRecommendedCookingItemsHandler_Handle_PrincipalMissing(t *testing.T) {
	testTable := []struct {
		name           string
		ctx            context.Context
		expectedStatus int
	}{
		{
			name:           "Principalがない場合、500を返すこと",
			ctx:            context.WithValue(context.Background(), "ctx-key-3", "ctx-3"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			usecaseMock := NewMockGetRecommendedCookingItemsUsecase(ctrl)
			usecaseMock.EXPECT().Exec(gomock.Any(), gomock.Any()).Times(0)
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/private/cooking-items/recommends", nil).WithContext(tt.ctx)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			handler := NewGetRecommendedCookingItemsHandler(usecaseMock)

			err := handler.Handle(c)
			require.Error(t, err)
			var httpErr *echo.HTTPError
			require.ErrorAs(t, err, &httpErr)
			assert.Equal(t, tt.expectedStatus, httpErr.Code)
		})
	}
}

func TestGetRecommendedCookingItemsHandler_Handle_MissingScope(t *testing.T) {
	testTable := []struct {
		name           string
		ctx            context.Context
		principal      dto.Principal
		expectedStatus int
	}{
		{
			name: "スコープ不足の場合、403を返すこと",
			ctx:  context.WithValue(context.Background(), "ctx-key-4", "ctx-4"),
			principal: dto.Principal{
				ActorCode: "actor-3",
				Scopes:    []string{"other-scope"},
			},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			usecaseMock := NewMockGetRecommendedCookingItemsUsecase(ctrl)
			usecaseMock.EXPECT().Exec(gomock.Any(), gomock.Any()).Times(0)
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/private/cooking-items/recommends", nil).WithContext(tt.ctx)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("principal", tt.principal)
			handler := NewGetRecommendedCookingItemsHandler(usecaseMock)

			err := handler.Handle(c)
			require.Error(t, err)
			var httpErr *echo.HTTPError
			require.ErrorAs(t, err, &httpErr)
			assert.Equal(t, tt.expectedStatus, httpErr.Code)
		})
	}
}

func TestGetRecommendedCookingItemsHandler_Handle_UsecaseError(t *testing.T) {
	testTable := []struct {
		name           string
		ctx            context.Context
		expectedStatus int
		expectedBody   string
		setupMock      func(t *testing.T, ctrl *gomock.Controller, ctx context.Context) usecase.IGetRecommendedCookingItemsUsecase
	}{
		{
			name:           "ユースケースが失敗した場合、200で空返却すること",
			ctx:            context.WithValue(context.Background(), "ctx-key-5", "ctx-5"),
			expectedStatus: http.StatusOK,
			expectedBody:   `{"recommended_cooking_items":[]}`,
			setupMock: func(t *testing.T, ctrl *gomock.Controller, ctx context.Context) usecase.IGetRecommendedCookingItemsUsecase {
				usecaseMock := NewMockGetRecommendedCookingItemsUsecase(ctrl)
				usecaseMock.EXPECT().
					Exec(ctx, usecase.ReccomendedCookingListFetchCondition{UserCode: "actor-4"}).
					Return(usecase.ReccomendedCookingListItemOutput{}, errors.New("usecase error"))
				return usecaseMock
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			usecaseMock := tt.setupMock(t, ctrl, tt.ctx)
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/private/cooking-items/recommends", nil).WithContext(tt.ctx)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("principal", dto.Principal{
				ActorCode: "actor-4",
				Scopes:    []string{dto.ScopeCookingItemsRead},
			})
			handler := NewGetRecommendedCookingItemsHandler(usecaseMock)

			err := handler.Handle(c)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
		})
	}
}
