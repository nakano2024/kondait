package util

type Principal struct {
	UserCode string
	Scopes   []string
}

type IPrincipalFetcher interface {
	FetchPrincipal(authToken string) (Principal, error)
}
