package usecase

import (
	"context"
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

func (m *MockAuthIntrospector) Introspect(ctx context.Context, token string) (auth.AuthIntrospectionResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Introspect", ctx, token)
	ret0 := ret[0].(auth.AuthIntrospectionResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockAuthIntrospectorMockRecorder) Introspect(ctx interface{}, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Introspect", reflect.TypeOf((*MockAuthIntrospector)(nil).Introspect), ctx, token)
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

func (m *MockActorRepository) FetchBySub(ctx context.Context, sub string) (*entity.Actor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchBySub", ctx, sub)
	ret0 := ret[0].(*entity.Actor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockActorRepositoryMockRecorder) FetchBySub(ctx interface{}, sub interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchBySub", reflect.TypeOf((*MockActorRepository)(nil).FetchBySub), ctx, sub)
}

func (m *MockActorRepository) Save(ctx context.Context, actor *entity.Actor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, actor)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockActorRepositoryMockRecorder) Save(ctx interface{}, actor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockActorRepository)(nil).Save), ctx, actor)
}

type MockUuidGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockUuidGeneratorMockRecorder
}

type MockUuidGeneratorMockRecorder struct {
	mock *MockUuidGenerator
}

func NewMockUuidGenerator(ctrl *gomock.Controller) *MockUuidGenerator {
	mock := &MockUuidGenerator{ctrl: ctrl}
	mock.recorder = &MockUuidGeneratorMockRecorder{mock}
	return mock
}

func (m *MockUuidGenerator) EXPECT() *MockUuidGeneratorMockRecorder {
	return m.recorder
}

func (m *MockUuidGenerator) Generate() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate")
	ret0 := ret[0].(string)
	return ret0
}

func (mr *MockUuidGeneratorMockRecorder) Generate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*MockUuidGenerator)(nil).Generate))
}

func TestGetPrincipalUsecase_Exec_Success(t *testing.T) {
	testTable := []struct {
		name      string
		ctx       context.Context
		input     GetPrincipalInput
		expected  PrincipalOutput
		setupMock func(t *testing.T, ctrl *gomock.Controller, ctx context.Context) (*MockAuthIntrospector, *MockActorRepository, *MockUuidGenerator)
	}{
		{
			name:  "トークンが有効でActorが取得できる場合、Principalが取得できること",
			ctx:   context.WithValue(context.Background(), "ctx-key-1", "ctx-1"),
			input: GetPrincipalInput{AuthToken: "token-1"},
			expected: PrincipalOutput{
				ActorCode: "actor-1",
				Scopes:    []string{"scope-a"},
			},
			setupMock: func(t *testing.T, ctrl *gomock.Controller, ctx context.Context) (*MockAuthIntrospector, *MockActorRepository, *MockUuidGenerator) {
				authIntrospector := NewMockAuthIntrospector(ctrl)
				actorRepo := NewMockActorRepository(ctrl)
				uuidGenerator := NewMockUuidGenerator(ctrl)
				authIntrospector.EXPECT().Introspect(ctx, "token-1").Return(auth.AuthIntrospectionResult{
					IsActive: true,
					Sub:      "sub-1",
					Scopes:   []string{"scope-a"},
				}, error(nil))
				actorRepo.EXPECT().FetchBySub(ctx, "sub-1").
					Return(&entity.Actor{Code: "actor-1", Sub: "sub-1"}, error(nil))
				actorRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Times(0)
				uuidGenerator.EXPECT().Generate().Times(0)
				return authIntrospector, actorRepo, uuidGenerator
			},
		},
		{
			name:  "Actorが取得できない場合、Actorを生成してPrincipalが取得できること",
			ctx:   context.WithValue(context.Background(), "ctx-key-2", "ctx-2"),
			input: GetPrincipalInput{AuthToken: "token-2"},
			expected: PrincipalOutput{
				ActorCode: "actor-2",
				Scopes:    []string{"scope-b"},
			},
			setupMock: func(t *testing.T, ctrl *gomock.Controller, ctx context.Context) (*MockAuthIntrospector, *MockActorRepository, *MockUuidGenerator) {
				authIntrospector := NewMockAuthIntrospector(ctrl)
				actorRepo := NewMockActorRepository(ctrl)
				uuidGenerator := NewMockUuidGenerator(ctrl)
				authIntrospector.EXPECT().Introspect(ctx, "token-2").Return(auth.AuthIntrospectionResult{
					IsActive: true,
					Sub:      "sub-2",
					Scopes:   []string{"scope-b"},
				}, error(nil))
				actorRepo.EXPECT().FetchBySub(ctx, "sub-2").
					Return((*entity.Actor)(nil), error(nil))
				uuidGenerator.EXPECT().Generate().Return("actor-2")
				actorRepo.EXPECT().Save(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, actor *entity.Actor) error {
					assert.Equal(t, "actor-2", actor.Code)
					assert.Equal(t, "sub-2", actor.Sub)
					return error(nil)
				})
				return authIntrospector, actorRepo, uuidGenerator
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			authIntrospector, actorRepo, uuidGenerator := tt.setupMock(t, ctrl, tt.ctx)
			usecase := NewGetPrincipalUsecase(authIntrospector, actorRepo, uuidGenerator)

			got, err := usecase.Exec(tt.ctx, tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestGetPrincipalUsecase_Exec_Failure(t *testing.T) {
	testTable := []struct {
		name      string
		ctx       context.Context
		input     GetPrincipalInput
		setupMock func(t *testing.T, ctx context.Context) (IGetPrincipalUsecase, func())
	}{
		{
			name:  "依存未設定の場合、エラーが返ること",
			ctx:   context.WithValue(context.Background(), "ctx-key-3", "ctx-3"),
			input: GetPrincipalInput{AuthToken: "token-1"},
			setupMock: func(t *testing.T, ctx context.Context) (IGetPrincipalUsecase, func()) {
				return NewGetPrincipalUsecase(nil, nil, nil), func() {}
			},
		},
		{
			name:  "イントロスペクションが失敗した場合、エラーが返ること",
			ctx:   context.WithValue(context.Background(), "ctx-key-4", "ctx-4"),
			input: GetPrincipalInput{AuthToken: "token-2"},
			setupMock: func(t *testing.T, ctx context.Context) (IGetPrincipalUsecase, func()) {
				ctrl := gomock.NewController(t)
				authIntrospector := NewMockAuthIntrospector(ctrl)
				actorRepo := NewMockActorRepository(ctrl)
				uuidGenerator := NewMockUuidGenerator(ctrl)
				authIntrospector.EXPECT().Introspect(ctx, "token-2").Return(auth.AuthIntrospectionResult{
					IsActive: true,
					Sub:      "sub-2",
					Scopes:   []string{"scope-a"},
				}, errors.New("introspection error"))
				usecase := NewGetPrincipalUsecase(authIntrospector, actorRepo, uuidGenerator)
				return usecase, ctrl.Finish
			},
		},
		{
			name:  "Actor保存に失敗した場合、エラーが返ること",
			ctx:   context.WithValue(context.Background(), "ctx-key-5", "ctx-5"),
			input: GetPrincipalInput{AuthToken: "token-4"},
			setupMock: func(t *testing.T, ctx context.Context) (IGetPrincipalUsecase, func()) {
				ctrl := gomock.NewController(t)
				authIntrospector := NewMockAuthIntrospector(ctrl)
				actorRepo := NewMockActorRepository(ctrl)
				uuidGenerator := NewMockUuidGenerator(ctrl)
				authIntrospector.EXPECT().Introspect(ctx, "token-4").Return(auth.AuthIntrospectionResult{
					IsActive: true,
					Sub:      "sub-4",
					Scopes:   []string{"scope-c"},
				}, error(nil))
				actorRepo.EXPECT().FetchBySub(ctx, "sub-4").
					Return((*entity.Actor)(nil), error(nil))
				uuidGenerator.EXPECT().Generate().Return("actor-4")
				actorRepo.EXPECT().Save(ctx, gomock.Any()).Return(errors.New("save error"))
				usecase := NewGetPrincipalUsecase(authIntrospector, actorRepo, uuidGenerator)
				return usecase, ctrl.Finish
			},
		},
		{
			name:  "Actorリポジトリエラーの場合、エラーが返ること",
			ctx:   context.WithValue(context.Background(), "ctx-key-6", "ctx-key-6"),
			input: GetPrincipalInput{AuthToken: "token-5"},
			setupMock: func(t *testing.T, ctx context.Context) (IGetPrincipalUsecase, func()) {
				ctrl := gomock.NewController(t)
				authIntrospector := NewMockAuthIntrospector(ctrl)
				actorRepo := NewMockActorRepository(ctrl)
				uuidGenerator := NewMockUuidGenerator(ctrl)
				authIntrospector.EXPECT().Introspect(ctx, "token-5").Return(auth.AuthIntrospectionResult{
					IsActive: true,
					Sub:      "sub-5",
					Scopes:   []string{"scope-d"},
				}, error(nil))
				actorRepo.EXPECT().FetchBySub(ctx, "sub-5").
					Return((*entity.Actor)(nil), errors.New("repository error"))
				usecase := NewGetPrincipalUsecase(authIntrospector, actorRepo, uuidGenerator)
				return usecase, ctrl.Finish
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			usecase, finish := tt.setupMock(t, tt.ctx)
			defer finish()
			got, err := usecase.Exec(tt.ctx, tt.input)
			require.Error(t, err)
			assert.Equal(t, PrincipalOutput{}, got)
			assert.Nil(t, got.Scopes)
		})
	}
}

func TestGetPrincipalUsecase_Exec_TokenInvalid(t *testing.T) {
	testTable := []struct {
		name  string
		ctx   context.Context
		input GetPrincipalInput
	}{
		{
			name:  "トークンが無効の場合、TokenInvalidErrorが返ること",
			ctx:   context.WithValue(context.Background(), "ctx-key-7", "ctx-7"),
			input: GetPrincipalInput{AuthToken: "token-3"},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			authIntrospector := NewMockAuthIntrospector(ctrl)
			actorRepo := NewMockActorRepository(ctrl)
			uuidGenerator := NewMockUuidGenerator(ctrl)
			authIntrospector.EXPECT().Introspect(tt.ctx, "token-3").Return(auth.AuthIntrospectionResult{
				IsActive: false,
				Sub:      "sub-3",
				Scopes:   []string{"scope-b"},
			}, error(nil))
			usecase := NewGetPrincipalUsecase(authIntrospector, actorRepo, uuidGenerator)

			got, err := usecase.Exec(tt.ctx, tt.input)
			require.Error(t, err)
			var tokenInvalidErr *TokenInvalidError
			require.ErrorAs(t, err, &tokenInvalidErr)
			assert.Equal(t, PrincipalOutput{}, got)
			assert.Nil(t, got.Scopes)
		})
	}
}
