package usecase

import (
	"context"
	"errors"

	"kondait-backend/application/auth"
	"kondait-backend/application/util"
	"kondait-backend/domain/entity"
	"kondait-backend/domain/repository"
)

type PrincipalOutput struct {
	ActorCode string
	Scopes    []string
}

type TokenInvalidError struct{}

func (err *TokenInvalidError) Error() string {
	return "token invalid"
}

type GetPrincipalInput struct {
	AuthToken string
}

type IGetPrincipalUsecase interface {
	Exec(ctx context.Context, input GetPrincipalInput) (PrincipalOutput, error)
}

type getPrincipalUsecase struct {
	authIntrospector auth.IAuthIntrospector
	actorRepo        repository.IActorRepository
	uuidGenerator    util.UuidGenerator
}

func NewGetPrincipalUsecase(
	authIntrospector auth.IAuthIntrospector,
	actorRepo repository.IActorRepository,
	uuidGenerator util.UuidGenerator,
) IGetPrincipalUsecase {
	return &getPrincipalUsecase{
		authIntrospector: authIntrospector,
		actorRepo:        actorRepo,
		uuidGenerator:    uuidGenerator,
	}
}

func (usecase *getPrincipalUsecase) Exec(ctx context.Context, input GetPrincipalInput) (PrincipalOutput, error) {
	if usecase.authIntrospector == nil || usecase.actorRepo == nil || usecase.uuidGenerator == nil {
		return PrincipalOutput{}, errors.New("getPrincipalUsecase: dependency not set")
	}

	introspection, err := usecase.authIntrospector.Introspect(ctx, input.AuthToken)
	if err != nil {
		return PrincipalOutput{}, err
	}
	if !introspection.IsActive {
		return PrincipalOutput{}, &TokenInvalidError{}
	}

	actor, err := usecase.actorRepo.FetchBySub(ctx, introspection.Sub)
	if err != nil {
		return PrincipalOutput{}, err
	}
	if actor != nil {
		return PrincipalOutput{
			ActorCode: actor.Code,
			Scopes:    introspection.Scopes,
		}, nil
	}

	actor = entity.NewActor(usecase.uuidGenerator.Generate(), introspection.Sub)
	if err := usecase.actorRepo.Save(ctx, actor); err != nil {
		return PrincipalOutput{}, err
	}
	return PrincipalOutput{
		ActorCode: actor.Code,
		Scopes:    introspection.Scopes,
	}, nil
}
