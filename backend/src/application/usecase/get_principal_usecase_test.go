package usecase

import (
	"errors"
	"reflect"
	"testing"

	"kondait-backend/application/auth"
	"kondait-backend/domain/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type MockAuthIntrospector struct {
	ctrl     *gomock.Controller
	recorder *MockAuthIntrospectorMockRecorder
}

type MockAuthIntrospectorMockRecorder struct {
	mock *MockAuthIntrospector
}

func NewMockAuthIntrospector(ctrl *gomock.Controller) *MockAuthIntrospector {
	mock := &MockAuthIntrospector{ctrl: ctrl}
	mock.recorder = &MockAuthIntrospectorMockRecorder{mock}
	return mock
}

func (m *MockAuthIntrospector) EXPECT() *MockAuthIntrospectorMockRecorder {
	return m.recorder
}

func (m *MockAuthIntrospector) Introspect(token string) (auth.AuthIntrospectionResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Introspect", token)
	ret0 := ret[0].(auth.AuthIntrospectionResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockAuthIntrospectorMockRecorder) Introspect(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Introspect", reflect.TypeOf((*MockAuthIntrospector)(nil).Introspect), token)
}

type MockActorRepository struct {
	ctrl     *gomock.Controller
	recorder *MockActorRepositoryMockRecorder
}

type MockActorRepositoryMockRecorder struct {
	mock *MockActorRepository
}

func NewMockActorRepository(ctrl *gomock.Controller) *MockActorRepository {
	mock := &MockActorRepository{ctrl: ctrl}
	mock.recorder = &MockActorRepositoryMockRecorder{mock}
	return mock
}

func (m *MockActorRepository) EXPECT() *MockActorRepositoryMockRecorder {
	return m.recorder
}

func (m *MockActorRepository) FetchBySub(sub string) (*entity.Actor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchBySub", sub)
	ret0 := ret[0].(*entity.Actor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockActorRepositoryMockRecorder) FetchBySub(sub interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchBySub", reflect.TypeOf((*MockActorRepository)(nil).FetchBySub), sub)
}

func TestGetPrincipalUsecase_Exec_Success(t *testing.T) {
	testTable := []struct {
		name      string
		input     GetPrincipalInput
		expected  PrincipalOutput
		setupMock func(t *testing.T, ctrl *gomock.Controller) (*MockAuthIntrospector, *MockActorRepository)
	}{
		{
			name:  "トークンが有効でActorが取得できる場合、Principalが取得できること",
			input: GetPrincipalInput{AuthToken: "token-1"},
			expected: PrincipalOutput{
				ActorCode: "actor-1",
				Scopes:    []string{"scope-a"},
			},
			setupMock: func(t *testing.T, ctrl *gomock.Controller) (*MockAuthIntrospector, *MockActorRepository) {
				authIntrospector := NewMockAuthIntrospector(ctrl)
				actorRepo := NewMockActorRepository(ctrl)
				authIntrospector.EXPECT().Introspect("token-1").Return(auth.AuthIntrospectionResult{
					IsActive: true,
					Sub:      "sub-1",
					Scopes:   []string{"scope-a"},
				}, error(nil))
				actorRepo.EXPECT().FetchBySub("sub-1").Return(&entity.Actor{Code: "actor-1"}, error(nil))
				return authIntrospector, actorRepo
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			authIntrospector, actorRepo := tt.setupMock(t, ctrl)
			usecase := NewGetPrincipalUsecase(authIntrospector, actorRepo)

			got, err := usecase.Exec(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestGetPrincipalUsecase_Exec_Failure(t *testing.T) {
	testTable := []struct {
		name      string
		input     GetPrincipalInput
		setupMock func(t *testing.T) (IGetPrincipalUsecase, func())
	}{
		{
			name:  "依存未設定の場合、エラーが返ること",
			input: GetPrincipalInput{AuthToken: "token-1"},
			setupMock: func(t *testing.T) (IGetPrincipalUsecase, func()) {
				return NewGetPrincipalUsecase(nil, nil), func() {}
			},
		},
		{
			name:  "イントロスペクションが失敗した場合、エラーが返ること",
			input: GetPrincipalInput{AuthToken: "token-2"},
			setupMock: func(t *testing.T) (IGetPrincipalUsecase, func()) {
				ctrl := gomock.NewController(t)
				authIntrospector := NewMockAuthIntrospector(ctrl)
				actorRepo := NewMockActorRepository(ctrl)
				authIntrospector.EXPECT().Introspect("token-2").Return(auth.AuthIntrospectionResult{
					IsActive: true,
					Sub:      "sub-2",
					Scopes:   []string{"scope-a"},
				}, errors.New("introspection error"))
				usecase := NewGetPrincipalUsecase(authIntrospector, actorRepo)
				return usecase, ctrl.Finish
			},
		},
		{
			name:  "トークンが無効の場合、エラーが返ること",
			input: GetPrincipalInput{AuthToken: "token-3"},
			setupMock: func(t *testing.T) (IGetPrincipalUsecase, func()) {
				ctrl := gomock.NewController(t)
				authIntrospector := NewMockAuthIntrospector(ctrl)
				actorRepo := NewMockActorRepository(ctrl)
				authIntrospector.EXPECT().Introspect("token-3").Return(auth.AuthIntrospectionResult{
					IsActive: false,
					Sub:      "sub-3",
					Scopes:   []string{"scope-b"},
				}, error(nil))
				usecase := NewGetPrincipalUsecase(authIntrospector, actorRepo)
				return usecase, ctrl.Finish
			},
		},
		{
			name:  "Actorが存在しない場合、エラーが返ること",
			input: GetPrincipalInput{AuthToken: "token-4"},
			setupMock: func(t *testing.T) (IGetPrincipalUsecase, func()) {
				ctrl := gomock.NewController(t)
				authIntrospector := NewMockAuthIntrospector(ctrl)
				actorRepo := NewMockActorRepository(ctrl)
				authIntrospector.EXPECT().Introspect("token-4").Return(auth.AuthIntrospectionResult{
					IsActive: true,
					Sub:      "sub-4",
					Scopes:   []string{"scope-c"},
				}, error(nil))
				actorRepo.EXPECT().FetchBySub("sub-4").Return((*entity.Actor)(nil), error(nil))
				usecase := NewGetPrincipalUsecase(authIntrospector, actorRepo)
				return usecase, ctrl.Finish
			},
		},
		{
			name:  "Actorリポジトリエラーの場合、エラーが返ること",
			input: GetPrincipalInput{AuthToken: "token-5"},
			setupMock: func(t *testing.T) (IGetPrincipalUsecase, func()) {
				ctrl := gomock.NewController(t)
				authIntrospector := NewMockAuthIntrospector(ctrl)
				actorRepo := NewMockActorRepository(ctrl)
				authIntrospector.EXPECT().Introspect("token-5").Return(auth.AuthIntrospectionResult{
					IsActive: true,
					Sub:      "sub-5",
					Scopes:   []string{"scope-d"},
				}, error(nil))
				actorRepo.EXPECT().FetchBySub("sub-5").Return((*entity.Actor)(nil), errors.New("repository error"))
				usecase := NewGetPrincipalUsecase(authIntrospector, actorRepo)
				return usecase, ctrl.Finish
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			usecase, finish := tt.setupMock(t)
			defer finish()
			got, err := usecase.Exec(tt.input)
			require.Error(t, err)
			assert.Equal(t, PrincipalOutput{}, got)
			assert.Nil(t, got.Scopes)
		})
	}
}
