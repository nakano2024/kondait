package auth

import (
	"context"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"kondait-backend/application/auth"
	"kondait-backend/infra/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type MockRoundTripper struct {
	ctrl     *gomock.Controller
	recorder *MockRoundTripperMockRecorder
}

type MockRoundTripperMockRecorder struct {
	mock *MockRoundTripper
}

func NewMockRoundTripper(ctrl *gomock.Controller) *MockRoundTripper {
	mock := &MockRoundTripper{ctrl: ctrl}
	mock.recorder = &MockRoundTripperMockRecorder{mock}
	return mock
}

func (m *MockRoundTripper) EXPECT() *MockRoundTripperMockRecorder {
	return m.recorder
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RoundTrip", req)
	ret0 := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockRoundTripperMockRecorder) RoundTrip(req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RoundTrip", reflect.TypeOf((*MockRoundTripper)(nil).RoundTrip), req)
}

func TestAuthIntrospector_Introspect_Success(t *testing.T) {
	testTable := []struct {
		name      string
		ctx       context.Context
		cfg       config.Config
		token     string
		setupMock func(t *testing.T, ctrl *gomock.Controller) *http.Client
		expected  auth.AuthIntrospectionResult
	}{
		{
			name:  "activeがtrueでスコープがある場合、結果が取得できること",
			ctx:   context.WithValue(context.Background(), "ctx-key-1", "ctx-1"),
			cfg:   config.Config{AuthServerUrl: "https://keycloak.example.com/realms/myrealm", ClientId: "client-1", ClientSecret: "secret-1"},
			token: "token-1",
			setupMock: func(t *testing.T, ctrl *gomock.Controller) *http.Client {
				transport := NewMockRoundTripper(ctrl)
				transport.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
					require.Equal(t, http.MethodPost, req.Method)
					assert.Equal(t, "https://keycloak.example.com/realms/myrealm/protocol/openid-connect/token/introspect", req.URL.String())
					assert.Equal(t, "application/x-www-form-urlencoded", req.Header.Get("Content-Type"))
					assert.Equal(t, "Basic Y2xpZW50LTE6c2VjcmV0LTE=", req.Header.Get("Authorization"))
					bodyBytes, err := io.ReadAll(req.Body)
					require.NoError(t, err)
					assert.Equal(t, "token=token-1", string(bodyBytes))
					resBody := `{"active":true,"sub":"user-id","scope":"openid profile email"}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(resBody)),
					}, error(nil)
				})
				return &http.Client{Transport: transport}
			},
			expected: auth.AuthIntrospectionResult{
				IsActive: true,
				Sub:      "user-id",
				Scopes:   []string{"openid", "profile", "email"},
			},
		},
		{
			name:  "activeがfalseでスコープが空の場合、空のスコープが取得できること",
			ctx:   context.WithValue(context.Background(), "ctx-key-2", "ctx-2"),
			cfg:   config.Config{AuthServerUrl: "https://keycloak.example.com/realms/myrealm/", ClientId: "client-2", ClientSecret: "secret-2"},
			token: "token-2",
			setupMock: func(t *testing.T, ctrl *gomock.Controller) *http.Client {
				transport := NewMockRoundTripper(ctrl)
				transport.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
					require.Equal(t, http.MethodPost, req.Method)
					assert.Equal(t, "https://keycloak.example.com/realms/myrealm/protocol/openid-connect/token/introspect", req.URL.String())
					assert.Equal(t, "Basic Y2xpZW50LTI6c2VjcmV0LTI=", req.Header.Get("Authorization"))
					resBody := `{"active":false,"sub":"user-id-2","scope":""}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(resBody)),
					}, error(nil)
				})
				return &http.Client{Transport: transport}
			},
			expected: auth.AuthIntrospectionResult{
				IsActive: false,
				Sub:      "user-id-2",
				Scopes:   []string{},
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			httpClient := tt.setupMock(t, ctrl)
			introspector := NewAuthIntrospector(tt.cfg, httpClient)

			got, err := introspector.Introspect(tt.ctx, tt.token)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestAuthIntrospector_Introspect_Failure(t *testing.T) {
	testTable := []struct {
		name        string
		ctx         context.Context
		cfg         config.Config
		token       string
		setupMock   func(t *testing.T, ctrl *gomock.Controller) *http.Client
		expected    auth.AuthIntrospectionResult
		expectedErr string
	}{
		{
			name:  "レスポンスが400の場合、エラーが返ること",
			ctx:   context.WithValue(context.Background(), "ctx-key-3", "ctx-3"),
			cfg:   config.Config{AuthServerUrl: "https://keycloak.example.com/realms/myrealm", ClientId: "client-3", ClientSecret: "secret-3"},
			token: "token-3",
			setupMock: func(t *testing.T, ctrl *gomock.Controller) *http.Client {
				transport := NewMockRoundTripper(ctrl)
				transport.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
					assert.Equal(t, "Basic Y2xpZW50LTM6c2VjcmV0LTM=", req.Header.Get("Authorization"))
					resBody := "bad request"
					return &http.Response{
						StatusCode: http.StatusBadRequest,
						Body:       io.NopCloser(strings.NewReader(resBody)),
					}, error(nil)
				})
				return &http.Client{Transport: transport}
			},
			expected:    auth.AuthIntrospectionResult{},
			expectedErr: "bad request",
		},
		{
			name:  "レスポンスが401の場合、エラーが返ること",
			ctx:   context.WithValue(context.Background(), "ctx-key-4", "ctx-4"),
			cfg:   config.Config{AuthServerUrl: "https://keycloak.example.com/realms/myrealm", ClientId: "client-4", ClientSecret: "secret-4"},
			token: "token-4",
			setupMock: func(t *testing.T, ctrl *gomock.Controller) *http.Client {
				transport := NewMockRoundTripper(ctrl)
				transport.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
					assert.Equal(t, "Basic Y2xpZW50LTQ6c2VjcmV0LTQ=", req.Header.Get("Authorization"))
					resBody := "unauthorized"
					return &http.Response{
						StatusCode: http.StatusUnauthorized,
						Body:       io.NopCloser(strings.NewReader(resBody)),
					}, error(nil)
				})
				return &http.Client{Transport: transport}
			},
			expected:    auth.AuthIntrospectionResult{},
			expectedErr: "unauthorized",
		},
		{
			name:  "http clientがnilの場合、エラーが返ること",
			ctx:   context.WithValue(context.Background(), "ctx-key-5", "ctx-5"),
			cfg:   config.Config{AuthServerUrl: "https://keycloak.example.com/realms/myrealm"},
			token: "token-5",
			setupMock: func(t *testing.T, ctrl *gomock.Controller) *http.Client {
				return nil
			},
			expected:    auth.AuthIntrospectionResult{},
			expectedErr: "authIntrospector: http client is nil",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			httpClient := tt.setupMock(t, ctrl)
			introspector := NewAuthIntrospector(tt.cfg, httpClient)

			got, err := introspector.Introspect(tt.ctx, tt.token)
			require.Error(t, err)
			assert.Equal(t, tt.expected, got)
			assert.Equal(t, tt.expectedErr, err.Error())
		})
	}
}
