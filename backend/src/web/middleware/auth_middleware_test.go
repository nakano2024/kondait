package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"kondait-backend/application/usecase"
	"kondait-backend/web/dto"
)

type MockGetPrincipalUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockGetPrincipalUsecaseMockRecorder
}

type MockGetPrincipalUsecaseMockRecorder struct {
	mock *MockGetPrincipalUsecase
}

func NewMockGetPrincipalUsecase(ctrl *gomock.Controller) *MockGetPrincipalUsecase {
	mock := &MockGetPrincipalUsecase{ctrl: ctrl}
	mock.recorder = &MockGetPrincipalUsecaseMockRecorder{mock}
	return mock
}

func (m *MockGetPrincipalUsecase) EXPECT() *MockGetPrincipalUsecaseMockRecorder {
	return m.recorder
}

func (m *MockGetPrincipalUsecase) Exec(input usecase.GetPrincipalInput) (usecase.PrincipalOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exec", input)
	ret0 := ret[0].(usecase.PrincipalOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockGetPrincipalUsecaseMockRecorder) Exec(input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockGetPrincipalUsecase)(nil).Exec), input)
}

func TestAuthMiddleware_Success(t *testing.T) {
	testTable := []struct {
		name       string
		authHeader string
		expected   dto.Principal
		setupMock  func(t *testing.T, ctrl *gomock.Controller) usecase.IGetPrincipalUsecase
	}{
		{
			name:       "有効なトークンの場合、次へ進むこと",
			authHeader: "Bearer token-1",
			expected: dto.Principal{
				ActorCode: "actor-1",
				Scopes:    []string{"scope-a"},
			},
			setupMock: func(t *testing.T, ctrl *gomock.Controller) usecase.IGetPrincipalUsecase {
				usecaseMock := NewMockGetPrincipalUsecase(ctrl)
				usecaseMock.EXPECT().Exec(usecase.GetPrincipalInput{AuthToken: "token-1"}).Return(usecase.PrincipalOutput{
					ActorCode: "actor-1",
					Scopes:    []string{"scope-a"},
				}, error(nil))
				return usecaseMock
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			usecaseMock := tt.setupMock(t, ctrl)
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/private", nil)
			req.Header.Set(echo.HeaderAuthorization, tt.authHeader)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			middleware := AuthMiddleware(usecaseMock)
			next := func(c echo.Context) error {
				c.Set("nextCalled", true)
				return c.NoContent(http.StatusNoContent)
			}

			err := middleware(next)(c)
			require.NoError(t, err)
			assert.Equal(t, http.StatusNoContent, rec.Code)
			assert.Equal(t, tt.expected, c.Get("principal"))
			assert.Equal(t, true, c.Get("nextCalled"))
		})
	}
}

func TestAuthMiddleware_Failure(t *testing.T) {
	testTable := []struct {
		name           string
		authHeader     string
		expectedStatus int
		setupMock      func(t *testing.T) (usecase.IGetPrincipalUsecase, func())
	}{
		{
			name:           "依存未設定の場合、500を返すこと",
			authHeader:     "Bearer token-1",
			expectedStatus: http.StatusInternalServerError,
			setupMock: func(t *testing.T) (usecase.IGetPrincipalUsecase, func()) {
				return nil, func() {}
			},
		},
		{
			name:           "Authorizationヘッダがない場合、401を返すこと",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			setupMock: func(t *testing.T) (usecase.IGetPrincipalUsecase, func()) {
				ctrl := gomock.NewController(t)
				usecaseMock := NewMockGetPrincipalUsecase(ctrl)
				usecaseMock.EXPECT().Exec(gomock.Any()).Times(0)
				return usecaseMock, ctrl.Finish
			},
		},
		{
			name:           "Authorizationヘッダの形式が不正な場合、401を返すこと",
			authHeader:     "Token token-2",
			expectedStatus: http.StatusUnauthorized,
			setupMock: func(t *testing.T) (usecase.IGetPrincipalUsecase, func()) {
				ctrl := gomock.NewController(t)
				usecaseMock := NewMockGetPrincipalUsecase(ctrl)
				usecaseMock.EXPECT().Exec(gomock.Any()).Times(0)
				return usecaseMock, ctrl.Finish
			},
		},
		{
			name:           "トークンが空の場合、401を返すこと",
			authHeader:     "Bearer ",
			expectedStatus: http.StatusUnauthorized,
			setupMock: func(t *testing.T) (usecase.IGetPrincipalUsecase, func()) {
				ctrl := gomock.NewController(t)
				usecaseMock := NewMockGetPrincipalUsecase(ctrl)
				usecaseMock.EXPECT().Exec(gomock.Any()).Times(0)
				return usecaseMock, ctrl.Finish
			},
		},
		{
			name:           "ユースケースが失敗した場合、401を返すこと",
			authHeader:     "Bearer token-3",
			expectedStatus: http.StatusUnauthorized,
			setupMock: func(t *testing.T) (usecase.IGetPrincipalUsecase, func()) {
				ctrl := gomock.NewController(t)
				usecaseMock := NewMockGetPrincipalUsecase(ctrl)
				usecaseMock.EXPECT().Exec(usecase.GetPrincipalInput{AuthToken: "token-3"}).
					Return(usecase.PrincipalOutput{}, errors.New("unauthorized"))
				return usecaseMock, ctrl.Finish
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			usecaseMock, finish := tt.setupMock(t)
			defer finish()
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/private", nil)
			req.Header.Set(echo.HeaderAuthorization, tt.authHeader)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			middleware := AuthMiddleware(usecaseMock)
			next := func(c echo.Context) error {
				return c.NoContent(http.StatusNoContent)
			}

			err := middleware(next)(c)
			require.Error(t, err)
			var httpErr *echo.HTTPError
			require.ErrorAs(t, err, &httpErr)
			assert.Equal(t, tt.expectedStatus, httpErr.Code)
		})
	}
}
