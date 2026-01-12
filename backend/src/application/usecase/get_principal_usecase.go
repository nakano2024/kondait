package usecase

import "kondait-backend/application/util"

type PrincipalOutput struct {
	UserCode string
	Scopes   []string
}

type GetPrincipalInput struct {
	AuthToken string
}

type IGetPrincipalUsecase interface {
	Exec(input GetPrincipalInput) (PrincipalOutput, error)
}

type getPrincipalUsecase struct {
	principalFecher util.IPrincipalFetcher
}

func NewGetPrincipalUsecase(principalFetcher util.IPrincipalFetcher) IGetPrincipalUsecase {
	return &getPrincipalUsecase{
		principalFecher: principalFetcher,
	}
}

func (usecase *getPrincipalUsecase) Exec(input GetPrincipalInput) (PrincipalOutput, error) {
	principal, err := usecase.principalFecher.FetchPrincipal(input.AuthToken)
	if err != nil {
		return PrincipalOutput{}, err
	}
	return PrincipalOutput{
		UserCode: principal.UserCode,
		Scopes:   principal.Scopes,
	}, nil
}
