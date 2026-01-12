package auth

type AuthIntrospectionResult struct {
	IsActive bool
	Sub      string
	Scopes   []string
}

type IAuthIntrospector interface {
	Introspect(token string) (AuthIntrospectionResult, error)
}
