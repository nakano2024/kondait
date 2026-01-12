package auth

import "kondait-backend/application/auth"

type authIntrospector struct{}

func (introspector *authIntrospector) Introspect(token string) (auth.AuthIntrospectionResult, error) {
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

func (introspector *authIntrospectorMock) Introspect(token string) (auth.AuthIntrospectionResult, error) {
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
