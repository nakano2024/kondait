package util

import (
	"kondait-backend/application/util"

	"gorm.io/gorm"
)

type principalFethcer struct {
	db *gorm.DB
}

func (fetcher *principalFethcer) FetchPrincipal(authToken string) (util.Principal, error) {
	return util.Principal{}, nil
}

func NewPrincipalFetcher() util.IPrincipalFetcher {
	return &principalFethcer{}
}

type principalFethcerMock struct{}

func (fetcher *principalFethcerMock) FetchPrincipal(authToken string) (util.Principal, error) {
	return util.Principal{
		UserCode: "11111111-1111-1111-1111-111111111111",
		Scopes: []string{
			"cooking-items.read",
		},
	}, nil
}

func NewPrincipalFetcherMock() util.IPrincipalFetcher {
	return &principalFethcerMock{}
}
