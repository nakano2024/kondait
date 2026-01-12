package usecase

import (
	"errors"

	"kondait-backend/application/auth"
	"kondait-backend/domain/repository"
)

type PrincipalOutput struct {
	ActorCode string
	Scopes    []string
}

type GetPrincipalInput struct {
	AuthToken string
}

type IGetPrincipalUsecase interface {
	Exec(input GetPrincipalInput) (PrincipalOutput, error)
}

type getPrincipalUsecase struct {
	authIntrospector auth.IAuthIntrospector
	actorRepo        repository.IActorRepository
}

func NewGetPrincipalUsecase(authIntrospector auth.IAuthIntrospector, actorRepo repository.IActorRepository) IGetPrincipalUsecase {
	return &getPrincipalUsecase{
		authIntrospector: authIntrospector,
		actorRepo:        actorRepo,
	}
}

func (usecase *getPrincipalUsecase) Exec(input GetPrincipalInput) (PrincipalOutput, error) {
	if usecase.authIntrospector == nil || usecase.actorRepo == nil {
		return PrincipalOutput{}, errors.New("getPrincipalUsecase: dependency not set")
	}

	introspection, err := usecase.authIntrospector.Introspect(input.AuthToken)
	if err != nil {
		return PrincipalOutput{}, err
	}
	if !introspection.IsActive {
		return PrincipalOutput{}, errors.New("inactive token")
	}

	actor, err := usecase.actorRepo.FetchBySub(introspection.Sub)
	if err != nil {
		return PrincipalOutput{}, err
	}
	if actor == nil {
		return PrincipalOutput{}, errors.New("actor not found")
	}
	return PrincipalOutput{
		ActorCode: actor.Code,
		Scopes:    introspection.Scopes,
	}, nil
}
