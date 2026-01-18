package auth

import (
	"context"

	"kondait-backend/application/auth"
	"kondait-backend/infra/config"
)

type authIntrospector struct {
	config config.Config
}

func (introspector *authIntrospector) Introspect(ctx context.Context, token string) (auth.AuthIntrospectionResult, error) {
	_ = ctx
	return auth.AuthIntrospectionResult{
		IsActive: true,
		Sub:      "",
		Scopes:   []string{},
	}, nil
}

func NewAuthIntrospector() auth.IAuthIntrospector {
	return &authIntrospector{}
}

type authIntrospectorMock struct{}

func (introspector *authIntrospectorMock) Introspect(ctx context.Context, token string) (auth.AuthIntrospectionResult, error) {
	_ = ctx
	return auth.AuthIntrospectionResult{
		IsActive: true,
		Sub:      "mock-sub",
		Scopes: []string{
			"cooking-items.read",
		},
	}, nil
}

func NewAuthIntrospectorMock() auth.IAuthIntrospector {
	return &authIntrospectorMock{}
}
