package auth

import "context"

type AuthIntrospectionResult struct {
	IsActive bool
	Sub      string
	Scopes   []string
}

type IAuthIntrospector interface {
	Introspect(ctx context.Context, token string) (AuthIntrospectionResult, error)
}
